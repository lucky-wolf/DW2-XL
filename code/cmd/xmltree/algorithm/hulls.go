package algorithm

import (
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
)

// we'll put hull related support algos here

// name a bay category (general, weapon, etc.) and the count of bays it should contain
// this either deletes from the end of the list, or copies the last entry into more entries
// it does NOT adjust the size of bays - only copies the tail node to add more of that type
// SUBTLE: "Engine" bays are always an odd number by adding +1 bay if you ask for an even number
// subtle: this allows the game to always have a center position for a single engine / odd number of engines within the allowed count
type ComponentType string
type BayCounts map[ComponentType]int
type BayCountsPerLevel struct {
	Tier int
	Bays BayCounts
}
type HullLevel = int
type BayCountsPerLevels map[HullLevel]BayCountsPerLevel

var (
	componentBaySchedules = map[string]BayCountsPerLevels{
		"FighterInterceptor": {
			0: {
				Tier: 0,
				Bays: BayCounts{
					"Weapon":  1,
					"Engine":  1,
					"Defense": 1,
					"General": 2,
				},
			},
			// +2 bays
			2: {
				Tier: 1,
				Bays: BayCounts{
					"Weapon":  2,
					"Engine":  2,
					"Defense": 1,
					"General": 2,
				},
			},
			// +2 bays
			4: {
				Tier: 2,
				Bays: BayCounts{
					"Weapon":  2,
					"Engine":  2,
					"Defense": 2,
					"General": 3,
				},
			},
			// +1 bay
			9: {
				Tier: 3,
				Bays: BayCounts{
					"Weapon":  2,
					"Engine":  3,
					"Defense": 2,
					"General": 3,
				},
			},
			// +2 bay
			13: {
				Tier: 4,
				Bays: BayCounts{
					"Weapon":  2,
					"Engine":  3,
					"Defense": 3,
					"General": 4,
				},
			},
			// +1 bay
			15: {
				Tier: 5,
				Bays: BayCounts{
					"Weapon":  3,
					"Engine":  3,
					"Defense": 3,
					"General": 4,
				},
			},
		},
	}
)

// we simply go through all of the component bays and number them sequentially
func RenumberComponentBays() {
}

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
		for _, e := range f.root.Elements.Elements() {

			err = assertIs(e, "ArrayOfShipHull")
			if err != nil {
				return
			}

			for _, e := range e.Elements() {

				// each of these is a ShipHull
				err = assertIs(e, "ShipHull")
				if err != nil {
					return
				}

				// see whether we have a schedule for this role
				roleName := e.Child("Role").StringValue()
				schedule, ok := componentBaySchedules[roleName]
				if ok {
					level := e.Child("Level").IntValue()
					bays, ok := schedule[level]
					if ok {
						err = j.applyComponentBaySchedule(e, bays)
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

func (j *Job) applyComponentBaySchedule(e *xmltree.XMLElement, bays BayCountsPerLevel) (err error) {

	// scan the component bays in order, inserting or deleting as needed

	// e is a ShipHull
	// b is a ComponentBays (list of ComponentBay)
	// c is a ComponentBay
	b := e.Child("ComponentBays")
	for _, c := range b.Elements() {
		// we have a funny situation...
		// this is an array of elements which are grouped by type
		// weapon, engine, defense, general (for strike craft)
		// we need to adjust the count of each group to conform to the caller's specifications
		// contraction is simply knowing which index range to delete
		// extension is knowing which index range to replicate (last element of group replicated N times after itself)
		// extension also requires fixing any <MeshName>#weapon0</MeshName> to use next available index (e.g. #weapon1)
		// both extension and contraction require that we fix up the <ComponentBayId> indexes as well to be fully sequential for total elements
		// this last step is probably easiest and more flex in another pass / separate algo -- which could be handy for hand-edits of the xml (insertion / deletion of components)
		c.Child("Type").StringValue()

	}

	// see if this hull is a level match for any
	return
}

type BayTypeGroups map[string]BayTypeIndexes

type BayTypeIndexes struct {
	start int
	count int
}

// renumbers all elements and returns the group counts (weapon, defense, etc.)
func (j *Job) renumberComponentBays(e *xmltree.XMLElement) (counts BayTypeGroups, err error) {

	// scan the component bays in order, inserting or deleting as needed

	// e is a ShipHull
	// b is a ComponentBays (list of ComponentBay)
	// c is a ComponentBay
	b := e.Child("ComponentBays")

	// initialize our map (golang will panic otherwise)
	counts = BayTypeGroups{}

	for i, c := range b.Elements() {
		// we have a funny situation...
		// this is an array of elements which are grouped by type
		// weapon, engine, defense, general (for strike craft)
		// we need to adjust the count of each group to conform to the caller's specifications
		// contraction is simply knowing which index range to delete
		// extension is knowing which index range to replicate (last element of group replicated N times after itself)
		// extension also requires fixing any <MeshName>#weapon0</MeshName> to use next available index (e.g. #weapon1)
		// both extension and contraction require that we fix up the <ComponentBayId> indexes as well to be fully sequential for total elements
		// this last step is probably easiest and more flex in another pass / separate algo -- which could be handy for hand-edits of the xml (insertion / deletion of components)

		// id is trivial - just linear numbering within this list
		c.Child("ComponentBayId").SetValue(i)

		// mesh name is harder...
		// first, figure out what type this is
		t := c.Child("Type").StringValue()

		// update our start index & count
		g := counts[t]
		if g.count == 0 {
			g.start = i
		}
		g.count++
		counts[t] = g

		if m := c.Child("MeshName"); m != nil {
			m.StringValue()

		}
	}

	// see if this hull is a level match for any
	return
}
