package algorithm

import (
	"log"
)

func IonShields(folder string) (err error) {

	log.Println("Updates core ion shields to have 5 levels off of a common data table")

	// load all component definition files
	j, err := LoadJobFor(folder, "ComponentDefinitions*")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyIonShields()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

func (j *Job) applyIonShields() (err error) {

	// t2..t10 = level 2..10
	components := map[string]ComponentData{
		"Ion Shield": {
			scaleToFtr: true,
			minLevel:   2,
			maxLevel:   10,
			fieldValues: ValuesTable{
				"ComponentIonDefense": func(level int) float64 { return 5 * float64(level) },
				"IonDamageDefense":    func(level int) float64 { return 10 * float64(level) },
				"CrewRequirement":     func(level int) float64 { return 5 },
				"StaticEnergyUsed":    func(level int) float64 { return 5 * float64(level) },
			},
		},
	}

	// apply stats for each component
	err = j.ApplyAll(components)

	return
}