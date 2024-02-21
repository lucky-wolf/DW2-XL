package algorithm

import (
	"fmt"
	"log"
	"lucky-wolf/DW2-XL/code/etc"
	"lucky-wolf/DW2-XL/code/xmltree"
	"math"
	"regexp"
)

// note: see globals.go for some standard level functions

type ComponentName = string
type WeaponFamilyName = string

// maps the name of a component attribute to a level function to derive its value
type ComponentStats = map[AttributeName]LevelFunc

// like component stats, but returns a simple level func which is more flexible (not always a number)
type FlexComponentStats = map[AttributeName]SimpleLevelFunc

// we take a "level" and return a numeric value
// typically we use this for component stats
type LevelFunc = func(level int) float64

// we return a simple value which can be applied to an XMLValue
type SimpleLevelFunc = func(level int) xmltree.SimpleValue

type ComponentData struct {
	values         map[AttributeName]xmltree.SimpleValue // these are applied to the base component
	minLevel       int
	maxLevel       int
	componentStats ComponentStats
	derivatives    []string // what other components are derived from this? (i.e. name of [Ftr] or [PD] components)
}

type ComponentIs struct {
	fighter bool
	weapon  bool
	pd      bool
	size    int
}

func MakeFixedLevelFunc(basis float64) LevelFunc {
	return func(level int) float64 { return basis }
}

func MakeLinearLevelFunc(basis float64, add float64) LevelFunc {
	return func(level int) float64 { return basis + (add * float64(level)) }
}

func MakeScalingLevelFunc(basis float64, scale float64) LevelFunc {
	return func(level int) float64 { return basis * (1.0 + scale*float64(level)) }
}

func MakeExpLevelFunc(basis float64, scale float64) LevelFunc {
	return func(level int) float64 { return basis * math.Pow(1.0+scale, float64(level)) }
}

func MakeScaledFuncLevelFunc(scale float64, levelfunc LevelFunc) LevelFunc {
	return func(level int) float64 { return scale * levelfunc(level) }
}

func MakeOffsetFuncLevelFunc(offset int, levelfunc LevelFunc) LevelFunc {
	return func(level int) float64 { return levelfunc(level + offset) }
}

func MakeIntegerLevelFunc(levelfunc LevelFunc) LevelFunc {
	return func(level int) float64 { return math.Round(levelfunc(level)) }
}

func ComposeComponentStats(fields ComponentStats, more ...ComponentStats) (result ComponentStats) {
	// clone the base stats
	result = ComponentStats{}
	for k, v := range fields {
		result[k] = v
	}
	// merge in each additional stats
	for _, more := range more {
		for k, v := range more {
			result[k] = v
		}
	}
	return
}

func GetComponentIsms(e *xmltree.XMLElement) (is ComponentIs) {
	is.fighter = e.HasChildWithValue("IsFighterOnly", "true")
	is.weapon = e.HasPrefix("Category", "Weapon")
	is.pd = !is.fighter && e.HasChildWithValue("Category", "WeaponIntercept")
	is.size = e.Child("Size").IntValue()

	return
}

// apply stats for each component
func (j *Job) ApplyComponentAll(components map[string]ComponentData) (err error) {
	for k, v := range components {
		err = j.ApplyComponent(k, v)
		if err != nil {
			return
		}
	}
	return
}

// applies stats for given component
func (j *Job) ApplyComponent(name string, data ComponentData) (err error) {

	// find this component definition
	e, f := j.FindElement("Name", name)
	if e == nil {
		return fmt.Errorf("%s not found", name)
	}
	statistics := &f.stats

	// apply any top-level values
	for key, sv := range data.values {
		if c := e.Child(key); c != nil {
			c.XMLValue.SetString(sv.String())
		}
	}

	// ensure we have correct number of component stats to update
	err = e.Child("Values").SetElementCountByCopyingFirstElementAsNeeded(1 + data.maxLevel - data.minLevel)
	if err != nil {
		return
	}

	// fill in the data from our data tables
	stats := e.Child("Values").Elements()
	for i, e := range stats {
		for key, f := range data.componentStats {
			e.Child(key).SetValue(f(data.minLevel + i))
		}
		statistics.elements++
		statistics.changed++
	}
	statistics.objects++

	// scale to fighter if required
	for _, name := range data.derivatives {
		err = j.DeriveFromComponentByName(e, name)
	}

	return
}

func (j *Job) DeriveFromComponentByName(source *xmltree.XMLElement, name string) (err error) {

	// figure out our target name
	e, f := j.FindElement("Name", name)
	if e == nil {
		err = fmt.Errorf("%s not found", name)
		return
	}

	// scale it
	err = j.DeriveFromComponent(f, source, e)

	return
}

func (j *Job) DeriveFromComponent(file *XFile, source *xmltree.XMLElement, e *xmltree.XMLElement) (err error) {

	statistics := &file.stats

	// distinguish what kind of target component we're dealing with
	is := GetComponentIsms(e)

	// scale component size
	if is.fighter {
		sourceSizeElem := source.Child("Size")
		if sourceSizeElem != nil {

			// attempt to extract the int64 value of the size element
			var value int64
			value, err = sourceSizeElem.GetInt64Value()
			if err != nil {
				return
			}

			// size must be an integer
			switch {
			case is.weapon:
				// 11 -> 5
				// 13 -> 10
				if value < 13 {
					e.Child("Size").SetValue(5)
					is.size = 5
				} else {
					e.Child("Size").SetValue(10)
					is.size = 10
				}
			default:
				// non-weapons = 33% size
				// (must be an integer value, rounded up, but cannot exceed size 10 for any slot, and reactors would exceed it)
				is.size = min(10, etc.MulDivRoundUp(int(value), 1, 3))

				e.Child("Size").SetValue(is.size)
			}

		}
	}

	// copy (and scale fighter) resource requirements
	if is.fighter {
		err = e.CopyAndVisitByTag("ResourcesRequired", source, func(e *xmltree.XMLElement) error { e.Child("Amount").ScaleBy(0.25); return nil })
	} else {
		err = e.CopyByTag("ResourcesRequired", source)
	}
	if err != nil {
		log.Println(err)
	}

	// copy component stats
	err = e.CopyByTag("Values", source)
	if err != nil {
		log.Println(err)
	}

	// now that we have our own copy of the component stats (same number of levels too)
	// we can update each of those to scale for [Ftr] version
	for _, e := range e.Child("Values").Elements() {

		// every element should be a component bay
		err = assertIs(e, "ComponentStats")
		if err != nil {
			return
		}

		if is.weapon {

			// "flatten" source volleys to 1 per shot but at 1/x fire rate (same dps, but distributed instead of burst firing)
			if va := e.Child("WeaponVolleyAmount").NumericValue(); va != 1 {
				e.Child("WeaponFireRate").ScaleBy(1.0 / va)
				e.Child("WeaponVolleyAmount").SetValue(1)
			}
			e.Child("WeaponVolleyFireRate").SetString("0")

			// scale standard fire relative to our source weapon
			err = ScaleFtrOrPDMainWeaponValues(e, is)
			if err != nil {
				return
			}

			// scale intercept function
			err = ScaleFtrOrPDInterceptValues(e, is)
			if err != nil {
				return
			}

			// fighters and PD never do bombard damage
			for _, e := range e.Matching(regexp.MustCompile("WeaponBombard.*")) {
				e.SetString("0")
			}
		}

		// scale down the ion defenses and offenses
		err = ScaleFtrOrPDIonValues(e, is)
		if err != nil {
			return
		}

		if is.fighter {
			// fighters never have crew requirements
			e.Child("CrewRequirement").SetString("0")

			// scale down static energy use
			e.Child("StaticEnergyUsed").ScaleBy(0.5)

			// scale armor values
			e.Child("ArmorBlastRating").ScaleBy(0.2)
			e.Child("ArmorReactiveRating").ScaleBy(0.2)

			// scale engine values
			e.Child("EngineMainCruiseThrust").ScaleBy(0.5)
			e.Child("EngineMainCruiseThrustEnergyUsage").ScaleBy(0.5)

			e.Child("EngineMainMaximumThrust").ScaleBy(0.6)
			e.Child("EngineMainMaximumThrustEnergyUsage").ScaleBy(0.6)

			e.Child("EngineVectoringThrust").ScaleBy(0.25)
			e.Child("EngineVectoringEnergyUsage").ScaleBy(0.25)

			// scale reactor values
			// note: seems like too much for early tech, but quickly isn't
			const reactorFactor = 0.25
			e.Child("ReactorEnergyOutputPerSecond").ScaleBy(reactorFactor)
			e.Child("ReactorEnergyStorageCapacity").ScaleBy(reactorFactor)
			e.Child("ReactorFuelUnitsForFullCharge").ScaleBy(reactorFactor)
			if value, err := e.Child("ReactorFuelUnitsForFullCharge").GetNumericValue(); err == nil {
				// set the fuel units to be enough for 10 recharges
				e.Child("FuelStorageCapacity").SetValue(value * 100)
			}

			// scale shield values
			e.Child("ShieldRechargeRate").ScaleBy(0.2)
			if c := e.Child("ShieldRechargeEnergyUsage"); c != nil {
				// warn: Golang claims it should be okay to call a nil -- yet windows explodes
				c.ScaleBy(0.5)
			}
			// e.Child("ShieldRechargeEnergyUsage").ScaleBy(0.2)
			e.Child("ShieldResistance").ScaleBy(0.2)
			e.Child("ShieldStrength").ScaleBy(0.2)
		}

		statistics.changed++
		statistics.elements++
	}

	statistics.objects++

	return
}

func ScaleFtrOrPDIonValues(e *xmltree.XMLElement, is ComponentIs) (err error) {

	// allow strike craft full ion defenses
	// if is.fighter {
	// e.Child("ComponentIonDefense").ScaleBy(IonFtrPDScaleFactor)
	// e.Child("IonDamageDefense").ScaleBy(IonFtrPDScaleFactor)
	// }

	// but they aren't just a ton of ion small weapons against ships (that would be OP I believe)
	if is.weapon {
		e.Child("WeaponIonEngineDamage").ScaleBy(IonFtrPDScaleFactor)
		e.Child("WeaponIonHyperDriveDamage").ScaleBy(IonFtrPDScaleFactor)
		e.Child("WeaponIonSensorDamage").ScaleBy(IonFtrPDScaleFactor)
		e.Child("WeaponIonShieldDamage").ScaleBy(IonFtrPDScaleFactor)
		e.Child("WeaponIonWeaponDamage").ScaleBy(IonFtrPDScaleFactor)
		e.Child("WeaponIonGeneralDamage").ScaleBy(IonFtrPDScaleFactor)
	}

	return
}

func FtrOrPDMainWeaponScaling(is ComponentIs) (rof float64, dmg float64) {

	if is.fighter && is.size > 5 {
		// make bombers essentially 1/2 the rof as compared to fighter intercept weapons
		// note: we could customize the scaling depending on e.Child("Family")
		// warn: but we cannot use that to determine "bomber" weapon or not (as we want to add bomber beams and gravitic or etc.)
		// 2 x .5 = 1x total output compared to source weapon (more highly negated by DR and the like)
		return 2, .5
	}

	// 4 x .375 = 1.5x total output
	return 4, .375
}

func FtrOrPDInterceptScaling(is ComponentIs) (rof float64, dmg float64) {

	// WARN: this is relative to already being scaled by FtrOrPDMainWeaponScaling()

	switch {
	case is.fighter && is.weapon:
		// for fighters scale intercept by...
		// previously we were at a net 4x dmg vs. fighters and 8x dmg vs. seeking
		// (8x rof, 1/2x dmg vs. fighters, and 8x rof, 1x dmg vs. seeking)
		// now we're at 5x4 = 20x rof vs. standard (was 32x)
		// and 5 x .4 = 200% total damage output compard to base, which is 1.5 normal = 300% total vs. standard weapon
		rof = 5
		dmg = .4

	case is.pd:

		switch is.size {
		case 13:
			// seeking based PD
			// we don't want to ramp up the fire rate all that much at all
			// we already flatten the fire rate and then 4x it (see FtrOrPDMainWeaponScaling)
			rof, dmg = FtrOrPDMainWeaponScaling(is)
			rof = 2

			// subtle: this should get us to 1/2 dmg of base missile type vs. ftr, and 1 dmg vs. missiles or torpedoes
			dmg = .5 / dmg

		default:
			// PD is 2 * 2 = 4x as effective as a ftr
			// the very high rof means we should get cool visuals (blasters are now approx 4/s)
			rof = 10
			dmg = .8
		}
	}

	return
}

func ScaleFtrOrPDMainWeaponValues(e *xmltree.XMLElement, is ComponentIs) (err error) {

	if is.weapon {
		// get appropriate scaling factors
		rof, dmg := FtrOrPDMainWeaponScaling(is)
		rof = 1 / rof

		// scale by our source weapon values
		e.Child("WeaponFireRate").ScaleBy(rof)
		e.Child("WeaponRawDamage").ScaleBy(dmg)
		e.Child("WeaponEnergyPerShot").ScaleBy(dmg)

		// range is 1/3 + 50% more rapid fall-off
		e.Child("WeaponRange").ScaleBy(0.3333333333)
		e.Child("WeaponDamageFalloffRatio").ScaleBy(1.5)

		// fighter & PD weapons generically get a +10% targeting across the board (very short range = enhanced accuracy)
		e.Child("ComponentTargetingBonus").AdjustValue(0.1)
	}

	return
}

func ScaleFtrOrPDInterceptValues(e *xmltree.XMLElement, is ComponentIs) (err error) {

	if is.weapon {

		// bombers will no longer have any intercept function at all
		// use interceptors (fighters) for that!
		if is.fighter && is.size == 10 {
			ZeroInterceptValues(e)
			return
		}

		// get appropriate scaling factors
		rof, dmg := FtrOrPDInterceptScaling(is)
		rof = 1 / rof

		// scale by our standard mode values
		e.ScaleChildToSiblingBy("WeaponInterceptFireRate", "WeaponFireRate", rof)
		e.ScaleChildToSiblingBy("WeaponInterceptDamageFighter", "WeaponRawDamage", dmg)
		e.ScaleChildToSiblingBy("WeaponInterceptDamageSeeking", "WeaponRawDamage", 2*dmg)
		e.ScaleChildToSiblingBy("WeaponInterceptEnergyPerShot", "WeaponEnergyPerShot", dmg)

		// currently we simply always set intercept range == base range for this weapon
		e.SetChildToSibling("WeaponInterceptRange", "WeaponRange")

		// PD must actually hit for it to be useful!
		e.SetChildToSibling("WeaponInterceptComponentTargetingBonus", "ComponentTargetingBonus")
		e.Child("WeaponInterceptComponentTargetingBonus").AdjustValue(0.1)

		// because the dw2 team is incredibly foolish, we have no direct way to know if a weapon is Ion or not
		// so, we'll look for Ion damage attribute and base it on being non-zero there
		// note: WeaponIonGeneralDamage is often zero in vanilla, but we've made it align with all other WeaponIon*Damage values in XL
		if e.Child("WeaponIonGeneralDamage").StringValue() != "0" {
			e.Child("WeaponInterceptIonDamageRatio").SetString("1")
		}
	}

	return
}

func ZeroInterceptValues(e *xmltree.XMLElement) (err error) {

	// scale by our standard mode values
	e.Child("WeaponInterceptFireRate").SetString("0")
	e.Child("WeaponInterceptDamageFighter").SetString("0")
	e.Child("WeaponInterceptDamageSeeking").SetString("0")
	e.Child("WeaponInterceptEnergyPerShot").SetString("0")

	// currently we simply always set intercept range == base range for this weapon
	e.Child("WeaponInterceptRange").SetString("0")

	// PD must actually hit for it to be useful!
	e.Child("WeaponInterceptComponentTargetingBonus").SetString("0")

	// not ionic
	e.Child("WeaponIonGeneralDamage").SetString("0")

	return
}

func GetFighterOrPointDefenseSourceName(targetName string, is ComponentIs) (sourceName string) {

	// find the corresponding small weapon by name
	// PD in particular uses asymmetric sources
	switch targetName {
	case "Buckler Repeating Blaster [PD]":
		sourceName = "Maxos Blaster [S]"
	case "Guardian Defense Grid [PD]":
		sourceName = "Omega Beam [S]"
	case "Maelstrom Defender [PD]":
		sourceName = "Titan Blaster [S]"
	case "Point Defense Cannon [PD]":
		sourceName = "Rail Gun [S]"
	case "Sentinel Multi-Beam Defense [PD]":
		sourceName = "Thuon Beam [S]"
	case "Interceptor Missile [PD]":
		sourceName = "Concussion Missile [S]"
	case "Aegis Missile Battery [PD]":
		sourceName = "Lightning Missile [S]"
	default:
		// simply use the component name [S] as our source component

		// ion cannon
		// ion rapid pulse array
		// impact assault blaster
		// terminator autocannon
		// swarm missile
		// reinforcing swarm battery

		if is.fighter {
			sourceName = targetName[:len(targetName)-len(" [Ftr]")] + " [S]"
		} else {
			sourceName = targetName[:len(targetName)-len(" [PD]")] + " [S]"
		}
	}

	return
}
