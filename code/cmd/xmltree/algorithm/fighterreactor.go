package algorithm

import (
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
	"strings"
)

func FighterReactors(folder string) (err error) {

	if !Quiet {
		log.Println("All strikecraft reactors will be set to:")
		log.Println("- 20% Power")
		log.Println("- 20% Energy Storage capacity")
		log.Println("- Fuel storage set to 2x energy storage")
	}

	// load all component definition files
	j, err := loadJobFor(folder, "ComponentDefinitions*")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyFighterReactors()
	if err != nil {
		return
	}

	// save them all
	j.save()

	return
}

func (j *job) applyFighterReactors() (err error) {

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

				// only reactors...
				if !e.Has("Category", "Reactor") {
					continue
				}

				// for fighters...
				targetName := e.Child("Name").StringValue()
				if !strings.HasSuffix(targetName, "[Ftr]") {
					continue
				}

				// find the corresponding ship reactors by same name
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
					e.Child("ReactorEnergyOutputPerSecond").ScaleBy(0.2)
					e.Child("ReactorEnergyStorageCapacity").ScaleBy(0.2)
					e.Child("ReactorFuelUnitsForFullCharge").ScaleBy(0.2)

					// set the fuel units to be enough for 10 recharges
					var value float64
					value, err = e.Child("ReactorFuelUnitsForFullCharge").GetNumericValue()
					e.Child("FuelStorageCapacity").SetValue(value * 100)

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
