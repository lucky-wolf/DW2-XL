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
				for _, e := range e.Child("Values").Elements() {

					// every element should be a component bay
					err = assertIs(e, "ComponentStats")
					if err != nil {
						return
					}

					// NOTE: do this before we scale the main guns so we're scaling off of the [S] value, independently of our standard output (below)
					e.Child("WeaponInterceptFireRate").SetValue(e.Child("WeaponFireRate").NumericValue() / 20) // when doing intercept duty (PD) we operate at a higher ROF

					// scale relative to [S]
					e.Child("WeaponEnergyPerShot").ScaleBy(0.5)
					e.Child("WeaponRawDamage").ScaleBy(0.5)
					e.Child("WeaponRange").ScaleBy(0.3)
					e.Child("WeaponDamageFalloffRatio").ScaleBy(1.5) // reduced range, more rapid fall-off
					e.Child("WeaponFireRate").ScaleBy(0.25)          // 4x rate of fire against ships compared to small weapons

					// all other intercept values are same scale as our standard output
					e.Child("WeaponInterceptRange").SetValue(e.Child("WeaponRange").StringValue())
					e.Child("WeaponInterceptDamageFighter").SetValue(e.Child("WeaponRawDamage").StringValue())
					e.Child("WeaponInterceptDamageSeeking").SetValue(e.Child("WeaponRawDamage").NumericValue() * 2)
					e.Child("WeaponInterceptEnergyPerShot").SetValue(e.Child("WeaponEnergyPerShot").StringValue())

					// never a crew requirement for fighter components
					e.Child("CrewRequirement").SetValue(0)
					e.Child("StaticEnergyUsed").SetValue(0)

					// fighter weapons generically get a +10% targeting across the board
					e.Child("ComponentTargetingBonus").AdjustValue(0.1)

					// fighters never do bombard damage
					for _, e := range e.Matching(regexp.MustCompile("WeaponBombard.*")) {
						e.SetValue(0)
					}
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
