package algorithm

import (
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
	"strings"
)

func FighterArmor(folder string) (err error) {

	if !Quiet {
		log.Println("All strikecraft armor will be set to:")
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
	err = j.applyFighterArmor()
	if err != nil {
		return
	}

	// save them all
	j.save()

	return
}

func (j *job) applyFighterArmor() (err error) {

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

				// only armor...
				if !e.Has("Category", "Armor") {
					continue
				}

				// for fighters...
				targetName := e.Child("Name").GetStringValue()
				if !strings.HasSuffix(targetName, "[Ftr]") {
					continue
				}

				// find the corresponding ship armor by same name
				sourceName := strings.TrimSpace(targetName[:len(targetName)-len("[Ftr]")])
				sourceDefinition, _ := j.find("Name", sourceName)
				if sourceDefinition == nil {
					log.Printf("element not found: %s for %s", sourceName, targetName)
					continue
				}

				// copy and scale resource requirements
				err = e.CopyAndVisitByTag("ResourcesRequired", sourceDefinition, func(e *xmltree.XMLElement) error { e.Child("Amount").ScaleBy(0.2); return nil })
				if err != nil {
					return
				}

				// copy values
				var targetValues *xmltree.XMLElement
				_, targetValues, err = e.CopyByTag("Values", sourceDefinition)
				if err != nil {
					return
				}

				// now that we have our own copy of the component stats (same number of levels too)
				// we can update each of those to scale for [Ftr] version
				for _, e := range targetValues.Elements() {

					// every element should be a component bay
					err = assertIs(e, "ComponentStats")
					if err != nil {
						return
					}

					// scale / modify the values for the component to match source
					e.Child("ArmorBlastRating").ScaleBy(0.2)
					e.Child("ArmorReactiveRating").ScaleBy(0.2)
					e.Child("IonDamageDefense").ScaleBy(0.2)
					e.Child("CrewRequirement").SetValue(0)

				}

				statistics.changed++
				statistics.elements++
				statistics.objects++
			}
		}
	}
	return
}
