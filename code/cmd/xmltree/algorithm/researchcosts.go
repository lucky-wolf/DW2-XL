package algorithm

import (
	"fmt"
	"log"
)

func ResearchCosts(folder string) (err error) {

	if !Quiet {
		log.Println("Applies size, expense, and resource costs to research based on column position")
	}

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

// columns 0..12+
var sizes = []int{33, 100, 300, 900, 2700, 5400, 10800, 21600, 32400, 48600, 72900, 109350, 60750}

// var costs = []int{0, 500, 2000, 4000, 8000, 16000, 24000, 32000, 48000, 64000, 96000, 128000, 192000}
// var amounts = []int{0, 50, 200, 400, 800, 1600, 2400, 3200, 4800, 6400, 9600, 12800, 19200}

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
				case "Study Degenerate Gizureans":
					col = 4
				case "Basic Vault Systems", "Basic Vault Structures":
					col = 5
				case "Restore Gizurean Hive Mind":
					col = 8
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
