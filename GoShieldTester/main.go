package main

/*
 * Go version of Down To Earth Astronomy's Shield Tester
 * Original version here:  https://github.com/DownToEarthAstronomy/D2EA_Shield_tester
 *
 * Go port by Andrew van der Stock, vanderaj@gmail.com
 *
 * Why this version? It's fast - about 15,500 times faster than the PowerShell version, so multi-threading is not necessary even with 8 utility slots and the full list
 */

import (
	"flag"
	"fmt"
	"time"
)

func processFlags() {
	flag.IntVar(&config.shieldBoosterCount, "boosters", config.shieldBoosterCount, "Number of Shield Boosters")
	flag.Float64Var(&config.explosiveDPS, "edps", config.explosiveDPS, "Explosive DPS percentage (use 0 for Thargoids)")
	flag.Float64Var(&config.kineticDPS, "kdps", config.kineticDPS, "Kinetic DPS percentage (use 0 for Thargoids)")
	flag.Float64Var(&config.thermalDPS, "tdps", config.thermalDPS, "Thermal DPS percentage (use 0 for Thargoids)")
	flag.Float64Var(&config.absoluteDPS, "adps", config.absoluteDPS, "Absolute DPS percentage (use 100 for Thargoids)")
	flag.Float64Var(&config.damageEffectiveness, "dmg", config.damageEffectiveness, "Damage effectiveness (use 0.1 for PvE, 0.5 for PvP, and 0.65 for Thargoids)")
	flag.Float64Var(&config.scbHitPoint, "scb", config.scbHitPoint, "SCB HitPoints (default 0)")
	flag.Float64Var(&config.guardianShieldHitPoint, "gshp", config.guardianShieldHitPoint, "Guardian HitPoints (default 0)")

	prismatics := flag.Bool("noprismatics", false, "Disable Prismatic shields")
	thargoid := flag.Bool("thargoid", false, "Useful Thargoid defaults")
	cucumber := flag.Bool("cucumber", false, "Useful Cucumber defaults")
	allboosters := flag.Bool("fullboost", false, "Load the full booster list")

	flag.Parse()

	if config.shieldBoosterCount < 0 {
		fmt.Println("Can't have negative shield boosters, setting to 0")
		config.shieldBoosterCount = 0
	}

	if config.shieldBoosterCount > 8 {
		fmt.Println("No current ship has more than 8 shield boosters, setting to 8")
		config.shieldBoosterCount = 8
	}

	if *prismatics {
		fmt.Println("Disabling Prismatics")
		config.prismatics = false
	}

	if *thargoid && *cucumber {
		fmt.Println("D2EA is not a Thargoid, loading only Cucumber")
		*thargoid = false
	}

	if *allboosters {
		fmt.Println("Loading all boosters")
		config.boosterFile = "../ShieldBoosterVariants.csv"
	}

	if *cucumber {
		fmt.Println("Loading Cucumber defenses")
		config.explosiveDPS = 0
		config.kineticDPS = 83
		config.thermalDPS = 47
		config.absoluteDPS = 0
		config.damageEffectiveness = 0.25
		config.shieldBoosterCount = 7
	}
	if *thargoid {
		fmt.Println("Loading Thargoid defenses")
		config.explosiveDPS = 0
		config.kineticDPS = 0
		config.thermalDPS = 0
		config.absoluteDPS = 200
		config.damageEffectiveness = 0.10
		config.shieldBoosterCount = 7
	}
}

func main() {
	fmt.Println("Down to Earth Astronomy's ShieldTester (https://github.com/DownToEarthAstronomy/D2EA_Shield_tester)")
	fmt.Println("Go port by Andrew van der Stock, vanderaj@gmail.com")
	fmt.Println()

	config = loadConfig()

	processFlags()

	fmt.Println("Test started at: ", time.Now().Format(time.RFC1123))

	var generators = loadGenerators()
	fmt.Println("Loaded", len(generators), "generator variants")

	var boosterVariants = loadboosterVariants()
	fmt.Println("Loaded", len(boosterVariants), "shield booster variants")

	var shieldBoosterLoadoutList = getBoosterLoadoutList(len(boosterVariants))

	fmt.Println("Loadout shield booster variations to be tested per generator: ", len(shieldBoosterLoadoutList))
	fmt.Println("Total loadouts to be tested: ", len(shieldBoosterLoadoutList)*len(generators))

	startTime := time.Now()

	var result = testGenerators(generators, boosterVariants, shieldBoosterLoadoutList)

	endTime := time.Now()

	fmt.Println("Time elapsed: [", endTime.Sub(startTime), "]")
	fmt.Println()

	showResults(result, boosterVariants)
}
