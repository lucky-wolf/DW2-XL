package algorithm

import (
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
	"math"
	"strings"
)

func FighterEngines(folder string) (err error) {

	if !Quiet {
		log.Println("All strikecraft engines will be set to:")
		log.Println("- 20% Blast rating")
		log.Println("- 20% Reactive rating")
		log.Println("- 100% Ion Defense")
	}

	// load all component definition files
	j, err := loadJobFor(folder, "ComponentDefinitions*")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyFighterEngines()
	if err != nil {
		return
	}

	// save them all
	j.save()

	return
}

func (j *job) applyFighterEngines() (err error) {

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

				// only engines...
				if !e.Has("Category", "Engine") {
					continue
				}

				// for fighters...
				targetName := e.Child("Name").GetStringValue()
				if !strings.HasSuffix(targetName, "[Ftr]") {
					continue
				}

				// find the corresponding ship engines by same name
				sourceName := strings.TrimSpace(targetName[:len(targetName)-len("[Ftr]")])
				sourceDefinition, _ := j.find("Name", sourceName)
				if sourceDefinition == nil {
					log.Printf("element not found: %s for %s", sourceName, targetName)
					continue
				}

				// scale our size
				size := sourceDefinition.Child("Size")
				if size != nil {
					var value float64
					value, err = size.GetNumericValue()
					if err != nil {
						return
					}
					// size must be an integer value
					e.Child("Size").SetValue(math.Round(value / 4))
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
					e.Child("EngineMainCruiseThrust").ScaleBy(0.5)
					e.Child("EngineMainCruiseThrustEnergyUsage").ScaleBy(0.5) // .SetValue(0.7)

					e.Child("EngineMainMaximumThrust").ScaleBy(0.6)
					e.Child("EngineMainMaximumThrustEnergyUsage").ScaleBy(0.6) //.SetValue(1)

					e.Child("EngineVectoringThrust").ScaleBy(0.25)
					e.Child("EngineVectoringEnergyUsage").ScaleBy(0.25) // .SetValue(1)

					// never a crew requirement for fighter components
					e.Child("CrewRequirement").SetValue(0)

					// no static energy usage for engines
					e.Child("StaticEnergyUsed").SetValue(0)
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
