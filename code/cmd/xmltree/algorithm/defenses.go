package algorithm

import (
	"fmt"
	"log"
)

func Defenses(folder string) (err error) {

	log.Println("Updates Defenses (Shields & Armor) off of a common data table")

	// load all component definition files
	j, err := LoadJobFor(folder, "ComponentDefinitions*")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyDefenses()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

func (j *Job) applyDefenses() (err error) {

	type schedule int
	const (
		avg schedule = iota
		good
		best
	)

	// data tables
	type ValueTable map[string][]float64
	data := []ValueTable{
		{
			"ComponentIonDefense":          []float64{2, 3, 4, 5, 6, 7, 8},
			"HyperDriveBlockingInsulation": []float64{10, 15, 20, 25, 30, 35, 40},
			"HyperDriveEnergyUsage":        []float64{80, 100, 120, 140, 160, 180, 200},
			"HyperDriveJumpRange":          []float64{160e6, 180e6, 200e6, 230e6, 260e6, 300e6, 340e6},
			"Defensespeed":                 []float64{300e3, 450e3, 600e3, 900e3, 1200e3, 1800e3, 2400e3},
			"HyperDriveJumpInitiationTime": []float64{25, 22.5, 20, 17.5, 15, 12.5, 10},
			"HyperDriveRechargeTime":       []float64{25, 22.5, 20, 17.5, 15, 12.5, 10},
			"HyperDriveJumpAccuracy":       []float64{2000, 2000, 2000, 2000, 2000, 2000, 2000},
		},
		{
			"HyperDriveEnergyUsage":        []float64{72, 90, 108, 126, 144, 162, 180}, // 10% discount vs. avg
			"HyperDriveJumpInitiationTime": []float64{18, 16, 14, 12, 10, 8, 6},
			"HyperDriveRechargeTime":       []float64{18, 16, 14, 12, 10, 8, 6},
			"Defensespeed":                 []float64{350e3, 525e3, 700e3, 1050e3, 1400e3, 2100e3, 2800e3},
		},
		{
			"HyperDriveEnergyUsage":        []float64{50, 70, 90, 110, 130, 150, 170},
			"HyperDriveJumpRange":          []float64{200e6, 240e6, 280e6, 330e6, 380e6, 440e6, 500e6},
			"Defensespeed":                 []float64{400e3, 600e3, 800e3, 1200e3, 1600e3, 2400e3, 3200e3},
			"HyperDriveJumpInitiationTime": []float64{12.5, 11.25, 10, 8.75, 7.5, 6.25, 5},
			"HyperDriveRechargeTime":       []float64{12.5, 11.25, 10, 8.75, 7.5, 6.25, 5},
		},
	}

	keys := []string{
		"ComponentIonDefense",
		"HyperDriveBlockingInsulation",
		"HyperDriveEnergyUsage",
		"HyperDriveJumpInitiationTime",
		"HyperDriveJumpRange",
		"HyperDriveRechargeTime",
		"Defensespeed",
		"HyperDriveJumpAccuracy",
	}

	componentSchedule := map[string]map[string][]float64{
		"Snap Drive": {
			"ComponentIonDefense":          data[avg]["ComponentIonDefense"],
			"HyperDriveBlockingInsulation": data[avg]["HyperDriveBlockingInsulation"],
			"HyperDriveEnergyUsage":        data[avg]["HyperDriveEnergyUsage"],
			"HyperDriveJumpInitiationTime": data[best]["HyperDriveJumpInitiationTime"],
			"HyperDriveRechargeTime":       data[best]["HyperDriveRechargeTime"],
			"Defensespeed":                 data[avg]["Defensespeed"],
			"HyperDriveJumpRange":          data[avg]["HyperDriveJumpRange"],
			"HyperDriveJumpAccuracy":       data[avg]["HyperDriveJumpAccuracy"],
		},
		"Sojourn Drive": {
			"ComponentIonDefense":          data[avg]["ComponentIonDefense"],
			"HyperDriveBlockingInsulation": data[avg]["HyperDriveBlockingInsulation"],
			"HyperDriveEnergyUsage":        data[avg]["HyperDriveEnergyUsage"],
			"HyperDriveJumpInitiationTime": data[avg]["HyperDriveJumpInitiationTime"],
			"HyperDriveRechargeTime":       data[avg]["HyperDriveRechargeTime"],
			"Defensespeed":                 data[avg]["Defensespeed"],
			"HyperDriveJumpRange":          data[best]["HyperDriveJumpRange"],
			"HyperDriveJumpAccuracy":       data[avg]["HyperDriveJumpAccuracy"],
		},
		"Hyperstream Drive": {
			"ComponentIonDefense":          data[avg]["ComponentIonDefense"],
			"HyperDriveBlockingInsulation": data[avg]["HyperDriveBlockingInsulation"],
			"HyperDriveEnergyUsage":        data[avg]["HyperDriveEnergyUsage"],
			"HyperDriveJumpInitiationTime": data[avg]["HyperDriveJumpInitiationTime"],
			"HyperDriveRechargeTime":       data[avg]["HyperDriveRechargeTime"],
			"Defensespeed":                 data[best]["Defensespeed"],
			"HyperDriveJumpRange":          data[avg]["HyperDriveJumpRange"],
			"HyperDriveJumpAccuracy":       data[avg]["HyperDriveJumpAccuracy"],
		},
		"Smart Drive": {
			"ComponentIonDefense":          data[avg]["ComponentIonDefense"],
			"HyperDriveBlockingInsulation": data[avg]["HyperDriveBlockingInsulation"],
			"HyperDriveEnergyUsage":        data[best]["HyperDriveEnergyUsage"],
			"HyperDriveJumpInitiationTime": data[avg]["HyperDriveJumpInitiationTime"],
			"HyperDriveRechargeTime":       data[avg]["HyperDriveRechargeTime"],
			"Defensespeed":                 data[avg]["Defensespeed"],
			"HyperDriveJumpRange":          data[best]["HyperDriveJumpRange"],
			"HyperDriveJumpAccuracy":       data[avg]["HyperDriveJumpAccuracy"],
		},
		"Velocity Drive": {
			"ComponentIonDefense":          data[avg]["ComponentIonDefense"],
			"HyperDriveBlockingInsulation": data[avg]["HyperDriveBlockingInsulation"],
			"HyperDriveEnergyUsage":        data[good]["HyperDriveEnergyUsage"],
			"HyperDriveJumpInitiationTime": data[good]["HyperDriveJumpInitiationTime"],
			"HyperDriveRechargeTime":       data[good]["HyperDriveRechargeTime"],
			"Defensespeed":                 data[good]["Defensespeed"],
			"HyperDriveJumpRange":          data[avg]["HyperDriveJumpRange"],
			"HyperDriveJumpAccuracy":       data[avg]["HyperDriveJumpAccuracy"],
		},
	}

	applyStats := func(name string) (err error) {

		// find this drive definition
		e, f := j.FindElement("Name", name)
		if e == nil {
			return fmt.Errorf("%s not found", name)
		}

		statistics := &f.stats

		componentValues, ok := componentSchedule[name]
		if !ok {
			return fmt.Errorf("component schedule for %s not found", name)
		}

		// fill in the data from our data tables
		stats := e.Child("Values").Elements()
		for i, e := range stats {
			for _, key := range keys {

				values, ok := componentValues[key]
				if !ok {
					err = fmt.Errorf("component values for %s not found", key)
					return
				}

				e.Child(key).SetValue(values[i])
			}
			statistics.elements++
			statistics.changed++
		}

		statistics.objects++

		return nil
	}

	for _, drive := range []string{"Snap Drive", "Sojourn Drive", "Hyperstream Drive", "Smart Drive", "Velocity Drive"} {
		err = applyStats(drive)
		if err != nil {
			return
		}
	}

	return
}
