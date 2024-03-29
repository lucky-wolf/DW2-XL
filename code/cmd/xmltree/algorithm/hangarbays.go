package algorithm

import (
	"fmt"
	"log"
)

func HangarBays(folder string) (err error) {

	log.Println("All hangar bays will be limited by hull classification")

	// load the ship hull files
	j, err := LoadJobFor(folder, "ShipHulls_*")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyHangarBays()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

func (j *Job) applyHangarBays() (err error) {

	for _, f := range j.xfiles {

		statistics := &f.stats

		// the root will result in a single ArrayOf[RootObjectType]
		for _, e := range f.root.Elements.Elements() {
			err = assertIs(e, "ArrayOfShipHull")
			if err != nil {
				return
			}

			for _, e := range e.Elements() {

				// each of these is a shiphull
				err = assertIs(e, "ShipHull")
				if err != nil {
					return
				}

				// what role does this hull have?
				role := e.Child("Role")
				if !Quiet {
					log.Println(role.XMLValue.String())
				}

				// default to auxiliary bays
				var size int
				switch role.XMLValue.String() {
				case "ConstructionShip":
					// constructors need a size 100 hangar for construction yard module!
					size = 100
				case "SpaceportSmall", "SpaceportMedium", "SpaceportLarge":
					// military & construction - both use 100
					size = 100
				case "DefensiveBase", "MonitoringStation":
					// military bases get huge bays
					size = 100
				case "Carrier", "MiningStation", "ResearchStation", "ResortBase":
					// carriers & civilian stations get full bays
					size = 50
				default:
					if e.Child("RaceId").StringValue() == "2" {
						// teekans special ability is to get
						size = 50
					} else {
						// all others just get aux bays if they have any hangar bays at all
						size = 25
					}
				}

				elements := statistics.elements

				// for all hangar bay modules, apply the new size
				for _, componentBay := range e.Child("ComponentBays").Elements() {

					// every element should be a component bay
					err = assertIs(componentBay, "ComponentBay")
					if err != nil {
						return
					}

					// for each component bay whose type is Hangar
					if componentBay.HasChildWithValue("Type", "Hangar") {

						// find and update the MaximumComponentSize
						c := componentBay.Child("MaximumComponentSize")
						v := fmt.Sprint(size)
						if c.StringValue() != v {
							c.SetString(v)
							statistics.changed++
						}

						statistics.elements++
					}
				}

				if statistics.elements != elements {
					statistics.objects++
				}
			}
		}
	}

	return
}
