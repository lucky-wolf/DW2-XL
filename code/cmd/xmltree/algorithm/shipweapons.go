package algorithm

import (
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
)

func KineticShipWeapons(folder string) (err error) {

	log.Println("All kinetic ship weapons will be scaled to XL values")

	// load all component definition files
	j, err := LoadJobFor(folder, "ComponentDefinitions*.xml")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.ApplyComponentAll(KineticWeaponData)
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

type ComponentName = string
type WeaponFamilyName = string

type WeaponBasis map[string]xmltree.SimpleValue

type SimpleValueFunc = func(level int) xmltree.SimpleValue
type SimpleValuesTable = map[string]SimpleValueFunc

// type ObjectDictionary map[ComponentName]

// subtle: level starts at one, NOT zero!

func MakeFixedLevelFunc(basis float64) LevelFunc {
	return func(level int) float64 { return basis }
}

func MakeLinearLevelFunc(basis float64, add float64) LevelFunc {
	return func(level int) float64 { return basis + (add * float64(level-1)) }
}

func MakeScalingLevelFunc(basis float64, scale float64) LevelFunc {
	return func(level int) float64 { return basis * (1.0 + scale*float64(level-1)) }
}

func MakeScaledFuncLevelFunc(scale float64, levelfunc LevelFunc) LevelFunc {
	return func(level int) float64 { return scale * levelfunc(level) }
}

func MakeOffsetFuncLevelFunc(offset int, levelfunc LevelFunc) LevelFunc {
	return func(level int) float64 { return levelfunc(level + offset) }
}

// 775, 800, ...

var (

	// warn: we number from 1..11 where 1 = t0, and 2,3,...,10 = t2,t3,...,t1l0
	KineticWeaponData = map[string]ComponentData{
		"Long Range Cannon [S]": {
			minLevel:    1,
			maxLevel:    1,
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
		"Rail Gun [L]": {
			minLevel:    2,
			maxLevel:    4,
			fieldValues: LargeKineticWeaponComponentStats,
		},
	}

	BasicKineticWeaponSpeed         = MakeFixedLevelFunc(750)
	BasicKineticWeaponDamage        = MakeScalingLevelFunc(19.2, 0.17)
	BasicKineticWeaponEnergyPerShot = MakeScaledFuncLevelFunc(0.5, BasicKineticWeaponDamage)

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
		"WeaponDamageFalloffRatio":          MakeFixedLevelFunc(0),       // rail guns don't lose dmg with distance
		"WeaponArmorBypass":                 MakeFixedLevelFunc(-0.3333), // rail guns are especially affected by armor
		"WeaponShieldBypass":                MakeFixedLevelFunc(0),       // rail guns have no special interaction with shields
		"WeaponRange":                       MakeScaledFuncLevelFunc(1.5, BasicKineticWeaponSpeed),
		"WeaponSpeed":                       BasicKineticWeaponSpeed,
		"WeaponFireRate":                    MakeFixedLevelFunc(12),
		"WeaponEnergyPerShot":               BasicKineticWeaponEnergyPerShot, // todo: drive this off of damage & family?
		"WeaponRawDamage":                   BasicKineticWeaponDamage,        // compounding per level
		"WeaponVolleyAmount":                MakeFixedLevelFunc(1),
		"WeaponVolleyFireRate":              MakeFixedLevelFunc(0),
	}

	SmallKineticWeaponSpeed          = MakeLinearLevelFunc(775, 25) // t2 = 800
	SmallKineticWeaponComponentStats = ExtendValuesTable(
		BasicKineticWeaponComponentStats,
		ComponentStats{
			"WeaponSpeed": SmallKineticWeaponSpeed,
			"WeaponRange": MakeScaledFuncLevelFunc(1.5, SmallKineticWeaponSpeed),
		},
	)

	MediumKineticWeaponSpeed          = MakeLinearLevelFunc(800, 50) // t2 = 850
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

	LargeKineticWeaponSpeed          = MakeLinearLevelFunc(800, 100) // t2 = 900
	LargeKineticWeaponEnergyPerShot  = MakeScaledFuncLevelFunc(2, MediumKineticWeaponEnergyPerShot)
	LargeKineticWeaponDamage         = MakeScaledFuncLevelFunc(2, MediumKineticWeaponDamage)
	LargeKineticWeaponComponentStats = ExtendValuesTable(
		SmallKineticWeaponComponentStats,
		ComponentStats{
			"CrewRequirement":     MakeFixedLevelFunc(12),
			"WeaponEnergyPerShot": LargeKineticWeaponEnergyPerShot,
			"WeaponRawDamage":     LargeKineticWeaponDamage,
			"WeaponSpeed":         LargeKineticWeaponSpeed,
			"WeaponRange":         MakeScaledFuncLevelFunc(2, LargeKineticWeaponSpeed),
		},
	)
)
