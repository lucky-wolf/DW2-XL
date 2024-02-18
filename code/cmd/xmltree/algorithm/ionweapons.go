package algorithm

import (
	"log"
)

func IonWeapons(folder string) (err error) {

	log.Println("Updates ion weapons off of a common core data table")

	// load all component definition files
	j, err := LoadJobFor(folder, "ComponentDefinitions*.xml")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyIonWeapons()
	if err != nil {
		return
	}

	// update derivatives
	err = j.applyFighterWeaponsAndPD()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

var (
	// WARN! please keep the multipliers of IonComponentDamage to a MINIMUM
	// HACK: everything is scaled off of this value - any increase means at equal tech with FULL defense, that "extra" amount will always get through
	// hack: and MOST EVERYTHING in the game will NOT have MAX ION DEFENSE!

	// note: ion weapons never have any bombard value (lighting in atmosphere is not a real issue)
	IonFieldProjector = ComponentStats{
		"ComponentCountermeasuresBonus":     DirectFireComponentCountermeasuresBonus,
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
		"WeaponEnergyPerShot":               func(level int) float64 { return 0.75 * IonWeaponRawDamage(level) },
		"WeaponFireRate":                    IonWeaponRateOfFire,
		"WeaponRawDamage":                   IonWeaponRawDamage,
		"WeaponIonEngineDamage":             IonWeaponComponentDamage,
		"WeaponIonHyperDriveDamage":         IonWeaponComponentDamage,
		"WeaponIonSensorDamage":             IonWeaponComponentDamage,
		"WeaponIonShieldDamage":             IonWeaponComponentDamage,
		"WeaponIonWeaponDamage":             IonWeaponComponentDamage,
		"WeaponIonGeneralDamage":            IonWeaponComponentDamage,
	}
	SmallIonCannon = IonFieldProjector

	// [M] is simply 2x fire [S]
	MediumIonCannon = ExtendValuesTable(
		SmallIonCannon,
		ComponentStats{
			"WeaponVolleyAmount":   func(level int) float64 { return 2 },
			"WeaponVolleyFireRate": func(level int) float64 { return 1 },
		},
	)

	// rapid is simply 1.25 rof
	RapidIonCannon = ExtendValuesTable(
		IonFieldProjector,
		ComponentStats{
			"WeaponFireRate": func(level int) float64 { return 0.8 * IonWeaponRateOfFire(level) }, // 25% faster, but no loss of ion output
		},
	)

	// heavy medium is more powerful per shot instead of double shot
	MediumHeavyIonCannon = ExtendValuesTable(
		SmallIonCannon,
		ComponentStats{
			"WeaponEnergyPerShot":       func(level int) float64 { return 1.5 * IonWeaponRawDamage(level) },
			"WeaponRange":               func(level int) float64 { return 800 + float64(level*125) },
			"WeaponDamageFalloffRatio":  func(level int) float64 { return 0.2 },
			"WeaponRawDamage":           func(level int) float64 { return 2 * IonWeaponRawDamage(level) },
			"WeaponIonEngineDamage":     func(level int) float64 { return 1.1 * IonWeaponComponentDamage(level) },
			"WeaponIonHyperDriveDamage": func(level int) float64 { return 1.1 * IonWeaponComponentDamage(level) },
			"WeaponIonSensorDamage":     func(level int) float64 { return 1.1 * IonWeaponComponentDamage(level) },
			"WeaponIonShieldDamage":     func(level int) float64 { return 1.1 * IonWeaponComponentDamage(level) },
			"WeaponIonWeaponDamage":     func(level int) float64 { return 1.1 * IonWeaponComponentDamage(level) },
			"WeaponIonGeneralDamage":    func(level int) float64 { return 1.1 * IonWeaponComponentDamage(level) },
		},
	)

	// heavy large is double heavy medium
	LargeHeavyIonCannon = ExtendValuesTable(
		SmallIonCannon,
		ComponentStats{
			"WeaponEnergyPerShot":       func(level int) float64 { return 3 * IonWeaponRawDamage(level) },
			"WeaponRange":               func(level int) float64 { return 800 + float64(level*150) },
			"WeaponDamageFalloffRatio":  func(level int) float64 { return 0.2 },
			"WeaponRawDamage":           func(level int) float64 { return 4 * IonWeaponRawDamage(level) },
			"WeaponIonEngineDamage":     func(level int) float64 { return 1.2 * IonWeaponComponentDamage(level) },
			"WeaponIonHyperDriveDamage": func(level int) float64 { return 1.2 * IonWeaponComponentDamage(level) },
			"WeaponIonSensorDamage":     func(level int) float64 { return 1.2 * IonWeaponComponentDamage(level) },
			"WeaponIonShieldDamage":     func(level int) float64 { return 1.2 * IonWeaponComponentDamage(level) },
			"WeaponIonWeaponDamage":     func(level int) float64 { return 1.2 * IonWeaponComponentDamage(level) },
			"WeaponIonGeneralDamage":    func(level int) float64 { return 1.2 * IonWeaponComponentDamage(level) },
		},
	)

	// Lance is slower with bigger per shot values, some targeting bonus, good range
	EMLance = ExtendValuesTable(
		SmallIonCannon,
		ComponentStats{
			"ComponentTargetingBonus":   func(level int) float64 { return float64(level-1) * 0.05 },
			"WeaponArmorBypass":         func(level int) float64 { return 0.15 },  // std +25
			"WeaponShieldBypass":        func(level int) float64 { return -0.15 }, // std -25
			"WeaponEnergyPerShot":       func(level int) float64 { return 3.75 * IonWeaponRawDamage(level) },
			"WeaponFireRate":            func(level int) float64 { return 1.25 * IonWeaponRateOfFire(level) }, // 25% slower
			"WeaponSpeed":               func(level int) float64 { return 5000 },                              // todo: drive this off of WeaponFireType or Family
			"WeaponRange":               func(level int) float64 { return 1000 + float64(level*200) },
			"WeaponDamageFalloffRatio":  func(level int) float64 { return 0.15 },
			"WeaponRawDamage":           func(level int) float64 { return 5 * IonWeaponRawDamage(level) },
			"WeaponIonEngineDamage":     func(level int) float64 { return 1.2 * IonWeaponComponentDamage(level) },
			"WeaponIonHyperDriveDamage": func(level int) float64 { return 1.2 * IonWeaponComponentDamage(level) },
			"WeaponIonSensorDamage":     func(level int) float64 { return 1.2 * IonWeaponComponentDamage(level) },
			"WeaponIonShieldDamage":     func(level int) float64 { return 1.2 * IonWeaponComponentDamage(level) },
			"WeaponIonWeaponDamage":     func(level int) float64 { return 1.2 * IonWeaponComponentDamage(level) },
			"WeaponIonGeneralDamage":    func(level int) float64 { return 1.2 * IonWeaponComponentDamage(level) },
		},
	)

	// Wave Lance has more armor bypass
	EMWaveLance = ExtendValuesTable(
		EMLance,
		ComponentStats{
			"WeaponArmorBypass": func(level int) float64 { return 0.33333 }, // std +25
		},
	)

	// Lance is slower with bigger per shot values, some targeting bonus, good range
	MediumIonBomb = ComponentStats{
		"ComponentCountermeasuresBonus":     func(level int) float64 { return 0.38 + 0.02*float64(level) },
		"ComponentTargetingBonus":           func(level int) float64 { return 0.2 },
		"WeaponAreaEffectRange":             func(level int) float64 { return 150 + float64(25*level) },
		"WeaponAreaBlastWaveSpeed":          func(level int) float64 { return 300 },
		"WeaponBombardDamageInfrastructure": func(level int) float64 { return 0 }, // zero bombard value
		"WeaponBombardDamageMilitary":       func(level int) float64 { return 0 },
		"WeaponBombardDamagePopulation":     func(level int) float64 { return 0 },
		"WeaponBombardDamageQuality":        func(level int) float64 { return 0 },
		"WeaponArmorBypass":                 func(level int) float64 { return 0.25 },  // std +25
		"WeaponShieldBypass":                func(level int) float64 { return -0.25 }, // std -25
		"WeaponSpeed":                       TorpedoSeekingSpeed,
		"WeaponRange":                       TorpedoSeekingRange,
		"WeaponDamageFalloffRatio":          func(level int) float64 { return 0.1 },
		"WeaponEnergyPerShot":               func(level int) float64 { return 1.5 * IonWeaponRawDamage(level) },
		"WeaponFireRate":                    func(level int) float64 { return 20 },
		"WeaponRawDamage":                   func(level int) float64 { return 2 * IonWeaponRawDamage(level) },
		"WeaponIonEngineDamage":             func(level int) float64 { return IonWeaponComponentDamage(level + 1) },
		"WeaponIonHyperDriveDamage":         func(level int) float64 { return IonWeaponComponentDamage(level + 1) },
		"WeaponIonSensorDamage":             func(level int) float64 { return IonWeaponComponentDamage(level + 1) },
		"WeaponIonShieldDamage":             func(level int) float64 { return IonWeaponComponentDamage(level + 1) },
		"WeaponIonWeaponDamage":             func(level int) float64 { return IonWeaponComponentDamage(level + 1) },
		"WeaponIonGeneralDamage":            func(level int) float64 { return IonWeaponComponentDamage(level + 1) },
	}
	LargeIonBomb = ExtendValuesTable(
		MediumIonBomb,
		ComponentStats{
			"WeaponEnergyPerShot":       func(level int) float64 { return 3 * IonWeaponRawDamage(level) },
			"WeaponRawDamage":           func(level int) float64 { return 4 * IonWeaponRawDamage(level) },
			"WeaponIonEngineDamage":     func(level int) float64 { return 1.2 * IonWeaponComponentDamage(level+1) },
			"WeaponIonHyperDriveDamage": func(level int) float64 { return 1.2 * IonWeaponComponentDamage(level+1) },
			"WeaponIonSensorDamage":     func(level int) float64 { return 1.2 * IonWeaponComponentDamage(level+1) },
			"WeaponIonShieldDamage":     func(level int) float64 { return 1.2 * IonWeaponComponentDamage(level+1) },
			"WeaponIonWeaponDamage":     func(level int) float64 { return 1.2 * IonWeaponComponentDamage(level+1) },
			"WeaponIonGeneralDamage":    func(level int) float64 { return 1.2 * IonWeaponComponentDamage(level+1) },
		},
	)

	// ion missiles have same ion and raw da,age output as medium cannons (scaled for slower rof)
	IonMissile = ComponentStats{
		"ComponentCountermeasuresBonus":     func(level int) float64 { return 0.48 + 0.02*float64(level+1) },
		"ComponentTargetingBonus":           func(level int) float64 { return 0.1 },
		"WeaponAreaEffectRange":             func(level int) float64 { return 0 },
		"WeaponAreaBlastWaveSpeed":          func(level int) float64 { return 0 },
		"WeaponBombardDamageInfrastructure": func(level int) float64 { return 0 }, // zero bombard value
		"WeaponBombardDamageMilitary":       func(level int) float64 { return 0 },
		"WeaponBombardDamagePopulation":     func(level int) float64 { return 0 },
		"WeaponBombardDamageQuality":        func(level int) float64 { return 0 },
		"WeaponArmorBypass":                 func(level int) float64 { return -0.1 },
		"WeaponShieldBypass":                func(level int) float64 { return 0 },
		"WeaponSpeed":                       MissileSeekingSpeed,
		"WeaponRange":                       MissileSeekingRange,
		"WeaponDamageFalloffRatio":          func(level int) float64 { return 0 },
		"WeaponEnergyPerShot":               func(level int) float64 { return IonWeaponRawDamage(level) },
		"WeaponFireRate":                    func(level int) float64 { return 20 },
		"WeaponRawDamage": func(level int) float64 {
			// todo: we really want to tie raw damage to say 80% of a standard missile, rather than to that of ion cannons
			return 20 / IonWeaponRateOfFire(level) * 2 * IonWeaponRawDamage(level)
		},
		"WeaponIonEngineDamage":     func(level int) float64 { return 1.1 * IonWeaponComponentDamage(level) },
		"WeaponIonHyperDriveDamage": func(level int) float64 { return 1.1 * IonWeaponComponentDamage(level) },
		"WeaponIonSensorDamage":     func(level int) float64 { return 1.1 * IonWeaponComponentDamage(level) },
		"WeaponIonShieldDamage":     func(level int) float64 { return 1.1 * IonWeaponComponentDamage(level) },
		"WeaponIonWeaponDamage":     func(level int) float64 { return 1.1 * IonWeaponComponentDamage(level) },
		"WeaponIonGeneralDamage":    func(level int) float64 { return 1.1 * IonWeaponComponentDamage(level) },
	}

	AdvIonMissile = ExtendValuesTable(
		IonMissile,
		ComponentStats{
			// todo: we really want to tie raw damage to say 80% of a standard missile, rather than to that of ion cannons
			"WeaponRawDamage": func(level int) float64 { return 20 / IonWeaponRateOfFire(level) * 2 * IonWeaponRawDamage(level+1) },
		},
	)

	UltraIonMissile = ExtendValuesTable(
		AdvIonMissile,
		ComponentStats{
			// todo: we really want to tie raw damage to say 80% of a standard missile, rather than to that of ion cannons
			"WeaponRawDamage": func(level int) float64 { return 20 / IonWeaponRateOfFire(level) * 4 * IonWeaponRawDamage(level+1) },
		},
	)

	// warn: we number from 1..11 where 1 = t0, and 2,3,...,10 = t2,t3,...,t1l0
	IonComponentData = map[string]ComponentData{
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
			minLevel:    4,
			maxLevel:    7,
			fieldValues: EMLance,
		},
		"Electromagnetic Wave Lance [L]": {
			minLevel:    8,
			maxLevel:    10,
			fieldValues: EMWaveLance,
		},

		"Ion Bomb [M]": {
			minLevel:    3,
			maxLevel:    6,
			fieldValues: MediumIonBomb,
		},
		"Ion Bomb [L]": {
			minLevel:    3,
			maxLevel:    6,
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
)

func (j *Job) applyIonWeapons() (err error) {

	// apply stats for each component
	err = j.ApplyComponentAll(IonComponentData)

	return
}
