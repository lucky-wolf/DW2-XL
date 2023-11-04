package algorithm

import (
	"fmt"
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
	"os"
	"path/filepath"
	"regexp"
)

var (
	Quiet bool
)

type Statistics struct {
	objects  int
	elements int
	changed  int
}

type LevelFunc = func(level int) float64
type ValuesTable = map[string]LevelFunc

func ExtendValuesTable(fields, more ValuesTable) (result ValuesTable) {
	result = ValuesTable{}
	for k, v := range fields {
		result[k] = v
	}
	for k, v := range more {
		result[k] = v
	}
	return
}

type ComponentData struct {
	// todo: could we ever look this up viz research tree for first occurrence there?
	//       that's just column in which it is listed for each level...
	scaleTo     []string // what other components are copies / scaled to this thing
	minLevel    int
	maxLevel    int
	fieldValues ValuesTable
}

const IonFtrPDScaleFactor = 0.5

func (statistics *Statistics) For(filename string) string {
	return fmt.Sprintf("%s: objects found: %d, elements updated: %d of %d", filename, statistics.objects, statistics.changed, statistics.elements)
}

func assertIs(e *xmltree.XMLElement, kind string) (err error) {
	// if !Quiet {
	// 	log.Println(e.Name.Local)
	// }
	if e.Name.Local != kind {
		err = fmt.Errorf("invalid file: expected %s but found %s", kind, e.Name.Local)
	}
	return
}

/////////////////////////////////////////////////////////////////////////

type Job struct {
	xfiles []*XFile
}

type XFile struct {
	filename string
	root     *xmltree.XMLTree
	stats    Statistics
}

func FindMatchingFiles(root, pattern string) (matches []string, err error) {

	walker := func(path string, fi os.FileInfo, pe error) (err error) {

		if pe != nil || fi.IsDir() {
			return
		}

		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			return
		}

		if matched {
			matches = append(matches, path)
		}

		return
	}

	err = filepath.Walk(root, walker)
	return
}

func LoadJobFor(root, pattern string) (j Job, err error) {

	// get the list of files applicable
	filenames, err := FindMatchingFiles(root, pattern)
	if err != nil {
		return
	}

	// load each and every one so we have access to all of them
	for i := range filenames {
		var root *xmltree.XMLTree
		root, err = xmltree.LoadFromFile(filenames[i])
		if err != nil {
			return
		}
		j.xfiles = append(j.xfiles, &XFile{filename: filenames[i], root: root})
	}

	return
}

func (j *Job) Save() {
	// save all of our files
	for _, f := range j.xfiles {
		err := f.root.WriteToFile(f.filename)
		switch err {
		case nil:
			log.Println(f.stats.For(f.filename))
		default:
			log.Printf("failed to write %s: %s", f.filename, err)
		}
	}
}

func (j *Job) FindElement(tag, value string) (element *xmltree.XMLElement, file *XFile) {
	for _, file = range j.xfiles {
		element, _ = file.root.Find(tag, value)
		if element != nil {
			return
		}
	}
	return
}

// apply stats for each component
func (j *Job) ApplyAll(components map[string]ComponentData) (err error) {
	for k, v := range components {
		err = j.Apply(k, v)
		if err != nil {
			return
		}
	}
	return
}

// applies stats for given component
func (j *Job) Apply(name string, data ComponentData) (err error) {

	// find this component definition
	e, f := j.FindElement("Name", name)
	if e == nil {
		return fmt.Errorf("%s not found", name)
	}
	statistics := &f.stats

	// ensure we have correct number of component stats to update
	err = e.Child("Values").SetElementCount(1 + data.maxLevel - data.minLevel)
	if err != nil {
		return
	}

	// fill in the data from our data tables
	stats := e.Child("Values").Elements()
	for i, e := range stats {
		for key, f := range data.fieldValues {
			e.Child(key).SetValue(f(data.minLevel + i))
		}
		statistics.elements++
		statistics.changed++
	}
	statistics.objects++

	// scale to fighter if required
	for _, name := range data.scaleTo {
		err = j.ScaleToComponentByName(e, name)
	}

	return
}

func (j *Job) ScaleToComponentByName(source *xmltree.XMLElement, name string) (err error) {

	// figure out our target name
	e, f := j.FindElement("Name", name)
	if e == nil {
		err = fmt.Errorf("%s not found", name)
		return
	}

	// scale it
	err = j.ScaleComponentToComponent(f, source, e)

	return
}

func (j *Job) ScaleComponentToComponent(file *XFile, source *xmltree.XMLElement, e *xmltree.XMLElement) (err error) {

	statistics := &file.stats

	// distinguish what kind of target component we're dealing with
	isFighter := e.Has("IsFighterOnly", "true")
	isPointDefense := !isFighter && e.Has("Category", "WeaponIntercept")
	isWeapon := e.HasPrefix("Category", "Weapon")

	// copy (and scale fighter) resource requirements
	if isFighter {
		err = e.CopyAndVisitByTag("ResourcesRequired", source, func(e *xmltree.XMLElement) error { e.Child("Amount").ScaleBy(0.25); return nil })
	} else {
		err = e.CopyByTag("ResourcesRequired", source)
	}
	if err != nil {
		return
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

		// "flatten" source volleys to 1 per shot but at 1/x fire rate (same dps, but distributed instead of burste firing)
		if va := e.Child("WeaponVolleyAmount").NumericValue(); va != 1 {
			e.Child("WeaponFireRate").ScaleBy(1.0 / va)
			e.Child("WeaponVolleyAmount").SetValue(1)
		}
		e.Child("WeaponVolleyFireRate").SetString("0")

		if isFighter || isPointDefense {

			// scale standard fire relative to our source weapon
			err = ScaleFtrOrPDMainWeaponValues(e, isFighter)
			if err != nil {
				return
			}

			// scale intercept function
			err = ScaleFtrOrPDInterceptValues(e, isFighter)
			if err != nil {
				return
			}

			// scale down the ion defenses and offenses
			err = ScaleFtrOrPDIonValues(e, isFighter)
			if err != nil {
				return
			}

			// fighters and PD never do bombard damage
			for _, e := range e.Matching(regexp.MustCompile("WeaponBombard.*")) {
				e.SetString("0")
			}

			if isFighter {
				// fighters never have crew requirements
				e.Child("CrewRequirement").SetString("0")

				if isWeapon {
					// fighter weapons have no static draw for ftr
					e.Child("StaticEnergyUsed").SetString("0")
				} else {
					// but other fighter components are simply scaled down
					e.Child("StaticEnergyUsed").ScaleBy(0.25)
				}
			}
		}

		statistics.changed++
		statistics.elements++
	}

	statistics.objects++

	return
}

func ScaleFtrOrPDIonValues(e *xmltree.XMLElement, isFighter bool) (err error) {

	e.Child("ComponentIonDefense").ScaleBy(IonFtrPDScaleFactor)
	e.Child("IonDamageDefense").ScaleBy(IonFtrPDScaleFactor)

	e.Child("WeaponIonEngineDamage").ScaleBy(IonFtrPDScaleFactor)
	e.Child("WeaponIonHyperDriveDamage").ScaleBy(IonFtrPDScaleFactor)
	e.Child("WeaponIonSensorDamage").ScaleBy(IonFtrPDScaleFactor)
	e.Child("WeaponIonShieldDamage").ScaleBy(IonFtrPDScaleFactor)
	e.Child("WeaponIonWeaponDamage").ScaleBy(IonFtrPDScaleFactor)
	e.Child("WeaponIonGeneralDamage").ScaleBy(IonFtrPDScaleFactor)

	return
}

func FtrOrPDMainWeaponScaling(isFighter bool) (rof float64, dmg float64) {
	// 4 x .375 = 1.5x total output
	return 4, .375
}

func FtrOrPDInterceptScaling(isFighter bool) (rof float64, dmg float64) {
	if isFighter {
		// for fighters scale intercept by...
		// previously we were at a net 4x dmg vs. fighters and 8x dmg vs. seeking
		// (8x rof, 1/2x dmg vs. fighters, and 8x rof, 1x dmg vs. seeking)
		// now we're at 5x4 = 20x rof vs. standard (was 32x)
		// and 5 x .4 = 200% total damage output compard to base, which is 1.5 normal = 300% total vs. standard weapon
		rof = 5
		dmg = .4
	} else {
		// PD is 2 * 2 = 4x as effective as a ftr
		// the very high rof means we should get cool visuals (blasters are now approx 4/s)
		// note: we might want to break this out by weapon type (super high for kinetic & blaster, less so for beams & missiles)
		rof = 10
		dmg = .8
	}
	return
}

func ScaleFtrOrPDMainWeaponValues(e *xmltree.XMLElement, isFighter bool) (err error) {

	// get appropriate scaling factors
	rof, dmg := FtrOrPDMainWeaponScaling(isFighter)

	// scale by our source weapon values
	e.Child("WeaponFireRate").ScaleBy(1 / rof)
	e.Child("WeaponRawDamage").ScaleBy(dmg)
	e.Child("WeaponEnergyPerShot").ScaleBy(dmg)

	// range is 1/3 + 50% more rapid fall-off
	e.Child("WeaponRange").ScaleBy(0.3333333333)
	e.Child("WeaponDamageFalloffRatio").ScaleBy(1.5)

	// fighter & PD weapons generically get a +10% targeting across the board (very short range = enhanced accuracy)
	e.Child("ComponentTargetingBonus").AdjustValue(0.1)

	return
}

func ScaleFtrOrPDInterceptValues(e *xmltree.XMLElement, isFighter bool) (err error) {

	// get appropriate scaling factors
	rof, dmg := FtrOrPDInterceptScaling(isFighter)

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

	return
}

func GetComponentSourceName(targetName string, isFighter bool) (sourceName string) {

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

		if isFighter {
			sourceName = targetName[:len(targetName)-len(" [Ftr]")] + " [S]"
		} else {
			sourceName = targetName[:len(targetName)-len(" [PD]")] + " [S]"
		}
	}

	return
}
