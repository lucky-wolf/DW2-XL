package algorithm

import (
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
	"regexp"
	"strings"
)

// TODO: we need to increase interceptor missiles (PD) speed so they can catch other missiles & fighters!
// maybe also slow down their ROF / increase dmg/hit (say by factor of 2x)

func FighterWeaponsAndPD(folder string) (err error) {

	log.Println("All strikecraft weapons and PD weapons will be scaled to ship components")

	// load all component definition files
	j, err := loadJobFor(folder, "ComponentDefinitions*")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyFighterWeaponsAndPD()
	if err != nil {
		return
	}

	// save them all
	j.save()

	return
}

func (j *job) applyFighterWeaponsAndPD() (err error) {

	for _, f := range j.xfiles {

		statistics := &f.stats

		// the root will result in a single ArrayOf[RootObjectType]
		for _, e := range f.root.Elements.Elements() {

			err = assertIs(e, "ArrayOfComponentDefinition")
			if err != nil {
				return
			}

			for _, e := range e.Elements() {

				// each of these is a componentdefinition
				err = assertIs(e, "ComponentDefinition")
				if err != nil {
					return
				}

				// only weapons...
				if !e.HasPrefix("Category", "Weapon") {
					continue
				}

				// for fighters...
				targetName := e.Child("Name").StringValue()
				if !strings.HasSuffix(targetName, "[Ftr]") && !strings.HasSuffix(targetName, "[PD]") {
					continue
				}

				// distinguish if this component is only for strike craft
				isFighterOnly := e.Has("IsFighterOnly", "true")

				// find the corresponding small weapon by name
				sourceName := getSourceOfPD(targetName, isFighterOnly)
				sourceDefinition, _ := j.find("Name", sourceName)
				if sourceDefinition == nil {
					log.Printf("Missing: %s (for %s)", sourceName, targetName)
					continue
				}

				// debug
				if !Quiet {
					log.Printf("%s from %s\n", targetName, sourceName)
				}

				// copy (and scale fighter) resource requirements
				if isFighterOnly {
					err = e.CopyAndVisitByTag("ResourcesRequired", sourceDefinition, func(e *xmltree.XMLElement) error { e.Child("Amount").ScaleBy(0.25); return nil })
				} else {
					err = e.CopyByTag("ResourcesRequired", sourceDefinition)
				}
				if err != nil {
					log.Printf("%s: %s from %s", err, targetName, sourceName)
				}

				// copy component stats
				err = e.CopyByTag("Values", sourceDefinition)
				if err != nil {
					log.Printf("%s: %s from %s", err, targetName, sourceName)
				}

				// now that we have our own copy of the component stats (same number of levels too)
				// we can update each of those to scale for [Ftr] version
				err = scaleFighterOrPDValues(e, isFighterOnly)
				if err != nil {
					log.Printf("%s: %s from %s", err, targetName, sourceName)
				}

				statistics.changed++
				statistics.elements++
				statistics.objects++
			}
		}
	}
	err = nil
	return
}

// scale all component stats from the copy from our [S] source weapon
func scaleFighterOrPDValues(e *xmltree.XMLElement, isFighterOnly bool) (err error) {

	for _, e := range e.Child("Values").Elements() {

		// every element should be a component bay
		err = assertIs(e, "ComponentStats")
		if err != nil {
			return
		}

		// "flatten" source volleys to 1 per shot but at 1/x fire rate (same dps, but distributed instead of burste firing)
		if va := e.Child("WeaponVolleyAmount").NumericValue(); va != 1 {
			e.Child("WeaponFireRate").ScaleBy(1.0 / va)
			e.Child("WeaponVolleyAmount").SetValue(1)
		}
		e.Child("WeaponVolleyFireRate").SetValue(0)

		scaleWeapon := func(rof float64, dmg float64) {
			// scale by our source weapon values
			e.Child("WeaponFireRate").ScaleBy(1 / rof)
			e.Child("WeaponRawDamage").ScaleBy(dmg)
			e.Child("WeaponEnergyPerShot").ScaleBy(dmg)

			// range is 1/3 + 50% more rapid fall-off
			e.Child("WeaponRange").ScaleBy(0.3333333333)
			e.Child("WeaponDamageFalloffRatio").ScaleBy(1.5)

			// fighter & PD weapons generically get a +10% targeting across the board (very short range = enhanced accuracy)
			e.Child("ComponentTargetingBonus").AdjustValue(0.1)
		}

		scaleIntercept := func(rof float64, dmg float64) {
			// scale by our standard mode values
			e.ScaleChildToSiblingBy("WeaponInterceptFireRate", "WeaponFireRate", 1/rof)
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

		// scale standard fire relative to our source weapon
		// 4 x .375 = 1.5x total output
		scaleWeapon(4, .375)

		if isFighterOnly {
			// for fighters scale intercept by...
			// previously we were at a net 4x dmg vs. fighters and 8x dmg vs. seeking
			// (8x rof, 1/2x dmg vs. fighters, and 8x rof, 1x dmg vs. seeking)
			// now we're at 5x4 = 20x rof vs. standard (was 32x)
			// and 5 x .4 = 200% total damage output compard to base, which is 1.5 normal = 300% total vs. standard weapon
			scaleIntercept(5, .4)
		} else {
			// PD is 2 * 2 = 4x as effective as a ftr
			// the very high rof means we should get cool visuals (blasters are now approx 4/s)
			// note: we might want to break this out by weapon type (super high for kinetic & blaster, less so for beams & missiles)
			scaleIntercept(10, .8)
		}

		// fighters and PD never do bombard damage
		for _, e := range e.Matching(regexp.MustCompile("WeaponBombard.*")) {
			e.SetValue(0)
		}

		// some things just don't apply to fighters
		if isFighterOnly {
			e.Child("CrewRequirement").SetValue(0)
			e.Child("StaticEnergyUsed").SetValue(0)
		}
	}

	return
}

func getSourceOfPD(targetName string, isFighterOnly bool) (sourceName string) {

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
		// hive missile battery
		// reinforcing swarm battery

		if isFighterOnly {
			sourceName = targetName[:len(targetName)-len(" [Ftr]")] + " [S]"
		} else {
			sourceName = targetName[:len(targetName)-len(" [PD]")] + " [S]"
		}
	}

	return
}
