package algorithm

import (
	"log"
)

func Engines(folder string) (err error) {

	log.Println("All engines will be scaled to XL values")

	// load all component definition files
	j, err := LoadJobFor(folder, "ComponentDefinitions*.xml")
	if err != nil {
		return
	}

	// update kinetic weapons
	err = j.applyStarshipEngines()
	if err != nil {
		return
	}

	// update derivatives
	err = j.applyFighterEngines()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

func (j *Job) applyStarshipEngines() (err error) {

	// apply stats for each component
	err = j.ApplyComponentAll(StarshipEngineData)

	return
}

var (
	StarshipEngineData = map[string]ComponentData{
		"Engines, Ion": {
			minLevel:    0,
			maxLevel:    1,
			fieldValues: IonStarshipEngineComponentStats,
		},
		"Engines, Pulsed Ion": {
			minLevel:    2,
			maxLevel:    5,
			fieldValues: PulsedIonEngineComponentStats,
		},
		"Engines, Compact Ion": {
			minLevel:    2,
			maxLevel:    5,
			fieldValues: CompactEngineComponentStats,
		},
		"Engines, Acceleros": {
			minLevel:    2,
			maxLevel:    5,
			fieldValues: AccelerosEngineComponentStats,
		},
		"Engines, Turbo Thruster": {
			minLevel:    2,
			maxLevel:    10,
			fieldValues: TurboThrusterEngineComponentStats,
		},
		"Engines, Vortex": {
			minLevel:    6,
			maxLevel:    10,
			fieldValues: VortexEngineComponentStats,
		},

		// TODO: super vector & super engines & mortalen line?
	}

	// this needs to correspond with engine maneuverability
	WeakStarshipEngineCountermeasureBonus   = MakeLinearLevelFunc(0, 0.025)
	MediumStarshipEngineCountermeasureBonus = MakeLinearLevelFunc(0, 0.05)
	StrongStarshipEngineCountermeasureBonus = MakeLinearLevelFunc(0, 0.075)

	// Base thrust is 12K, with 10% compounding increase per level
	// Medium is 25% better, and Strong is 50% better
	WeakStarshipCruiseThrust   = MakeExpLevelFunc(12000, .1)
	MediumStarshipCruiseThrust = MakeScaledFuncLevelFunc(1.25, WeakStarshipCruiseThrust)
	StrongStarshipCruiseThrust = MakeScaledFuncLevelFunc(1.5, WeakStarshipCruiseThrust)

	// 33% boost for max thrust
	WeakStarshipMaxThrust   = MakeScaledFuncLevelFunc(1.33333333, WeakStarshipCruiseThrust)
	MediumStarshipMaxThrust = MakeScaledFuncLevelFunc(1.33333333, MediumStarshipCruiseThrust)
	StrongStarshipMaxThrust = MakeScaledFuncLevelFunc(1.33333333, StrongStarshipCruiseThrust)

	WeakStarshipCruiseThrustEnergy   = MakeScaledFuncLevelFunc(0.0002, WeakStarshipCruiseThrust)
	MediumStarshipCruiseThrustEnergy = MakeScaledFuncLevelFunc(0.0002, MediumStarshipCruiseThrust)
	StrongStarshipCruiseThrustEnergy = MakeScaledFuncLevelFunc(0.0002, StrongStarshipCruiseThrust)

	WeakStarshipMaxThrustEnergy   = MakeScaledFuncLevelFunc(0.0002, WeakStarshipMaxThrust)
	MediumStarshipMaxThrustEnergy = MakeScaledFuncLevelFunc(0.0002, MediumStarshipMaxThrust)
	StrongStarshipMaxThrustEnergy = MakeScaledFuncLevelFunc(0.0002, StrongStarshipMaxThrust)

	// Base vector is 60, with 18% compounding increase per level
	// Medium is 20% better, and Strong is 40% better
	WeakStarshipVectorThrust   = MakeExpLevelFunc(60, .18)
	MediumStarshipVectorThrust = MakeScaledFuncLevelFunc(1.25, WeakStarshipVectorThrust)
	StrongStarshipVectorThrust = MakeScaledFuncLevelFunc(1.5, WeakStarshipVectorThrust)

	// vector energy is 1/50
	WeakStarshipVectorEnergy   = MakeScaledFuncLevelFunc(0.02, WeakStarshipVectorThrust)
	MediumStarshipVectorEnergy = MakeScaledFuncLevelFunc(0.02, MediumStarshipVectorThrust)
	StrongStarshipVectorEnergy = MakeScaledFuncLevelFunc(0.02, StrongStarshipVectorThrust)

	StarshipEngineBaseStats = ComponentStats{
		"ComponentIonDefense": DefaultComponentIonDefense,
		"CrewRequirement":     MakeFixedLevelFunc(10),
		"StaticEnergyUsed":    MakeFixedLevelFunc(1),
	}

	WeakStarshipEngineThrust = ComponentStats{
		"EngineMainCruiseThrust":             WeakStarshipCruiseThrust,
		"EngineMainCruiseThrustEnergyUsage":  WeakStarshipCruiseThrustEnergy,
		"EngineMainMaximumThrust":            WeakStarshipMaxThrust,
		"EngineMainMaximumThrustEnergyUsage": WeakStarshipMaxThrustEnergy,
	}

	MediumStarshipEngineThrust = ComponentStats{
		"EngineMainCruiseThrust":             MediumStarshipCruiseThrust,
		"EngineMainCruiseThrustEnergyUsage":  MediumStarshipCruiseThrustEnergy,
		"EngineMainMaximumThrust":            MediumStarshipMaxThrust,
		"EngineMainMaximumThrustEnergyUsage": MediumStarshipMaxThrustEnergy,
	}

	StrongStarshipEngineThrust = ComponentStats{
		"EngineMainCruiseThrust":             StrongStarshipCruiseThrust,
		"EngineMainCruiseThrustEnergyUsage":  StrongStarshipCruiseThrustEnergy,
		"EngineMainMaximumThrust":            StrongStarshipMaxThrust,
		"EngineMainMaximumThrustEnergyUsage": StrongStarshipMaxThrustEnergy,
	}

	WeakStarshipEngineVector = ComponentStats{
		"CountermeasuresBonus":       WeakStarshipEngineCountermeasureBonus,
		"EngineVectoringThrust":      WeakStarshipVectorThrust,
		"EngineVectoringEnergyUsage": WeakStarshipVectorEnergy,
	}

	MediumStarshipEngineVector = ComponentStats{
		"CountermeasuresBonus":       MediumStarshipEngineCountermeasureBonus,
		"EngineVectoringThrust":      MediumStarshipVectorThrust,
		"EngineVectoringEnergyUsage": MediumStarshipVectorEnergy,
	}

	StrongStarshipEngineVector = ComponentStats{
		"EngineVectoringThrust":      StrongStarshipVectorThrust,
		"EngineVectoringEnergyUsage": StrongStarshipVectorEnergy,
		"CountermeasuresBonus":       StrongStarshipEngineCountermeasureBonus,
	}

	// basic
	IonStarshipEngineComponentStats = ExtendValuesTable(
		StarshipEngineBaseStats,
		WeakStarshipEngineThrust,
		WeakStarshipEngineVector,
	)

	// nimble
	PulsedIonEngineComponentStats = ExtendValuesTable(
		StarshipEngineBaseStats,
		MediumStarshipEngineThrust,
		StrongStarshipEngineVector,
	)

	// compact
	CompactEngineComponentStats = IonStarshipEngineComponentStats

	// acceleros
	AccelerosEngineComponentStats = ExtendValuesTable(
		StarshipEngineBaseStats,
		StrongStarshipEngineThrust,
		MediumStarshipEngineVector,
	)

	// powerful, nimble, efficient
	TurboThrusterEngineComponentStats = ExtendValuesTable(
		StarshipEngineBaseStats,
		ComponentStats{
			"EngineMainCruiseThrust":             StrongStarshipCruiseThrust,
			"EngineMainCruiseThrustEnergyUsage":  MakeScaledFuncLevelFunc(0.75, StrongStarshipCruiseThrustEnergy),
			"EngineMainMaximumThrust":            StrongStarshipMaxThrust,
			"EngineMainMaximumThrustEnergyUsage": MakeScaledFuncLevelFunc(0.75, StrongStarshipMaxThrustEnergy),
			"EngineVectoringThrust":              StrongStarshipVectorThrust,
			"EngineVectoringEnergyUsage":         MakeScaledFuncLevelFunc(0.75, StrongStarshipVectorEnergy),
			"CountermeasuresBonus":               StrongStarshipEngineCountermeasureBonus,
		},
	)

	// strong at everything (for smallest size)
	VortexEngineComponentStats = ExtendValuesTable(
		StarshipEngineBaseStats,
		StrongStarshipEngineThrust,
		StrongStarshipEngineVector,
	)
)
