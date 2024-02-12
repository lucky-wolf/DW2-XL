package algorithm

import (
	"log"
)

func Troops(folder string) (err error) {

	log.Println("All strikecraft armor will be scaled to ship components")

	// load all component definition files
	j, err := LoadJobFor(folder, "TroopDefinitions*.xml")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyTroops()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

type Type = string
type Key = string
type Value = string
type TroopValue map[Key]Value

// hack: fuck - so because battlebots don't have a type, this won't really work (they're the one exception that breaks everything)
// warn: also Quameno have no true infantry - only a copy of those battlebots
// warn: and then Teekan infantry + spec-ops are size 4K, instead of the usual 5k
var TroopValues = map[Type]TroopValue{
	"Infantry": TroopValue{
		"Size":            "2000",
		"RecruitmentCost": "crap",
	},
}

func (j *Job) applyTroops() (err error) {

	for _, f := range j.xfiles {

		// the root will result in a single ArrayOf[RootObjectType]
		for _, e := range f.root.Elements.Elements() {

			err = assertIs(e, "ArrayOfTroopDefinition")
			if err != nil {
				return
			}

			for _, e := range e.Elements() {

				// each of these is a componentdefinition
				err = assertIs(e, "TroopDefinition")
				if err != nil {
					return
				}

				// only armor...
				t := e.Child("Type")
				if t == nil {
					continue
				}

				// todo: apply values
				// values, ok := TroopValues[t.StringValue()]
				// if !ok {
				// 	continue
				// }
			}
		}
	}

	return
}
