package algorithm

import "log"

func All(folder string) (err error) {

	log.Println("All components and research values will be updated for all algorithms to-date")

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
