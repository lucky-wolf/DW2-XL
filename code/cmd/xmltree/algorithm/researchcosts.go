package algorithm

import (
	"fmt"
	"log"
	"strings"
)

func ResearchCosts(folder string) (err error) {

	log.Println("Research size, expense, and resource costs will be scaled by tech level (column)")

	// load all research definition files
	j, err := loadJobFor(folder, "ResearchProjectDefinitions*.xml")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyResearchCosts()
	if err != nil {
		return
	}

	// save them all
	j.save()

	return
}

// columns 0..12

// x3...
// var sizes = []int{33, 100, 300, 900, 2700, 8100, 24300, 72900, 218700, 656100, 1968300, 5904900, 400000}

// x3...x2          x3   x3   x3   x3    x2    x2     x2     x2     x2     x2      x2      ~t9
var sizes = []int{33, 100, 300, 900, 2700, 5400, 10800, 21600, 43200, 86400, 172800, 345600, 86400}

// more aggressive
// // x3...x2            x3   x3   x3   x3    x3    x2.5   x2.5   x2      x2      x1.5    x1.5    ~t9
// var sizes = []int{33, 100, 300, 900, 2700, 8100, 20250, 50625, 101250, 202500, 303750, 455625, 202500}

// XL up to 1.18.2
// // x3...x2...x1.5     x3   x3   x3   x3    x2    x2     x2     x3/2   x3/2   x3/2   x3/2   ~t9.5
// var sizes = []int{33, 100, 300, 900, 2700, 5400, 10800, 21600, 32400, 48600, 72900, 109350, 60750}

func (j *job) applyResearchCosts() (err error) {

	for _, f := range j.xfiles {

		statistics := &f.stats

		// the root will result in a single ArrayOf[RootObjectType]
		for _, e := range f.root.Elements.Elements() {

			err = assertIs(e, "ArrayOfResearchProjectDefinition")
			if err != nil {
				return
			}

			for _, e := range e.Elements() {

				// each of these is a researchprojectdefinition
				err = assertIs(e, "ResearchProjectDefinition")
				if err != nil {
					return
				}

				// get our nominal column
				col := e.Child("Column").IntValue()
				if col >= len(sizes) {
					err = fmt.Errorf("Column %d exceeds maximum: %s", col, e.Child("Name").StringValue())
					return
				}

				// get the tech name
				techName := e.Child("Name").StringValue()
				if !Quiet {
					log.Println(techName)
				}

				// use that to see if there is a cost override / exception
				switch techName {
				case "Assault Pods", "Bombardment Weapons", "Regenerating Hull Splinters":
					col = 1
				case "Cure Degenerate Gizureans", "Cure Shakturi Psionic Virus",
					"Puzzle Pirate Culture Research", "Shakturi Design and Behavior":
					col = 2
				case "Study Degenerate Gizureans", "Restore Gizurean Hive Mind":
					col = 3
				case "Basic Vault Systems", "Basic Vault Structures":
					col = 3
				default:
					if strings.HasPrefix(techName, "Ancient Guardian") {
						col = 3
					} else if strings.HasPrefix(techName, "Xeno Studies:") && col == 0 {
						col = 1
					}
				}

				// set size from that
				e.Child("Size").SetValue(sizes[col])

				// set our initiation cost and resource amounts (if present)
				if cost := e.Child("InitiationCost"); cost != nil {
					cost.Child("Money").SetValue(sizes[col] * 5)
					resources := cost.Child("Resources")
					if resources != nil {
						for _, e := range resources.Elements() {
							e.Child("Amount").SetValue(sizes[col] / 2)
						}
					}
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
