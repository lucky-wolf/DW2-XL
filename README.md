# DW2-XL
Distant Worlds 2 - XL

Author: Mordachai (lucky-wolf)

Version: v1.1.5

- [DW2-XL](#dw2-xl)
	- [Latest Changes](#latest-changes)
	- [Guiding Principles](#guiding-principles)
	- [Research (Tech Tree)](#research-tech-tree)
	- [Ship Components](#ship-components)
	- [Fleets](#fleets)
	- [Ships](#ships)
	- [Crew & Star Barracks](#crew--star-barracks)
	- [Sensors](#sensors)
	- [Armor](#armor)
	- [Weapons](#weapons)
	- [Kinetic Weapons](#kinetic-weapons)
	- [Energy Torpedos](#energy-torpedos)
	- [Fighters and Bombers](#fighters-and-bombers)
	- [Hyper Drives](#hyper-drives)
	- [Reactors](#reactors)
	- [Targeting and Countermeasures](#targeting-and-countermeasures)
	- [Colonization](#colonization)
	- [Planetary Facilities](#planetary-facilities)
	- [Vanilla Bug Fixes](#vanilla-bug-fixes)
	- [Mod Bug Fixes](#mod-bug-fixes)

## Latest Changes

- Tweaked Long Hault hyper drives to make them a little more attractive
- Moved ion defenses to appear near shields && reorganized point deflectors (right next to them)
  - both are rooted in advanced deflectors rather than in weapons techs
- Separated scanners and jammers into their own tech-lines (independent of each other)
- Further tweaked hyperdrives to make long-haul drives truely efficient in multiple analysis (also the most accurate drive)
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

## Guiding Principles
Principly this mod aims to create a better player experience while playing games of Distant Wordls 2.
It is not a completely different game than vanilla, rather it's "just better" in every way I had time or insight to make something better than it was. YMMV, but this is my take on making it better.

## Research (Tech Tree)
- Some aspects of the tech tree were simplified or reorganized to make it more logical and reasonable
  - Removes unnecessary dependencies in the tech tree to make it easier to research what you want
    - For example - it is not necessary to research all of the "standard" tech in addition to your racial bonus technology in order to get to some upper level technology.  Rather, yours will either eventually give you that opportunity directly, or it has been extended to give you an equivalent end-game version that is every bit as powerful.
  - Removes all of the research bonus requirements from tech, which gives a more even playing field for all races.
    - For example, the Boskara have a -20% research malus, which slows them down, but doesn't randomly preclude them from overcoming that with perseverence.
- Moves most of the real decisions about which direction to go in research starting around the third level of tech, instead of the first two levels.
  - This lends itself to starting games where empires have some worlds already, but aren't then gifted a massive volume of tech.  Rather, they have enough to have expanded, but the core decisions about each empire are yet to be encountered, making playing a game with this advanced start more replayable and fun.
- Tech costs scale better as the game progresses
  - giving you an early research period where you're just exploring your local system and its immediate neighbors
  - a middle-game where tech increases apace with empire expansion
  - a late-game where the most envy-worthy techs require a massive empire to obtain and field

## Ship Components
- Ship Components have been extended to provide additional levels of various ships systems
  - This extends some of the decisions you make to be more functional and therefore viable, allowing you to play through using efficent or energy hog designs without bumping up against as many of vanilla's arbitrary limitations, or without feeling like it's a non-choice because ultiamtely you're forced to choose the one and only one viable approach after a few tech levels anyway.

## Fleets
- Slightly increased the size of standard fleets to include more and heavier ship classes
- By default, fleets will use the 33% of fuel range as their operational theater
  - Hopefully this makes them less idiotically prone to constantly running out of fuel
  - And keeps them nearer to their operational theater

## Ships
- Many of the ship classes have been enlarged
  - this offers you more flexibility in designing your ships,
  - allows you to play with some of the more exotic technologies, and
  - generally creating more balanced designs (the AI is helped by this change as well)
- Ship designs have a more natural progression
  - first you learn how to field all of the core ship classes shy of capital ships
  - then you're given the ability to research advanced variants of those ship classes or expand into capital ships in whatever order works for your empire's growth.
- Overall, ships are bigger in this mod than in vanilla, giving you the ability to load them up with all the fun things you've researched to have the mega battles we're all hoping for!
- Several ship classes have been renamed

## Crew & Star Barracks
- Marine Barracks always branch off of related crew tech
- Each level of Marine Barracks supports the equivalent amount of crew as the corresponding crew module
- There are 4 levels of Marine Barracks to keep parity
- This makes the ship designer "just work" correctly at every level

## Sensors
- Small sensors give some small amount of targeting bonus
- Large sensors give a small fleet targeting bonus
- Small sensors are size 5
- Large sensors are size 40

## Armor
- Separated out ionic armor from armor-materials
  - Ionic require you to research the base material, and then an additional tech to realize its Ionic enhanced derivative
  - Ionic armors don't upgrade in situ.  Like other armor types, you must refit your ships to gain their advantages.
  - Hence Ionics are a more powerful armor type that you pay for in terms of research and construction, not a weird alternative that plays by different rules.
- Paced out the armor for a more even game progression
- Simplified Boskaran armors to achieve the highest levels of armor without having to re-research the standard ones
  - Also gave them the special facilities at high levels of their research tree so that again, they don't need to also research standard armor to gain access

## Weapons
- Weapon sizes are normalized
  - Allows more predictable ship designs
  - Gives you the ability to reason about trade-offs between different weapon systems using their many other factors, such as power use, range, alpha-damage v. sustained damage, etc.
  - Offers you and the AIs to design better ships overall (and typically more weapons at any given ship class)
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
- Do not degrate with distance (slugs don't slow down in space, or lose energy)
  - But they continue to be quite inaccurate with range
- Kinetic PD does more raw damage than blaster or beam PD, but it's less accurate

## Energy Torpedos
- Boskaran's Firestorm torpedo line is deeper and gives you a large weapon sooner
- Mortalen's torpedos line is split into ftr+small and med+lrg
  - The med/large line also gives the planetary torpedo facilities
  - all of these techs have been extended by 4 additional levels each to keep them viable throughout the game

## Fighters and Bombers
- The Fighter and Bomber craft portion of the tech tree has been made more elegant and rational.
- Fighter Bays themselves have been simplified to hold 4/8/16 craft for S/M/L bays.
  - Rather than holding more units as you tech up, they instead simply get faster at replenishing them.
    - The point being to make fighters & bombers a great tool on the battlefield, but not an automatic "I win" button.
  - So to get more fighters & bombers, you need more bays, larger ships, etc., which is a more natural flow.

## Hyper Drives
- Experimental Warp Fields allow for nearby exploration, making for a more gradual early game
- Have had the various types extended and balanced a bit so that every path is viable, with obvious trade-offs between lag, in flight speed, and balanced performance
- Will make another pass at this at some point to really make this "sing"

## Reactors
- Have had the various types extended and balanced so that every major type is in-principle viable and competitive, depending on your empire's strategic situation (has enough fuel to feed the engines of destruction)
- Will make another pass at this at some point to really make this "sing"

## Targeting and Countermeasures
- Are size 5
- Have been somewhat nerfed to keep these systems from becoming too much of a "win" for high-tech fleets.  They help, they're still critical, but they're not a lock-out against your opponents.
- Sensors have gained some targeting value
- Long range sensors have gained some fleet targeting value

## Colonization
- This tech tree is more freely explorable, not requiring you to research every type of habitat in order to get to the higher levels of technology in those biomes you care about.
  - This should give you and the AIs the ability to make good use out of whatever resources can be found in your region of the galaxy, supplementing tech to make your local area colonizable, rather than needing the perfect map to have a viable and fun game.
- There is a 4th level of colonization modules
  - 25M, 50M, 75M, 100M capacity
  - Static energy is higher than vanilla
  - Static energy cost scales with capacity (so you're unlikely to be able to have more than 1 on a single ship)

## Planetary Facilities
- Maintenance is much reduced to make it possible to not wreck yours (or the AI's) economy.
  - At some point I will probably increase these costs, but nowhere near what they were in vanilla.

## Vanilla Bug Fixes
- Technocracy is only available to Ackdarians, as stated in the game's messaging
- Giving crew to Star Marine Barracks fixes those components to work as expected
- Limit the Troop Academy to 1 per Empire, plugging an exploit for the player (AIs didn't build more than one)

## Mod Bug Fixes
Fixes for bugs I introduced in previous versions of this mod, and later fixed in this version:
- Fixed a bug in Ackdarian Fast Interceptors and Bombers - event now works properly.
- Fixed a bug in Mysterious Plague event -- now works properly.