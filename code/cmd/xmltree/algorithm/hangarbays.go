package algorithm

import (
	"fmt"
	"log"
)

func HangarBays(folder string) (err error) {

	if !Quiet {
		log.Println("All stations will have size 100 hangar bays")
		log.Println("All carriers will have size 50 hangar bays")
		log.Println("Everything else will have size 25 hangar bays")
	}

	// load the ship hull files
	j, err := loadJobFor(folder, "ShipHulls_*")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyHangarBays()
	if err != nil {
		return
	}

	// save them all
	j.save()

	return
}

func (j *job) applyHangarBays() (err error) {

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

				// default to auxillary bays
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
					// all others just get aux bays if they have any hangar bays at all
					size = 25
				}

				elements := statistics.elements

				// for all hangar bay modules, apply the new size
				bays := e.Child("ComponentBays")
				if bays != nil {
					for _, componentBay := range bays.Elements() {

						// every element should be a component bay
						err = assertIs(componentBay, "ComponentBay")
						if err != nil {
							return
						}

						// for each component bay whose type is Hangar
						if componentBay.Has("Type", "Hangar") {

							// find and update the MaximumComponentSize
							c := componentBay.Child("MaximumComponentSize")
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
	}

	return
}
