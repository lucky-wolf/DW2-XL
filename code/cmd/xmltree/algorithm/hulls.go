package algorithm

import (
	"fmt"
	"log"
	"lucky-wolf/DW2-XL/code/cmd/etc"
	"lucky-wolf/DW2-XL/code/xmltree"
	"regexp"
)

// we'll put hull related support algos here

// name a bay category (general, weapon, etc.) and the count of bays it should contain
// this either deletes from the end of the list, or copies the last entry into more entries
// it does NOT adjust the size of bays - only copies the tail node to add more of that type
// SUBTLE: "Engine" bays are always an odd number by adding +1 bay if you ask for an even number
// subtle: this allows the game to always have a center position for a single engine / odd number of engines within the allowed count
type HullLevel = int
type ComponentType = string
type RoleName = string
type AttributeName = string
type BayCounts map[ComponentType]int
type BayCountsPerLevels map[HullLevel]BayCounts
type HullTiers map[HullLevel]int
type HullRoleDefinition struct {
	HullTiers          HullTiers
	ValuesTable        ValuesTable
	BayCountsPerLevels BayCountsPerLevels
}
type HullBaySchedule map[RoleName]HullRoleDefinition

type BayTypeGroups map[ComponentType]BayTypeIndexes
type BayTypeIndexes struct {
	start int
	count int
}

const (
	Human     = 0
	Ackdarian = 1
	Teekan    = 2
	Haakonish = 3
	Mortalen  = 4
	Ikkuro    = 5
	Boskara   = 6
	Zenox     = 7
	Wekkarus  = 8
	Atuuk     = 9
	Dhayut    = 10
	Gizurean  = 11
	Ketarov   = 12
	Kiadian   = 13
	Naxxilian = 14
	Quameno   = 15
	Securan   = 16
	Shandar   = 17
	Sluken    = 18
	Ugnari    = 19
	Shakturi  = 20
)

type RaceID = int
type RacialQuirk struct {
	bomberCenterIsTwo bool
}
type RacialQuirks map[RaceID]RacialQuirk

// because some races use 0,1,2 for left,right,center, but others use 1,2,0 for same, we have to do stupid shit here
func GetWeaponMeshIndex(shiphull *xmltree.XMLElement, index int, desiredCounts BayCounts) int {

	role := shiphull.Child("Role").StringValue()
	if role != "FighterBomber" {
		return index
	}

	race := shiphull.Child("RaceId").IntValue()
	zeroCenter := !racialQuirks[race].bomberCenterIsTwo
	if zeroCenter {
		if desiredCounts["Weapon"] == 1 {
			// only a single weapon means use the center
			return 0
		}
		switch index {
		case 0:
			return 1
		case 1:
			return 2
		case 2:
			return 0
		}
	} else if desiredCounts["Weapon"] == 1 {
		// only a single weapon means use the center mesh
		return 2
	}
	return index
}

var (
	// racial quirks (we might just want to supply translator functions for various mesh indexes)
	racialQuirks = RacialQuirks{
		Ackdarian: {bomberCenterIsTwo: true},
		Ikkuro:    {bomberCenterIsTwo: true},
		Teekan:    {bomberCenterIsTwo: true},
	}

	// required order of component bays types
	componentBayOrder        = []ComponentType{"Weapon", "Engine", "Defense", "General"}
	reverseComponentBayOrder = etc.Reverse(componentBayOrder)

	// what we wish it were
	fighterHullSchedule = HullBaySchedule{
		"FighterInterceptor": {
			HullTiers: HullTiers{0: 1, 2: 2, 4: 3, 9: 4, 13: 5, 15: 6},
			ValuesTable: ValuesTable{
				"ArmorReactiveRating":  func(tier int) float64 { return float64(2 * tier) },
				"IonDefenseRating":     func(tier int) float64 { return float64(4 * tier) },
				"CountermeasuresBonus": func(tier int) float64 { return []float64{.26, .32, .38, .44, .50, .56}[tier-1] },
				"TargetingBonus":       func(tier int) float64 { return []float64{.03, .06, .09, .12, .15, .18}[tier-1] },
				"ManeuveringBonus":     func(tier int) float64 { return []float64{.08, .16, .24, .32, .40, .48}[tier-1] },
			},
			BayCountsPerLevels: BayCountsPerLevels{
				// Fighter I
				0: {
					"Weapon":  1,
					"Engine":  1,
					"Defense": 1,
					"General": 2,
				},
				// Fighter II
				2: {
					"Weapon":  1,
					"Engine":  2,
					"Defense": 2,
					"General": 3,
				},
				// Fighter III
				4: {
					"Weapon":  2,
					"Engine":  2,
					"Defense": 2,
					"General": 3,
				},
				// Fighter IV
				9: {
					"Weapon":  2,
					"Engine":  3,
					"Defense": 2,
					"General": 4,
				},
				// Fighter V, Strike
				13: {
					"Weapon":  2,
					"Engine":  3,
					"Defense": 3,
					"General": 4,
				},
				// Fighter VI, Superiority
				15: {
					"Weapon":  2,
					"Engine":  3,
					"Defense": 3,
					"General": 5,
				},
			},
		},
		"FighterBomber": {
			HullTiers: HullTiers{0: 1, 1: 2, 2: 3, 4: 4, 5: 5, 6: 6},
			ValuesTable: ValuesTable{
				"ArmorReactiveRating":  func(tier int) float64 { return float64(2 * tier) },
				"IonDefenseRating":     func(tier int) float64 { return float64(4 * tier) },
				"CountermeasuresBonus": func(tier int) float64 { return []float64{.13, .16, .19, .22, .25, .28}[tier-1] },
				"TargetingBonus":       func(tier int) float64 { return []float64{.06, .09, .12, .15, .18, .21}[tier-1] },
				"ManeuveringBonus":     func(tier int) float64 { return []float64{.00, .08, .16, .24, .32, .40}[tier-1] },
			},
			BayCountsPerLevels: BayCountsPerLevels{
				// Bomber I
				0: {
					"Weapon":  1,
					"Engine":  1,
					"Defense": 1,
					"General": 2,
				},
				// Bomber II
				1: {
					"Weapon":  1,
					"Engine":  2,
					"Defense": 2,
					"General": 3,
				},
				// Bomber III
				2: {
					"Weapon":  2,
					"Engine":  2,
					"Defense": 2,
					"General": 3,
				},
				// Bomber IV
				4: {
					"Weapon":  2,
					"Engine":  2,
					"Defense": 3,
					"General": 4,
				},
				// Bomber V, Strike
				5: {
					"Weapon":  2,
					"Engine":  3,
					"Defense": 3,
					"General": 4,
				},
				// Bomber X
				6: {
					"Weapon":  3,
					"Engine":  3,
					"Defense": 3,
					"General": 5,
				},
			},
		},
	}
)

func Hulls(folder string) (err error) {

	log.Println("All strikecraft component bay counts will be adjusted to match desired schedule")

	// load all ship hull definition files
	j, err := LoadJobFor(folder, "ShipHulls*.xml")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyFighterHulls()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

func (j *Job) applyFighterHulls() (err error) {
	return j.applyComponentBays(fighterHullSchedule)
}

func (j *Job) applyComponentBays(schedule HullBaySchedule) (err error) {

	for _, f := range j.xfiles {

		// the root will result in a single ArrayOf[RootObjectType]
		for _, shiphulls := range f.root.Elements.Elements() {

			err = assertIs(shiphulls, "ArrayOfShipHull")
			if err != nil {
				return
			}

			for _, shiphull := range shiphulls.Elements() {

				// each of these is a ShipHull
				err = assertIs(shiphull, "ShipHull")
				if err != nil {
					return
				}

				// see whether we have a schedule for this role
				roleName := shiphull.Child("Role").StringValue()
				hullDefn, ok := schedule[roleName]
				if ok {
					level := shiphull.Child("Level").IntValue()
					desiredCounts, ok := hullDefn.BayCountsPerLevels[level]
					if ok {
						// set the slot counts
						err = j.applyComponentBaySchedule(shiphull, desiredCounts)
						if err != nil {
							return
						}

						// convert arbitrary engine level to more useful tier
						tier := hullDefn.HullTiers[level]

						// set maximum maxSize for all components combined
						componentSize := j.totalComponentBaySize(shiphull.Child("ComponentBays"))

						// set engine limit
						if count, ok := desiredCounts["Engine"]; ok {
							shiphull.Child("EngineLimit").SetValue(count)
							if count%2 == 0 {
								// hack: fucking crap engine!
								// remove extra engine size if it was auto-added
								ebay := shiphull.Child("ComponentBays").FindRecurse("Type", "Engine")
								componentSize -= ebay.Child("MaximumComponentSize").IntValue()
							}
						}

						maxSize := shiphull.Child("Size").IntValue() + etc.MulDivRoundUp(componentSize, 5, 8)
						shiphull.Child("MaximumSize").SetValue(maxSize)

						// todo: we could have a schedule for Size and DisplaySize if we wish

						// set some purely tier based numbers
						shiphull.SetChildValueIfExists("ArmorReactiveRating", hullDefn.ValuesTable["ArmorReactiveRating"](tier))
						shiphull.SetChildValueIfExists("IonDefense", hullDefn.ValuesTable["IonDefenseRating"](tier))
						shiphull.SetChildValueIfExists("CountermeasuresBonus", hullDefn.ValuesTable["CountermeasuresBonus"](tier))
						shiphull.SetChildValueIfExists("TargetingBonus", hullDefn.ValuesTable["TargetingBonus"](tier))

						// maneuvering bonsues
						if bonuses := shiphull.Child("Bonuses"); bonuses != nil {
							if bonus := bonuses.FindRecurse("Type", "ShipManeuvering"); bonus != nil {
								bonus.Child("Amount").SetValue(hullDefn.ValuesTable["ManeuveringBonus"](tier))
							}
						}
					}
				}
			}
		}
	}

	return
}

func (j *Job) applyComponentBaySchedule(shiphull *xmltree.XMLElement, desiredCounts BayCounts) (err error) {

	// scan the component bays by type to figure out where they each begin (and their counts)
	indexes, err := j.getComponentBayIndexes(shiphull.Child("ComponentBays"))
	if err != nil {
		return
	}

	// insert or remove elements from collection at the now known indexes (in reverse order)
	for _, componentBayType := range reverseComponentBayOrder {

		desired, ok := desiredCounts[componentBayType]
		if !ok {
			continue
		}
		actual := indexes[componentBayType].count

		switch componentBayType {

		case "Engine":
			if desiredCounts["Engine"]%2 == 0 {
				// note: always make engine bays odd to allow for the game to balance them
				// note: but we don't allocate space for it, and we set the EngineLimit according to the schedule
				desired += 1
			}
			// case "General":
			// 	// as of 1205 we should be able to add only 2 general slots for strike craft
			// 	desired = max(3, desired)
		}

		switch {
		case desired > actual:
			// append copies
			i := indexes[componentBayType].start + actual
			e := shiphull.Child("ComponentBays").Elements()[i-1]
			for c, d := 0, desired-actual; c < d; c++ {
				err = shiphull.Child("ComponentBays").InsertAt(i+c, e.Clone())
				if err != nil {
					return
				}
			}
		case actual > desired:
			// delete unneeded elements
			err = shiphull.Child("ComponentBays").RemoveSpan(indexes[componentBayType].start+desired, actual-desired)
			if err != nil {
				return
			}
		}
	}

	// renumber everything to ensure it's coherent
	err = j.renumberComponentBays(shiphull, desiredCounts)

	return
}

// renumbers all elements and returns the group counts (weapon, defense, etc.)
func (j *Job) getComponentBayIndexes(componentbays *xmltree.XMLElement) (counts BayTypeGroups, err error) {

	// scan the component bays in order, inserting or deleting as needed

	// initialize our map (golang will panic otherwise)
	counts = BayTypeGroups{}

	for i, c := range componentbays.Elements() {

		// figure out what type this is
		t := c.Child("Type").StringValue()

		// update our start index & count
		g := counts[t]
		if g.count == 0 {
			g.start = i
		}
		g.count++
		counts[t] = g
	}

	return
}

// matches #weapon0 and the like...
var meshRegex = regexp.MustCompile(`(#\w+)\d+`)

// renumbers all elements and returns the group counts (weapon, defense, etc.)
func (j *Job) renumberComponentBays(shiphull *xmltree.XMLElement, desiredCounts BayCounts) (err error) {

	// track our counts as we renumber
	counts := BayCounts{}

	// scan the component bays in order
	for i, c := range shiphull.Child("ComponentBays").Elements() {

		// id is trivial - just linear numbering within this list
		c.Child("ComponentBayId").SetValue(i)

		// mesh name is harder...
		// determine what type of component bay this is (we count mesh names independently)
		t := c.Child("Type").StringValue()

		// get the current count of this type (starts at zero)
		if m := c.Child("MeshName"); m != nil && meshRegex.MatchString(m.StringValue()) {

			ci := counts[t]

			switch t {
			case "Weapon":
				// get the translated mesh index for this weapon for this race
				ci = GetWeaponMeshIndex(shiphull, ci, desiredCounts)
			}

			oldname := m.StringValue()
			newname := meshRegex.ReplaceAllString(oldname, fmt.Sprintf("${1}%d", ci))

			// update this one to match our count
			m.SetString(newname)

			// update our count for the next one
			counts[t] += 1
		}
	}

	return
}

// adds up all the bay sizes
func (j *Job) totalComponentBaySize(componentbays *xmltree.XMLElement) (size int) {

	// scan the component bays in order
	for _, c := range componentbays.Elements() {
		size += c.Child("MaximumComponentSize").IntValue()
	}

	return
}
