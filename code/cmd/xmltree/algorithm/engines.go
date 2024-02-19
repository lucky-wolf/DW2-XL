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
	err = j.applyEngines()
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

func (j *Job) applyEngines() (err error) {

	// apply stats for each component
	err = j.ApplyComponentAll(StarshipEngineData)

	return
}

const (
	EngineIncreaseExp       = 0.15                        // compounding increase (level over level)
	TurboEfficiencyRatio    = 0.66666                     // 33% reduced energy cost
	MaxThrustRatio          = 1.33333                     // combat only
	CruiseThrustEnergyRatio = 0.0001                      // 1/10K
	MaxThrustEnergyRatio    = 2 * CruiseThrustEnergyRatio // 1/5K
	VectorThrustEnergyRatio = 0.02                        // vector energy is 1/50
)

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
		"Engines, Infinite Flux": {
			minLevel:    11,
			maxLevel:    11,
			fieldValues: InfiniteFluxEngineComponentStats,
		},
		"Inertialess Thruster": {
			minLevel:    11,
			maxLevel:    11,
			fieldValues: InertialessEngineComponentStats,
		},
	}

	// Cruise Thrust
	WeakStarshipCruiseThrust   = MakeExpLevelFunc(12000, EngineIncreaseExp)
	MediumStarshipCruiseThrust = MakeScaledFuncLevelFunc(1.25, WeakStarshipCruiseThrust)
	StrongStarshipCruiseThrust = MakeScaledFuncLevelFunc(1.5, WeakStarshipCruiseThrust)

	WeakStarshipCruiseThrustEnergy   = MakeScaledFuncLevelFunc(CruiseThrustEnergyRatio, WeakStarshipCruiseThrust)
	MediumStarshipCruiseThrustEnergy = MakeScaledFuncLevelFunc(CruiseThrustEnergyRatio, MediumStarshipCruiseThrust)
	StrongStarshipCruiseThrustEnergy = MakeScaledFuncLevelFunc(CruiseThrustEnergyRatio, StrongStarshipCruiseThrust)

	// Max Thrust
	WeakStarshipMaxThrust   = MakeScaledFuncLevelFunc(MaxThrustRatio, WeakStarshipCruiseThrust)
	MediumStarshipMaxThrust = MakeScaledFuncLevelFunc(MaxThrustRatio, MediumStarshipCruiseThrust)
	StrongStarshipMaxThrust = MakeScaledFuncLevelFunc(MaxThrustRatio, StrongStarshipCruiseThrust)

	WeakStarshipMaxThrustEnergy   = MakeScaledFuncLevelFunc(MaxThrustEnergyRatio, WeakStarshipMaxThrust)
	MediumStarshipMaxThrustEnergy = MakeScaledFuncLevelFunc(MaxThrustEnergyRatio, MediumStarshipMaxThrust)
	StrongStarshipMaxThrustEnergy = MakeScaledFuncLevelFunc(MaxThrustEnergyRatio, StrongStarshipMaxThrust)

	// Vectoring Thrust
	WeakStarshipVectorThrust   = MakeExpLevelFunc(60, EngineIncreaseExp)
	MediumStarshipVectorThrust = MakeScaledFuncLevelFunc(1.25, WeakStarshipVectorThrust)
	StrongStarshipVectorThrust = MakeScaledFuncLevelFunc(1.5, WeakStarshipVectorThrust)

	WeakStarshipVectorEnergy   = MakeScaledFuncLevelFunc(VectorThrustEnergyRatio, WeakStarshipVectorThrust)
	MediumStarshipVectorEnergy = MakeScaledFuncLevelFunc(VectorThrustEnergyRatio, MediumStarshipVectorThrust)
	StrongStarshipVectorEnergy = MakeScaledFuncLevelFunc(VectorThrustEnergyRatio, StrongStarshipVectorThrust)

	// Vectoring gives countermeasure bonus
	WeakStarshipEngineCountermeasureBonus   = MakeLinearLevelFunc(0, 0.02)
	MediumStarshipEngineCountermeasureBonus = MakeLinearLevelFunc(0, 0.035)
	StrongStarshipEngineCountermeasureBonus = MakeLinearLevelFunc(0, 0.05)

	StarshipEngineBaseStats = ComponentStats{
		"ComponentIonDefense": HardenedComponentIonDefense, // engines are a hardened component
		"CrewRequirement":     MakeFixedLevelFunc(5),
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
			"EngineMainCruiseThrustEnergyUsage":  MakeScaledFuncLevelFunc(TurboEfficiencyRatio, StrongStarshipCruiseThrustEnergy),
			"EngineMainMaximumThrust":            StrongStarshipMaxThrust,
			"EngineMainMaximumThrustEnergyUsage": MakeScaledFuncLevelFunc(TurboEfficiencyRatio, StrongStarshipMaxThrustEnergy),
			"EngineVectoringThrust":              StrongStarshipVectorThrust,
			"EngineVectoringEnergyUsage":         MakeScaledFuncLevelFunc(TurboEfficiencyRatio, StrongStarshipVectorEnergy),
			"CountermeasuresBonus":               StrongStarshipEngineCountermeasureBonus,
		},
	)

	// strong at everything (for smallest size)
	VortexEngineComponentStats = ExtendValuesTable(
		StarshipEngineBaseStats,
		StrongStarshipEngineThrust,
		StrongStarshipEngineVector,
	)

	// super engine
	InfiniteFluxEngineComponentStats = VortexEngineComponentStats

	// super vector thurster
	InertialessEngineComponentStats = VortexEngineComponentStats
)
