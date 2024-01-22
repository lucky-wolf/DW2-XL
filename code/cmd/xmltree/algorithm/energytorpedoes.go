package algorithm

import (
	"log"
)

func TorpedoWeapons(folder string) (err error) {

	log.Println("Updates torpedo weapons off of a common core data table")

	// load all component definition files
	j, err := LoadJobFor(folder, "ComponentDefinitions*.xml")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyTorpedoWeapons()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

// standard weapon countermeasure schedule (by tech level)
func ComponentCountermeasuresBonus(level int) float64 {
	return 0.6 + float64(level)*0.02
}

// t0 ... t10
// 300, 325, 350, 375, 400, 425, 450, 475, 500, 525, 550
// todo: drive this off of WeaponFireType or Family
func TorpedoSeekingSpeed(level int) float64 { return 300 + 25*float64(level) }

func (j *Job) applyTorpedoWeapons() (err error) {

	// // this is 50% slower (2/3 of) blasters
	// WeaponRateOfFire := func(level int) float64 { return 13.5 }

	// // standard damage is based on pulsed blasters-ish, but at about 2/3 the ROF, so 2/3 the DPS
	// WeaponRawDamage := func(level int) float64 {
	// 	// this gives us 20 at (t0) and a gain of 20% per level beyond that
	// 	return 20 * (1 + 0.2*float64(level-1))
	// }

	// // note: torpedo weapons never have any bombard value (lighting in atmosphere is not a real issue)
	// EpsilonTorpedo := ValuesTable{
	// 	"ComponentCountermeasuresBonus":     ComponentCountermeasuresBonus,
	// 	"ComponentTargetingBonus":           func(level int) float64 { return 0 },
	// 	"WeaponBombardDamageInfrastructure": func(level int) float64 { return 0 },
	// 	"WeaponBombardDamageMilitary":       func(level int) float64 { return 0 },
	// 	"WeaponBombardDamagePopulation":     func(level int) float64 { return 0 },
	// 	"WeaponBombardDamageQuality":        func(level int) float64 { return 0 },
	// 	"WeaponShieldBypass":                func(level int) float64 { return -0.4 },
	// 	"WeaponArmorBypass":                 func(level int) float64 { return 0.4 },
	// 	"WeaponSpeed":                       TorpedoSeekingSpeed,
	// 	"WeaponRange":                       func(level int) float64 { return 800 + float64(level*100) },
	// 	"WeaponDamageFalloffRatio":          func(level int) float64 { return 0.25 },
	// 	"WeaponEnergyPerShot":               func(level int) float64 { return 0.75 * WeaponRawDamage(level) },
	// 	"WeaponFireRate":                    WeaponRateOfFire,
	// 	"WeaponRawDamage":                   WeaponRawDamage,
	// }

	// apply stats for each component
	// err = j.ApplyComponentAll(components)

	return
}
