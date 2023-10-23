package main

// trivial test harness for our custom xmltree lib
// simply want to read & write without scrambling the file or losing comments

import (
	"flag"
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
	"os"

	"lucky-wolf/DW2-XL/code/cmd/xmltree/algorithm"
)

type Transformation = algorithm.Transformation

var (
	source   string
	target   string
	function string
)

func main() {

	log.SetFlags(0)

	flag.StringVar(&source, "source", "", "source filename")
	flag.StringVar(&target, "target", "", "target filename")
	flag.StringVar(&function, "algorithm", "", "algorithm to apply")
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
		err = run(algorithm.HangarBays)
	case "FighterShields":
		err = run(algorithm.FighterShields)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func run(transformer Transformation) (err error) {

	// any xml file should be readable
	values, err := xmltree.LoadFromFile(source)
	if err != nil {
		return
	}

	// apply requested transformation (if any)
	if transformer != nil {
		var statistics algorithm.Statistics

		statistics, err = transformer(values)
		if err != nil {
			return
		}

		log.Println(statistics.For(target))
	}

	// convert it to output
	err = values.WriteToFile(target)

	return
}
