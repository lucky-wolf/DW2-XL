package algorithm

import (
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
)

func Armors(folder string) (err error) {

	log.Println("All armors will be scaled to XL values")

	// load all component definition files
	j, err := LoadJobFor(folder, "ComponentDefinitions*.xml")
	if err != nil {
		return
	}

	// update kinetic weapons
	err = j.applyArmors()
	if err != nil {
		return
	}

	// update derivatives
	err = j.applyFighterArmors()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

func (j *Job) applyArmors() (err error) {

	// apply stats for each component
	err = j.ApplyComponentAll(StarshipArmorData)

	return
}

// general trend is currently that weapons increase in output faster than armors can resist it
// subtle: but note that starting values matter immensely - if weapons outpace armor to begin with - they'll only spread that delta w/o any extra exp delta
const (
	ArmorSizeSmall    = 10
	ArmorSizeStandard = 14
	ArmorSizeLarge    = 18
)

var (
	WeakArmorBlastRating     = MakeExpLevelFunc(ArmorStrengthBasis, ArmorStrengthIncreaseExp)
	StandardArmorBlastRating = MakeScaledFuncLevelFunc(1.5, WeakArmorBlastRating)
	StrongArmorBlastRating   = MakeScaledFuncLevelFunc(2.0, WeakArmorBlastRating)

	WeakArmorReactiveRating     = MakeLinearLevelFunc(2, 1)
	StandardArmorReactiveRating = MakeLinearLevelFunc(2, 1.5)
	StrongArmorReactiveRating   = MakeLinearLevelFunc(2, 2)

	WeakArmorIonDefense   = MakeFixedLevelFunc(0)
	MediumArmorIonDefense = MakeLinearLevelFunc(0, 2)
	StrongArmorIonDefense = MakeLinearLevelFunc(0, 4)

	InertArmorStaticEnergy  = MakeFixedLevelFunc(0)
	ActiveArmorStaticEnergy = MakeLinearLevelFunc(.5, .5)

	InertArmorCrew  = MakeFixedLevelFunc(1)
	ActiveArmorCrew = MakeIntegerLevelFunc(MakeLinearLevelFunc(1, .25))

	SmallArmorValues = map[AttributeName]xmltree.SimpleValue{
		"Size": xmltree.CreateInt(ArmorSizeSmall),
	}
	StandardArmorValues = map[AttributeName]xmltree.SimpleValue{
		"Size": xmltree.CreateInt(ArmorSizeStandard),
	}
	LargeArmorValues = map[AttributeName]xmltree.SimpleValue{
		"Size": xmltree.CreateInt(ArmorSizeLarge),
	}

	StarshipArmorData = map[string]ComponentData{

		"Armor Plating": {
			values:         StandardArmorValues,
			minLevel:       0,
			maxLevel:       1,
			componentStats: ArmorPlatingComponentStats,
		},

		"Ion Sheath Armor": {
			values:         StandardArmorValues,
			minLevel:       2,
			maxLevel:       3,
			componentStats: IonArmorComponentStats,
		},
		"Enhanced Ion Sheath Armor": {
			values:         StandardArmorValues,
			minLevel:       4,
			maxLevel:       5,
			componentStats: IonArmorComponentStats,
		},
		"Hardened Ion Sheath Armor": {
			values:         StandardArmorValues,
			minLevel:       6,
			maxLevel:       7,
			componentStats: IonArmorComponentStats,
		},
		"Ultra-Dense Ion Sheath Armor": {
			values:         StandardArmorValues,
			minLevel:       8,
			maxLevel:       9,
			componentStats: IonArmorComponentStats,
		},
		"Absorbing Ion Sheath Armor": {
			values:         StandardArmorValues,
			minLevel:       10,
			maxLevel:       10,
			componentStats: IonArmorComponentStats,
		},

		"Heavy Armor": {
			values:         StandardArmorValues,
			minLevel:       2,
			maxLevel:       3,
			componentStats: HeavyArmorComponentStats,
		},
		"Enhanced Heavy Armor": {
			values:         StandardArmorValues,
			minLevel:       4,
			maxLevel:       5,
			componentStats: HeavyArmorComponentStats,
		},
		"Hardened Heavy Armor": {
			values:         StandardArmorValues,
			minLevel:       6,
			maxLevel:       7,
			componentStats: HeavyArmorComponentStats,
		},
		"Ultra-Dense Heavy Armor": {
			values:         StandardArmorValues,
			minLevel:       8,
			maxLevel:       9,
			componentStats: HeavyArmorComponentStats,
		},
		"Absorbing Heavy Armor": {
			values:         StandardArmorValues,
			minLevel:       10,
			maxLevel:       10,
			componentStats: HeavyArmorComponentStats,
		},

		"Reactive Armor": {
			values:         StandardArmorValues,
			minLevel:       2,
			maxLevel:       3,
			componentStats: ReactiveArmorComponentStats,
		},
		"Enhanced Reactive Armor": {
			values:         StandardArmorValues,
			minLevel:       4,
			maxLevel:       5,
			componentStats: ReactiveArmorComponentStats,
		},
		"Hardened Reactive Armor": {
			values:         StandardArmorValues,
			minLevel:       6,
			maxLevel:       7,
			componentStats: ReactiveArmorComponentStats,
		},
		"Ultra-Dense Reactive Armor": {
			values:         StandardArmorValues,
			minLevel:       8,
			maxLevel:       9,
			componentStats: ReactiveArmorComponentStats,
		},
		"Absorbing Reactive Armor": {
			values:         StandardArmorValues,
			minLevel:       10,
			maxLevel:       10,
			componentStats: ReactiveArmorComponentStats,
		},

		"Flux Sheath Armor": {
			values:         StandardArmorValues,
			minLevel:       2,
			maxLevel:       3,
			componentStats: FluxArmorComponentStats,
		},
		"Flux Enhanced Armor": {
			values:         StandardArmorValues,
			minLevel:       4,
			maxLevel:       5,
			componentStats: FluxArmorComponentStats,
		},
		"Flux Hardened Armor": {
			values:         StandardArmorValues,
			minLevel:       6,
			maxLevel:       7,
			componentStats: FluxArmorComponentStats,
		},
		"Flux Ultra-Dense Armor": {
			values:         StandardArmorValues,
			minLevel:       8,
			maxLevel:       9,
			componentStats: FluxArmorComponentStats,
		},
		"Flux Absorbing Armor": {
			values:         StandardArmorValues,
			minLevel:       10,
			maxLevel:       10,
			componentStats: FluxArmorComponentStats,
		},

		"Hex Armor": {
			values:         StandardArmorValues,
			minLevel:       2,
			maxLevel:       3,
			componentStats: HexArmorComponentStats,
		},
		"Reactive Hex Armor": {
			values:         StandardArmorValues,
			minLevel:       4,
			maxLevel:       5,
			componentStats: HexArmorComponentStats,
		},
		"Dense Hex Armor": {
			values:         StandardArmorValues,
			minLevel:       6,
			maxLevel:       7,
			componentStats: HexArmorComponentStats,
		},
		"Multi-Dimensional Hex Armor": {
			values:         StandardArmorValues,
			minLevel:       8,
			maxLevel:       9,
			componentStats: HexArmorComponentStats,
		},
		"Infinite Hex Armor": {
			values:         StandardArmorValues,
			minLevel:       10,
			maxLevel:       10,
			componentStats: HexArmorComponentStats,
		},

		"Stellar Armor": {
			values:         LargeArmorValues,
			minLevel:       10,
			maxLevel:       10,
			componentStats: StellarArmorComponentStats,
		},
	}

	ArmorPlatingComponentStats = ComponentStats{
		"ComponentIonDefense": HardenedComponentIonDefense, // armors are hardened
		"CrewRequirement":     InertArmorCrew,
		"StaticEnergyUsed":    InertArmorStaticEnergy,
		"ArmorBlastRating":    StandardArmorBlastRating,
		"ArmorReactiveRating": WeakArmorReactiveRating,
		"IonDamageDefense":    WeakArmorIonDefense,
	}

	ReactiveArmorComponentStats = ComponentStats{
		"ComponentIonDefense": HardenedComponentIonDefense, // armors are hardened
		"CrewRequirement":     ActiveArmorCrew,
		"StaticEnergyUsed":    ActiveArmorStaticEnergy,
		"ArmorBlastRating":    WeakArmorBlastRating,
		"ArmorReactiveRating": StrongArmorReactiveRating,
		"IonDamageDefense":    WeakArmorIonDefense,
	}

	HeavyArmorComponentStats = ComponentStats{
		"ComponentIonDefense": HardenedComponentIonDefense, // armors are hardened
		"CrewRequirement":     InertArmorCrew,
		"StaticEnergyUsed":    InertArmorStaticEnergy,
		"ArmorBlastRating":    StrongArmorBlastRating,
		"ArmorReactiveRating": WeakArmorReactiveRating,
		"IonDamageDefense":    WeakArmorIonDefense,
	}

	IonArmorComponentStats = ComponentStats{
		"ComponentIonDefense": HardenedComponentIonDefense, // armors are hardened
		"CrewRequirement":     ActiveArmorCrew,
		"StaticEnergyUsed":    ActiveArmorStaticEnergy,
		"ArmorBlastRating":    StandardArmorBlastRating,
		"ArmorReactiveRating": WeakArmorReactiveRating,
		"IonDamageDefense":    StrongArmorIonDefense,
	}

	FluxArmorComponentStats = ComponentStats{
		"ComponentIonDefense": HardenedComponentIonDefense, // armors are hardened
		"CrewRequirement":     ActiveArmorCrew,
		"StaticEnergyUsed":    ActiveArmorStaticEnergy,
		"ArmorBlastRating":    StandardArmorBlastRating,
		"ArmorReactiveRating": StandardArmorReactiveRating,
		"IonDamageDefense":    MediumArmorIonDefense,
	}

	HexArmorComponentStats = ComponentStats{
		"ComponentIonDefense": HardenedComponentIonDefense, // armors are hardened
		"CrewRequirement":     ActiveArmorCrew,
		"StaticEnergyUsed":    ActiveArmorStaticEnergy,
		"ArmorBlastRating":    StandardArmorBlastRating,
		"ArmorReactiveRating": StrongArmorReactiveRating,
		"IonDamageDefense":    WeakArmorIonDefense,
	}

	StellarArmorComponentStats = ComponentStats{
		"ComponentIonDefense": HardenedComponentIonDefense, // armors are hardened
		"CrewRequirement":     ActiveArmorCrew,
		"StaticEnergyUsed":    ActiveArmorStaticEnergy,
		"ArmorBlastRating":    StrongArmorBlastRating,
		"ArmorReactiveRating": StrongArmorReactiveRating,
		"IonDamageDefense":    StrongArmorIonDefense,
	}
)
