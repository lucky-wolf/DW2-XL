package algorithm

import (
	"log"
)

func IonShields(folder string) (err error) {

	log.Println("Updates ion shields off of a common data table (ship & ftr)")

	// load all component definition files
	j, err := LoadJobFor(folder, "ComponentDefinitions*.xml")
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
			derivatives: []string{"Ion Shield [F/B]"},
			minLevel:    2,
			maxLevel:    10,
			componentStats: ComponentStats{
				"ComponentIonDefense": HardenedComponentIonDefense,
				"IonDamageDefense":    IonShieldIonDamageDefense,
				"CrewRequirement":     func(level int) float64 { return 5 },
				"StaticEnergyUsed":    func(level int) float64 { return 2.5 * float64(level) },
			},
		},
	}

	// apply stats for each component
	err = j.ApplyComponentAll(components)

	return
}
