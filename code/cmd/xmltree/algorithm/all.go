package algorithm

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
