# DW2-XL
Distant Worlds 2 - XL

Author: Mordachai (lucky-wolf)

- [DW2-XL](#dw2-xl)
	- [Guiding Principles](#guiding-principles)
	- [Mod Highlights](#mod-highlights)
	- [Latest Changes](#latest-changes)
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
- Removed "Robo-" from Zennox troop names (now they're are Ice whatever) since they're not robotic, so it made no sense!

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

## Ship Components
- Ship Components have been extended to provide additional levels of various ships systems.
  - This extends some of the decisions you make to be more functional and therefore viable, allowing you to play through using efficient or energy hog designs without bumping up against as many of vanilla's arbitrary limitations, or without feeling like it's a non-choice because ultimately you're forced to choose the one and only one viable approach after a few tech levels anyway.

## Fleets
- Slightly increased the size of standard fleets to include more and heavier ship classes.
- By default, fleets will use the 33% of fuel range as their operational theater (instead of vanilla 50% default).
  - Hopefully this makes them less prone to constantly running out of fuel.
  - And keeps them nearer to their operational theater (so that they can actually perform their durty in their assigned area of space).

## Ships
- Many of the ship classes have been enlarged.
  - this offers you more flexibility in designing your ships...
  - allows you to play with some of the more exotic technologies...
  - generally creating more balanced designs (the AI is helped by this change as well).
- Ship designs have a more natural progression.
  - first you learn how to field all of the core ship classes shy of capital ships
  - then you're given the ability to research advanced variants of those ship classes or expand into capital ships in whatever order works for your empire's growth.
- Overall, ships are bigger in this mod than in vanilla, giving you the ability to load them up with more of the fun things you've researched.

## Crew & Star Barracks
- Marine Barracks supports the equivalent amount of crew as the corresponding crew module, at one tech level higher (thence 2x the tech cost).
- There are 4 levels of Marine Barracks to keep parity with standard crew modules.
- This makes the ship designer "just work" correctly at every level (overcomes a vinalla bug in ship designer).

## Sensors
- Short range sensors give some targeting bonus.
  - This should substantially improve default and AI designs that always want a short range sensor on them.
- Long range sensors give a fleet targeting bonus.
- Reduced sizes to allow for more components on your ship designs:
  - Short range sensors are size 5
  - Long range sensors are size 40

## Armor
- Added levels of ionic armor to reach parity with standard armor.
  - Ionic is alwas one tech level higher, hence 2x the tech cost of the standard armor equivalent.
  - Ionic armors don't upgrade without a refit: like other armor types, you must refit your ships to gain their advantages.
- Paced out the armor for a more even game progression.
- Expanded Boskaran armors to achieve parity with the highest levels of standard armor in the game, but at a lower tech-level.

## Weapons
- Weapon sizes are normalized across all types.
  - Allows more predictable ship designs.
  - Gives you the ability to reason about trade-offs between different weapon systems using their many other factors, such as power use, range, alpha-damage v. sustained damage, etc.
  - Offers you and the AIs to design better ships overall (and typically more weapons at any given ship class).
  - Direct fire weapons:
    - Small = 12
    - Medium = 25
    - Large = 50
  - Tracking weapons:
    - Small = 15
    - Medium = 30
    - Large = 60
  - Area and Bombardment weapons:
    - Medium = 25
    - Large = 50
    - Mines = 60

## Kinetic Weapons
- Do not lose damage with distance (slugs don't slow down in space, or lose energy).
  - But they continue to be quite inaccurate with range.
- Kinetic PD does more raw damage than blaster or beam PD, but it's less accurate.

## Energy Torpedoes
- Boskaran's Firestorm torpedo line is deeper and gives you a large weapon sooner.
- Mortalen's torpedoes line is split into fighter/small and medium/large lines.
  - all of these techs have been extended by 4 additional levels each to keep them viable throughout the game.

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

## Reactors
- Have had the various types extended and balanced so that every major type is in-principle viable and competitive, depending on your empire's strategic situation.

## Targeting and Countermeasures
- Are size 5 for standard ones, and 15 for fleet ones.
- Have been somewhat nerfed to keep these systems from becoming too much of a "win" for high-tech fleets.  They help, they're still critical, but they're not a lock-out against your opponents.
- Sensors have gained some targeting value.
- Long range sensors have gained some fleet targeting value.

## Colonization
- This tech tree is more freely explorable, not requiring you to research every type of habitat in order to get to the higher levels of technology in those biomes you care about.
  - This should give you and the AIs the ability to make good use out of whatever resources can be found in your region of the galaxy, supplementing tech to make your local area colonizable, rather than needing the perfect map to have a viable and fun game.
- There is a 4th level of colonization modules
  - 50M, 100M, 150M, 200M capacity.
    - The larger capacities give the AI a much better chance of expanding at a more profitable rate without the micromanagement a Human player might indulge in.
  - Static energy requirements are much higher than vanilla.
  - Static energy cost scales with capacity (so you're unlikely to be able to have more than 1 on a single ship).
    - note: you can still manage it in some situations - however, it's expensive, and not an easy way to out-colonize the AI as compared to vanilla.

## Planetary Facilities
- Maintenance is greatly reduced.
  - This is to make it possible to not wreck yours or the AI's economy.

## Bug Fixes: Base Game
- Technocracy is only available to Ackdarians, as stated in the game's original messaging.
  - Although they've updated vanilla to give Zenox technocracy, the Zenox already start with a +20 to all research on their homeworld.  This is already obscene, and the do NOT need to double that with technocracy.  So NO, it is not available to them in this mod.
- Giving crew to Star Marine Barracks fixes those components to work as expected in the ship designer.
- Troop Academy is limited to 1 per Empire, plugging an exploit for the player (AIs didn't build more than one).

## Bug Fixes: Early versions of this mod
Fixes for bugs I introduced in previous versions of this mod, and later fixed in this version:
- Fixed a bug in Ackdarian Fast Interceptors and Bombers - event now works properly.
- Fixed a bug in Ackdarians who often get an event for maneuvering thrusters.  This now works properly.
- Fixed a bug in Mysterious Plague event -- now works properly.
