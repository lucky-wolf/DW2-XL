package algorithm

import (
	"log"
)

func IonWeapons(folder string) (err error) {

	log.Println("Updates ion weapons off of a common core data table")

	// load all component definition files
	j, err := LoadJobFor(folder, "ComponentDefinitions*")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyIonWeapons()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

func (j *Job) applyIonWeapons() (err error) {

	// warn: we number from 1..11 where 1 = t0, and 2,3,...,10 = t2,t3,...,t1l0
	// simply 11/level, so 12..120 for a [S] ion cannon
	IonComponentDamage := func(level int) float64 {
		return float64(level) * 12
	}

	// this is 50% slower (2/3 of) blasters
	WeaponRateOfFire := func(level int) float64 { return 13.5 }

	// standard damage is based on pulsed blasters-ish, but at about 2/3 the ROF, so 2/3 the DPS
	WeaponRawDamage := func(level int) float64 {
		// this gives us 20 at (t0) and a gain of 15% per level beyond that
		return 20 * (1 + 0.15*float64(level-1))
	}

	// this seems more or less pointless -- what's going to counter a blaster or a beam or a blast?
	ComponentCountermeasuresBonus := func(level int) float64 { return 0.6 + float64(level)*0.02 }

	// note: ion weapons never have any bombard value (lighting in atmosphere is not a real issue)
	IonFieldProjector := ValuesTable{
		"ComponentCountermeasuresBonus":     ComponentCountermeasuresBonus,
		"ComponentTargetingBonus":           func(level int) float64 { return 0 },
		"WeaponBombardDamageInfrastructure": func(level int) float64 { return 0 }, // zero bombard value
		"WeaponBombardDamageMilitary":       func(level int) float64 { return 0 },
		"WeaponBombardDamagePopulation":     func(level int) float64 { return 0 },
		"WeaponBombardDamageQuality":        func(level int) float64 { return 0 },
		"WeaponArmorBypass":                 func(level int) float64 { return 0.25 },  // std +25
		"WeaponShieldBypass":                func(level int) float64 { return -0.25 }, // std -25
		"WeaponSpeed":                       func(level int) float64 { return 2200 },  // todo: drive this off of WeaponFireType or Family
		"WeaponRange":                       func(level int) float64 { return 800 + float64(level*100) },
		"WeaponDamageFalloffRatio":          func(level int) float64 { return 0.25 },
		"WeaponEnergyPerShot":               func(level int) float64 { return 0.75 * WeaponRawDamage(level) },
		"WeaponFireRate":                    WeaponRateOfFire,
		"WeaponRawDamage":                   WeaponRawDamage,
		"WeaponIonEngineDamage":             IonComponentDamage,
		"WeaponIonHyperDriveDamage":         IonComponentDamage,
		"WeaponIonSensorDamage":             IonComponentDamage,
		"WeaponIonShieldDamage":             IonComponentDamage,
		"WeaponIonWeaponDamage":             IonComponentDamage,
		"WeaponIonGeneralDamage":            IonComponentDamage,
	}
	SmallIonCannon := ExtendValuesTable(
		IonFieldProjector,
		ValuesTable{},
	)

	// [M] is simply 2x fire [S]
	MediumIonCannon := ExtendValuesTable(
		SmallIonCannon,
		ValuesTable{
			"WeaponVolleyAmount":   func(level int) float64 { return 2 },
			"WeaponVolleyFireRate": func(level int) float64 { return 1 },
		},
	)

	// rapid is simply 1.25 rof
	RapidIonCannon := ExtendValuesTable(
		IonFieldProjector,
		ValuesTable{
			"WeaponFireRate": func(level int) float64 { return 0.8 * WeaponRateOfFire(level) }, // 25% faster, but no loss of ion output
		},
	)

	// heavy medium is more powerful per shot instead of double shot
	MediumHeavyIonCannon := ExtendValuesTable(
		SmallIonCannon,
		ValuesTable{
			"WeaponEnergyPerShot":       func(level int) float64 { return 1.5 * WeaponRawDamage(level) },
			"WeaponRange":               func(level int) float64 { return 800 + float64(level*125) },
			"WeaponDamageFalloffRatio":  func(level int) float64 { return 0.2 },
			"WeaponRawDamage":           func(level int) float64 { return 2 * WeaponRawDamage(level) },
			"WeaponIonEngineDamage":     func(level int) float64 { return 1.25 * IonComponentDamage(level) },
			"WeaponIonHyperDriveDamage": func(level int) float64 { return 1.25 * IonComponentDamage(level) },
			"WeaponIonSensorDamage":     func(level int) float64 { return 1.25 * IonComponentDamage(level) },
			"WeaponIonShieldDamage":     func(level int) float64 { return 1.25 * IonComponentDamage(level) },
			"WeaponIonWeaponDamage":     func(level int) float64 { return 1.25 * IonComponentDamage(level) },
			"WeaponIonGeneralDamage":    func(level int) float64 { return 1.25 * IonComponentDamage(level) },
		},
	)

	// heavy large is double heavy medium
	LargeHeavyIonCannon := ExtendValuesTable(
		SmallIonCannon,
		ValuesTable{
			"WeaponEnergyPerShot":       func(level int) float64 { return 3 * WeaponRawDamage(level) },
			"WeaponRange":               func(level int) float64 { return 800 + float64(level*150) },
			"WeaponDamageFalloffRatio":  func(level int) float64 { return 0.2 },
			"WeaponRawDamage":           func(level int) float64 { return 4 * WeaponRawDamage(level) },
			"WeaponIonEngineDamage":     func(level int) float64 { return 1.5 * IonComponentDamage(level) },
			"WeaponIonHyperDriveDamage": func(level int) float64 { return 1.5 * IonComponentDamage(level) },
			"WeaponIonSensorDamage":     func(level int) float64 { return 1.5 * IonComponentDamage(level) },
			"WeaponIonShieldDamage":     func(level int) float64 { return 1.5 * IonComponentDamage(level) },
			"WeaponIonWeaponDamage":     func(level int) float64 { return 1.5 * IonComponentDamage(level) },
			"WeaponIonGeneralDamage":    func(level int) float64 { return 1.5 * IonComponentDamage(level) },
		},
	)

	// Lance is slower with bigger per shot values, some targeting bonus, good range
	EMLance := ExtendValuesTable(
		SmallIonCannon,
		ValuesTable{
			"ComponentTargetingBonus":   func(level int) float64 { return float64(level-1) * 0.05 },
			"WeaponArmorBypass":         func(level int) float64 { return 0.15 },  // std +25
			"WeaponShieldBypass":        func(level int) float64 { return -0.15 }, // std -25
			"WeaponEnergyPerShot":       func(level int) float64 { return 3.75 * WeaponRawDamage(level) },
			"WeaponFireRate":            func(level int) float64 { return 1.25 * WeaponRateOfFire(level) }, // 25% slower
			"WeaponSpeed":               func(level int) float64 { return 5000 },                           // todo: drive this off of WeaponFireType or Family
			"WeaponRange":               func(level int) float64 { return 1000 + float64(level*200) },
			"WeaponDamageFalloffRatio":  func(level int) float64 { return 0.15 },
			"WeaponRawDamage":           func(level int) float64 { return 5 * WeaponRawDamage(level) },
			"WeaponIonEngineDamage":     func(level int) float64 { return 1.5 * IonComponentDamage(level) },
			"WeaponIonHyperDriveDamage": func(level int) float64 { return 1.5 * IonComponentDamage(level) },
			"WeaponIonSensorDamage":     func(level int) float64 { return 1.5 * IonComponentDamage(level) },
			"WeaponIonShieldDamage":     func(level int) float64 { return 1.5 * IonComponentDamage(level) },
			"WeaponIonWeaponDamage":     func(level int) float64 { return 1.5 * IonComponentDamage(level) },
			"WeaponIonGeneralDamage":    func(level int) float64 { return 1.5 * IonComponentDamage(level) },
		},
	)

	// Wave Lance has more armor bypass
	EMWaveLance := ExtendValuesTable(
		EMLance,
		ValuesTable{
			"WeaponArmorBypass": func(level int) float64 { return 0.33333 }, // std +25
		},
	)

	// Lance is slower with bigger per shot values, some targeting bonus, good range
	MediumIonBomb := ValuesTable{
		"ComponentCountermeasuresBonus":     func(level int) float64 { return 0.38 + 0.02*float64(level) },
		"ComponentTargetingBonus":           func(level int) float64 { return 0.2 },
		"WeaponAreaEffectRange":             func(level int) float64 { return 150 + float64(25*level) },
		"WeaponAreaBlastWaveSpeed":          func(level int) float64 { return 300 },
		"WeaponBombardDamageInfrastructure": func(level int) float64 { return 0 }, // zero bombard value
		"WeaponBombardDamageMilitary":       func(level int) float64 { return 0 },
		"WeaponBombardDamagePopulation":     func(level int) float64 { return 0 },
		"WeaponBombardDamageQuality":        func(level int) float64 { return 0 },
		"WeaponArmorBypass":                 func(level int) float64 { return 0.25 },                    // std +25
		"WeaponShieldBypass":                func(level int) float64 { return -0.25 },                   // std -25
		"WeaponSpeed":                       func(level int) float64 { return 325 + float64(25*level) }, // todo: drive this off of WeaponFireType or Family
		"WeaponRange":                       func(level int) float64 { return float64(735 * level) },
		"WeaponDamageFalloffRatio":          func(level int) float64 { return 0.1 },
		"WeaponEnergyPerShot":               func(level int) float64 { return 1.5 * WeaponRawDamage(level) },
		"WeaponFireRate":                    func(level int) float64 { return 20 },
		"WeaponRawDamage":                   func(level int) float64 { return 2 * WeaponRawDamage(level) },
		"WeaponIonEngineDamage":             func(level int) float64 { return IonComponentDamage(level + 1) },
		"WeaponIonHyperDriveDamage":         func(level int) float64 { return IonComponentDamage(level + 1) },
		"WeaponIonSensorDamage":             func(level int) float64 { return IonComponentDamage(level + 1) },
		"WeaponIonShieldDamage":             func(level int) float64 { return IonComponentDamage(level + 1) },
		"WeaponIonWeaponDamage":             func(level int) float64 { return IonComponentDamage(level + 1) },
		"WeaponIonGeneralDamage":            func(level int) float64 { return IonComponentDamage(level + 1) },
	}
	LargeIonBomb := ExtendValuesTable(
		MediumIonBomb,
		ValuesTable{
			"WeaponEnergyPerShot":       func(level int) float64 { return 3 * WeaponRawDamage(level) },
			"WeaponRawDamage":           func(level int) float64 { return 4 * WeaponRawDamage(level) },
			"WeaponIonEngineDamage":     func(level int) float64 { return 1.5 * IonComponentDamage(level+1) },
			"WeaponIonHyperDriveDamage": func(level int) float64 { return 1.5 * IonComponentDamage(level+1) },
			"WeaponIonSensorDamage":     func(level int) float64 { return 1.5 * IonComponentDamage(level+1) },
			"WeaponIonShieldDamage":     func(level int) float64 { return 1.5 * IonComponentDamage(level+1) },
			"WeaponIonWeaponDamage":     func(level int) float64 { return 1.5 * IonComponentDamage(level+1) },
			"WeaponIonGeneralDamage":    func(level int) float64 { return 1.5 * IonComponentDamage(level+1) },
		},
	)

	IonMissile := ValuesTable{
		"ComponentCountermeasuresBonus":     func(level int) float64 { return 0.38 + 0.02*float64(level) },
		"ComponentTargetingBonus":           func(level int) float64 { return 0.1 },
		"WeaponAreaEffectRange":             func(level int) float64 { return 0 },
		"WeaponAreaBlastWaveSpeed":          func(level int) float64 { return 0 },
		"WeaponBombardDamageInfrastructure": func(level int) float64 { return 0 }, // zero bombard value
		"WeaponBombardDamageMilitary":       func(level int) float64 { return 0 },
		"WeaponBombardDamagePopulation":     func(level int) float64 { return 0 },
		"WeaponBombardDamageQuality":        func(level int) float64 { return 0 },
		"WeaponArmorBypass":                 func(level int) float64 { return -0.1 },
		"WeaponShieldBypass":                func(level int) float64 { return 0 },
		"WeaponSpeed":                       func(level int) float64 { return 375 },
		"WeaponRange":                       func(level int) float64 { return 2000 + float64(level*100) },
		"WeaponDamageFalloffRatio":          func(level int) float64 { return 0 },
		"WeaponEnergyPerShot":               func(level int) float64 { return 1 * WeaponRawDamage(level) },
		"WeaponFireRate":                    func(level int) float64 { return 20 },
		"WeaponRawDamage":                   func(level int) float64 { return 2 * WeaponRawDamage(level) },
		"WeaponIonEngineDamage":             func(level int) float64 { return IonComponentDamage(level + 1) },
		"WeaponIonHyperDriveDamage":         func(level int) float64 { return IonComponentDamage(level + 1) },
		"WeaponIonSensorDamage":             func(level int) float64 { return IonComponentDamage(level + 1) },
		"WeaponIonShieldDamage":             func(level int) float64 { return IonComponentDamage(level + 1) },
		"WeaponIonWeaponDamage":             func(level int) float64 { return IonComponentDamage(level + 1) },
		"WeaponIonGeneralDamage":            func(level int) float64 { return IonComponentDamage(level + 1) },
	}

	AdvIonMissile := ExtendValuesTable(
		IonMissile,
		ValuesTable{
			"WeaponSpeed": func(level int) float64 { return 450 },
			"WeaponRange": func(level int) float64 { return 2000 + float64(level*200) },
		},
	)

	UltraIonMissile := ExtendValuesTable(
		AdvIonMissile,
		ValuesTable{
			"WeaponSpeed": func(level int) float64 { return 550 },
			"WeaponRange": func(level int) float64 { return 2000 + float64(level*400) },
		},
	)

	// warn: we number from 1..11 where 1 = t0, and 2,3,...,10 = t2,t3,...,t1l0
	components := map[string]ComponentData{
		"Ion Field Projector [S]": {
			minLevel:    1,
			maxLevel:    1,
			fieldValues: SmallIonCannon,
		},
		"Ion Cannon [S]": {
			minLevel:    2,
			maxLevel:    5,
			fieldValues: SmallIonCannon,
		},
		"Ion Cannon [M]": {
			minLevel:    3,
			maxLevel:    5,
			fieldValues: MediumIonCannon,
		},

		"Rapid Ion Cannon [S]": {
			minLevel:    6,
			maxLevel:    10,
			fieldValues: RapidIonCannon,
		},

		"Heavy Ion Cannon [M]": {
			minLevel:    6,
			maxLevel:    10,
			fieldValues: MediumHeavyIonCannon,
		},
		"Heavy Ion Cannon [L]": {
			minLevel:    6,
			maxLevel:    10,
			fieldValues: LargeHeavyIonCannon,
		},

		"Electromagnetic Lance [L]": {
			minLevel:    3,
			maxLevel:    6,
			fieldValues: EMLance,
		},
		"Electromagnetic Wave Lance [L]": {
			minLevel:    8,
			maxLevel:    10,
			fieldValues: EMWaveLance,
		},

		"Ion Bomb [M]": {
			minLevel:    2,
			maxLevel:    5,
			fieldValues: MediumIonBomb,
		},
		"Ion Bomb [L]": {
			minLevel:    2,
			maxLevel:    5,
			fieldValues: LargeIonBomb,
		},

		"Ion Pulse [M]": {
			minLevel:    7,
			maxLevel:    10,
			fieldValues: MediumIonBomb,
		},
		"Ion Pulse [L]": {
			minLevel:    7,
			maxLevel:    10,
			fieldValues: LargeIonBomb,
		},

		"Ion Missile [M]": {
			minLevel:    2,
			maxLevel:    5,
			fieldValues: IonMissile,
		},
		"Advanced Ion Missile [M]": {
			minLevel:    6,
			maxLevel:    9,
			fieldValues: AdvIonMissile,
		},
		"Ultra Ion Missile [M]": {
			minLevel:    10,
			maxLevel:    10,
			fieldValues: UltraIonMissile,
		},
	}

	// apply stats for each component
	err = j.ApplyAll(components)

	return
}
