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
type BayCounts map[ComponentType]int
type BayCountsPerLevels map[HullLevel]BayCounts
type HullBaySchedule map[string]BayCountsPerLevels

type BayTypeGroups map[ComponentType]BayTypeIndexes
type BayTypeIndexes struct {
	start int
	count int
}

var (
	// relationship of linear tier from assigned level (0 = invalid)
	fighterTier = map[HullLevel]int{0: 1, 2: 2, 4: 3, 9: 4, 13: 5, 15: 6}

	// required order of component bays types
	componentBayOrder        = []ComponentType{"Weapon", "Engine", "Defense", "General"}
	reverseComponentBayOrder = etc.Reverse(componentBayOrder)

	// what we wish it were
	fighterHullSchedule = HullBaySchedule{
		"FighterInterceptor": {
			// Fighter I
			0: {
				"Weapon":  1,
				"Engine":  1,
				"Defense": 1,
				"General": 2, // warn: their dumpster fire code can't handle a fighter with only 2 general slots (morons)
			},
			// Fighter II +1 engine
			2: {
				"Weapon":  1,
				"Engine":  2,
				"Defense": 1,
				"General": 2,
			},
			// Fighter III +1 weapon +1 defense
			4: {
				"Weapon":  2,
				"Engine":  2,
				"Defense": 2,
				"General": 3,
			},
			// Fighter IV +1 general
			9: {
				"Weapon":  2,
				"Engine":  2,
				"Defense": 2,
				"General": 4,
			},
			// Fighter V +1 engine
			13: {
				"Weapon":  2,
				"Engine":  3,
				"Defense": 2,
				"General": 4,
			},
			// Fighter X +1 defense
			15: {
				"Weapon":  2,
				"Engine":  3,
				"Defense": 3,
				"General": 4,
			},
		},
	}

	fighterComponentIdealSizes = BayCounts{
		"Weapon":  5,
		"Engine":  10,
		"Defense": 10,
		"General": 10,
	}
	bomberComponentIdealSizes = BayCounts{
		"Weapon":  10,
		"Engine":  10,
		"Defense": 10,
		"General": 10,
	}

	// these are simply based on tier
	ftrArmorReactiveRatings   = []int{2, 4, 6, 8, 12, 16}
	ftrIonDefenseRatings      = []int{4, 8, 12, 16, 20, 24}
	ftrCountermeasuresBonuses = []float64{.4, .5, .6, .7, .8, .9}
	ftrTargetingBonuses       = []float64{.4, .5, .6, .7, .8, .9}
	ftrManeuveringBonuses     = []float64{.08, .16, .24, .32, .40, .48}
)

// gives the space required to maximize all slots
func FighterComponentIdealSize(schedule BayCounts) (size int) {
	for k, v := range schedule {
		size += fighterComponentIdealSizes[k] * v
	}
	return
}

// gives the space required to maximize all slots
func BomberComponentIdealSize(schedule BayCounts) (size int) {
	for k, v := range schedule {
		size += bomberComponentIdealSizes[k] * v
	}
	return
}

func FighterHulls(folder string) (err error) {

	log.Println("All strikecraft component bay counts will be adjusted to match desired schedule")

	// load all ship hull definition files
	j, err := LoadJobFor(folder, "ShipHulls*.xml")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyComponentBays(fighterHullSchedule)
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
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
				schedule, ok := schedule[roleName]
				if ok {
					level := shiphull.Child("Level").IntValue()
					desiredCounts, ok := schedule[level]
					if ok {
						// set the slot counts
						err = j.applyComponentBaySchedule(shiphull, desiredCounts)
						if err != nil {
							return
						}

						// todo: we could have a schedule for Size and DisplaySize if we wish

						// set maximum size (90% of ideal)
						size := 9 * FighterComponentIdealSize(desiredCounts) / 10
						size += shiphull.Child("Size").IntValue()
						shiphull.Child("MaximumSize").SetValue(size)

						// set engine limit
						if count, ok := desiredCounts["Engine"]; ok {
							shiphull.Child("EngineLimit").SetValue(count)
						}

						// set some purely tier based numbers
						tier := fighterTier[level]
						shiphull.Child("ArmorReactiveRating").SetValue(ftrArmorReactiveRatings[tier])
						shiphull.Child("IonDefense").SetValue(ftrIonDefenseRatings[tier])
						shiphull.Child("CountermeasuresBonus").SetValue(ftrCountermeasuresBonuses[tier])
						shiphull.Child("TargetingBonus").SetValue(ftrTargetingBonuses[tier])

						// maneuvering bonsues
						if bonuses := shiphull.Child("Bonuses"); bonuses != nil {
							if bonus := bonuses.FindRecurse("Type", "ShipManeuvering"); bonus != nil {
								bonus.Child("Amount").SetValue(ftrManeuveringBonuses[tier])
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
		case "General":
			// warn: their dumpster fire code can't handle a fighter with only 2 general slots (morons)
			// warn: it turns out they prioritize an ion shield over reactors, so it fails the design due to no reactor
			// hack: so we enforce 3 physical slots, but no space allocation for it
			desired = max(3, desired)
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
	err = j.renumberComponentBays(shiphull.Child("ComponentBays"))

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
func (j *Job) renumberComponentBays(componentbays *xmltree.XMLElement) (err error) {

	// track our counts as we renumber
	counts := BayCounts{}

	// scan the component bays in order
	for i, c := range componentbays.Elements() {

		// id is trivial - just linear numbering within this list
		c.Child("ComponentBayId").SetValue(i)

		// mesh name is harder...
		// determine what type of component bay this is (we count mesh names independently)
		t := c.Child("Type").StringValue()

		// get the current count of this type (starts at zero)
		if m := c.Child("MeshName"); m != nil && meshRegex.MatchString(m.StringValue()) {
			// get the count of mesh names of this type so far
			ci := counts[t]

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
