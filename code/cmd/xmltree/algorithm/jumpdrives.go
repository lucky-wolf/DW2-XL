package algorithm

import (
	"fmt"
	"log"
)

func HyperDrives(folder string) (err error) {

	log.Println("TODO...")

	// load all component definition files
	j, err := loadJobFor(folder, "ComponentDefinitions*")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyHyperDrives()
	if err != nil {
		return
	}

	// save them all
	j.save()

	return
}

func (j *job) applyHyperDrives() (err error) {

	type schedule int
	const (
		avg schedule = iota
		good
	)

	// data tables
	type ValueTable map[string][]float64
	data := []ValueTable{
		{
			"ComponentIonDefense":          []float64{2, 3, 4, 5, 6, 7, 8},
			"HyperDriveBlockingInsulation": []float64{10, 15, 20, 25, 30, 35, 40},
			"HyperDriveEnergyUsage":        []float64{80, 100, 120, 140, 160, 180, 200},
			"HyperDriveJumpRange":          []float64{160e6, 180e6, 200e6, 230e6, 260e6, 300e6, 340e6},
			"HyperDriveDriveSpeed":         []float64{300e3, 450e3, 600e3, 900e3, 1200e3, 1800e3, 2400e3},
			"HyperDriveJumpInitiationTime": []float64{25, 22.5, 20, 17.5, 15, 12.5, 10},
			"HyperDriveRechargeTime":       []float64{25, 22.5, 20, 17.5, 15, 12.5, 10},
			"HyperDriveDriveJumpAccuracy":  []float64{2000, 2000, 2000, 2000, 2000, 2000, 2000},
		},
		{
			"HyperDriveEnergyUsage":        []float64{50, 70, 90, 110, 130, 150, 170},
			"HyperDriveJumpRange":          []float64{200e6, 240e6, 280e6, 330e6, 380e6, 440e6, 500e6},
			"HyperDriveDriveSpeed":         []float64{400e3, 600e3, 800e3, 1200e3, 1600e3, 2400e3, 3200e3},
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
		"HyperDriveDriveSpeed",
		"HyperDriveDriveJumpAccuracy",
	}

	componentSchedule := map[string]map[string]schedule{
		"Snap Drive": {
			"ComponentIonDefense":          avg,
			"HyperDriveBlockingInsulation": avg,
			"HyperDriveEnergyUsage":        avg,
			"HyperDriveJumpInitiationTime": good,
			"HyperDriveRechargeTime":       good,
			"HyperDriveDriveSpeed":         avg,
			"HyperDriveJumpRange":          avg,
			"HyperDriveDriveJumpAccuracy":  avg,
		},
		"Sojourn Drive": {
			"ComponentIonDefense":          avg,
			"HyperDriveBlockingInsulation": avg,
			"HyperDriveEnergyUsage":        avg,
			"HyperDriveJumpInitiationTime": avg,
			"HyperDriveRechargeTime":       avg,
			"HyperDriveDriveSpeed":         avg,
			"HyperDriveJumpRange":          good,
			"HyperDriveDriveJumpAccuracy":  avg,
		},
		"Hyperstream Drive": {
			"ComponentIonDefense":          avg,
			"HyperDriveBlockingInsulation": avg,
			"HyperDriveEnergyUsage":        avg,
			"HyperDriveJumpInitiationTime": avg,
			"HyperDriveRechargeTime":       avg,
			"HyperDriveDriveSpeed":         good,
			"HyperDriveJumpRange":          avg,
			"HyperDriveDriveJumpAccuracy":  avg,
		},
		"Smart Drive": {
			"ComponentIonDefense":          avg,
			"HyperDriveBlockingInsulation": avg,
			"HyperDriveEnergyUsage":        good,
			"HyperDriveJumpInitiationTime": avg,
			"HyperDriveRechargeTime":       avg,
			"HyperDriveDriveSpeed":         avg,
			"HyperDriveJumpRange":          good,
			"HyperDriveDriveJumpAccuracy":  avg,
		},
	}

	statistics := &j.xfiles[0].stats

	applyStats := func(name string) error {

		// find this drive definition
		e, _ := j.find("Name", name)
		if e == nil {
			return fmt.Errorf("%s not found", name)
		}

		componentSchedule, ok := componentSchedule[name]
		if !ok {
			return fmt.Errorf("component schedule for %s not found", name)
		}

		// fill in the data from our data tables
		stats := e.Child("Values").Elements()
		for i, e := range stats {
			for _, key := range keys {
				e.Child(key).SetValue(data[componentSchedule[key]][key][i])
			}
			statistics.elements++
			statistics.changed++
		}

		return nil
	}

	for _, drive := range []string{"Snap Drive", "Sojourn Drive", "Hyperstream Drive", "Smart Drive"} {
		err = applyStats(drive)
		if err != nil {
			return
		}
		statistics.objects++
	}

	return
}
