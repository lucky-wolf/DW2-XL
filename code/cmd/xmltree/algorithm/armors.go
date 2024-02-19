package algorithm

import (
	"log"
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
	ArmorStrengthIncreaseExp = 0.15 // compounding increase (level over level)
	WeaponDmgIncreaseExp     = 0.18 // % level over level
)

var (
	StarshipArmorData = map[string]ComponentData{
		"Armor Plating": {
			minLevel:    0,
			maxLevel:    1,
			fieldValues: ArmorPlatingComponentStats,
		},
		"Ion Sheath Armor": {
			minLevel:    2,
			maxLevel:    3,
			fieldValues: IonArmorComponentStats,
		},
		"Heavy Armor": {
			minLevel:    2,
			maxLevel:    3,
			fieldValues: HeavyArmorComponentStats,
		},
		"Reactive Armor": {
			minLevel:    2,
			maxLevel:    3,
			fieldValues: ReactiveArmorComponentStats,
		},
		"Enhanced Ion Sheath Armor": {
			minLevel:    4,
			maxLevel:    5,
			fieldValues: IonArmorComponentStats,
		},
		"Enhanced Heavy Armor": {
			minLevel:    4,
			maxLevel:    5,
			fieldValues: HeavyArmorComponentStats,
		},
		"Enhanced Reactive Armor": {
			minLevel:    4,
			maxLevel:    5,
			fieldValues: ReactiveArmorComponentStats,
		},
		"Hardened Ion Sheath Armor": {
			minLevel:    6,
			maxLevel:    7,
			fieldValues: IonArmorComponentStats,
		},
		"Hardened Heavy Armor": {
			minLevel:    6,
			maxLevel:    7,
			fieldValues: HeavyArmorComponentStats,
		},
		"Hardened Reactive Armor": {
			minLevel:    6,
			maxLevel:    7,
			fieldValues: ReactiveArmorComponentStats,
		},
		"Ultra-Dense Ion Sheath Armor": {
			minLevel:    8,
			maxLevel:    9,
			fieldValues: IonArmorComponentStats,
		},
		"Ultra-Dense Heavy Armor": {
			minLevel:    8,
			maxLevel:    9,
			fieldValues: HeavyArmorComponentStats,
		},
		"Ultra-Dense Reactive Armor": {
			minLevel:    8,
			maxLevel:    9,
			fieldValues: ReactiveArmorComponentStats,
		},
		"Absorbing Ion Sheath Armor": {
			minLevel:    10,
			maxLevel:    10,
			fieldValues: IonArmorComponentStats,
		},
		"Absorbing Heavy Armor": {
			minLevel:    10,
			maxLevel:    10,
			fieldValues: HeavyArmorComponentStats,
		},
		"Absorbing Reactive Armor": {
			minLevel:    10,
			maxLevel:    10,
			fieldValues: ReactiveArmorComponentStats,
		},

		"Flux Sheath Armor": {
			minLevel:    2,
			maxLevel:    3,
			fieldValues: FluxArmorComponentStats,
		},
		"Flux Enhanced Armor": {
			minLevel:    4,
			maxLevel:    5,
			fieldValues: FluxArmorComponentStats,
		},
		"Flux Hardened Armor": {
			minLevel:    6,
			maxLevel:    7,
			fieldValues: FluxArmorComponentStats,
		},
		"Flux Ultra-Dense Armor": {
			minLevel:    8,
			maxLevel:    9,
			fieldValues: FluxArmorComponentStats,
		},
		"Flux Absorbing Armor": {
			minLevel:    10,
			maxLevel:    10,
			fieldValues: FluxArmorComponentStats,
		},

		"Hex Armor": {
			minLevel:    2,
			maxLevel:    3,
			fieldValues: HexArmorComponentStats,
		},
		"Reactive Hex Armor": {
			minLevel:    4,
			maxLevel:    5,
			fieldValues: HexArmorComponentStats,
		},
		"Dense Hex Armor": {
			minLevel:    6,
			maxLevel:    7,
			fieldValues: HexArmorComponentStats,
		},
		"Multi-Dimensional Hex Armor": {
			minLevel:    8,
			maxLevel:    9,
			fieldValues: HexArmorComponentStats,
		},
		"Infinite Hex Armor": {
			minLevel:    10,
			maxLevel:    10,
			fieldValues: HexArmorComponentStats,
		},

		"Stellar Armor": {
			minLevel:    10,
			maxLevel:    10,
			fieldValues: StellarArmorComponentStats,
		},
	}

	WeakArmorBlastRating     = MakeExpLevelFunc(15*BlasterBaseDamage, ArmorStrengthIncreaseExp)
	StandardArmorBlastRating = MakeScaledFuncLevelFunc(1.3333333, WeakArmorBlastRating)
	StrongArmorBlastRating   = MakeScaledFuncLevelFunc(1.6666666, WeakArmorBlastRating)

	WeakArmorReactiveRating     = MakeLinearLevelFunc(2, 1)
	StandardArmorReactiveRating = MakeLinearLevelFunc(2, 1.5)
	StrongArmorReactiveRating   = MakeLinearLevelFunc(3, 2)

	WeakArmorIonDefense   = MakeFixedLevelFunc(0)
	MediumArmorIonDefense = MakeLinearLevelFunc(0, 2)
	StrongArmorIonDefense = MakeLinearLevelFunc(0, 4)

	InertArmorStaticEnergy  = MakeFixedLevelFunc(0)
	ActiveArmorStaticEnergy = MakeLinearLevelFunc(.25, .25)

	InertArmorCrew  = MakeFixedLevelFunc(1)
	ActiveArmorCrew = MakeIntegerLevelFunc(MakeLinearLevelFunc(1, .25))

	ArmorPlatingComponentStats = ComponentStats{
		"ComponentIonDefense": HardenedComponentIonDefense, // armors are hardened
		"CrewRequirement":     InertArmorCrew,
		"StaticEnergyUsed":    InertArmorStaticEnergy,
		"ArmorBlastRating":    StandardArmorBlastRating,
		"ArmorReactiveRating": StandardArmorReactiveRating,
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
		"ArmorBlastRating":    WeakArmorBlastRating,
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
