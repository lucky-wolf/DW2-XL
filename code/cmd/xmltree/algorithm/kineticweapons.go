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
	KineticBaseDamageFactor = 1.1 // base damage is 10% better than blasters
	KineticBaseRateOfFire   = 12
	KineticEnergyRatio      = 0.5
)

var (
	KineticWeaponData = map[string]ComponentData{
		"Long Range Cannon [S]": {
			minLevel:       0,
			maxLevel:       0,
			componentStats: BasicKineticWeaponComponentStats,
		},
		"Rail Gun [S]": {
			minLevel:       2,
			maxLevel:       4,
			componentStats: SmallKineticWeaponComponentStats,
		},
		"Rail Gun [M]": {
			minLevel:       2,
			maxLevel:       4,
			componentStats: MediumKineticWeaponComponentStats,
		},
		"Heavy Rail Gun [L]": {
			minLevel:       3,
			maxLevel:       4,
			componentStats: LargeKineticWeaponComponentStats,
		},
		"Forge Rail Battery [M]": {
			minLevel:       5,
			maxLevel:       10,
			componentStats: MediumForgeRailWeaponComponentStats,
		},
		"Forge Rail Battery [L]": {
			minLevel:       5,
			maxLevel:       10,
			componentStats: LargeForgeRailWeaponComponentStats,
		},
		"Terminator Autocannon [S]": {
			minLevel:       5,
			maxLevel:       10,
			componentStats: AutocannonWeaponComponentStats,
		},
		"Hail Cannon [S]": {
			minLevel:       5,
			maxLevel:       10,
			componentStats: HailCannonWeaponComponentStats,
		},

		// TODO: in order to make this work, we have to either have a facility per level
		// todo: or a level-stride (which is always even)
		// todo: or use different scaling functions that take the "stride" into account
		// todo: or refer to levels as they fall on the tech tree (v1, v4, v7)
		"Planetary Forge Battery": {
			minLevel:       3,
			maxLevel:       9,
			componentStats: PlanetaryForgeBatteryComponentStats,
		},
	}

	BasicKineticWeaponArmorBypass   = MakeFixedLevelFunc(-1. / 3.)
	BasicKineticWeaponSpeed         = MakeLinearLevelFunc(750, 25)
	BasicKineticWeaponDamage        = MakeExpLevelFunc(KineticBaseDamageFactor*BlasterBaseDamage*KineticBaseRateOfFire/BlasterRateOfFire, WeaponDamageIncreaseExp)
	BasicKineticWeaponEnergyPerShot = MakeScaledFuncLevelFunc(KineticEnergyRatio, BasicKineticWeaponDamage)
	BasicKineticWeaponROF           = MakeFixedLevelFunc(KineticBaseRateOfFire)

	// note: ion weapons never have any bombard value (lighting in atmosphere is not a real issue)
	BasicKineticWeaponComponentStats = ComponentStats{
		// "ComponentCountermeasuresBonus":     MakeLevelFunc(WeaponFamilyValues["Rail Guns"]["ComponentCountermeasuresBonus"], .02),
		"ComponentCountermeasuresBonus":     DirectFireComponentCountermeasuresBonus,
		"ComponentTargetingBonus":           MakeFixedLevelFunc(0),
		"CrewRequirement":                   SmallCrewRequirements,
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
	SmallKineticWeaponComponentStats = ComposeComponentStats(
		BasicKineticWeaponComponentStats,
		ComponentStats{
			"WeaponSpeed": SmallKineticWeaponSpeed,
			"WeaponRange": MakeScaledFuncLevelFunc(1.5, SmallKineticWeaponSpeed),
		},
	)

	MediumKineticWeaponSpeed          = MakeLinearLevelFunc(750, 50) // t2 = 850
	MediumKineticWeaponDamage         = MakeScaledFuncLevelFunc(2, BasicKineticWeaponDamage)
	MediumKineticWeaponEnergyPerShot  = MakeScaledFuncLevelFunc(2, BasicKineticWeaponEnergyPerShot)
	MediumKineticWeaponComponentStats = ComposeComponentStats(
		SmallKineticWeaponComponentStats,
		ComponentStats{
			"CrewRequirement":     MediumCrewRequirements,
			"WeaponEnergyPerShot": MediumKineticWeaponEnergyPerShot,
			"WeaponRawDamage":     MediumKineticWeaponDamage,
			"WeaponSpeed":         MediumKineticWeaponSpeed,
			"WeaponRange":         MakeScaledFuncLevelFunc(1.75, MediumKineticWeaponSpeed),
		},
	)

	LargeKineticWeaponSpeed          = MakeLinearLevelFunc(700, 100) // t2 = 900
	LargeKineticWeaponDamage         = MakeScaledFuncLevelFunc(2, MediumKineticWeaponDamage)
	LargeKineticWeaponEnergyPerShot  = MakeScaledFuncLevelFunc(2, MediumKineticWeaponEnergyPerShot)
	LargeKineticWeaponComponentStats = ComposeComponentStats(
		MediumKineticWeaponComponentStats,
		ComponentStats{
			"CrewRequirement":     LargeCrewRequirements,
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

	MediumForgeRailWeaponDamage         = MakeScaledFuncLevelFunc(1.5, MediumKineticWeaponDamage)
	MediumForgeRailWeaponEnergyPerShot  = MakeScaledFuncLevelFunc(KineticEnergyRatio, MediumForgeRailWeaponDamage)
	MediumForgeRailWeaponComponentStats = ComposeComponentStats(
		MediumKineticWeaponComponentStats,
		ComponentStats{
			"WeaponEnergyPerShot": MediumForgeRailWeaponEnergyPerShot,
			"WeaponRawDamage":     MediumForgeRailWeaponDamage,
			"WeaponArmorBypass":   MakeScaledFuncLevelFunc(.5, BasicKineticWeaponArmorBypass), // forge rail guns are less bothered by armor
			"WeaponFireRate":      MakeScaledFuncLevelFunc(1.5, BasicKineticWeaponROF),        // 50% slower
		},
	)

	LargeForgeRailWeaponDamage         = MakeScaledFuncLevelFunc(2, MediumForgeRailWeaponDamage)
	LargeForgeRailWeaponEnergyPerShot  = MakeScaledFuncLevelFunc(KineticEnergyRatio, LargeForgeRailWeaponDamage)
	LargeForgeRailWeaponROF            = MakeScaledFuncLevelFunc(1.5, BasicKineticWeaponROF)
	LargeForgeRailWeaponComponentStats = ComposeComponentStats(
		LargeKineticWeaponComponentStats,
		ComponentStats{
			"WeaponEnergyPerShot": LargeForgeRailWeaponEnergyPerShot,
			"WeaponRawDamage":     LargeForgeRailWeaponDamage,
			"WeaponArmorBypass":   MakeScaledFuncLevelFunc(.5, BasicKineticWeaponArmorBypass), // forge rail guns are less bothered by armor
			"WeaponFireRate":      LargeForgeRailWeaponROF,                                    // 50% slower
		},
	)

	AutocannonWeaponDamage         = MakeScaledFuncLevelFunc(0.5, BasicKineticWeaponDamage)
	AutocannonWeaponEnergyPerShot  = MakeScaledFuncLevelFunc(KineticEnergyRatio, AutocannonWeaponDamage)
	AutocannonRateOfFire           = MakeScaledFuncLevelFunc(.5, BasicKineticWeaponROF)
	AutocannonWeaponComponentStats = ComposeComponentStats(
		BasicKineticWeaponComponentStats,
		ComponentStats{
			"WeaponRawDamage":     AutocannonWeaponDamage,
			"WeaponEnergyPerShot": AutocannonWeaponEnergyPerShot,
			"WeaponFireRate":      AutocannonRateOfFire, // 2x rof
		},
	)

	HailCannonRateOfFire           = MakeScaledFuncLevelFunc(.8, AutocannonRateOfFire) // 25% faster than autocannons
	HailCannonWeaponComponentStats = ComposeComponentStats(
		AutocannonWeaponComponentStats,
		ComponentStats{
			"WeaponArmorBypass":  MakeFixedLevelFunc(0), // hail cannons are armor neutral
			"WeaponShieldBypass": MakeFixedLevelFunc(0), // rail guns have no special interaction with shields
			"WeaponFireRate":     HailCannonRateOfFire,  // 25% faster than autocannons
		},
	)

	// note: we only care about a few of these (v1, v4, v7)
	PlanetaryRailWeaponDamage           = MakeScaledFuncLevelFunc(2, LargeForgeRailWeaponDamage)
	PlanetaryRailWeaponEnergyPerShot    = MakeScaledFuncLevelFunc(KineticEnergyRatio, PlanetaryRailWeaponDamage)
	PlanetaryForgeBatteryComponentStats = ComposeComponentStats(
		LargeForgeRailWeaponComponentStats,
		ComponentStats{
			"CrewRequirement":               PlanetaryCrewRequirements,
			"ComponentCountermeasuresBonus": MakeScaledFuncLevelFunc(2, DirectFireComponentCountermeasuresBonus),
			"ComponentTargetingBonus":       MakeScaledFuncLevelFunc(2, DirectFireComponentCountermeasuresBonus),
			"WeaponEnergyPerShot":           PlanetaryRailWeaponEnergyPerShot,
			"WeaponRawDamage":               PlanetaryRailWeaponDamage,
			"WeaponRange":                   MakeScaledFuncLevelFunc(4, LargeKineticWeaponSpeed),
			"WeaponFireRate":                MakeScaledFuncLevelFunc(.5, LargeForgeRailWeaponROF), // 100% faster
		},
	)
)
