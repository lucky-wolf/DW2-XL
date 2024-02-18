package algorithm

import (
	"log"
)

func KineticWeapons(folder string) (err error) {

	log.Println("All kinetic weapons will be scaled to XL values")

	// load all component definition files
	j, err := LoadJobFor(folder, "ComponentDefinitions*.xml")
	if err != nil {
		return
	}

	// update kinetic weapons
	err = j.applyKineticWeapons()
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

func (j *Job) applyKineticWeapons() (err error) {

	// apply stats for each component
	err = j.ApplyComponentAll(KineticWeaponData)

	return
}

const (
	KineticEnergyPerShotDamageRatio = 0.5
)

var (

	// warn: we number from 1..11 where 1 = t0, and 2,3,...,10 = t2,t3,...,t1l0
	KineticWeaponData = map[string]ComponentData{
		"Long Range Cannon [S]": {
			minLevel:    0,
			maxLevel:    0,
			fieldValues: BasicKineticWeaponComponentStats,
		},
		"Rail Gun [S]": {
			minLevel:    2,
			maxLevel:    4,
			fieldValues: SmallKineticWeaponComponentStats,
		},
		"Rail Gun [M]": {
			minLevel:    2,
			maxLevel:    4,
			fieldValues: MediumKineticWeaponComponentStats,
		},
		"Heavy Rail Gun [L]": {
			minLevel:    3,
			maxLevel:    4,
			fieldValues: LargeKineticWeaponComponentStats,
		},
		"Forge Rail Battery [M]": {
			minLevel:    5,
			maxLevel:    8,
			fieldValues: MediumForgeRailWeaponComponentStats,
		},
		"Forge Rail Battery [L]": {
			minLevel:    5,
			maxLevel:    8,
			fieldValues: LargeForgeRailWeaponComponentStats,
		},
		"Terminator Autocannon [S]": {
			minLevel:    5,
			maxLevel:    8,
			fieldValues: AutocannonWeaponComponentStats,
		},
		"Hail Cannon [S]": {
			minLevel:    5,
			maxLevel:    8,
			fieldValues: HailCannonWeaponComponentStats,
		},
		// "Planetary Cannon Batteries": {
		// 	minLevel:    5,
		// 	maxLevel:    8,
		// 	fieldValues: HailCannonWeaponComponentStats,
		// },
	}

	BasicKineticWeaponArmorBypass   = MakeFixedLevelFunc(-1. / 3.)
	BasicKineticWeaponSpeed         = MakeLinearLevelFunc(750, 25)
	BasicKineticWeaponDamage        = MakeExpLevelFunc(19.2, 0.18)
	BasicKineticWeaponEnergyPerShot = MakeScaledFuncLevelFunc(KineticEnergyPerShotDamageRatio, BasicKineticWeaponDamage)
	BasicKineticWeaponROF           = MakeFixedLevelFunc(12)

	// note: ion weapons never have any bombard value (lighting in atmosphere is not a real issue)
	BasicKineticWeaponComponentStats = ComponentStats{
		// "ComponentCountermeasuresBonus":     MakeLevelFunc(WeaponFamilyValues["Rail Guns"]["ComponentCountermeasuresBonus"], .02),
		"ComponentCountermeasuresBonus":     DirectFireComponentCountermeasuresBonus,
		"ComponentTargetingBonus":           MakeFixedLevelFunc(0),
		"CrewRequirement":                   MakeFixedLevelFunc(5),
		"StaticEnergyUsed":                  MakeFixedLevelFunc(1),
		"WeaponAreaEffectRange":             MakeFixedLevelFunc(0),
		"WeaponAreaBlastWaveSpeed":          MakeFixedLevelFunc(0),
		"WeaponBombardDamageInfrastructure": MakeFixedLevelFunc(0),
		"WeaponBombardDamageMilitary":       MakeFixedLevelFunc(0),
		"WeaponBombardDamagePopulation":     MakeFixedLevelFunc(0),
		"WeaponBombardDamageQuality":        MakeFixedLevelFunc(0),
		"WeaponIonEngineDamage":             MakeFixedLevelFunc(0),
		"WeaponIonHyperDriveDamage":         MakeFixedLevelFunc(0),
		"WeaponIonSensorDamage":             MakeFixedLevelFunc(0),
		"WeaponIonShieldDamage":             MakeFixedLevelFunc(0),
		"WeaponIonWeaponDamage":             MakeFixedLevelFunc(0),
		"WeaponIonGeneralDamage":            MakeFixedLevelFunc(0),
		"WeaponDamageFalloffRatio":          MakeFixedLevelFunc(0),         // rail guns don't lose dmg with distance
		"WeaponArmorBypass":                 BasicKineticWeaponArmorBypass, // rail guns are especially affected by armor
		"WeaponShieldBypass":                MakeFixedLevelFunc(0),         // rail guns have no special interaction with shields
		"WeaponRange":                       MakeScaledFuncLevelFunc(1.5, BasicKineticWeaponSpeed),
		"WeaponSpeed":                       BasicKineticWeaponSpeed,
		"WeaponFireRate":                    BasicKineticWeaponROF,
		"WeaponEnergyPerShot":               BasicKineticWeaponEnergyPerShot, // todo: drive this off of damage & family?
		"WeaponRawDamage":                   BasicKineticWeaponDamage,        // compounding per level
		"WeaponVolleyAmount":                MakeFixedLevelFunc(1),
		"WeaponVolleyFireRate":              MakeFixedLevelFunc(0),
	}

	SmallKineticWeaponSpeed          = BasicKineticWeaponSpeed // t2 = 800
	SmallKineticWeaponComponentStats = ExtendValuesTable(
		BasicKineticWeaponComponentStats,
		ComponentStats{
			"WeaponSpeed": SmallKineticWeaponSpeed,
			"WeaponRange": MakeScaledFuncLevelFunc(1.5, SmallKineticWeaponSpeed),
		},
	)

	MediumKineticWeaponSpeed          = MakeLinearLevelFunc(750, 50) // t2 = 850
	MediumKineticWeaponDamage         = MakeScaledFuncLevelFunc(2, BasicKineticWeaponDamage)
	MediumKineticWeaponEnergyPerShot  = MakeScaledFuncLevelFunc(2, BasicKineticWeaponEnergyPerShot)
	MediumKineticWeaponComponentStats = ExtendValuesTable(
		SmallKineticWeaponComponentStats,
		ComponentStats{
			"CrewRequirement":     MakeFixedLevelFunc(8),
			"WeaponEnergyPerShot": MediumKineticWeaponEnergyPerShot,
			"WeaponRawDamage":     MediumKineticWeaponDamage,
			"WeaponSpeed":         MediumKineticWeaponSpeed,
			"WeaponRange":         MakeScaledFuncLevelFunc(1.75, MediumKineticWeaponSpeed),
		},
	)

	LargeKineticWeaponSpeed          = MakeLinearLevelFunc(700, 100) // t2 = 900
	LargeKineticWeaponEnergyPerShot  = MakeScaledFuncLevelFunc(2, MediumKineticWeaponEnergyPerShot)
	LargeKineticWeaponDamage         = MakeScaledFuncLevelFunc(2, MediumKineticWeaponDamage)
	LargeKineticWeaponComponentStats = ExtendValuesTable(
		MediumKineticWeaponComponentStats,
		ComponentStats{
			"CrewRequirement":     MakeFixedLevelFunc(12),
			"WeaponEnergyPerShot": LargeKineticWeaponEnergyPerShot,
			"WeaponRawDamage":     LargeKineticWeaponDamage,
			"WeaponSpeed":         LargeKineticWeaponSpeed,
			"WeaponRange":         MakeScaledFuncLevelFunc(2, LargeKineticWeaponSpeed),
			// "WeaponBombardDamageInfrastructure": MakeScalingLevelFunc(0),
			// "WeaponBombardDamageMilitary":       MakeScalingLevelFunc(0),
			// "WeaponBombardDamagePopulation":     MakeScalingLevelFunc(0),
			// "WeaponBombardDamageQuality":        MakeScalingLevelFunc(0),
		},
	)

	// note: Forge Batteries fire 50% slower, and do 50% more damage per stroke
	// note they also reduce their armor penalty by 50%

	MediumForgeRailWeaponEnergyPerShot  = MakeScaledFuncLevelFunc(KineticEnergyPerShotDamageRatio, MediumForgeRailWeaponDamage)
	MediumForgeRailWeaponDamage         = MakeScaledFuncLevelFunc(1.5, MediumKineticWeaponDamage)
	MediumForgeRailWeaponComponentStats = ExtendValuesTable(
		MediumKineticWeaponComponentStats,
		ComponentStats{
			"WeaponEnergyPerShot": MediumForgeRailWeaponEnergyPerShot,
			"WeaponRawDamage":     MediumForgeRailWeaponDamage,
			"WeaponArmorBypass":   MakeScaledFuncLevelFunc(.5, BasicKineticWeaponArmorBypass), // forge rail guns are less bothered by armor
			"WeaponFireRate":      MakeScaledFuncLevelFunc(1.5, BasicKineticWeaponROF),        // 50% slower
		},
	)

	LargeForgeRailWeaponEnergyPerShot  = MakeScaledFuncLevelFunc(KineticEnergyPerShotDamageRatio, LargeForgeRailWeaponDamage)
	LargeForgeRailWeaponDamage         = MakeScaledFuncLevelFunc(2, MediumForgeRailWeaponDamage)
	LargeForgeRailWeaponComponentStats = ExtendValuesTable(
		LargeKineticWeaponComponentStats,
		ComponentStats{
			"WeaponEnergyPerShot": LargeForgeRailWeaponEnergyPerShot,
			"WeaponRawDamage":     LargeForgeRailWeaponDamage,
			"WeaponArmorBypass":   MakeScaledFuncLevelFunc(.5, BasicKineticWeaponArmorBypass), // forge rail guns are less bothered by armor
			"WeaponFireRate":      MakeScaledFuncLevelFunc(1.5, BasicKineticWeaponROF),        // 50% slower
		},
	)

	AutocannonWeaponDamage         = MakeScaledFuncLevelFunc(0.5, BasicKineticWeaponDamage)
	AutocannonWeaponEnergyPerShot  = MakeScaledFuncLevelFunc(KineticEnergyPerShotDamageRatio, AutocannonWeaponDamage)
	AutocannonRateOfFire           = MakeScaledFuncLevelFunc(.5, BasicKineticWeaponROF)
	AutocannonWeaponComponentStats = ExtendValuesTable(
		BasicKineticWeaponComponentStats,
		ComponentStats{
			"WeaponRawDamage":     AutocannonWeaponDamage,
			"WeaponEnergyPerShot": AutocannonWeaponEnergyPerShot,
			"WeaponFireRate":      AutocannonRateOfFire, // 2x rof
		},
	)

	HailCannonRateOfFire           = MakeScaledFuncLevelFunc(.8, AutocannonRateOfFire) // 25% faster than autocannons
	HailCannonWeaponComponentStats = ExtendValuesTable(
		AutocannonWeaponComponentStats,
		ComponentStats{
			"WeaponArmorBypass":  MakeFixedLevelFunc(0), // hail cannons are armor neutral
			"WeaponShieldBypass": MakeFixedLevelFunc(0), // rail guns have no special interaction with shields
			"WeaponFireRate":     HailCannonRateOfFire,  // 25% faster than autocannons
		},
	)

	// todo: planetary facility
	// PlanetaryRailWeaponEnergyPerShot  = MakeScaledFuncLevelFunc(KineticEnergyPerShotDamageRatio, PlanetaryRailWeaponDamage)
	// PlanetaryRailWeaponDamage         = MakeScaledFuncLevelFunc(2, MediumForgeRailWeaponDamage)
	// PlanetaryRailWeaponComponentStats = ExtendValuesTable(
	// 	LargeKineticWeaponComponentStats,
	// 	ComponentStats{
	// 		"WeaponEnergyPerShot": PlanetaryRailWeaponEnergyPerShot,
	// 		"WeaponRawDamage":     PlanetaryRailWeaponDamage,
	// 		"WeaponArmorBypass":   MakeScaledFuncLevelFunc(.5, BasicKineticWeaponArmorBypass), // forge rail guns are less bothered by armor
	// 		"WeaponFireRate":      MakeScaledFuncLevelFunc(1.5, BasicKineticWeaponROF),        // 50% slower
	// 	},
	// )
)
