package algorithm

import (
	"fmt"
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
)

func IonShields(folder string) (err error) {

	log.Println("Updates core ion shields to have 5 levels off of a common data table")

	// load all component definition files
	j, err := loadJobFor(folder, "ComponentDefinitions*")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyIonShields()
	if err != nil {
		return
	}

	// save them all
	j.save()

	return
}

func (j *job) applyIonShields() (err error) {

	// data tables
	type ValueTable map[string]func(level int) float64
	data := ValueTable{
		"ComponentIonDefense": func(level int) float64 { return 10 + float64(2*level) },
		"IonDamageDefense":    func(level int) float64 { return 8 + float64(2*level) },
		"CrewRequirement":     func(level int) float64 { return 5 },
		"StaticEnergyUsed":    func(level int) float64 { return 4 + float64(2*level) },
	}

	applyStats := func(name string) (err error) {

		// find this drive definition
		e, f := j.find("Name", name)
		if e == nil {
			return fmt.Errorf("%s not found", name)
		}

		statistics := &f.stats

		// fill in the data from our data tables
		stats := e.Child("Values").Elements()
		for i, e := range stats {
			for key, f := range data {
				e.Child(key).SetValue(f(i))
			}
			statistics.elements++
			statistics.changed++
		}

		statistics.objects++

		// apply to [ftr]
		sourceDefinition := e
		name += " [Ftr]"
		e, f = j.find("Name", name)
		if e == nil {
			err = fmt.Errorf("%s not found", name)
			return
		}

		// copy and scale resource requirements
		err = e.CopyAndVisitByTag("ResourcesRequired", sourceDefinition, func(e *xmltree.XMLElement) error { e.Child("Amount").ScaleBy(0.2); return nil })
		if err != nil {
			log.Println(err)
		}

		// copy component stats
		err = e.CopyByTag("Values", sourceDefinition)
		if err != nil {
			log.Println(err)
		}

		// now that we have our own copy of the component stats (same number of levels too)
		// we can update each of those to scale for [Ftr] version
		for _, e := range e.Child("Values").Elements() {

			// every element should be a component bay
			err = assertIs(e, "ComponentStats")
			if err != nil {
				return
			}

			// scale down the ion defenses (or ion PD / ftr weapons will never penetrate)
			e.Child("ComponentIonDefense").ScaleBy(0.2)
			e.Child("IonDamageDefense").ScaleBy(0.2)

			// never a crew requirement for fighter components
			e.Child("CrewRequirement").SetValue(0)

			statistics.changed++
			statistics.elements++
		}

		statistics.objects++

		return
	}

	// apply to ship & ftr component
	err = applyStats("Ion Shield")
	if err != nil {
		return
	}

	return
}
