package algorithm

var (
	cruiserNames = []string{"Cruiser", "Cruiser, Fast", "Cruiser, Shield", "Cruiser, Fleet", "Cruiser, Heavy", "Cruiser, Galaxy", "Cruiser, Patrol", "Cruiser, Command", "Cruiser, Battle"}

	// what we wish it were
	cruiserHullSchedule = HullBaySchedule{
		"Cruiser": {
			StringsTable: StringsTable{
				"Name": func(tier Tier) string { return cruiserNames[tier] },
			},
			ValuesTable: ComponentStats{
				"ArmorReactiveRating":  func(tier Tier) float64 { return float64(2 * tier) },
				"IonDefenseRating":     func(tier Tier) float64 { return float64(4 * tier) },
				"CountermeasuresBonus": func(tier Tier) float64 { return []float64{.26, .32, .38, .44, .50, .56}[tier-1] },
				"TargetingBonus":       func(tier Tier) float64 { return []float64{.03, .06, .09, .12, .15, .18}[tier-1] },
				"ManeuveringBonus":     func(tier Tier) float64 { return []float64{.08, .16, .24, .32, .40, .48}[tier-1] },
			},
			BayCountsPerLevels: BayCountsPerLevels{
				// Cruiser
				0: {},
				// Cruiser, Fast
				1: {
					"Weapon":  1,
					"Engine":  2,
					"Defense": 2,
					"General": 3,
				},
				// Cruiser, Shield
				2: {
					"Weapon":  2,
					"Engine":  2,
					"Defense": 2,
					"General": 3,
				},
				// Cruiser, Fleet
				3: {
					"Weapon":  2,
					"Engine":  3,
					"Defense": 2,
					"General": 4,
				},
				// Cruiser, Heavy
				4: {
					"Weapon":  2,
					"Engine":  3,
					"Defense": 3,
					"General": 4,
				},
				// Cruiser, Galaxy
				5: {
					"Weapon":  2,
					"Engine":  3,
					"Defense": 3,
					"General": 5,
				},
			},
		},
	}
)
