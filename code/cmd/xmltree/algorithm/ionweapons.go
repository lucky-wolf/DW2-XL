package algorithm

import (
	"fmt"
	"log"
)

func IonWeapons(folder string) (err error) {

	log.Println("Updates core ion shields to have 5 levels off of a common data table")

	// load all component definition files
	j, err := loadJobFor(folder, "ComponentDefinitions*")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyIonWeapons()
	if err != nil {
		return
	}

	// save them all
	j.save()

	return
}

func (j *job) applyIonWeapons() (err error) {

	type LevelFunc = func(level int) float64

	WeaponRawDamage := func(level int) float64 {
		return float64(12 * level)
	}
	IonComponentDamage := func(level int) float64 {
		return WeaponRawDamage(level) * 1.5
	}
	EnergyPerShot := func(level int) float64 { return 14 + float64(level*2) }

	type FieldLookup = map[string]LevelFunc

	mergeFields := func(fields, more FieldLookup) (result FieldLookup) {
		for k, v := range fields {
			result[k] = v
		}
		for k, v := range more {
			result[k] = v
		}
		return
	}

	type ComponentData struct {
		// todo: could we ever look this up viz research tree for first occurrence there?
		//       that's just column in which it is listed for each level...
		minLevel int
		maxLevel int
		fields   FieldLookup
	}

	// note: ion weapons never have any bombard value (lighting in atmosphere is not a real issue)
	IonFieldProjector := FieldLookup{
		"ComponentCountermeasuresBonus":     func(level int) float64 { return 0 }, // can't really shoot them down anyway
		"ComponentTargetingBonus":           func(level int) float64 { return 0 },
		"WeaponBombardDamageInfrastructure": func(level int) float64 { return 0 }, // zero bombard value
		"WeaponBombardDamageMilitary":       func(level int) float64 { return 0 },
		"WeaponBombardDamagePopulation":     func(level int) float64 { return 0 },
		"WeaponBombardDamageQuality":        func(level int) float64 { return 0 },
		"WeaponArmorBypass":                 func(level int) float64 { return 0.25 },  // std +25
		"WeaponShieldBypass":                func(level int) float64 { return -0.25 }, // std -25
		"WeaponSpeed":                       func(level int) float64 { return 2200 },  // todo: drive this off of WeaponFireType or Family
		"WeaponRange":                       func(level int) float64 { return 800 + float64(level*100) },
		"WeaponDamageFalloffRatio":          func(level int) float64 { return 0.2 },
		"WeaponEnergyPerShot":               EnergyPerShot,
		"WeaponRawDamage":                   WeaponRawDamage,
		"WeaponIonEngineDamage":             IonComponentDamage,
		"WeaponIonHyperDriveDamage":         IonComponentDamage,
		"WeaponIonSensorDamage":             IonComponentDamage,
		"WeaponIonShieldDamage":             IonComponentDamage,
		"WeaponIonWeaponDamage":             IonComponentDamage,
		"WeaponIonGeneralDamage":            IonComponentDamage,
	}
	SmallIonCannon := mergeFields(
		IonFieldProjector,
		FieldLookup{},
	)

	// [M] is simply 2x fire [S]
	MediumIonCannon := mergeFields(
		SmallIonCannon,
		FieldLookup{
			"WeaponVolleyAmount":   func(level int) float64 { return 2 },
			"WeaponVolleyFireRate": func(level int) float64 { return 1 },
		},
	)

	// heavy is 2x powerful, not a double shot
	MediumHeavyIonCannon := mergeFields(
		SmallIonCannon,
		FieldLookup{
			"WeaponEnergyPerShot":       func(level int) float64 { return EnergyPerShot(level) * 2 },
			"WeaponRange":               func(level int) float64 { return 800 + float64(level*125) },
			"WeaponDamageFalloffRatio":  func(level int) float64 { return 0.175 },
			"WeaponRawDamage":           func(level int) float64 { return 2 * WeaponRawDamage(level) },
			"WeaponIonEngineDamage":     func(level int) float64 { return 2 * IonComponentDamage(level) },
			"WeaponIonHyperDriveDamage": func(level int) float64 { return 2 * IonComponentDamage(level) },
			"WeaponIonSensorDamage":     func(level int) float64 { return 2 * IonComponentDamage(level) },
			"WeaponIonShieldDamage":     func(level int) float64 { return 2 * IonComponentDamage(level) },
			"WeaponIonWeaponDamage":     func(level int) float64 { return 2 * IonComponentDamage(level) },
			"WeaponIonGeneralDamage":    func(level int) float64 { return 2 * IonComponentDamage(level) },
		},
	)

	// heavy large is 4x powerful
	LargeHeavyIonCannon := mergeFields(
		SmallIonCannon,
		FieldLookup{
			"WeaponEnergyPerShot":       func(level int) float64 { return EnergyPerShot(level) * 4 },
			"WeaponRange":               func(level int) float64 { return 800 + float64(level*150) },
			"WeaponDamageFalloffRatio":  func(level int) float64 { return 0.15 },
			"WeaponRawDamage":           func(level int) float64 { return 4 * WeaponRawDamage(level) },
			"WeaponIonEngineDamage":     func(level int) float64 { return 4 * IonComponentDamage(level) },
			"WeaponIonHyperDriveDamage": func(level int) float64 { return 4 * IonComponentDamage(level) },
			"WeaponIonSensorDamage":     func(level int) float64 { return 4 * IonComponentDamage(level) },
			"WeaponIonShieldDamage":     func(level int) float64 { return 4 * IonComponentDamage(level) },
			"WeaponIonWeaponDamage":     func(level int) float64 { return 4 * IonComponentDamage(level) },
			"WeaponIonGeneralDamage":    func(level int) float64 { return 4 * IonComponentDamage(level) },
		},
	)

	// Lance is 4x powerful
	EMLance := mergeFields(
		SmallIonCannon,
		FieldLookup{
			"ComponentTargetingBonus":   func(level int) float64 { return float64(level-1) * 0.025 },
			"WeaponEnergyPerShot":       func(level int) float64 { return EnergyPerShot(level) * 4 },
			"WeaponRange":               func(level int) float64 { return 800 + float64(level*150) },
			"WeaponDamageFalloffRatio":  func(level int) float64 { return 0.2 },
			"WeaponRawDamage":           func(level int) float64 { return 4 * WeaponRawDamage(level) },
			"WeaponIonEngineDamage":     func(level int) float64 { return 4 * IonComponentDamage(level) },
			"WeaponIonHyperDriveDamage": func(level int) float64 { return 4 * IonComponentDamage(level) },
			"WeaponIonSensorDamage":     func(level int) float64 { return 4 * IonComponentDamage(level) },
			"WeaponIonShieldDamage":     func(level int) float64 { return 4 * IonComponentDamage(level) },
			"WeaponIonWeaponDamage":     func(level int) float64 { return 4 * IonComponentDamage(level) },
			"WeaponIonGeneralDamage":    func(level int) float64 { return 4 * IonComponentDamage(level) },
		},
	)

	components := map[string]ComponentData{
		"Ion Field Projector [S]": {
			fields: SmallIonCannon,
		},
		"Ion Cannon [S]": {
			minLevel: 2,
			maxLevel: 5,
			fields:   SmallIonCannon,
		},
		"Ion Cannon [M]": {
			minLevel: 3,
			maxLevel: 5,
			fields:   MediumIonCannon,
		},

		"Rapid Ion Cannon [S]": {
			minLevel: 6,
			maxLevel: 10,
			fields:   SmallIonCannon,
		},

		"Heavy Ion Cannon [M]": {
			minLevel: 6,
			maxLevel: 10,
			fields:   MediumHeavyIonCannon,
		},
		"Heavy Ion Cannon [L]": {
			minLevel: 6,
			maxLevel: 10,
			fields:   LargeHeavyIonCannon,
		},

		"Electromagnetic Lance [L]": {
			minLevel: 3,
			maxLevel: 10,
			fields:   EMLance,
		},
	}

	applyStats := func(name string, data ComponentData) (err error) {

		// find this component definition
		e, f := j.find("Name", name)
		if e == nil {
			return fmt.Errorf("%s not found", name)
		}
		statistics := &f.stats

		// ensure we have as many stats as we need
		count := 1 + data.maxLevel - data.minLevel
		stats := e.Child("Values").Elements()
		if len(stats) == 0 {
			return fmt.Errorf("no ComponentStats found for %s", name)
		}
		for i := len(stats); i < count; i++ {
			e.Child("Values").Append(stats[len(stats)-1])
		}

		// fill in the data from our data tables
		for i, e := range e.Child("Values").Elements() {
			for key, f := range data.fields {
				e.Child(key).SetValue(f(data.minLevel + i))
			}
			statistics.elements++
			statistics.changed++
		}

		statistics.objects++

		return
	}

	// apply stats for each component
	for k, v := range components {
		err = applyStats(k, v)
		if err != nil {
			return
		}
	}

	return
}
