package algorithm

import (
	"log"
	"strings"
)

func FighterReactors(folder string) (err error) {

	log.Println("All strikecraft reactors will be scaled to ship components")

	// load all component definition files
	j, err := LoadJobFor(folder, "ComponentDefinitions*.xml")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyFighterReactors()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

func (j *Job) applyFighterReactors() (err error) {

	for _, f := range j.xfiles {

		// the root will result in a single ArrayOf[RootObjectType]
		for _, e := range f.root.Elements.Elements() {

			err = assertIs(e, "ArrayOfComponentDefinition")
			if err != nil {
				return
			}

			for _, e := range e.Elements() {

				// each of these is a componentdefinition
				err = assertIs(e, "ComponentDefinition")
				if err != nil {
					return
				}

				// only reactors...
				if !e.HasChildWithValue("Category", "Reactor") {
					continue
				}

				// for fighters...
				targetName := e.Child("Name").StringValue()
				if !strings.HasSuffix(targetName, "[Ftr]") {
					continue
				}

				// find the corresponding ship reactors by same name
				sourceName := strings.TrimSpace(targetName[:len(targetName)-len("[Ftr]")])
				source, _ := j.FindElement("Name", sourceName)
				if source == nil {
					log.Printf("element not found: %s for %s", sourceName, targetName)
					continue
				}

				// debug
				if !Quiet {
					log.Printf("%s from %s\n", targetName, sourceName)
				}

				// do it
				err = j.DeriveFromComponent(f, source, e)
				if err != nil {
					return
				}
			}
		}
	}

	return
}
