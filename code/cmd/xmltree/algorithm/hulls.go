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
	componentBaySchedules = map[string]BayCountsPerLevels{
		"FighterInterceptor": {
			0: {
				"Weapon":  1,
				"Engine":  1,
				"Defense": 1,
				"General": 2,
			},
			// +2 bays
			2: {
				"Weapon":  2,
				"Engine":  2,
				"Defense": 1,
				"General": 2,
			},
			// +2 bays
			4: {
				"Weapon":  2,
				"Engine":  2,
				"Defense": 2,
				"General": 3,
			},
			// +1 bay
			9: {
				"Weapon":  2,
				"Engine":  3,
				"Defense": 2,
				"General": 3,
			},
			// +2 bay
			13: {
				"Weapon":  2,
				"Engine":  3,
				"Defense": 3,
				"General": 4,
			},
			// +1 bay
			15: {
				"Weapon":  3,
				"Engine":  3,
				"Defense": 3,
				"General": 4,
			},
		},
	}
)

func AdjustComponentBays(folder string) (err error) {

	log.Println("All strikecraft component bay counts will be adjusted to match desired schedule")

	// load all ship hull definition files
	j, err := LoadJobFor(folder, "ShipHulls*.xml")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyComponentBays()
	if err != nil {
		return
	}

	// save them all
	j.Save()

	return
}

func (j *Job) applyComponentBays() (err error) {

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
				schedule, ok := componentBaySchedules[roleName]
				if ok {
					level := shiphull.Child("Level").IntValue()
					bays, ok := schedule[level]
					if ok {
						err = j.applyComponentBaySchedule(shiphull, bays)
						if err != nil {
							return
						}
					}
				}
			}
		}
	}

	return
}

func (j *Job) applyComponentBaySchedule(shiphull *xmltree.XMLElement, desiredCounts BayCounts) (err error) {

	componentbays := shiphull.Child("ComponentBays")

	// scan the component bays by type to figure out where they each begin (and their counts)
	indexes, err := j.getComponentBayIndexes(componentbays)
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

		switch {
		case desired > actual:
			// append copies
			err = componentbays.ExtendAt(indexes[componentBayType].start+actual-1, desired-actual)
			if err != nil {
				return
			}
		case actual > desired:
			// delete unneeded elements
			err = componentbays.RemoveSpan(indexes[componentBayType].start+desired, actual-desired)
			if err != nil {
				return
			}
		}
	}

	// renumber everything to ensure it's coherent
	err = j.renumberComponentBays(componentbays)

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
var meshRegex = regexp.MustCompile("(#\\w+)\\d+")

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
			i := counts[t]

			// update this one to match our count
			m.SetString(meshRegex.ReplaceAllString(m.StringValue(), fmt.Sprintf("$1%d", i)))

			// update our count for the next one
			counts[t] += 1
		}
	}

	// see if this hull is a level match for any
	return
}
