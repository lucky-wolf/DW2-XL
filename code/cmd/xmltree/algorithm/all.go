package algorithm

import "log"

func All(folder string) (err error) {

	log.Println("All components and research values will be updated for all algorithms to-date")

	// load all component and research definition files
	j, err := LoadJobFor(folder, "ComponentDefinitions*.xml", "ResearchProjectDefinitions*.xml")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyComponents()
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyResearchCosts()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}
