package algorithm

import (
	"log"
	"strings"
)

// TODO: we need to increase interceptor missiles (PD) speed so they can catch other missiles & fighters!
// maybe also slow down their ROF / increase dmg/hit (say by factor of 2x)

func FighterWeaponsAndPD(folder string) (err error) {

	log.Println("All strikecraft weapons and PD weapons will be scaled to ship components")

	// load all component definition files
	j, err := LoadJobFor(folder, "ComponentDefinitions*")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyFighterWeaponsAndPD()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

func (j *Job) applyFighterWeaponsAndPD() (err error) {

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

				// only weapons...
				if !e.HasPrefix("Category", "Weapon") {
					continue
				}

				// for fighters and PD...
				targetName := e.Child("Name").StringValue()
				if !strings.HasSuffix(targetName, "[Ftr]") && !strings.HasSuffix(targetName, "[PD]") {
					continue
				}

				// find the corresponding small weapon by name
				sourceName := GetComponentSourceName(targetName, e.Has("IsFighterOnly", "true"))
				source, _ := j.FindElement("Name", sourceName)
				if source == nil {
					log.Printf("Missing: %s (for %s)", sourceName, targetName)
					continue
				}

				// debug
				if !Quiet {
					log.Printf("%s from %s\n", targetName, sourceName)
				}

				// do it
				err = j.ScaleComponentToComponent(f, source, e)
				if err != nil {
					return
				}
			}
		}
	}

	return
}
