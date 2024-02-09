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
		log.Println(cwd)
	}

	var err error
	switch function {
	case "":
		log.Println("no algorithm selected: will simply copy source to target")
	case "All":
		err = algorithm.All(folder)
	case "Components":
		err = algorithm.Components(folder)
	case "FighterArmor":
		err = algorithm.FighterArmor(folder)
	case "FighterEngines":
		err = algorithm.FighterEngines(folder)
	case "FighterHulls":
		err = algorithm.Hulls(folder)
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
	case "HyperDrives":
		err = algorithm.HyperDrives(folder)
	case "IonShields":
		err = algorithm.IonShields(folder)
	case "IonWeapons":
		err = algorithm.IonWeapons(folder)
	case "ScalePlanetFrequencies":
		err = algorithm.ScalePlanetFrequencies(folder, scale)
	case "NamesFirst":
		err = algorithm.NamesFirst(folder)
	default:
		err = fmt.Errorf("unknown algorithm: %s", function)
	}
	if err != nil {
		log.Fatal(err)
	}
}
