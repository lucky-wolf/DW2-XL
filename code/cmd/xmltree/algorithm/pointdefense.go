package algorithm

import (
	"log"
	"strings"
)

func PointDefense(folder string) (err error) {

	// load all component definition files
	j, err := loadJobFor(folder, "ComponentDefinitions*")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyPointDefense()
	if err != nil {
		return
	}

	// save them all
	j.save()

	return
}

func (j *job) applyPointDefense() (err error) {

	for _, f := range j.xfiles {

		statistics := &f.stats

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

				// only intercept weapons...
				if !e.Has("Category", "WeaponIntercept") {
					continue
				}

				// for point defense...
				targetName := e.Child("Name").StringValue()
				if !strings.HasSuffix(targetName, "[PD]") {
					continue
				}

				// find the corresponding small weapon by name
				sourceName := getSourceOfPD(targetName)
				sourceDefinition, _ := j.find("Name", sourceName)
				if sourceDefinition == nil {
					log.Printf("Missing: %s (for %s)", sourceName, targetName)
					continue
				}

				// debug
				if !Quiet {
					log.Printf("%s from %s\n", targetName, sourceName)
				}

				// copy resource requirements
				err = e.CopyByTag("ResourcesRequired", sourceDefinition)
				if err != nil {
					log.Printf("%s: %s from %s", err, targetName, sourceName)
				}

				// copy component stats
				err = e.CopyByTag("Values", sourceDefinition)
				if err != nil {
					log.Printf("%s: %s from %s", err, targetName, sourceName)
				}

				// scale PD nearly identically to fighter weapons
				err = scaleFighterOrPDValues(e)
				if err != nil {
					log.Printf("%s: %s from %s", err, targetName, sourceName)
				}

				statistics.changed++
				statistics.elements++
				statistics.objects++
			}
		}
	}
	err = nil
	return
}

func getSourceOfPD(targetName string) string {
	targetName = targetName[:len(targetName)-len(" [PD]")]
	switch targetName {
	case "Buckler Repeating Blaster":
		return "Maxos Blaster [S]"
	case "Guardian Defense Grid":
		return "Omega Beam [S]"
	case "Maelstrom Defender":
		return "Titan Blaster [S]"
	case "Point Defense Cannon":
		return "Rail Gun [S]"
	case "Sentinel Multi-Beam Defense":
		return "Thuon Beam [S]"
	case "Interceptor Missile":
		return "Concussion Missile [S]"
	case "Aegis Missile Battery":
		return "Lightning Missile [S]"
	default:
		// ion cannon
		// ion rapid pulse array
		// impact assault blaster
		// terminator autocannon
		// hive missile battery
		// reinforcing swarm battery
		return targetName + " [S]"
	}
}
