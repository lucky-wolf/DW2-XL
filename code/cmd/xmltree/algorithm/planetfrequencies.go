package algorithm

import (
	"log"
)

func ScalePlanetFrequencies(folder string, factor float64) (err error) {

	log.Println("Planets will be updated to desired relative frequencies")

	// load all orbtype files
	j, err := LoadJobFor(folder, "OrbTypes*.xml")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.scalePlanetFrequencies(factor)
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

func (j *Job) scalePlanetFrequencies(factor float64) (err error) {

	for _, f := range j.xfiles {

		statistics := &f.stats

		// the root will result in a single ArrayOf[RootObjectType]
		for _, e := range f.root.Elements.Elements() {

			err = assertIs(e, "ArrayOfOrbType")
			if err != nil {
				return
			}

			for _, e := range e.Elements() {

				// each of these is a OrbType
				err = assertIs(e, "OrbType")
				if err != nil {
					return
				}

				// we're adjusting stars
				// if e.Child("Category").StringValue() != "Star" {
				// 	continue
				// }

				// go through all children and apply our adjustment there
				for _, f := range e.Child("ChildTypes").Elements() {

					err = assertIs(f, "OrbTypeFactor")
					if err != nil {
						return
					}

					// we want to adjust those which are not colonizable (things that have them as children)
					switch f.Child("OrbTypeId").IntValue() {
					case 7, 8, 9, 10, 11, 12, 17, 18, 19, 20, 21, 22, 23, 27, 29, 30:
						err = f.ScaleChildBy("Factor", factor)
						if err != nil {
							return
						}
						statistics.changed++
					}
				}

				statistics.elements++
			}

			statistics.objects++
		}
	}
	err = nil
	return
}
