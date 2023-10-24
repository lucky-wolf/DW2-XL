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
	flag.StringVar(&function, "algorithm", "", "algorithm to apply")
	flag.StringVar(&folder, "folder", "XL", "folder to apply changes to")
	flag.BoolVar(&algorithm.Quiet, "quiet", false, "set if you don't want debug output")
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
	case "HangarBays":
		err = algorithm.HangarBays(folder)
	case "FighterArmor":
		err = algorithm.FighterArmor(folder)
	case "FighterShields":
		err = algorithm.FighterShields(folder)
	case "FighterReactor":
		err = algorithm.FighterReactors(folder)
	default:
		err = fmt.Errorf("unknown algorithm: %s", function)
	}
	if err != nil {
		log.Fatal(err)
	}
}
