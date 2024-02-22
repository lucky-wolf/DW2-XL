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

const (
	IonWeaponRawDamageRatio = 0.8 // ion weapons do ~80% raw damage vs. other direct fire weapons
	IonWeaponDamageBasis    = IonWeaponRawDamageRatio * WeaponDamageBasis
	IonBaseRateOfFire       = 1.5 * BlasterRateOfFire // ion is 50% slower than blasters

	IonFtrPDScaleFactor = 0.75
)

var (
	IonWeaponRawDamage  = MakeExpLevelFunc(IonWeaponDamageBasis*IonBaseRateOfFire/BlasterRateOfFire, WeaponDamageIncreaseExp)
	IonWeaponRateOfFire = MakeFixedLevelFunc(IonBaseRateOfFire)

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
	MediumIonCannon = ComposeComponentStats(
		SmallIonCannon,
		ComponentStats{
			"WeaponVolleyAmount":   func(level int) float64 { return 2 },
			"WeaponVolleyFireRate": func(level int) float64 { return 1 },
		},
	)

	// rapid is simply 1.25 rof
	RapidIonCannon = ComposeComponentStats(
		IonFieldProjector,
		ComponentStats{
			"WeaponFireRate": func(level int) float64 { return 0.8 * IonWeaponRateOfFire(level) }, // 25% faster, but no loss of ion output
		},
	)

	// heavy medium is more powerful per shot instead of double shot
	MediumHeavyIonCannon = ComposeComponentStats(
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
	LargeHeavyIonCannon = ComposeComponentStats(
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

	// planetary large is double large heavy
	PlanetaryHeavyIonCannon = ComposeComponentStats(
		LargeHeavyIonCannon,
		ComponentStats{
			"CrewRequirement":               PlanetaryCrewRequirements, // meaningless, but doesn't hurt
			"ComponentCountermeasuresBonus": MakeScaledFuncLevelFunc(2, DirectFireComponentCountermeasuresBonus),
			"ComponentTargetingBonus":       MakeScaledFuncLevelFunc(2, DirectFireComponentCountermeasuresBonus),
			"WeaponRawDamage":               MakeScaledFuncLevelFunc(2, LargeHeavyIonCannon["WeaponRawDamage"]),
			"WeaponRange":                   MakeScaledFuncLevelFunc(3, LargeHeavyIonCannon["WeaponRange"]),
			"WeaponDamageFalloffRatio":      MakeFixedLevelFunc(0.1),
			"WeaponIonEngineDamage":         func(level int) float64 { return 1.33333 * IonWeaponComponentDamage(level) },
			"WeaponIonHyperDriveDamage":     func(level int) float64 { return 1.33333 * IonWeaponComponentDamage(level) },
			"WeaponIonSensorDamage":         func(level int) float64 { return 1.33333 * IonWeaponComponentDamage(level) },
			"WeaponIonShieldDamage":         func(level int) float64 { return 1.33333 * IonWeaponComponentDamage(level) },
			"WeaponIonWeaponDamage":         func(level int) float64 { return 1.33333 * IonWeaponComponentDamage(level) },
			"WeaponIonGeneralDamage":        func(level int) float64 { return 1.33333 * IonWeaponComponentDamage(level) },
		},
	)

	// Lance is slower with bigger per shot values, some targeting bonus, good range
	EMLance = ComposeComponentStats(
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
	EMWaveLance = ComposeComponentStats(
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
	LargeIonBomb = ComposeComponentStats(
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

	AdvIonMissile = ComposeComponentStats(
		IonMissile,
		ComponentStats{
			// todo: we really want to tie raw damage to say 80% of a standard missile, rather than to that of ion cannons
			"WeaponRawDamage": func(level int) float64 { return 20 / IonWeaponRateOfFire(level) * 2 * IonWeaponRawDamage(level+1) },
		},
	)

	UltraIonMissile = ComposeComponentStats(
		AdvIonMissile,
		ComponentStats{
			// todo: we really want to tie raw damage to say 80% of a standard missile, rather than to that of ion cannons
			"WeaponRawDamage": func(level int) float64 { return 20 / IonWeaponRateOfFire(level) * 4 * IonWeaponRawDamage(level+1) },
		},
	)

	// warn: we treat the t0 tech as t1
	IonComponentData = map[string]ComponentData{
		"Ion Field Projector [S]": {
			minLevel:       1,
			maxLevel:       1,
			componentStats: SmallIonCannon,
		},
		"Ion Cannon [S]": {
			minLevel:       2,
			maxLevel:       5,
			componentStats: SmallIonCannon,
		},
		"Ion Cannon [M]": {
			minLevel:       3,
			maxLevel:       5,
			componentStats: MediumIonCannon,
		},

		"Rapid Ion Cannon [S]": {
			minLevel:       6,
			maxLevel:       10,
			componentStats: RapidIonCannon,
		},

		"Heavy Ion Cannon [M]": {
			minLevel:       6,
			maxLevel:       10,
			componentStats: MediumHeavyIonCannon,
		},
		"Heavy Ion Cannon [L]": {
			minLevel:       6,
			maxLevel:       10,
			componentStats: LargeHeavyIonCannon,
		},

		"Electromagnetic Lance [L]": {
			minLevel:       4,
			maxLevel:       7,
			componentStats: EMLance,
		},
		"Electromagnetic Wave Lance [L]": {
			minLevel:       8,
			maxLevel:       10,
			componentStats: EMWaveLance,
		},

		"Ion Bomb [M]": {
			minLevel:       3,
			maxLevel:       6,
			componentStats: MediumIonBomb,
		},
		"Ion Bomb [L]": {
			minLevel:       3,
			maxLevel:       6,
			componentStats: LargeIonBomb,
		},

		"Ion Pulse [M]": {
			minLevel:       7,
			maxLevel:       10,
			componentStats: MediumIonBomb,
		},
		"Ion Pulse [L]": {
			minLevel:       7,
			maxLevel:       10,
			componentStats: LargeIonBomb,
		},

		"Ion Missile [M]": {
			minLevel:       2,
			maxLevel:       5,
			componentStats: IonMissile,
		},
		"Advanced Ion Missile [M]": {
			minLevel:       6,
			maxLevel:       9,
			componentStats: AdvIonMissile,
		},
		"Ultra Ion Missile [M]": {
			minLevel:       10,
			maxLevel:       10,
			componentStats: UltraIonMissile,
		},

		"Planetary Ion Cannon": {
			minLevel:       3,
			maxLevel:       9,
			componentStats: PlanetaryHeavyIonCannon,
		},
	}
)

func (j *Job) applyIonWeapons() (err error) {

	// apply stats for each component
	err = j.ApplyComponentAll(IonComponentData)

	return
}
