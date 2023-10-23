package algorithm

import (
	"fmt"
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
	"strings"
)

func FighterShields(root *xmltree.XMLTree) (statistics Statistics, err error) {

	if !Quiet {
		log.Println("All strikecraft shields will be set to:")
		log.Println("- 20% Strength (and energy cost)")
		log.Println("- 20% Recharge (and energy cost)")
		log.Println("- 100% Ion Defense and Component Defense")
		log.Println("- 100% Resistance and Component Defense")
	}

	// the root will result in a single ArrayOf[RootObjectType]
	for _, e := range root.Elements.Elements() {
		err = AssertIs(e, "ArrayOfComponentDefinition")
		if err != nil {
			return
		}

		for _, e := range e.Elements() {

			// each of these is a componentdefinition
			err = AssertIs(e, "ComponentDefinition")
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
			sourceDefinition, _ := root.Find("Name", sourceName)
			if sourceDefinition == nil {
				log.Printf("element not found: %s for %s", sourceName, targetName)
				continue
			}

			// copy and modify the values into our ftr shield
			sourceValues := sourceDefinition.Child("Values")
			targetValues := e.Child("Values")
			if targetValues == nil {
				panic(fmt.Errorf("%s doesn't have a values node?!", targetName))
				// e.AddChild(sourceValues.Clone())
			} else {
				targetValues.SetContents(sourceValues.CloneContents())
			}

			// for _, componentStats := range sourceValues.Elements() {

			// 	// every element should be a component bay
			// 	err = AssertIs(componentStats, "componentStats")
			// 	if err != nil {
			// 		return
			// 	}

			// 	// for each component bay whose type is Hangar
			// 	if componentStats.Find("Type", "Hangar") != nil {

			// 		// find and update the MaximumComponentSize
			// 		c := componentStats.Child("MaximumComponentSize")
			// 		v := fmt.Sprint(size)
			// 		if c.GetStringValue() != v {
			// 			c.SetString(v)
			// 			statistics.changed++
			// 		}

			// 		statistics.elements++
			// 	}
			// }

			statistics.objects++
		}
	}

	return
}
