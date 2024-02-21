package algorithm

import (
	"fmt"
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
)

func PartialOrdering(folder string) (err error) {

	log.Println("All xml objects will have their ID and name placed first")

	// do all file types we can handle
	j, err := LoadJobFor(folder, "ComponentDefinitions*.xml", "OrbTypes*.xml", "PlanetaryFacilityDefinitions*.xml", "Races*.xml", "ResearchProjectDefinitions*.xml", "ShipHulls*.xml", "TroopDefinitions*.xml")
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
				err = j.applyPartialOrderingTo(f, arrayOf, "ComponentId", "Name", "Description", "ImageFilename", "Size", "Category", "IsFighterOnly", "Family")
			case "ArrayOfOrbType":
				err = j.applyPartialOrderingTo(f, arrayOf, "OrbTypeId", "Category", "Name", "Description", "ImageFilename")
			case "ArrayOfPlanetaryFacilityDefinition":
				err = j.applyPartialOrderingTo(f, arrayOf, "PlanetaryFacilityDefinitionId", "Name", "Description", "ImageFilename", "Size", "IsRuins", "BuildCost", "MaintenanceCost", "FacilityFamilyId", "FacilityFamilyLevel", "ExclusiveWithinFamily")
			case "ArrayOfRace":
				err = j.applyPartialOrderingTo(f, arrayOf, "RaceId", "Name", "Description", "ImageFilename", "BundleName")
			case "ArrayOfResearchProjectDefinition":
				err = j.applyPartialOrderingTo(f, arrayOf, "ResearchProjectId", "Name", "Description", "ImageFilename")
			case "ArrayOfShipHull":
				err = j.applyPartialOrderingTo(f, arrayOf, "ShipHullId", "Name", "Description", "ImageFilename")
			case "ArrayOfTroopDefinition":
				err = j.applyPartialOrderingTo(f, arrayOf, "TroopDefinitionId", "Name", "Description", "ImageFilename", "RaceId", "Type", "Size")
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
	f.stats.objects = len(elements)
	for _, object := range elements {

		// for each first, find it, and insert at next top position
		to := object.ZeroElementIndex()
		for _, tag := range firsts {

			// see if this element exists
			from := object.ChildIndex(tag)
			if from == -1 {
				continue
			}

			// yes, then see if we need to reorder it (not already in correct spot)
			f.stats.elements++
			if from != to {
				object.Reorder(from, to)
				f.stats.changed++
			}

			// we'll place the next element one further down
			to++
		}
	}

	return
}
