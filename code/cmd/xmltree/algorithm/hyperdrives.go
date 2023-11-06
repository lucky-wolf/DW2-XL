package algorithm

import (
	"fmt"
	"log"
)

func HyperDrives(folder string) (err error) {

	log.Println("Updates core hyperdrives to have 7 levels off of a common data table")

	// load all component definition files
	j, err := LoadJobFor(folder, "ComponentDefinitions*")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyHyperDrives()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

func (j *Job) applyHyperDrives() (err error) {

	type schedule int
	const (
		avg schedule = iota
		good
		best
		worst
	)

	// data tables
	type ValueTable map[string][]float64
	data := []ValueTable{
		{ // avg
			"ComponentIonDefense":          []float64{2, 3, 4, 5, 6, 7, 8},
			"HyperDriveBlockingInsulation": []float64{10, 15, 20, 25, 30, 35, 40},
			"HyperDriveEnergyUsage":        []float64{70, 90, 110, 130, 150, 170, 190},
			// "HyperDriveJumpRange":          []float64{160e6, 180e6, 200e6, 230e6, 260e6, 300e6, 340e6},
			"HyperDriveJumpRange":          []float64{120e6, 150e6, 187.5e6, 234.375e6, 292.96875e6, 366.2109375e6, 457.763672e6}, // 1.25% compounding
			"HyperDriveSpeed":              []float64{300e3, 450e3, 600e3, 900e3, 1200e3, 1800e3, 2400e3},
			"HyperDriveJumpInitiationTime": []float64{25, 22.5, 20, 17.5, 15, 12.5, 10},
			"HyperDriveRechargeTime":       []float64{25, 22.5, 20, 17.5, 15, 12.5, 10},
			"HyperDriveJumpAccuracy":       []float64{2000, 2000, 2000, 2000, 2000, 2000, 2000},
		},
		{ // good
			"HyperDriveEnergyUsage":        []float64{60, 80, 100, 120, 140, 160, 180}, // -10 vs. avg
			"HyperDriveJumpInitiationTime": []float64{18, 16, 14, 12, 10, 8, 6},
			"HyperDriveRechargeTime":       []float64{18, 16, 14, 12, 10, 8, 6},
			"HyperDriveSpeed":              []float64{350e3, 525e3, 700e3, 1050e3, 1400e3, 2100e3, 2800e3},
		},
		{ // best
			"HyperDriveEnergyUsage": []float64{50, 70, 90, 110, 130, 150, 170}, // -20 vs. avg (-10 vs. good)
			// "HyperDriveJumpRange":          []float64{200e6, 240e6, 280e6, 330e6, 380e6, 440e6, 500e6},
			"HyperDriveJumpRange":          []float64{180e6, 225e6, 281.25e6, 351.5625e6, 439.453125e6, 549.3164063e6, 686.645508e6}, // 50% better than avg at each level
			"HyperDriveSpeed":              []float64{400e3, 600e3, 800e3, 1200e3, 1600e3, 2400e3, 3200e3},
			"HyperDriveJumpInitiationTime": []float64{12.5, 11.25, 10, 8.75, 7.5, 6.25, 5},
			"HyperDriveRechargeTime":       []float64{12.5, 11.25, 10, 8.75, 7.5, 6.25, 5},
		},
		{ // worst
			"HyperDriveEnergyUsage": []float64{80, 100, 120, 140, 160, 180, 200}, // +10 vs. avg
		},
	}

	keys := []string{
		"ComponentIonDefense",
		"HyperDriveBlockingInsulation",
		"HyperDriveEnergyUsage",
		"HyperDriveJumpInitiationTime",
		"HyperDriveJumpRange",
		"HyperDriveRechargeTime",
		"HyperDriveSpeed",
		"HyperDriveJumpAccuracy",
	}

	componentSchedule := map[string]map[string]schedule{
		"Snap Drive": {
			"ComponentIonDefense":          avg,
			"HyperDriveBlockingInsulation": avg,
			"HyperDriveEnergyUsage":        avg,
			"HyperDriveJumpInitiationTime": best,
			"HyperDriveRechargeTime":       best,
			"HyperDriveSpeed":              avg,
			"HyperDriveJumpRange":          avg,
			"HyperDriveJumpAccuracy":       avg,
		},
		"Sojourn Drive": {
			"ComponentIonDefense":          avg,
			"HyperDriveBlockingInsulation": avg,
			"HyperDriveEnergyUsage":        avg,
			"HyperDriveJumpInitiationTime": avg,
			"HyperDriveRechargeTime":       avg,
			"HyperDriveSpeed":              avg,
			"HyperDriveJumpRange":          best,
			"HyperDriveJumpAccuracy":       avg,
		},
		"Hyperstream Drive": {
			"ComponentIonDefense":          avg,
			"HyperDriveBlockingInsulation": avg,
			"HyperDriveEnergyUsage":        worst,
			"HyperDriveJumpInitiationTime": avg,
			"HyperDriveRechargeTime":       avg,
			"HyperDriveSpeed":              best,
			"HyperDriveJumpRange":          avg,
			"HyperDriveJumpAccuracy":       avg,
		},
		"Smart Drive": {
			"ComponentIonDefense":          avg,
			"HyperDriveBlockingInsulation": avg,
			"HyperDriveEnergyUsage":        best,
			"HyperDriveJumpInitiationTime": avg,
			"HyperDriveRechargeTime":       avg,
			"HyperDriveSpeed":              avg,
			"HyperDriveJumpRange":          best,
			"HyperDriveJumpAccuracy":       avg,
		},
		"Velocity Drive": {
			"ComponentIonDefense":          avg,
			"HyperDriveBlockingInsulation": avg,
			"HyperDriveEnergyUsage":        good,
			"HyperDriveJumpInitiationTime": good,
			"HyperDriveRechargeTime":       good,
			"HyperDriveSpeed":              good,
			"HyperDriveJumpRange":          avg,
			"HyperDriveJumpAccuracy":       avg,
		},
	}

	applyStats := func(name string) (err error) {

		// find this drive definition
		e, f := j.FindElement("Name", name)
		if e == nil {
			return fmt.Errorf("%s not found", name)
		}

		statistics := &f.stats

		componentSchedules, ok := componentSchedule[name]
		if !ok {
			return fmt.Errorf("component schedule for %s not found", name)
		}

		// fill in the data from our data tables
		stats := e.Child("Values").Elements()
		for i, e := range stats {
			for _, key := range keys {

				values, ok := data[componentSchedules[key]][key]
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
