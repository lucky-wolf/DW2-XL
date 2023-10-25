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
var costs = []int{0, 500, 1000, 2000, 4500, 8000, 12500, 18000, 24500, 32000, 48000, 64000, 90000}
var amounts = []int{0, 50, 100, 200, 450, 800, 1250, 1800, 2450, 3200, 4800, 6400, 9000}

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

				// get our column
				col := e.Child("Column").IntValue()
				if col >= len(sizes) {
					err = fmt.Errorf("Column %d exceeds maximum: %s", col, e.Child("Name").StringValue())
					return
				}

				// set size from that
				e.Child("Size").SetValue(sizes[col])

				// set our initiation cost and resource amounts (if present)
				if cost := e.Child("InitiationCost"); cost != nil {
					cost.Child("Money").SetValue(costs[col])
					for _, e := range cost.Child("Resources").Elements() {
						e.Child("Amount").SetValue(amounts[col])
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
