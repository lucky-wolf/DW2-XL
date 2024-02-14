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

func (j *Job) applyTorpedoWeapons() (err error) {

	// // note: medium+ torpedo weapons do have bombard values
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
