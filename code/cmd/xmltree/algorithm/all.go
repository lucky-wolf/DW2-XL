package algorithm

import "log"

func All(folder string) (err error) {

	// do components first
	err = Components(folder)
	if err != nil {
		return
	}

	// now do research costs
	err = ResearchCosts(folder)
	if err != nil {
		return
	}

	return
}

func Components(folder string) (err error) {

	log.Println("All components will be updated for all algorithms to-date")

	// load all component definition files
	j, err := LoadJobFor(folder, "ComponentDefinitions*.xml")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyComponents()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

func (j *Job) applyComponents() (err error) {

	// primary ship components first
	err = j.applyHyperDrives()
	if err != nil {
		return
	}
	err = j.applyIonShields()
	if err != nil {
		return
	}
	err = j.applyIonWeapons()
	if err != nil {
		return
	}

	// then derivative components
	err = j.applyFighterArmor()
	if err != nil {
		return
	}
	err = j.applyFighterEngines()
	if err != nil {
		return
	}
	err = j.applyFighterReactors()
	if err != nil {
		return
	}
	err = j.applyFighterShields()
	if err != nil {
		return
	}
	err = j.applyFighterWeaponsAndPD()
	if err != nil {
		return
	}

	return
}
