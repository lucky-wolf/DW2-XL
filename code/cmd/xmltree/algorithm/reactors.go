package algorithm

import (
	"log"
)

func Reactors(folder string) (err error) {

	log.Println("All reactor will be scaled to XL values")

	// load all component definition files
	j, err := LoadJobFor(folder, "ComponentDefinitions*.xml")
	if err != nil {
		return
	}

	// update kinetic weapons
	err = j.applyReactors()
	if err != nil {
		return
	}

	// update derivatives
	err = j.applyFighterReactors()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

func (j *Job) applyReactors() (err error) {

	// apply stats for each component
	err = j.ApplyComponentAll(StarshipReactorData)

	return
}

func MakeFuelUnitsLevelFunc(efficiency, capacity LevelFunc) LevelFunc {
	return func(level int) float64 { return efficiency(level) / 1000 * capacity(level) }
}

const (
	LowReactorCapacityRatio  = 1.3333333333333
	MedReactorCapacityRatio  = 1.6666666666666
	HighReactorCapacityRatio = 2.0
)

var (
	StarshipReactorData = map[string]ComponentData{
		"Reactor, Caslon": {
			minLevel:       0,
			maxLevel:       1,
			componentStats: CaslonReactorComponentStats,
		},
		"Reactor, Caslon Fusion Prototype": {
			minLevel:       1,
			maxLevel:       1,
			componentStats: FusionReactorComponentStats,
		},
		"Reactor, Novacore": {
			minLevel:       2,
			maxLevel:       10,
			componentStats: NovaReactorComponentStats,
		},
		"Reactor, Caslon Fusion": {
			minLevel:       2,
			maxLevel:       4,
			componentStats: FusionReactorComponentStats,
		},
		"Reactor, Harmonic Caslon": {
			minLevel:       2,
			maxLevel:       4,
			componentStats: HarmonicReactorComponentStats,
		},
		"Reactor, Plasmatic Caslon": {
			minLevel:       2,
			maxLevel:       4,
			componentStats: PlasmaticReactorComponentStats,
		},
		"Reactor, Caslon Hyperfusion": {
			minLevel:       5,
			maxLevel:       7,
			componentStats: HyperFusionReactorComponentStats,
		},
		"Reactor, Resonant Caslon": {
			minLevel:       5,
			maxLevel:       7,
			componentStats: HyperHarmonicReactorComponentStats,
		},
		"Reactor, Plasmatic Caslon Cycling": {
			minLevel:       5,
			maxLevel:       7,
			componentStats: HyperPlasmaticReactorComponentStats,
		},
		"Reactor, Caslon Ultrafusion": {
			minLevel:       8,
			maxLevel:       10,
			componentStats: UltraFusionReactorComponentStats,
		},
		"Reactor, Zero Point": {
			minLevel:       8,
			maxLevel:       10,
			componentStats: ZeroPointReactorComponentStats,
		},
		"Reactor, Dark Star": {
			minLevel:       11,
			maxLevel:       11,
			componentStats: ZeroPointReactorComponentStats,
		},
	}

	// each reactor type is about 10% more efficient than the nearest neighbor
	ExcellentEfficiency = func(level int) float64 { return 2.3 - (.1 * float64(level)) }
	GoodEfficiency      = func(level int) float64 { return 2.6 - (.1 * float64(level)) }
	PoorEfficiency      = func(level int) float64 { return 2.9 - (.1 * float64(level)) }

	CaslonReactorOutput    = MakeLinearLevelFunc(60, 18)
	FusionReactorOutput    = MakeLinearLevelFunc(80, 20)
	HarmonicReactorOutput  = MakeLinearLevelFunc(100, 22)
	PlasmaticReactorOutput = MakeLinearLevelFunc(120, 24)
	NovaReactorOutput      = MakeLinearLevelFunc(100, 33)
	ZeroPointReactorOutput = MakeLinearLevelFunc(0, 48)

	HyperFusionReactorOutput    = MakeOffsetFuncLevelFunc(1, FusionReactorOutput)
	HyperHarmonicReactorOutput  = MakeOffsetFuncLevelFunc(1, HarmonicReactorOutput)
	HyperPlasmaticReactorOutput = MakeOffsetFuncLevelFunc(1, PlasmaticReactorOutput)

	UltraFusionReactorOutput = MakeOffsetFuncLevelFunc(2, FusionReactorOutput)

	// basic
	CaslonReactorCapacity       = MakeScaledFuncLevelFunc(LowReactorCapacityRatio, CaslonReactorOutput)
	CaslonEfficiency            = GoodEfficiency
	CaslonFuelUnits             = MakeFuelUnitsLevelFunc(CaslonEfficiency, CaslonReactorCapacity)
	CaslonReactorComponentStats = ComponentStats{
		"ComponentIonDefense":           HardenedComponentIonDefense, // reactors are a hardened component
		"CrewRequirement":               MakeFixedLevelFunc(5),
		"ReactorEnergyOutputPerSecond":  CaslonReactorOutput,
		"ReactorEnergyStorageCapacity":  CaslonReactorCapacity,
		"ReactorFuelUnitsForFullCharge": CaslonFuelUnits,
		"StaticEnergyUsed":              MakeFixedLevelFunc(0),
	}

	// novacore (Quameno)
	NovaReactorCapacity       = MakeScaledFuncLevelFunc(MedReactorCapacityRatio, NovaReactorOutput)
	NovaEfficiency            = ExcellentEfficiency
	NovaFuelUnits             = MakeFuelUnitsLevelFunc(NovaEfficiency, NovaReactorCapacity)
	NovaReactorComponentStats = ComposeComponentStats(
		CaslonReactorComponentStats,
		ComponentStats{
			"ReactorEnergyOutputPerSecond":  NovaReactorOutput,
			"ReactorEnergyStorageCapacity":  NovaReactorCapacity,
			"ReactorFuelUnitsForFullCharge": NovaFuelUnits,
		},
	)

	// fusion
	FusionReactorCapacity       = MakeScaledFuncLevelFunc(HighReactorCapacityRatio, FusionReactorOutput)
	FusionEfficiency            = ExcellentEfficiency
	FusionFuelUnits             = MakeFuelUnitsLevelFunc(FusionEfficiency, FusionReactorCapacity)
	FusionReactorComponentStats = ComposeComponentStats(
		CaslonReactorComponentStats,
		ComponentStats{
			"ReactorEnergyOutputPerSecond":  FusionReactorOutput,
			"ReactorEnergyStorageCapacity":  FusionReactorCapacity,
			"ReactorFuelUnitsForFullCharge": FusionFuelUnits,
		},
	)

	// harmonic
	HarmonicReactorCapacity       = MakeScaledFuncLevelFunc(MedReactorCapacityRatio, HarmonicReactorOutput)
	HarmonicEfficiency            = GoodEfficiency
	HarmonicFuelUnits             = MakeFuelUnitsLevelFunc(HarmonicEfficiency, HarmonicReactorCapacity)
	HarmonicReactorComponentStats = ComposeComponentStats(
		CaslonReactorComponentStats,
		ComponentStats{
			"ReactorEnergyOutputPerSecond":  HarmonicReactorOutput,
			"ReactorEnergyStorageCapacity":  HarmonicReactorCapacity,
			"ReactorFuelUnitsForFullCharge": HarmonicFuelUnits,
		},
	)

	// plasmatic
	PlasmaticReactorCapacity       = MakeScaledFuncLevelFunc(LowReactorCapacityRatio, PlasmaticReactorOutput)
	PlasmaticEfficiency            = PoorEfficiency
	PlasmaticFuelUnits             = MakeFuelUnitsLevelFunc(PlasmaticEfficiency, PlasmaticReactorCapacity)
	PlasmaticReactorComponentStats = ComposeComponentStats(
		CaslonReactorComponentStats,
		ComponentStats{
			"ReactorEnergyOutputPerSecond":  PlasmaticReactorOutput,
			"ReactorEnergyStorageCapacity":  PlasmaticReactorCapacity,
			"ReactorFuelUnitsForFullCharge": PlasmaticFuelUnits,
		},
	)

	// hyper fusion
	HyperFusionReactorCapacity       = MakeScaledFuncLevelFunc(HighReactorCapacityRatio, HyperFusionReactorOutput)
	HyperFusionEfficiency            = ExcellentEfficiency
	HyperFusionFuelUnits             = MakeFuelUnitsLevelFunc(HyperFusionEfficiency, HyperFusionReactorCapacity)
	HyperFusionReactorComponentStats = ComposeComponentStats(
		CaslonReactorComponentStats,
		ComponentStats{
			"ReactorEnergyOutputPerSecond":  HyperFusionReactorOutput,
			"ReactorEnergyStorageCapacity":  HyperFusionReactorCapacity,
			"ReactorFuelUnitsForFullCharge": HyperFusionFuelUnits,
		},
	)

	// hyper harmonic
	HyperHarmonicReactorCapacity       = MakeScaledFuncLevelFunc(MedReactorCapacityRatio, HyperHarmonicReactorOutput)
	HyperHarmonicEfficiency            = GoodEfficiency
	HyperHarmonicFuelUnits             = MakeFuelUnitsLevelFunc(HyperHarmonicEfficiency, HyperHarmonicReactorCapacity)
	HyperHarmonicReactorComponentStats = ComposeComponentStats(
		CaslonReactorComponentStats,
		ComponentStats{
			"ReactorEnergyOutputPerSecond":  HyperHarmonicReactorOutput,
			"ReactorEnergyStorageCapacity":  HyperHarmonicReactorCapacity,
			"ReactorFuelUnitsForFullCharge": HyperHarmonicFuelUnits,
		},
	)

	// hyper plasmatic
	HyperPlasmaticReactorCapacity       = MakeScaledFuncLevelFunc(LowReactorCapacityRatio, HyperPlasmaticReactorOutput)
	HyperPlasmaticEfficiency            = PoorEfficiency
	HyperPlasmaticFuelUnits             = MakeFuelUnitsLevelFunc(HyperPlasmaticEfficiency, HyperPlasmaticReactorCapacity)
	HyperPlasmaticReactorComponentStats = ComposeComponentStats(
		CaslonReactorComponentStats,
		ComponentStats{
			"ReactorEnergyOutputPerSecond":  HyperPlasmaticReactorOutput,
			"ReactorEnergyStorageCapacity":  HyperPlasmaticReactorCapacity,
			"ReactorFuelUnitsForFullCharge": HyperPlasmaticFuelUnits,
		},
	)

	// ultra fusion
	UltraFusionReactorCapacity       = MakeScaledFuncLevelFunc(HighReactorCapacityRatio, UltraFusionReactorOutput)
	UltraFusionEfficiency            = ExcellentEfficiency
	UltraFusionFuelUnits             = MakeFuelUnitsLevelFunc(UltraFusionEfficiency, UltraFusionReactorCapacity)
	UltraFusionReactorComponentStats = ComposeComponentStats(
		CaslonReactorComponentStats,
		ComponentStats{
			"ReactorEnergyOutputPerSecond":  UltraFusionReactorOutput,
			"ReactorEnergyStorageCapacity":  UltraFusionReactorCapacity,
			"ReactorFuelUnitsForFullCharge": UltraFusionFuelUnits,
		},
	)

	// zero point
	ZeroPointReactorCapacity       = MakeScaledFuncLevelFunc(MedReactorCapacityRatio, ZeroPointReactorOutput)
	ZeroPointEfficiency            = GoodEfficiency
	ZeroPointFuelUnits             = MakeFuelUnitsLevelFunc(ZeroPointEfficiency, ZeroPointReactorCapacity)
	ZeroPointReactorComponentStats = ComposeComponentStats(
		CaslonReactorComponentStats,
		ComponentStats{
			"ReactorEnergyOutputPerSecond":  ZeroPointReactorOutput,
			"ReactorEnergyStorageCapacity":  ZeroPointReactorCapacity,
			"ReactorFuelUnitsForFullCharge": ZeroPointFuelUnits,
		},
	)
)
