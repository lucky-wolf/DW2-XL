package algorithm

import (
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
)

func NamesFirst(folder string) (err error) {

	log.Println("All xml objects will have their ID and name placed first")

	// start with ship hulls
	j, err := LoadJobFor(folder, "ShipHulls_Ackdarian.xml")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyNamesFirst()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

func (j *Job) applyNamesFirst() (err error) {

	for _, f := range j.xfiles {

		// the root will result in a single ArrayOf[RootObjectType]
		for _, e := range f.root.Elements.Elements() {

			switch e.Name.Local {
			case "ArrayOfShipHull":
				err = j.applyFieldOrdering(e, "ShipHullId", "Name")
			}
			if err != nil {
				return
			}
		}
	}

	err = nil
	return
}

func (j *Job) applyFieldOrdering(collection *xmltree.XMLElement, firsts ...string) (err error) {

	elements := collection.Elements()
	count := len(elements)
	for _, object := range elements {

		// for each first, find it, and insert at next top position
		to := 0
		for _, tag := range firsts {
			from := object.ChildIndex(tag)
			if from != count {
				object.Reorder(from, to)
				to++
			}
		}
	}

	return
}
