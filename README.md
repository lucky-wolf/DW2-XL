# DW2-XL
Distant Worlds 2 - XL

Author: Mordachai (lucky-wolf)

- [DW2-XL](#dw2-xl)
	- [Guiding Principles](#guiding-principles)
	- [Mod Highlights](#mod-highlights)
	- [Latest Changes](#latest-changes)
		- [v1.7.0](#v170)
		- [v1.6.3](#v163)
		- [v1.6.2](#v162)
		- [v1.6.1](#v161)
		- [v1.6.0](#v160)
		- [v1.5.0](#v150)
		- [v1.4.4](#v144)
		- [v1.4.3](#v143)
		- [v1.4.2](#v142)
		- [v1.4.1](#v141)
		- [v1.4.0](#v140)
		- [v1.3.7](#v137)
		- [v1.3.6](#v136)
		- [v1.3.5](#v135)
		- [v1.3.4](#v134)
		- [v1.3.3](#v133)
		- [v1.3.2](#v132)
		- [v1.3.1](#v131)
		- [v1.3.0](#v130)
		- [v1.2.3](#v123)
		- [v1.2.2](#v122)
		- [v1.2.1](#v121)
		- [v1.2.0](#v120)
		- [v1.1.10](#v1110)
		- [v1.1.9](#v119)
		- [v1.1.8](#v118)
		- [v1.1.7](#v117)
		- [v1.1.6](#v116)
		- [v1.1.5](#v115)
		- [v1.1.4](#v114)
	- [Research (Tech Tree)](#research-tech-tree)
	- [Ship Components](#ship-components)
	- [Fleets](#fleets)
	- [Ships](#ships)
	- [Crew & Star Barracks](#crew--star-barracks)
	- [Sensors](#sensors)
	- [Armor](#armor)
	- [Weapons](#weapons)
	- [Kinetic Weapons](#kinetic-weapons)
	- [Energy Torpedoes](#energy-torpedoes)
	- [Fighters and Bombers](#fighters-and-bombers)
	- [Hyper Drives](#hyper-drives)
	- [Reactors](#reactors)
	- [Targeting and Countermeasures](#targeting-and-countermeasures)
	- [Colonization](#colonization)
	- [Planetary Facilities](#planetary-facilities)
	- [Bug Fixes: Base Game](#bug-fixes-base-game)
	- [Bug Fixes: Early versions of this mod](#bug-fixes-early-versions-of-this-mod)
	- [DW2-XL Hull Sizes](#dw2-xl-hull-sizes)

## Guiding Principles
Principally this mod aims to create a better player experience while playing games of Distant Worlds 2.
It is not a completely different game than vanilla, rather it's "just better" in every way I had time or insight to make something better than it was. YMMV, but this is my take on making it better.

## Mod Highlights
- Colonization tech is pick & choose (you're not required to become an expert in every planetary biome, you can specialize in those that make sense for the game you're playing as it evolves)
- Specialized tech is viable throughout the game
  - This mod makes sure that you don't have to also research the common techs in order to access some of the later features or facilities, rather they're available to everyone
- More specializations to choose from
  - Engine technology lines are more heavily specialized - e.g. Acceleros have poor maneuverability, but best-in-class acceleration; whereas Proton engines are all about maneuverability
  - Warp technology is now more strongly differentiated
    - Kaldos are super quick to initiate and fast point-to-point, but have very limited range, requiring multiple zig-zagging hops to get anywhere distant.
    - Gerax are a steady climber in all categories - gaining speed, range, fuel efficiency, and quickness steadily throughout their entire lineage.
    - Equinox are power-hungry version of Gerax with similar range and slightly worse efficiency, gaining over time, but at increased power and fuel consumption.
    - Calista-Dal are "long-haul" drives - slower, but much more efficient than any of the others, and able to sustain their warp-bubbles for greater distance than any of the others, making them ideal for long-haul operations or non-combat ships (albeit they're an interesting choice and very playable for your warships depending on your galaxy, and how open or obstructed it is with nebulas)
- Tech lines that bog down the game with useless non-choices are removed or simplified or deepened to have an appreciable impact on the game
  - Vectoring thrusters only slowed down the game and saddled it with non-choices and confusion for the AIs to get stuck on, so I removed them (hid them, they're still in game, so as not to break other mods or game-events).
  - The couple of specialized interceptors for the Arkadians are now hidden, and only unlock if you get the in-game story events to do so.  Since they end very quickly, they're a dead-end tech that is otherwise not helpful to the flow of the game.
- Basic tech is cheap & basic
  - The first two levels of tech are pretty much applicable to all games, all species.
  - They are cheaper than in vanilla.
  - Early warp techs are more viable, allowing limited early game expansion or small local empires to be founded.
- Game-changing technologies were pushed back to mostly start at the 3rd level of the tech tree
  - The game currently requires any nation to be gifted at least the first two levels of tech when starting a game with more than one initial colony...
  - By pushing back much of the more interesting decisions and more powerful technologies, those gifted free techs remains just what everyone needs, while leaving character building techs for later.
    - This makes game-starts with more than one colony much more viable and enjoyable to play, in my opinion
- There aren't any one-per-game facilities, rather, they're all one-per-empire.
  - This helps to reduce the winner run-away effect, snatching up all "wonders" and leaving everyone else permanently in the dust
  - Hopefully, this will make AIs more viable (once their code is improved in general in the game engine itself)
    - Or if you're coming from behind, should allow you to catch the leaders if you can keep yourself alive to reach and build these facilities yourself.

## Latest Changes

### v1.7.0
- Fixed ship hulls, names, sizes & related tech names to be as intended (was broken by 1.0.8.x builds)
  - Fixed ship hull sizes for all player-races to: [DW2-XL Hull Sizes](#dw2-xl-hull-sizes)
  - Fixed interceptors to always use size 5 weapons, bombers are size 10 (will force the ship designer to actually use both slots, and to define the roles of the two more honestly)
  - Added the .shmd files for all of this to force the game to actually respect these values!
- Tweaked tech tree to hopefully make the mid-game flow a bit better

### v1.6.3
- Reduced "efficient" thrusters to size 16, giving them some possible utility
- Added Hail Cannons \[M\] for Humans
  - Tweaked Hail Cannons and Terminator Cannons to be more sensible, and to give the Humans a slight advantage in their racial tech
- Fixed missing diplomatic specialization center (Universal Peace Center)
- Gave wormhole drive better range (1000M)
- Rebalanced all \[Ftr\] weapons to make better sense in the current tech tree across all weapon-lines and levels
  - Many weapons had their ROF reduced and their dmg/shot increased for a lighter load on the CPU, but effectively same damage output
  - Which also reduces the effectiveness of reflective armor and shields, since more dmg will penetrate, make fighters much more powerful (they were basically gnats to be ignored otherwise)
- Reduced the reflective value of Ion Armors, so that there is more of a trade off between standard armor being both tougher (25%), and better at reflecting damage (2x) vs. ion armors being able to replace ion defense generators (considerable space-savings)

### v1.6.2
- Fixed Talassos Shields v3 and v4 were crazy expensive
- Tweaked all shield recharge rates, strength progression
- Tweaked all reactor power levels and progressions

### v1.6.1
- Zenox can choose technocracy now that the bug with giving them a +20% all research on hw is fixed
- Integrated 1.0.8.3 Beta Changes from Matrix Games
  - This includes a very large number of minor component tweaks since this also includes 1.0.8.0 component changes

### v1.6.0
- Tons of improvements, new techs, new components, new facilities
- Fixed bug in Zenox that was giving them a +20 all research at their homeworld (this is also fixed in the next beta from Matrix Games)
- Revamped how diplomatic facilities are obtained
  - added another level to diplomatic facilities
  - you now need multiple diplomatic techs of a certain level to obtain access to diplomatic focus and diplomatic specialization
  - moved espionage facility to its own tech in the diplomatic section
  - it requires a certain number of basic diplomacy with other major powers
  - improved Espionage Facility to improve espionage (go figure!)
- Expanded Zenox star beams
  - they start at tier 2, like other racial weapons in this mod
  - they bifurcate into a fighter-small line, and a medium-large line (also like others in this mod)
  - there is more tier for all star beams at the high end, now
- Rebalanced the shield tech tree
  - shield tech branches at tier 2, so player decides right after deflectors which path to take
  - all shields have been rebalanced against each other
  - advanced shields are available to all research paths at tier 6
- Reorganized much of the weapons, shields and armor techs so they're all together at the start of the tech tree
  - in general, it goes direct fire to seeking and area weapons
- Bombardment weapons now have the designation \[B\] and are of medium slot size (same as a medium seeking weapon)
- Fixed interceptors to definitely always have light weapons
- Improved the base troop templates to use a much more balanced selection of troops
- Updated all planetary facilities to cost 1% operational maintenance
  - this doubles the maintenance cost of many of them
- Increased abhorrence of Boskarans by other races
- Increased diplomacy malus against hive mind by other governments
- Added restrictions so that Boskarans cannot specialize in diplomacy
- Removed all research bonsues from all races, and replaced them with more specific bonuses for that race
- Mildly nerfed technocracy to have only a 15% bonus to all tech (was 20%)
- Fixed all techs so that no groups are split up (all techs of a category are together, in a single block)
- Updated policies to default to capturing all bases instead of only those in your territory

### v1.5.0
- rebalanced the weapons tech tree
  - size \[X\] weapons all are size 150
  - added more tiers for some weapons & point defenses
  - tons of work here filling in & filling out & balancing
- moved up planetary defenses and some other planetary facilities to tier 3
- updated game events so that vector engines are never generated through this path

### v1.4.4
- Increased minimum deposits to at least 10%
- Fixed size of small Velocity Torps
- Updated races.xml for 1.0.7.9

### v1.4.3
- Fighter-interceptors use small weapons for all mounts
- Re-normalized weapons into 11/22/44 for direct fire weapons, and 13/26/52 for seeking, bombing, and AOE weapons
- Stealth is now size 30 (was 50)
- Tweaked Construction research tree
- Moved up research, research buildings, and terraforming
- Moderated the terraforming facilities
- Moved up terraforming to be a bit earlier in the tech tree
- Reduced tech cost overall

### v1.4.2
- Reduced size of Mortalen's Multi-Targeting from 10 -> 5 (same as other targeting computers)
  - Added them as alternative prerequisite technologies for advanced short range sensors

### v1.4.1
- Extended short range sensors +2 more levels (crystal also)
- Made the tech tree a little more organized by block

### v1.4.0
- Updated for DW2 1.7.1 Beta
- Random tech tree now works
- Significant changes to Sensors
  - Fleet targeting & Fleet countermeasures are now size 10 (was 15)
  - Short Range Sensors are now size 10 (was 5)
  - Long Range Sensors are now size 20 (was 40)
  - Deep Sensor Arrays remain size 40
  - Sensors require the targeting tech that enables that level of targeting bonus (so you need the targeting tech first)
  - There are +2 new, higher level short range sensors (basic, standard, improved, enhanced, advanced)
  - Sensors no longer give you more targeting bonus than the equivalent targeting tech, and they take up 2x ship space (so they continue to fulfil their purpose, are relatively small, but the actual targeting tech has more of a niche to fill, and you can't ignore that in favor of pure sensor tech)
  - Long Range Sensors are one level deeper in the tech tree (and are no longer a requirement for stealth)
  - Fleet Targeting and Countermeasures now start at level 2, instead of 3

### v1.3.7
- Added Heavy Rail Gun [M] v2
  - This gives parity for heavy and light rail guns so that upgraded designs are sensible
- Updated to 1.0.6.9 Beta
  -  Added dual tech requirements to support newest beta engine (no longer defaults to "all required")

### v1.3.6
- Swapped Gerax & Equinox sizes (22<->17) (Equinox is now the larger of the two)
- Increased Equinox speeds from 650/975/1950/3800 -> 700/1050/2100/4200

### v1.3.5
- Tiny tweaks to Automation Policies

### v1.3.4
- Tweaked fleet templates
  - Added +1/+1 capital ship and carrier to attack template
  - Removed capital ship from raid template
  - Removed all tankers from fleets (engine just cannot use them sensibly still)
- Adjusted AIs
  - Increased constructor production to match explorer (high) to help them grow their economies

### v1.3.3
- Fixed bug in Boskaran's Firestorm v2 tech - was costing 4x too much!
- Tweaked the default fleet templates
  - Added back in some fuel tankers (now that they're arguably marginally less buggy)
  - Gave invasion fleets a Carrier, and Raid fleets a Battleship
  - Switch the defaults so that fleet policies override ship policies
  - Changed the invade to "immediately" to try to avoid situation where they land sporadically (this is fundamentally still a bug in the game itself)
  - Changed fleets to be more willing to retreat - hopefully to give the AI a better ability to trade blows over time and not disintegrate at the first clash
- Increased visibility of nebulas (fixed 1.0.6.6's settings so that they actually work)
- Planetary Facilities
  - Buffed shield facilities
  - Buffed fighter bases
- Adjusted all AIs to try to
  - Better balance their military construction for what they actually need
  - Put more resources into economy (exploration & construction)
  - Avoid building too many fixed defenses (which are a poor investment)
  - Focus on a single tech at a time
  - Grow to worlds that are a little bit borderline, but which they should be able to terraform or tech into good worlds over time

### v1.3.2
- Adjusted Boskarans
  - Adjusted tech level of their specialized torpedo weapons to be on the same cadence as their plasma beams
    - thus ship designs are in parity (latest is latest of each)
  - Added ion defense to flux armor
    - not as good as true Ionic Armor, but better than vanilla
- Buffed planetary defense facilities
  - Doubled planetary weapons volley size (e.g. 3 torpedoes -> 6 torpedoes per shot)
  - +160, +180, +200 targeting and countermeasures (they will hit, and are extremely hard to shoot down)
  - Added a secondary tech that gives planetary beam defense weapons (from pulsed weapons)
- Adjusted multi-role ships to come one tech level sooner

### v1.3.1
- Adjusted AI policies to allocate more cashflow to troops & facilities

### v1.3.0
- Expanded the reactors line to 6-deep for all 3 main techs (fusion, fission, quantum)
  - fusion is the most compact but the lowest output and capacity
  - fission is 50% larger for 50% more output and 20% more capacity
  - quantum is 100% larger than fusion for 100% more output and 80% greater capacity
  - fusion is the most fuel efficient, fission in the middle, and quantum is the worst
- Added AI policies to help them be more competitive

### v1.2.3
- Applied 1.0.6.4 fixes to TroopDefinitions.xml

### v1.2.2
- Fixed vanilla bug in phasers: weapon speed 1500
- Relaxed the tech lines of ionic armor and star marine barracks to be linear, not requiring the normal armor nor the normal crew quarters, but being their own line of research
  - Makes upgrades less of a pain in the butt (and they already cost 2x research vs. their normal variants for the same level of base armor or crew capacity)
- Added one more level of ion sheath armor (equivalent to impervious)
- Upgraded super armor to have ion sheath properties

### v1.2.1
- Normalized all hyperdrives to increase speed of hyperdrives in general (starting above base Gerax)
  - Gerax/Torrent gives incremental improvement in all areas for higher research cost overall
  - Kaldos offers the quickest init & recharge times plus best accuracy, but with limited range and 2nd worst efficiency
  - Equinox offers the highest speeds, but with the worst efficiency and poor accuracy
  - Calista-Dal offers the longest ranged jumps combined with highly efficient engines and good accuracy, but with somewhat reduced speed and poorest initialize and recharge times
- Super Drives (Flux) now have respectable vector thrust
- +5 Energy to Ftr Reactors (were a bit anemic)

### v1.2.0
- Standardized static energy on colony modules to 1E/1M capacity (e.g. 50M capacity = 50 static energy)
- Improved sizes to be what they were with double modules (50,100,150,200)
- Luxury Barracks give a maintenance savings of 10/15/20/25 at various levels
- Marine Barracks no longer give a maintenance savings (they give boarding defense bonus)
- Applied fixes from DW2 build 1.0.5.8

### v1.1.10
- Incorporated minor race planetary habitability fixes from DW2 build 1.0.5.6

### v1.1.9
- Fixed an issue with Starship Maneuvering where getting that tech through an event would crash the game (mod bug)

### v1.1.8
- Normalized many of the planetary defense facility techs so they become available around the same tech-level
- Updated engines to have meaningful vector thrust values (70-375)

### v1.1.7
- Increased fighter engine output
  - base by 1.25
  - max by 1.5 x base
  - should allow fighters to advance into combat ahead of main fleet & hurry back to carriers when leaving
- Further organized construction, parts of shields, ...
- Tweaked Quick Jump hyperdrives again (a bit slower for less insane energy requirements)

### v1.1.6
- Applied DW2 1.0.5.1 Fixes
  - Kaidian spelling
  - Skip drive sound effects
- Reorganized Troops
  - All troop related tech and facilities is right there under troops
- Reorganized Boarding Pods
  - All boarding related tech is under boarding pots

### v1.1.5
- Revamped engines tech tree & components
  - three distinct choices: efficient, nimble, and maximum speed
  - increased energy consumption of all engines so that this actually makes a difference
  - improved vectoring of nimble, advanced, and Ackdarian engines
  - engine choice is not level 2 (was 3) - makes a nicer tree & decision point
  - extended "ion thrusters" one level to make this work nicely
- Fixed some minor inconsistencies in hyper drive deny blocking
- Improved the hyper deny research tree
- Removed "Robo-" from Zenox troop names (now they're are Ice whatever) since they're not robotic, so it made no sense!

### v1.1.4
- Tweaked Long Haul hyper drives to make them a little more attractive
- Moved ion defenses to appear near shields && reorganized point deflectors (right next to them)
  - both are rooted in advanced deflectors rather than in weapons techs
- Separated scanners and jammers into their own tech-lines (independent of each other)
- Further tweaked hyperdrives to make long-haul drives truly efficient in multiple analysis (also the most accurate drive)
- Tweaked Kaldos further to decrease range and increase energy consumption
- Fixed gerax line to make its accuracy improve continually through the whole series
- Updated hyperdrives to consume similar amounts of energy for what they do, with Calista-Dal being most efficient, and Equinox least, but Kaldos now has to pay for its super-speed short-hops.
- Tweaked the Kaldos line to be fixed range (1.4M) but to be much, much faster (so short, fast hops)
- Fixed: re-disabled the interceptors line (shouldn't have been visible - pointless duplicate of vanilla)
- Got rid of a bunch of cycles in tech tree
- Restored the efficient hyperdrive line
- Made Hyperdrive lines all usable to endgame
  - Kalista-Dal are efficient, slower, but have extreme range (long-haul hyperdrives)
  - Kaldos have more limited range, but become crazy fast at turnaround times
- Increased the cost of late game tech (and rebalanced the scaling to be pure doubling)
- Rooted the various planetary weapon facilities earlier in their respective trees so that all races should be able to build them without having to research pointless branches just to find that...
- Simplified a lot of the earlier changes I made now that I understand how the tech .xml files work with the engine to use the simplest approach that works reliably

## Research (Tech Tree)
- Some aspects of the tech tree were simplified or reorganized to make it more logical and reasonable.
  - Removes unnecessary dependencies in the tech tree to make it easier to research what you want.
    - For example - it is not necessary to research all of the "standard" tech in addition to your racial bonus technology in order to get to some upper level technology.  Rather, yours will either eventually give you that opportunity directly, or it has been extended to give you an equivalent end-game version that is every bit as powerful.
  - Removes all of the research bonus requirements from tech, which gives a more even playing field for all races.
    - For example, the Boskara have a -20% research malus, which slows them down, but doesn't randomly preclude them from overcoming that with perseverance.
- Moves most of the real decisions about which direction to go in research starting around the third level of tech, instead of the first two levels.
  - This lends itself to starting games where empires have some worlds already, but aren't then gifted a massive volume of tech.  Rather, they have enough to have expanded, but the core decisions about each empire are yet to be encountered, making playing a game with this advanced start more replayable and fun.
- Tech costs scale better as the game progresses.
  - giving you an early research period where you're just exploring your local system and its immediate neighbors.
  - a middle-game where tech increases apace with empire expansion.
  - a late-game where the most envy-worthy techs require a massive empire to obtain and field.
- NOTE: it is nigh-impossible to design a tech-cost scale that scales perfectly with all of the galaxy-creation options & settings.
  - turn up the difficulty, adjust tech speed, or give the AIs more relative head-starts (or both), and play with various galaxy creation options to find a good game to develop.

## Ship Components
- Ship Components have been extended to provide additional levels of various ships systems.
  - This extends some of the decisions you make to be more functional and therefore viable, allowing you to play through using efficient or energy hog designs without bumping up against as many of vanilla's arbitrary limitations, or without feeling like it's a non-choice because ultimately you're forced to choose the one and only one viable approach after a few tech levels anyway.
  - Which offers more replayability, given more distinct game-paths and strategies to explore & play-through.

## Fleets
- Slightly increased the size of standard fleets to include more and heavier ship classes.
- By default, fleets will use the 33% of fuel range as their operational theater (instead of vanilla 50% default).
  - This is tied to your personal preferences, so as soon as you override them, you can save those for future games.
  - Hopefully this makes them less prone to constantly running out of fuel.
  - And keeps them nearer to their assigned operational theater (so that they can actually perform their duty in their assigned area of space).

## Ships
- Base Cruisers were also bumped up in size to make them more significant, a work-horse of your fleets.
- Heavy ship variants have been made larger, and have a uniform increase in size over their base models.
- Ship sizes can be found here: [DW2-XL Hull Sizes](#dw2-xl-hull-sizes)
  - This offers you more fun with later tech designs, allowing you to play with some of the more exotic technologies...
  - Generally allows you and the AIs to design more balanced designs that include sensors, targeting, and a balanced array of weapons and defenses.
- Ship designs have a more natural progression.
  - First you learn how to field all of the core ship classes up to and including capital ships.
  - Then you're given the ability to research advanced variants of those ship classes in whatever order works for your empire's growth.

## Crew Quarters & Star Marine Barracks
- Marine Barracks supports the equivalent amount of crew as the corresponding crew module, at one tech level higher (thence 2x the tech cost).
  - This makes the ship designer "just work" correctly at every level (overcomes a vanilla bug in ship designer).
  - There are 4 levels of Marine Barracks to keep parity with standard crew modules.
  - They are always lower in capcity than the standard crew quarters for that tech-level (currently 80%).
  - Standard crew quarters also has a better damage reduction, and even ship cost reduction (luxury crew quarters), which marine do not have.

## Sensors
- Short range sensors give some targeting bonus.
  - This should substantially improve default and AI designs that always want a short range sensor on them.
- Long range sensors give a fleet targeting bonus.
- Reduced sizes to allow for more components on your ship designs:
  - Short range sensors are size 5
  - Long range sensors are size 40

## Armor
- Added levels of ionic armor to reach parity with standard armor.
  - Ionic is always one tech level higher, hence 2x the tech cost of the standard armor equivalent.
  - Ionic armors don't upgrade without a refit: like other armor types, you must refit your ships to gain their advantages.
- Paced out the armor for a more even game progression.
- Expanded Boskaran armors to achieve parity with the highest levels of standard armor in the game, but at a lower tech-level.

## Weapons
- Weapon sizes are normalized across all types.
  - Allows more predictable ship designs.
  - Gives you the ability to reason about trade-offs between different weapon systems using their many other factors, such as power use, range, alpha-damage v. sustained damage, etc.
  - Offers you and the AIs to design better ships overall (and typically more weapons at any given ship class).
  - Direct fire weapons:
    - Small = 11
    - Medium = 22
    - Large = 44
  - Tracking weapons:
    - Small = 13
    - Medium = 26
    - Large = 52
  - Area weapons:
    - Large = 52
  - Bombard weapons:
    - Medium = 26

## Kinetic Weapons
- Do not lose damage with distance (slugs don't slow down in space, or lose energy).
  - But they continue to be quite inaccurate with range.
- Kinetic PD does more raw damage than blaster or beam PD, but it's less accurate.
- I overhauled this line of weapons to give a better progression and to make them more viable for non-humans.
  - Humans still have racial locked techs in this line, and they are just better than the standard models.

## Energy Torpedoes
- Boskaran's Firestorm torpedo line is deeper and gives you a large weapon sooner.
  - The large varient requires a separate line of research to unlock.
- Mortalen's torpedoes line is split into fighter/small and medium/large lines.
  - All of these techs have been extended by 4 additional levels each to keep them viable throughout the game.

## Fighters and Bombers
- The Fighter and Bomber craft portion of the tech tree has been made more elegant and rational.
- Fighter Bays themselves have been simplified to hold 4/8/16 craft for S/M/L bays.
  - Rather than holding more units as you tech up, they instead simply get faster at replenishing them.
    - The point being to make fighters & bombers a great tool on the battlefield, but not an automatic "I win" button.
  - So to get more fighters & bombers, you need more bays, larger ships, etc., which is a more natural flow.
    - Vanilla made them an exponentially improving weapon choice, but this mod makes them an incrementally improving weapon choice, like all other weapon platforms.

## Hyper Drives
- Experimental Warp Fields allow for nearby exploration, making for a more gradual early game.
- Have had the various types extended and balanced a bit so that every path is viable, with obvious trade-offs between lag, in flight speed, and balanced performance.
  - Kaldos drives are all about snappy response, but sacrifice range for a single hop.
  - Equinox drives are about flight-speed, and have a middling range per hop.
  - Calista Dal drives are all about range per hop, but sacrifice response time, being quite slow to spin up, as well as significantly slower than either Kaldos or Equinox.  However, they're much lower energy consumption so make for an intriguing option well suited to fusion reactors.

## Reactors
- Have had the various types extended and balanced so that every major type is in-principle viable and competitive, depending on your empire's strategic situation.
  - Fusion reactors are small sized, low-power, efficient with fuel consumpiton.
  - Fission reactors are medium sized, mid-power, and balanced with fuel consumpiton.
  - Antimatter reactors are large sized, high-power, fuel-hogs.
- All reactor types improve on all ares with tech, but fusion remains the gas-sipping king, fission the middle ground, and antimatter the power king.
- All reactor tech lines have been greatly extended to higher levels of tech.

## Targeting and Countermeasures
- Are size 5 for standard ones, and 15 for fleet ones.
- Have been somewhat nerfed to keep these systems from becoming too much of a "win" for high-tech fleets.  They help, they're still critical, but they're not a lock-out against your opponents.
- Sensors have gained some targeting value.
- Long range sensors have gained some fleet targeting value.
- Sensors, targeting computers, and ship classes have been linked by dependencies to help ensure AI researches these (and it doesn't hurt the player to have some decent targeting computers for their ships, likely helps new players as well).

## Colonization
- This tech tree is more freely explorable, not requiring you to research every type of habitat in order to get to the higher levels of technology in those biomes you care about.
  - This should give you and the AIs the ability to make good use out of whatever resources can be found in your region of the galaxy, supplementing tech to make your local area colonizable, rather than needing the perfect map to have a viable and fun game.
- There is a 4th level of colonization modules
  - 50M, 100M, 150M, 200M capacity.
    - The larger capacities give the AI a much better chance of expanding at a more profitable rate without the micromanagement a Human player might indulge in.
    - Currently, the denser modules give a discount on resources required per pop, 0, 10%, 20%, and 30%.
  - Static energy requirements are much higher than vanilla.
  - Static energy cost scales with capacity (so you're unlikely to be able to have more than 1 on a single ship).
  - It is still entirely possible to use more than one module on a given colonizer, boosting your capacity further.  This is the same as vanilla, but can create ships that are hard for your economy to stomach, so use wisely.

## Planetary Facilities
- Maintenance is greatly reduced.
  - This is to make it possible to not wreck yours or the AI's economy.
  - Typically maintenance is 10% of initial buld cost.

## Bug Fixes: Base Game
- Star Marine Barracks has crew which fixes those components to work as expected in the ship designer.
- Troop Academy is limited to 1 per Empire, plugging an exploit for the player (AIs didn't build more than one).

## Bug Fixes: Early versions of this mod
Fixes for bugs I introduced in previous versions of this mod, and later fixed in this version:
- Fixed a bug in Ackdarian Fast Interceptors and Bombers - event now works properly.
- Fixed a bug in Ackdarians who often get an event for maneuvering thrusters.  This now works properly.
- Fixed a bug in Mysterious Plague event -- now works properly.

## DW2-XL Hull Sizes

| type      | standard | type-increase | specialized +20% | advanced +40% | super +60% |
| --------- | -------- | ------------- | ---------------- | ------------- | ---------- |
| Escort    | 375      |               | 450              |               |            |
| Frigate   | 450      | 1.2           | 540              |               |            |
| Destroyer | 600      | 1.33          | 720              |               |            |
| Cruiser   | 800      | 1.33          | 960              | 1120          |            |
| BB+CV     | 1200     | 1.5           | 1440             | 1680          | 1920       |

Base sizes for everything except Cruiser are the same as vanilla, but cruisers were increased from 750 to 800, which gives +33% over Destroyers, which are +33% over Frigates.  Vanilla basically badly shortchanges cruisers in general.

This also allows cruisers to be your early fleet command ships - able to give up a little firepower for fleet targeting & countermeasures, including long range scanners if you're willing to give over the space.

All upgrades in DW2-XL follow a +20% over the previous technology level.  So specialized destroyers are 20% larger than their base size, allowing for a significant increase in firepower (moving from mediums to large weapons where possible, and maximizing speed, defenses, etc.)
