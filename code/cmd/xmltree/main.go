package main

// trivial test harness for our custom xmltree lib
// simply want to read & write without scrambling the file or losing comments

import (
	"flag"
	"fmt"
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
	"os"
)

type Transformation = func(data *xmltree.XMLTree) (statistics Statistics, err error)

type Statistics struct {
	objects  uint
	elements uint
	changed  uint
}

func (statistics *Statistics) For(filename string) string {
	return fmt.Sprintf("%s: objects found: %d, elements updated: %d of %d", filename, statistics.objects, statistics.changed, statistics.elements)
}

var (
	source    string
	target    string
	algorithm string
	quiet     bool
)

func main() {

	log.SetFlags(0)

	flag.StringVar(&source, "source", "", "source filename")
	flag.StringVar(&target, "target", "", "target filename")
	flag.StringVar(&algorithm, "algorithm", "", "algorithm to apply")
	flag.BoolVar(&quiet, "quiet", false, "set if you don't want debug output")
	flag.Parse()

	if !quiet {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		log.Println(cwd)
	}

	var algo Transformation
	switch algorithm {
	case "":
		log.Println("no algorithm selected: will simply copy source to target")
	case "HangarBays":
		algo = HangarBays
	}

	err := run(algo)
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
		var statistics Statistics

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

func AssertIs(e *xmltree.XMLElement, kind string) (err error) {
	if !quiet {
		log.Println(e.Name.Local)
	}
	if e.Name.Local != kind {
		err = fmt.Errorf("invalid file: expected %s but found %s", kind, e.Name.Local)
	}
	return
}

func HangarBays(root *xmltree.XMLTree) (statistics Statistics, err error) {

	if !quiet {
		log.Println("All stations will have size 100 hangar bays")
		log.Println("All carriers will have size 50 hangar bays")
		log.Println("Everything else will have size 25 hangar bays")
	}

	// the root will result in a single ArrayOf[RootObjectType]
	for _, e := range root.Elements.Elements() {
		err = AssertIs(e, "ArrayOfShipHull")
		if err != nil {
			return
		}

		for _, e := range e.Elements() {

			// each of these is a shiphull
			err = AssertIs(e, "ShipHull")
			if err != nil {
				return
			}

			// what role does this hull have?
			role := Get(e.Elements(), "Role")
			if !quiet {
				log.Println(role.XMLValue.String())
			}

			// default to auxillary bays
			size := 25
			switch role.XMLValue.String() {
			case "Carrier", "MiningStation", "ResearchStation", "ResortBase":
				// carriers & civilian stations get full bays
				size = 50
			case "DefensiveBase", "MonitoringStation", "SpaceportSmall", "SpaceportMedium", "SpaceportLarge":
				// military bases get huge bays
				size = 100
			}

			elements := statistics.elements

			// for all hangar bay modules, apply the new size
			bays := Get(e.Elements(), "ComponentBays")
			if bays != nil {
				for _, componentBay := range bays.Elements() {

					// every element should be a component bay
					err = AssertIs(componentBay, "ComponentBay")
					if err != nil {
						return
					}

					t := Get(componentBay.Elements(), "Type")
					// s, _ := t.StringValue()
					// log.Printf("%s=%s", t.Name.Local, s)

					// for each component bay whose type is Hangar
					if t.StringValueEquals("Hangar") {

						// find and update the MaximumComponentSize
						c := Get(componentBay.Elements(), "MaximumComponentSize")
						v := fmt.Sprint(size)
						if c.GetStringValue() != v {
							c.SetString(v)
							statistics.changed++
						}

						statistics.elements++
					}
				}
			}

			if statistics.elements != elements {
				statistics.objects++
			}
		}
	}

	return
}

// returns the first matching element from the list of elements based on tag (name)
func Get(elements []*xmltree.XMLElement, tag string) (e *xmltree.XMLElement) {
	for _, e = range elements {
		if e.Name.Local == tag {
			return
		}
	}
	e = nil
	return
}
