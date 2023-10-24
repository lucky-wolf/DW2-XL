package algorithm

import (
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
	"strings"
)

func FighterWeapons(folder string) (err error) {

	if !Quiet {
		log.Println("All strikecraft weapons will be set to:")
		log.Println("- 50% Energy / shot")
		log.Println("- 10x Rate of fire")
		log.Println("- 25% Range")
		log.Println("- 20% Damage / shot")
	}

	// load all component definition files
	j, err := loadJobFor(folder, "ComponentDefinitions*")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyFighterWeapons()
	if err != nil {
		return
	}

	// save them all
	j.save()

	return
}

func (j *job) applyFighterWeapons() (err error) {

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
				targetName := e.Child("Name").GetStringValue()
				if !strings.HasSuffix(targetName, "[Ftr]") {
					continue
				}

				// find the corresponding small weapon by name
				sourceName := strings.Replace(targetName, "[Ftr]", "[S]", 1)
				sourceDefinition, _ := j.find("Name", sourceName)
				if sourceDefinition == nil {
					log.Printf("element not found: %s for %s", sourceName, targetName)
					continue
				}

				// debug
				log.Printf("processing %s from %s...\n", targetName, sourceName)

				// copy and scale resource requirements
				err = e.CopyAndVisitByTag("ResourcesRequired", sourceDefinition, func(e *xmltree.XMLElement) error { e.Child("Amount").ScaleBy(0.25); return nil })
				if err != nil {
					return
				}

				// copy component stats
				err = e.CopyByTag("Values", sourceDefinition)
				if err != nil {
					return
				}

				// now that we have our own copy of the component stats (same number of levels too)
				// we can update each of those to scale for [Ftr] version
				for _, e := range e.Child("Values").Elements() {

					// every element should be a component bay
					err = assertIs(e, "ComponentStats")
					if err != nil {
						return
					}

					// scale / modify the values for the component to match source
					e.Child("WeaponEnergyPerShot").ScaleBy(0.5)
					e.Child("WeaponFireRate").ScaleBy(10)
					e.Child("WeaponRange").ScaleBy(0.25)
					e.Child("WeaponRawDamage").ScaleBy(0.2)

					// never a crew requirement for fighter components
					e.Child("CrewRequirement").SetValue(0)
					e.Child("StaticEnergyUsed").SetValue(0)

					// copy our own values into intercept capabilities
					// WeaponInterceptRange
				}

				statistics.changed++
				statistics.elements++
				statistics.objects++
			}
		}
	}
	return
}
