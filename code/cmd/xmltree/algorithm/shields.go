package algorithm

import (
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
)

func Shields(folder string) (err error) {

	log.Println("All shields will be scaled to XL values")

	// load all component definition files
	j, err := LoadJobFor(folder, "ComponentDefinitions*.xml")
	if err != nil {
		return
	}

	// update kinetic weapons
	err = j.applyShields()
	if err != nil {
		return
	}

	// update derivatives
	err = j.applyFighterShields()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

func (j *Job) applyShields() (err error) {

	// apply stats for each component
	err = j.ApplyComponentAll(StarshipShieldData)

	return
}

func ShieldStaticEnergy(shieldStrength LevelFunc) LevelFunc {
	return func(level int) float64 { return ShieldStaticEnergyCoefficient * shieldStrength(level) }
}

func ShieldRechargeEnergy(shieldRecharge LevelFunc) LevelFunc {
	return func(level int) float64 { return ShieldRechargeEnergyCoefficient * shieldRecharge(level) }
}

func QuickShieldRecharge(shieldStrength LevelFunc) LevelFunc {
	return func(level int) float64 { return QuickShieldRechargeCoefficient * shieldStrength(level) }
}

func StandardShieldRecharge(shieldStrength LevelFunc) LevelFunc {
	return func(level int) float64 { return StandardShieldRechargeCoefficient * shieldStrength(level) }
}

func SlowShieldRecharge(shieldStrength LevelFunc) LevelFunc {
	return func(level int) float64 { return SlowShieldRechargeCoefficient * shieldStrength(level) }
}

// 30%, 25%, 20%, 15%, 10%, 5%, 0...
func ShieldPenetrationChance(level int) float64 {
	return max(0, .3-float64(level)*.05)
}

// 60%, 50%, 40%, 30%, 20%, 10%, 0...
func ShieldPenetrationRatio(level int) float64 {
	return max(0, .6-float64(level)*.1)
}

func GenerateShieldStatsFor(shieldStrength LevelFunc, rechargeRate func(LevelFunc) LevelFunc) ComponentStats {
	return ComponentStats{
		"ShieldStrength":            shieldStrength,
		"ShieldRechargeRate":        rechargeRate(shieldStrength),
		"ShieldRechargeEnergyUsage": ShieldRechargeEnergy(rechargeRate(shieldStrength)),
		"StaticEnergyUsed":          ShieldStaticEnergy(shieldStrength),
	}
}

const (
	ShieldSizeSmall    = 8
	ShieldSizeStandard = 12
	ShieldSizeLarge    = 16

	ShieldStaticEnergyCoefficient   = 0.0125
	ShieldRechargeEnergyCoefficient = 0.0125

	QuickShieldRechargeCoefficient    = 0.0100
	StandardShieldRechargeCoefficient = 0.0050
	SlowShieldRechargeCoefficient     = 0.0025
)

var (
	StarshipShieldData = map[string]ComponentData{

		// basic
		"Deflectors": {
			values:         StandardShieldValues,
			minLevel:       0,
			maxLevel:       1,
			componentStats: BasicShieldComponentStats,
		},

		// mid-game
		"Corvidian Shields": {
			values:         StandardShieldValues,
			minLevel:       2,
			maxLevel:       5,
			componentStats: BalancedShieldComponentStats,
		},
		"Talassos Shields": {
			values:         StandardShieldValues,
			minLevel:       2,
			maxLevel:       5,
			componentStats: QuickShieldComponentStats,
		},
		"Deucalios Shields": {
			values:         StandardShieldValues,
			minLevel:       2,
			maxLevel:       5,
			componentStats: StrongShieldComponentStats,
		},

		// advanced
		"Meridian Shields": {
			values:         StandardShieldValues,
			minLevel:       6,
			maxLevel:       10,
			componentStats: AdvancedShieldComponentStats,
		},

		// super
		"Citadel Shields": {
			values:         LargeShieldValues,
			minLevel:       11,
			maxLevel:       11,
			componentStats: AdvancedShieldComponentStats,
		},

		// zenox
		"Megatron Z4 Shields": {
			values:         SmallShieldValues,
			minLevel:       2,
			maxLevel:       10,
			componentStats: ZenoxShieldComponentStats,
		},

		// quameno
		"Tortoise Shields": {
			values:         StandardShieldValues,
			minLevel:       2,
			maxLevel:       10,
			componentStats: QuamenoShieldComponentStats,
		},
	}

	SmallShieldValues = map[AttributeName]xmltree.SimpleValue{
		"Size": xmltree.CreateInt(ShieldSizeSmall),
	}
	StandardShieldValues = map[AttributeName]xmltree.SimpleValue{
		"Size": xmltree.CreateInt(ShieldSizeStandard),
	}
	LargeShieldValues = map[AttributeName]xmltree.SimpleValue{
		"Size": xmltree.CreateInt(ShieldSizeLarge),
	}

	WeakShieldStrength     = MakeExpLevelFunc(ShieldStrengthBasis, ShieldStrengthIncreaseExp)
	StandardShieldStrength = MakeScaledFuncLevelFunc(1.3333333, WeakShieldStrength)
	StrongShieldStrength   = MakeScaledFuncLevelFunc(1.6666666, WeakShieldStrength)

	WeakShieldResistance     = MakeLinearLevelFunc(0, 1)
	StandardShieldResistance = MakeLinearLevelFunc(0, 1.5)
	StrongShieldResistance   = MakeLinearLevelFunc(0, 2)

	WeakShieldIonDefense     = MakeLinearLevelFunc(0, 1)
	StandardShieldIonDefense = MakeLinearLevelFunc(0, 2)
	StrongShieldIonDefense   = MakeLinearLevelFunc(0, 3)

	CoreShieldStats = ComponentStats{
		"ComponentIonDefense":     HardenedComponentIonDefense, // shields are hardened
		"CrewRequirement":         SmallCrewRequirements,
		"ShieldPenetrationChance": ShieldPenetrationChance,
		"ShieldPenetrationRatio":  ShieldPenetrationRatio,
	}

	BasicShieldComponentStats = ComposeComponentStats(
		CoreShieldStats,
		GenerateShieldStatsFor(WeakShieldStrength, SlowShieldRecharge),
		ComponentStats{
			"ShieldResistance": WeakShieldResistance,
			"IonDamageDefense": WeakShieldIonDefense,
		},
	)

	BalancedShieldComponentStats = ComposeComponentStats(
		CoreShieldStats,
		GenerateShieldStatsFor(StandardShieldStrength, StandardShieldRecharge),
		ComponentStats{
			"ShieldResistance": StandardShieldResistance,
			"IonDamageDefense": StandardShieldIonDefense,
		},
	)

	QuickShieldComponentStats = ComposeComponentStats(
		CoreShieldStats,
		GenerateShieldStatsFor(WeakShieldStrength, QuickShieldRecharge),
		ComponentStats{
			"ShieldResistance": StandardShieldResistance,
			"IonDamageDefense": StandardShieldIonDefense,
		},
	)

	StrongShieldComponentStats = ComposeComponentStats(
		CoreShieldStats,
		GenerateShieldStatsFor(StrongShieldStrength, SlowShieldRecharge),
		ComponentStats{
			"ShieldResistance": StandardShieldResistance,
			"IonDamageDefense": StandardShieldIonDefense,
		},
	)

	AdvancedShieldComponentStats = ComposeComponentStats(
		CoreShieldStats,
		GenerateShieldStatsFor(StrongShieldStrength, QuickShieldRecharge),
		ComponentStats{
			"ShieldResistance": StrongShieldResistance,
			"IonDamageDefense": StrongShieldIonDefense,
		},
	)

	ZenoxShieldComponentStats = ComposeComponentStats(
		CoreShieldStats,
		GenerateShieldStatsFor(StandardShieldStrength, QuickShieldRecharge),
		ComponentStats{
			"ShieldResistance": StandardShieldResistance,
			"IonDamageDefense": StandardShieldIonDefense,
		},
	)

	QuamenoShieldComponentStats = ComposeComponentStats(
		CoreShieldStats,
		GenerateShieldStatsFor(StrongShieldStrength, QuickShieldRecharge),
		ComponentStats{
			"ShieldResistance": StandardShieldResistance,
			"IonDamageDefense": WeakShieldIonDefense,
		},
	)
)
