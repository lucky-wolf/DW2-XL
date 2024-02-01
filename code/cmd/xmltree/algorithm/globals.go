package algorithm

// we have some standard global functions defined here
// to allow ourselves to build specific schedules with
// references to common global schedules
//
// the aim is to allow some global basis for all of our data
// then we can make things relative to each other
// tweak some global base values, and you rearrange the entire web of dependencies

const IonFtrPDScaleFactor = 0.75

// standard weapon countermeasure schedule (by tech level)
func ComponentCountermeasuresBonus(level int) float64 {
	return 0.6 + float64(level)*0.02
}

// arbitrary
func BlasterWeaponRateOfFire(level int) float64 {
	return 9
}

// ion is 50% slower than blasters
func IonWeaponRateOfFire(level int) float64 {
	return 1.5 * BlasterWeaponRateOfFire(level)
}

// 12, 12, 24, 36, 48, 60, 72, 84, 96, 108, 120
func IonWeaponComponentDamage(level int) float64 {
	// treat level 0 and 1 as the same for our purposes (currently)
	switch level {
	case 0:
		level = 1
	}
	return float64(level) * 12
}

// 0,  8, 16, 24, 32
func IonShieldIonDamageDefense(level int) float64 {
	return 8 * float64(level)
}

// standard component Ion defense
func StandardComponentIonDefense(level int) float64 {
	return float64(level) * 2
}

// hardened component Ion defense
func HardenedComponentIonDefense(level int) float64 {
	return float64(level) * 4
}

// standard damage is based on pulsed blasters-ish, but at about 2/3 the ROF, so 2/3 the DPS
func IonWeaponRawDamage(level int) float64 {
	// this gives us 20 at (t0) and a gain of 20% per level beyond that
	return 20 * (1 + 0.2*float64(level-1))
}

// 300, 325, 350, 375, 400, 425, 450, 475, 500, 525, 550
func TorpedoSeekingSpeed(level int) float64 {
	return 300 + 25*float64(level)
}

func TorpedoSeekingRange(level int) float64 {
	return 10 * TorpedoSeekingSpeed(level)
}

// 300, 325, 350, 375, 400, 425, 450, 475, 500, 525, 550
func MissileSeekingSpeed(level int) float64 {
	return 300 + 25*float64(level)
}

// 1k, 2k, 3k, ... 11k
func MissileSeekingRange(level int) float64 {
	return 1000 * float64(level+1)
}
