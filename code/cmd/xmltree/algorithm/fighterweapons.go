package algorithm

import (
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
	"regexp"
	"strings"
)

func FighterWeapons(folder string) (err error) {

	// load all component definition files
	j, err := loadJobFor(folder, "ComponentDefinitions*")
	if err != nil {
		return
	}

	// apply this transformation
	err = j.applyFighterWeapons()
	if err != nil {
		return
	}

	// save them all
	j.save()

	return
}

func (j *job) applyFighterWeapons() (err error) {

	for _, f := range j.xfiles {

		statistics := &f.stats

		// the root will result in a single ArrayOf[RootObjectType]
		for _, e := range f.root.Elements.Elements() {

			err = assertIs(e, "ArrayOfComponentDefinition")
			if err != nil {
				return
			}

			for _, e := range e.Elements() {

				// each of these is a componentdefinition
				err = assertIs(e, "ComponentDefinition")
				if err != nil {
					return
				}

				// only weapons...
				if !e.HasPrefix("Category", "Weapon") {
					continue
				}

				// for fighters...
				targetName := e.Child("Name").StringValue()
				if !strings.HasSuffix(targetName, "[Ftr]") {
					continue
				}

				// find the corresponding small weapon by name
				sourceName := strings.Replace(targetName, "[Ftr]", "[S]", 1)
				sourceDefinition, _ := j.find("Name", sourceName)
				if sourceDefinition == nil {
					log.Printf("Missing: %s (for %s)", sourceName, targetName)
					continue
				}

				// debug
				if !Quiet {
					log.Printf("%s from %s\n", targetName, sourceName)
				}

				// copy and scale resource requirements
				err = e.CopyAndVisitByTag("ResourcesRequired", sourceDefinition, func(e *xmltree.XMLElement) error { e.Child("Amount").ScaleBy(0.25); return nil })
				if err != nil {
					log.Printf("%s: %s from %s", err, targetName, sourceName)
				}

				// copy component stats
				err = e.CopyByTag("Values", sourceDefinition)
				if err != nil {
					log.Printf("%s: %s from %s", err, targetName, sourceName)
				}

				// now that we have our own copy of the component stats (same number of levels too)
				// we can update each of those to scale for [Ftr] version
				err = scaleFighterOrPDValues(e)
				if err != nil {
					log.Printf("%s: %s from %s", err, targetName, sourceName)
				}

				statistics.changed++
				statistics.elements++
				statistics.objects++
			}
		}
	}
	err = nil
	return
}

// scale all component stats from the copy from our [S] source weapon
func scaleFighterOrPDValues(e *xmltree.XMLElement) (err error) {

	for _, e := range e.Child("Values").Elements() {

		// every element should be a component bay
		err = assertIs(e, "ComponentStats")
		if err != nil {
			return
		}

		isFighter := e.Has("IsFighterOnly", "true")

		// "flatten" source volleys to 1 per shot but at 1/x fire rate (same dps, but distributed instead of burste firing)
		if va := e.Child("WeaponVolleyAmount").NumericValue(); va != 1 {
			e.Child("WeaponFireRate").ScaleBy(1.0 / va)
			e.Child("WeaponVolleyAmount").SetValue(1)
		}
		e.Child("WeaponVolleyFireRate").SetValue(0)

		// scale standard fire relative to [S] identically to FighterWeapons
		e.Child("WeaponFireRate").ScaleBy(0.25) // 4x rate of fire against ships compared to small weapons
		e.Child("WeaponRawDamage").ScaleBy(0.25)
		e.Child("WeaponEnergyPerShot").ScaleBy(0.25)
		e.Child("WeaponRange").ScaleBy(0.3)
		e.Child("WeaponDamageFalloffRatio").ScaleBy(1.5) // reduced range, more rapid fall-off

		// all other intercept values are same scale as our standard output
		e.Child("WeaponInterceptFireRate").SetValue(e.Child("WeaponFireRate").NumericValue() / 4)           // 4x standard action (which is currently 4x source gun rof)
		e.Child("WeaponInterceptDamageFighter").SetValue(e.Child("WeaponRawDamage").NumericValue() / 2)     // x4/2 = x2 effective dps vs. fighters
		e.Child("WeaponInterceptDamageSeeking").SetValue(e.Child("WeaponRawDamage").NumericValue() * 1)     // x4/1 = x4 effective dps vs. seeking ordinance
		e.Child("WeaponInterceptEnergyPerShot").SetValue(e.Child("WeaponEnergyPerShot").NumericValue() / 2) // x4/2 = x2 energy cost during intercept mode
		e.Child("WeaponInterceptRange").SetValue(e.Child("WeaponRange").StringValue())

		// fighter & PD weapons generically get a +10% targeting across the board
		e.Child("ComponentTargetingBonus").AdjustValue(0.1)

		// fighters and PD never do bombard damage
		for _, e := range e.Matching(regexp.MustCompile("WeaponBombard.*")) {
			e.SetValue(0)
		}

		// some things just don't apply to fighters (but do to PD)
		if isFighter {
			e.Child("CrewRequirement").SetValue(0)
			e.Child("StaticEnergyUsed").SetValue(0)
		}
	}

	return
}
