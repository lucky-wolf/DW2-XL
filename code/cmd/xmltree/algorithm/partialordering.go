package algorithm

import (
	"fmt"
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
)

func PartialOrdering(folder string) (err error) {

	log.Println("All xml objects will have their ID and name placed first")

	// do all file types we can handle
	j, err := LoadJobFor(folder, "ComponentDefinitions*.xml") //, "OrbTypes*.xml", "PlanetaryFacilityDefinitions*.xml", "Races*.xml", "ResearchProjectDefinitions*.xml", "ShipHulls*.xml")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyPartialOrdering()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

func (j *Job) applyPartialOrdering() (err error) {

	for _, f := range j.xfiles {

		// the root will result in a single ArrayOf[RootObjectType]
		for _, arrayOf := range f.root.Elements.Elements() {

			switch arrayOf.Name.Local {
			case "ArrayOfComponentDefinition":
				err = j.applyPartialOrderingTo(f, arrayOf, "ComponentId", "Name", "Description", "ImageFilename")
			case "ArrayOfOrbType":
				err = j.applyPartialOrderingTo(f, arrayOf, "OrbTypeId", "Name", "Description", "ImageFilename", "Category")
			case "ArrayOfPlanetaryFacilityDefinition":
				err = j.applyPartialOrderingTo(f, arrayOf, "PlanetaryFacilityDefinitionId", "Name", "Description", "ImageFilename")
			case "ArrayOfRace":
				err = j.applyPartialOrderingTo(f, arrayOf, "RaceId", "Name", "Description", "ImageFilename", "BundleName")
			case "ArrayOfResearchProjectDefinition":
				err = j.applyPartialOrderingTo(f, arrayOf, "ResearchProjectId", "Name", "Description", "ImageFilename")
			case "ArrayOfShipHull":
				err = j.applyPartialOrderingTo(f, arrayOf, "ShipHullId", "Name", "ImageFilename")
			case "ArrayOfTroopDefinition":
				err = j.applyPartialOrderingTo(f, arrayOf, "TroopDefinitionId", "RaceId", "Name", "ImageFilename", "Type")
			default:
				err = fmt.Errorf("no ordering defined for: %s", arrayOf.Name.Local)
			}
			if err != nil {
				return
			}
		}
	}

	err = nil
	return
}

func (j *Job) applyPartialOrderingTo(f *XFile, arrayOf *xmltree.XMLElement, firsts ...string) (err error) {

	elements := arrayOf.Elements()
	count := len(elements)
	f.stats.objects = count
	for _, object := range elements {

		// for each first, find it, and insert at next top position
		to := 0
		for _, tag := range firsts {
			from := object.ChildIndex(tag)
			if from != count {
				f.stats.elements++
				if from != to {
					object.Reorder(from, to)
					f.stats.changed++
				}
				to++
			}
		}
	}

	return
}
