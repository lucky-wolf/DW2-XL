package algorithm

import (
	"fmt"
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
)

func HangarBays(root *xmltree.XMLTree) (statistics Statistics, err error) {

	if !Quiet {
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
					err = AssertIs(componentBay, "ComponentBay")
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

	return
}
