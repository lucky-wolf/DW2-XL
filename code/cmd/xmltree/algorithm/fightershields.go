package algorithm

import (
	"fmt"
	"log"
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
				targetName := e.Child("Name").GetStringValue()
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

				// copy values into our ftr shield
				sourceValues := sourceDefinition.Child("Values")
				targetValues := e.Child("Values")
				if targetValues == nil {
					panic(fmt.Errorf("%s doesn't have a values node?!", targetName))
					// e.AddChild(sourceValues.Clone())
				} else {
					targetValues.SetContents(sourceValues.CloneContents())
				}

				// now that we have our own copy of the component stats (same number of levels too)
				// we can update each of those to scale for [Ftr] version
				for _, componentStats := range targetValues.Elements() {

					// every element should be a component bay
					err = assertIs(componentStats, "ComponentStats")
					if err != nil {
						return
					}

					// scale / modify the values for the component to match source
					componentStats.Child("CrewRequirement").SetString("0")
					componentStats.Child("ShieldRechargeRate").ScaleBy(0.2)
					componentStats.Child("ShieldRechargeEnergyUsage").ScaleBy(0.2)
					componentStats.Child("ShieldResistance").ScaleBy(0.2)
					componentStats.Child("ShieldStrength").ScaleBy(0.2)
					componentStats.Child("StaticEnergyUsed").ScaleBy(0.2)

					statistics.changed++
					statistics.elements++
				}

				statistics.objects++
			}
		}
	}
	return
}
