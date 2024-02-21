package algorithm

import (
	"log"
)

func BlasterWeapons(folder string) (err error) {

	log.Println("All blaster weapons will be scaled to XL values")

	// load all component definition files
	j, err := LoadJobFor(folder, "ComponentDefinitions*.xml")
	if err != nil {
		return
	}

	// update blaster weapons
	err = j.applyBlasterWeapons()
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

func (j *Job) applyBlasterWeapons() (err error) {

	// apply stats for each component
	err = j.ApplyComponentAll(BlasterWeaponData)

	return
}

const (
	BlasterWeaponSpeed    = 2200
	BlasterBaseDamage     = 12
	BlasterVolleyFireRate = .5

	BlasterEnergyRatio = .8
	ImpactEnergyRatio  = .7
	TitanEnergyRatio   = .9
	PhaserEnergyRatio  = .6
	PlasmaEnergyRatio  = .9 // note: high damage = high multiplier already

	BlasterFalloff       = .25
	ImpactBlasterFalloff = .2
	TitanBlasterFalloff  = .15
	PhaserBlasterFalloff = .1
	PlasmaBlasterFalloff = .25

	BlasterRateOfFire       = 9
	ImpactBlasterRateOfFire = 11 // fires a bit slower with larger bolts
	TitanBlasterRateOfFire  = 10
	PhaserBlasterRateOfFire = 12
	PlasmaBlasterRateOfFire = 10

	// plasma interacts with shields and armor
	PlasmaShieldBypass = 0
	PlasmaArmorBypass  = 0.3333333

	// phasers interact with shields and armor
	PhaserShieldBypass = 0.666666
	PhaserArmorBypass  = 0.3333333
)

var (
	BlasterWeaponData = map[string]ComponentData{
		"Laser Cannon [S]": {
			minLevel:       0,
			maxLevel:       0,
			componentStats: BasicBlasterWeaponComponentStats,
		},
		"Maxos Blaster [S]": {
			minLevel:       2,
			maxLevel:       4,
			componentStats: SmallBlasterWeaponComponentStats,
		},
		"Maxos Blaster [M]": {
			minLevel:       2,
			maxLevel:       4,
			componentStats: MediumBlasterWeaponComponentStats,
		},

		"Impact Assault Blaster [S]": {
			minLevel:       5,
			maxLevel:       8,
			componentStats: SmallImpactBlasterWeaponComponentStats,
		},
		"Impact Assault Blaster [M]": {
			minLevel:       5,
			maxLevel:       8,
			componentStats: MediumImpactBlasterWeaponComponentStats,
		},
		"Impact Assault Blaster [L]": {
			minLevel:       5,
			maxLevel:       8,
			componentStats: LargeImpactBlasterWeaponComponentStats,
		},

		"Titan Blaster [S]": {
			minLevel:       9,
			maxLevel:       10,
			componentStats: SmallTitanBlasterWeaponComponentStats,
		},
		"Titan Blaster [M]": {
			minLevel:       9,
			maxLevel:       10,
			componentStats: MediumTitanBlasterWeaponComponentStats,
		},
		"Titan Blaster [L]": {
			minLevel:       9,
			maxLevel:       10,
			componentStats: LargeTitanBlasterWeaponComponentStats,
		},

		"Phaser Blaster [S]": {
			minLevel:       9,
			maxLevel:       10,
			componentStats: SmallPhaserBlasterWeaponComponentStats,
		},
		"Phaser Blaster [M]": {
			minLevel:       9,
			maxLevel:       10,
			componentStats: MediumPhaserBlasterWeaponComponentStats,
		},
		"Phaser Blaster [L]": {
			minLevel:       9,
			maxLevel:       10,
			componentStats: LargePhaserBlasterWeaponComponentStats,
		},

		// Boskaran Plasma Cannons
		"Plasma Cannon [S]": {
			minLevel:       2,
			maxLevel:       10,
			componentStats: SmallPlasmaBlasterWeaponComponentStats,
		},
		"Plasma Cannon [M]": {
			minLevel:       2,
			maxLevel:       10,
			componentStats: MediumPlasmaBlasterWeaponComponentStats,
		},
		"Plasma Cannon [L]": {
			minLevel:       2,
			maxLevel:       10,
			componentStats: LargePlasmaBlasterWeaponComponentStats,
		},

		"Planetary Blaster Emplacement": {
			minLevel:       3,
			maxLevel:       9,
			componentStats: PlanetaryBlasterComponentStats,
		},
	}

	BasicBlasterFalloffRatio        = MakeFixedLevelFunc(BlasterFalloff)
	BasicBlasterWeaponDamage        = MakeExpLevelFunc(BlasterBaseDamage, WeaponDmgIncreaseExp)
	BasicBlasterWeaponEnergyPerShot = MakeScaledFuncLevelFunc(BlasterEnergyRatio, BasicBlasterWeaponDamage)
	BasicBlasterWeaponRange         = MakeExpLevelFunc(1000, .05) // 5% level over level
	BasicBlasterWeaponROF           = MakeFixedLevelFunc(BlasterRateOfFire)
	BasicBlasterWeaponSpeed         = MakeFixedLevelFunc(BlasterWeaponSpeed)

	// warn: blasters do NOT have bombard values because they're always just repeating small weapons

	BasicBlasterWeaponComponentStats = ComponentStats{
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
		"WeaponDamageFalloffRatio":          BasicBlasterFalloffRatio,
		"WeaponArmorBypass":                 MakeFixedLevelFunc(0),
		"WeaponShieldBypass":                MakeFixedLevelFunc(0),
		"WeaponRange":                       BasicBlasterWeaponRange,
		"WeaponSpeed":                       BasicBlasterWeaponSpeed,
		"WeaponFireRate":                    BasicBlasterWeaponROF,
		"WeaponEnergyPerShot":               BasicBlasterWeaponEnergyPerShot,
		"WeaponRawDamage":                   BasicBlasterWeaponDamage,
		"WeaponVolleyAmount":                MakeFixedLevelFunc(1),
		"WeaponVolleyFireRate":              MakeFixedLevelFunc(0),
	}

	SmallBlasterWeaponComponentStats = ExtendValuesTable(
		BasicBlasterWeaponComponentStats,
	)

	MediumBlasterWeaponComponentStats = ExtendValuesTable(
		SmallBlasterWeaponComponentStats,
		ComponentStats{
			"CrewRequirement":      MediumCrewRequirements,
			"WeaponVolleyAmount":   MakeFixedLevelFunc(2),
			"WeaponVolleyFireRate": MakeFixedLevelFunc(BlasterVolleyFireRate),
		},
	)

	LargeBlasterWeaponComponentStats = ExtendValuesTable(
		SmallBlasterWeaponComponentStats,
		ComponentStats{
			"CrewRequirement":      LargeCrewRequirements,
			"WeaponVolleyAmount":   MakeFixedLevelFunc(4),
			"WeaponVolleyFireRate": MakeFixedLevelFunc(BlasterVolleyFireRate),
		},
	)

	// Impact Assault are a 20% improvement
	ImpactBlasterFalloffRatio        = MakeFixedLevelFunc(ImpactBlasterFalloff)
	ImpactBlasterWeaponDamage        = MakeScaledFuncLevelFunc(1.1*ImpactBlasterRateOfFire/BlasterRateOfFire, BasicBlasterWeaponDamage)
	ImpactBlasterWeaponEnergyPerShot = MakeScaledFuncLevelFunc(ImpactEnergyRatio, ImpactBlasterWeaponDamage)
	ImpactBlasterWeaponRange         = BasicBlasterWeaponRange // it's already compounding
	ImpactBlasterWeaponROF           = MakeFixedLevelFunc(ImpactBlasterRateOfFire)

	SmallImpactBlasterWeaponComponentStats = ExtendValuesTable(
		BasicBlasterWeaponComponentStats,
		ComponentStats{
			"WeaponDamageFalloffRatio": ImpactBlasterFalloffRatio,
			"WeaponRange":              ImpactBlasterWeaponRange,
			"WeaponFireRate":           ImpactBlasterWeaponROF,
			"WeaponEnergyPerShot":      ImpactBlasterWeaponEnergyPerShot,
			"WeaponRawDamage":          ImpactBlasterWeaponDamage,
		},
	)

	MediumImpactBlasterWeaponComponentStats = ExtendValuesTable(
		SmallImpactBlasterWeaponComponentStats,
		ComponentStats{
			"CrewRequirement":      MediumCrewRequirements,
			"WeaponVolleyAmount":   MakeFixedLevelFunc(2),
			"WeaponVolleyFireRate": MakeFixedLevelFunc(BlasterVolleyFireRate),
		},
	)

	LargeImpactBlasterWeaponComponentStats = ExtendValuesTable(
		SmallImpactBlasterWeaponComponentStats,
		ComponentStats{
			"CrewRequirement":      LargeCrewRequirements,
			"WeaponVolleyAmount":   MakeFixedLevelFunc(4),
			"WeaponVolleyFireRate": MakeFixedLevelFunc(BlasterVolleyFireRate),
		},
	)

	// Titan Assault are a 20% improvement
	TitanBlasterFalloffRatio        = MakeFixedLevelFunc(TitanBlasterFalloff)
	TitanBlasterWeaponDamage        = MakeScaledFuncLevelFunc(1.1*TitanBlasterRateOfFire/ImpactBlasterRateOfFire, ImpactBlasterWeaponDamage)
	TitanBlasterWeaponEnergyPerShot = MakeScaledFuncLevelFunc(TitanEnergyRatio, TitanBlasterWeaponDamage)
	TitanBlasterWeaponRange         = BasicBlasterWeaponRange // it's already compounding
	TitanBlasterWeaponROF           = MakeFixedLevelFunc(TitanBlasterRateOfFire)

	SmallTitanBlasterWeaponComponentStats = ExtendValuesTable(
		BasicBlasterWeaponComponentStats,
		ComponentStats{
			"WeaponDamageFalloffRatio": TitanBlasterFalloffRatio,
			"WeaponRange":              TitanBlasterWeaponRange,
			"WeaponFireRate":           TitanBlasterWeaponROF,
			"WeaponEnergyPerShot":      TitanBlasterWeaponEnergyPerShot,
			"WeaponRawDamage":          TitanBlasterWeaponDamage,
		},
	)

	MediumTitanBlasterWeaponComponentStats = ExtendValuesTable(
		SmallTitanBlasterWeaponComponentStats,
		ComponentStats{
			"CrewRequirement":      MediumCrewRequirements,
			"WeaponVolleyAmount":   MakeFixedLevelFunc(2),
			"WeaponVolleyFireRate": MakeFixedLevelFunc(BlasterVolleyFireRate),
		},
	)

	LargeTitanBlasterWeaponComponentStats = ExtendValuesTable(
		SmallTitanBlasterWeaponComponentStats,
		ComponentStats{
			"CrewRequirement":      LargeCrewRequirements,
			"WeaponVolleyAmount":   MakeFixedLevelFunc(4),
			"WeaponVolleyFireRate": MakeFixedLevelFunc(BlasterVolleyFireRate),
		},
	)

	// Phaser Assault are a 20% improvement
	PhaserBlasterFalloffRatio        = MakeFixedLevelFunc(PhaserBlasterFalloff)
	PhaserBlasterWeaponDamage        = MakeScaledFuncLevelFunc(0.75*PhaserBlasterRateOfFire/ImpactBlasterRateOfFire, ImpactBlasterWeaponDamage)
	PhaserBlasterWeaponEnergyPerShot = MakeScaledFuncLevelFunc(PhaserEnergyRatio, PhaserBlasterWeaponDamage)
	PhaserBlasterWeaponRange         = BasicBlasterWeaponRange // it's already compounding
	PhaserBlasterWeaponROF           = MakeFixedLevelFunc(PhaserBlasterRateOfFire)

	SmallPhaserBlasterWeaponComponentStats = ExtendValuesTable(
		BasicBlasterWeaponComponentStats,
		ComponentStats{
			"WeaponDamageFalloffRatio": PhaserBlasterFalloffRatio,
			"WeaponRange":              PhaserBlasterWeaponRange,
			"WeaponFireRate":           PhaserBlasterWeaponROF,
			"WeaponEnergyPerShot":      PhaserBlasterWeaponEnergyPerShot,
			"WeaponRawDamage":          PhaserBlasterWeaponDamage,
			"WeaponShieldBypass":       MakeFixedLevelFunc(PhaserShieldBypass),
			"WeaponArmorBypass":        MakeFixedLevelFunc(PhaserArmorBypass),
		},
	)

	MediumPhaserBlasterWeaponComponentStats = ExtendValuesTable(
		SmallPhaserBlasterWeaponComponentStats,
		ComponentStats{
			"CrewRequirement":      MediumCrewRequirements,
			"WeaponVolleyAmount":   MakeFixedLevelFunc(2),
			"WeaponVolleyFireRate": MakeFixedLevelFunc(BlasterVolleyFireRate),
		},
	)

	LargePhaserBlasterWeaponComponentStats = ExtendValuesTable(
		SmallPhaserBlasterWeaponComponentStats,
		ComponentStats{
			"CrewRequirement":      LargeCrewRequirements,
			"WeaponVolleyAmount":   MakeFixedLevelFunc(4),
			"WeaponVolleyFireRate": MakeFixedLevelFunc(BlasterVolleyFireRate),
		},
	)

	// Boskaran Plasma Cannons
	PlasmaBlasterFalloffRatio        = MakeFixedLevelFunc(PlasmaBlasterFalloff)
	PlasmaBlasterWeaponDamage        = MakeScaledFuncLevelFunc(1.2*PlasmaBlasterRateOfFire/BlasterRateOfFire, BasicBlasterWeaponDamage)
	PlasmaBlasterWeaponEnergyPerShot = MakeScaledFuncLevelFunc(PlasmaEnergyRatio, PlasmaBlasterWeaponDamage)
	PlasmaBlasterWeaponRange         = BasicBlasterWeaponRange // it's already compounding
	PlasmaBlasterWeaponROF           = MakeFixedLevelFunc(PlasmaBlasterRateOfFire)

	SmallPlasmaBlasterWeaponComponentStats = ExtendValuesTable(
		BasicBlasterWeaponComponentStats,
		ComponentStats{
			"WeaponDamageFalloffRatio": PlasmaBlasterFalloffRatio,
			"WeaponRange":              PlasmaBlasterWeaponRange,
			"WeaponFireRate":           PlasmaBlasterWeaponROF,
			"WeaponEnergyPerShot":      PlasmaBlasterWeaponEnergyPerShot,
			"WeaponRawDamage":          PlasmaBlasterWeaponDamage,
			"WeaponShieldBypass":       MakeFixedLevelFunc(PlasmaShieldBypass),
			"WeaponArmorBypass":        MakeFixedLevelFunc(PlasmaArmorBypass),
		},
	)

	MediumPlasmaBlasterWeaponComponentStats = ExtendValuesTable(
		SmallPlasmaBlasterWeaponComponentStats,
		ComponentStats{
			"CrewRequirement":      MediumCrewRequirements,
			"WeaponVolleyAmount":   MakeFixedLevelFunc(2),
			"WeaponVolleyFireRate": MakeFixedLevelFunc(BlasterVolleyFireRate),
		},
	)

	LargePlasmaBlasterWeaponComponentStats = ExtendValuesTable(
		SmallPlasmaBlasterWeaponComponentStats,
		ComponentStats{
			"CrewRequirement":      LargeCrewRequirements,
			"WeaponVolleyAmount":   MakeFixedLevelFunc(4),
			"WeaponVolleyFireRate": MakeFixedLevelFunc(BlasterVolleyFireRate),
		},
	)

	// we'll make these derivative of impact assault blasters
	PlanetaryBlasterWeaponDamage   = MakeScaledFuncLevelFunc(8, ImpactBlasterWeaponDamage) // 2x large (where large is 4x small)
	PlanetaryBlasterComponentStats = ExtendValuesTable(
		SmallImpactBlasterWeaponComponentStats,
		ComponentStats{
			"CrewRequirement":               PlanetaryCrewRequirements, // meaningless, but doesn't hurt
			"ComponentCountermeasuresBonus": MakeScaledFuncLevelFunc(2, DirectFireComponentCountermeasuresBonus),
			"ComponentTargetingBonus":       MakeScaledFuncLevelFunc(2, DirectFireComponentCountermeasuresBonus),
			"WeaponRawDamage":               PlanetaryBlasterWeaponDamage,
			"WeaponEnergyPerShot":           MakeScaledFuncLevelFunc(BlasterEnergyRatio, PlanetaryBlasterWeaponDamage), // meaningless, but doesn't hurt
			"WeaponDamageFalloffRatio":      MakeFixedLevelFunc(0.15),                                                  // reduce the falloff so that we can shoot further as a planetary class weapon
			"WeaponRange":                   MakeScaledFuncLevelFunc(3, ImpactBlasterWeaponRange),                      // todo: write a limiter that takes into account falloff ratio so we don't hit negative output or even < 50% starting dmg
			"WeaponFireRate":                MakeScaledFuncLevelFunc(.5, ImpactBlasterWeaponROF),                       // two shots per interval
		},
	)
)
