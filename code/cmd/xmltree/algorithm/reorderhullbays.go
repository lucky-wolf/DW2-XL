package algorithm

import (
	"log"
)

func RenumberHullComponentBays(folder string) (err error) {

	log.Println("All component bay indexes will be fixed to a simple incremental index")

	// load all ship hull definition files
	j, err := LoadJobFor(folder, "ShipHulls*.xml")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.renumberComponentBays()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

// renumbers all elements and returns the group counts (weapon, defense, etc.)
func (j *Job) renumberComponentBays() (err error) {

	for _, f := range j.xfiles {

		// the root will result in a single ArrayOf[RootObjectType]
		for _, shiphulls := range f.root.Elements.Elements() {

			err = assertIs(shiphulls, "ArrayOfShipHull")
			if err != nil {
				return
			}

			for _, shiphull := range shiphulls.Elements() {

				// each of these is a ShipHull
				err = assertIs(shiphull, "ShipHull")
				if err != nil {
					return
				}

				f.stats.objects++

				// scan the component bays in order
				for i, c := range shiphull.Child("ComponentBays").Elements() {

					f.stats.elements++

					// id is trivial - just linear numbering within this list
					e := c.Child("ComponentBayId")
					if e.IntValue() != i {
						e.SetValue(i)
						f.stats.changed++
					}

					// todo: to do mesh names, we'd need to have a master list of meshes per model
					// todo: (ordered from what should be given first to last)
					// todo: for now, this is too complex and doesn't really make sense
					// todo: what we can do is take the "archetype" hull component bays definition - and trim that somehow
					// todo: rather than trying to grow anything from a lesser model
					// todo: this per-race, per 3d model, basically, this is super hard to do algorithmically
					// todo: so for now, we can renumber component bays after you edit a ship hull definition
					// todo: but it's up to you to manage the mesh identities
				}
			}
		}
	}

	return
}
