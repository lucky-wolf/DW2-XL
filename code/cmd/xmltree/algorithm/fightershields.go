package algorithm

import (
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
	"strings"
)

func FighterShields(folder string) (err error) {

	if !Quiet {
		log.Println("All strikecraft shields will be set to:")
		log.Println("- 20% Strength (and energy cost)")
		log.Println("- 20% Recharge (and energy cost)")
		log.Println("- 100% Ion Defense and Component Defense")
		log.Println("- 100% Resistance and Component Defense")
	}

	// load all component definition files
	j, err := loadJobFor(folder, "ComponentDefinitions*")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyFighterShields()
	if err != nil {
		return
	}

	// save them all
	j.save()

	return
}

func (j *job) applyFighterShields() (err error) {

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

				// only shields...
				if !e.Has("Category", "Shields") {
					continue
				}

				// for fighters...
				targetName := e.Child("Name").StringValue()
				if !strings.HasSuffix(targetName, "[Ftr]") {
					continue
				}

				// find the corresponding ship shields by same name
				sourceName := strings.TrimSpace(targetName[:len(targetName)-len("[Ftr]")])
				sourceDefinition, _ := j.find("Name", sourceName)
				if sourceDefinition == nil {
					log.Printf("element not found: %s for %s", sourceName, targetName)
					continue
				}

				// copy and scale resource requirements
				err = e.CopyAndVisitByTag("ResourcesRequired", sourceDefinition, func(e *xmltree.XMLElement) error { e.Child("Amount").ScaleBy(0.2); return nil })
				if err != nil {
					log.Println(err)
				}

				// copy component stats
				err = e.CopyByTag("Values", sourceDefinition)
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

					// scale / modify the values for the component to match source
					e.Child("ShieldRechargeRate").ScaleBy(0.2)
					e.Child("ShieldRechargeEnergyUsage").ScaleBy(0.2)
					e.Child("ShieldResistance").ScaleBy(0.2)
					e.Child("ShieldStrength").ScaleBy(0.2)
					e.Child("StaticEnergyUsed").ScaleBy(0.2)

					// never a crew requirement for fighter components
					e.Child("CrewRequirement").SetValue(0)

					statistics.changed++
					statistics.elements++
				}

				statistics.objects++
			}
		}
	}
	err = nil
	return
}