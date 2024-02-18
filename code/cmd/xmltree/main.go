package main

// trivial test harness for our custom xmltree lib
// simply want to read & write without scrambling the file or losing comments

import (
	"flag"
	"fmt"
	"log"
	"os"

	"lucky-wolf/DW2-XL/code/cmd/xmltree/algorithm"
)

func main() {

	log.SetFlags(0)

	var err error
	var function, folder string
	var scale float64

	flag.StringVar(&function, "algorithm", "", "algorithm to apply")
	flag.StringVar(&folder, "folder", "XL", "folder to apply changes to")
	flag.BoolVar(&algorithm.Quiet, "quiet", false, "set if you don't want debug output")
	flag.Float64Var(&scale, "scale", 1.0, "scale factor to apply")
	flag.Parse()

	if !algorithm.Quiet {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		log.Println(fmt.Sprintf("Active Path=%s", cwd))
	}

	var algos = []struct {
		name string
		desc string
	}{
		{"All", "Runs (all) Components, ResearchCosts, and Hulls"},
		{"Components", "Runs all component algorithms"},
		{"  Engines", "Updates engine components to XL data table"},
		{"  HyperDrives", "Updates hyperdrive components to have 7 levels off of a common data table"},
		{"  IonShields", "Updates ion shield components off of a common data table (ship & ftr)"},
		{"  IonWeapons", "Updates ion weapon components off of a common core data table"},
		{"  KineticWeapons", "Updates kinetic weapon components off of a common core data table"},
		{"  FighterArmor", "Fighter armor components are derived from ship armors"},
		{"  FighterEngines", "Fighter engine components are derived from ship engines"},
		{"  FighterReactors", "Fighter reactor components are derived from ship reactors"},
		{"  FighterShields", "Fighter shield components are derived from ship shields"},
		{"  FighterWeaponsAndPD", "[Ftr] and [PD] weapon components are derived from ship weapons"},
		{"ResearchCosts", "All Research costs are set to conform to their columnar position (with a few exceptions)"},
		{"HangarBays", "All ship hangarbay slots will be size-limited by role"},
		{"FighterHulls", "All strikecraft component slots will be adjusted to match desired schedule"},
		{"ScalePlanetFrequencies", "Planet frequencies will be scaled by your input"},
		{"PartialOrdering", "All xml objects will have their ID and name and a few other fields placed first"},
		{"RenumberHullComponentBays", "All component bay indexes will be fixed to a simple incremental index"},
	}

	switch function {
	case "":
		flag.PrintDefaults()
		log.Println("Possible algorithms include:")
		for _, v := range algos {
			log.Println(fmt.Sprintf("  %-30s%s", v.name, v.desc))
		}
	case "All":
		err = algorithm.All(folder)
	case "Components":
		err = algorithm.Components(folder)
	case "FighterArmor":
		err = algorithm.FighterArmor(folder)
	case "FighterEngines":
		err = algorithm.FighterEngines(folder)
	case "FighterHulls":
		err = algorithm.FighterHulls(folder)
	case "FighterReactors":
		err = algorithm.FighterReactors(folder)
	case "FighterShields":
		err = algorithm.FighterShields(folder)
	case "FighterWeaponsAndPD":
		err = algorithm.FighterWeaponsAndPD(folder)
	case "HangarBays":
		err = algorithm.HangarBays(folder)
	case "ResearchCosts":
		err = algorithm.ResearchCosts(folder)
	case "Engines":
		err = algorithm.Engines(folder)
	case "HyperDrives":
		err = algorithm.HyperDrives(folder)
	case "IonShields":
		err = algorithm.IonShields(folder)
	case "IonWeapons":
		err = algorithm.IonWeapons(folder)
	case "KineticWeapons":
		err = algorithm.KineticWeapons(folder)
	case "ScalePlanetFrequencies":
		err = algorithm.ScalePlanetFrequencies(folder, scale)
	case "PartialOrdering":
		err = algorithm.PartialOrdering(folder)
	case "RenumberHullComponentBays":
		err = algorithm.RenumberHullComponentBays(folder)
	default:
		err = fmt.Errorf("unknown algorithm: %s", function)
	}
	if err != nil {
		log.Fatal(err)
	}
}
