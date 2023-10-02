package utility

import (
	"fmt"
	"github.com/couchbase/go-slab"
	"sort"
	"strings"
)

// PrintStats Support method to print the statistics in an ordered and structured manner
func PrintStats(arena *slab.Arena, stats map[string]int64) {

	arena.Stats(stats)
	slabStats := make(map[string]map[string]int)
	generalStats := make(map[string]int)

	fmt.Println("STATS:")
	for key, value := range stats {
		if strings.HasPrefix(key, "slabClass") {
			parts := strings.SplitN(key, "-", 3)
			slabClass := parts[0] + "-" + parts[1]
			statName := parts[2]

			if _, ok := slabStats[slabClass]; !ok {
				slabStats[slabClass] = make(map[string]int)
			}
			slabStats[slabClass][statName] = int(value)
		} else {
			generalStats[key] = int(value)
		}
	}

	// Print information about each slabClass
	for slabClass, stats := range slabStats {
		fmt.Println(slabClass + ":")
		for statName, value := range stats {
			fmt.Printf("    %s: %d\n", statName, value)
		}
		fmt.Println()
	}

	// Print general information
	fmt.Println("General statistics:")
	for statName, value := range generalStats {
		fmt.Printf("    %s: %d\n", statName, value)
	}
}

func PrintStats2(arena *slab.Arena, stats map[string]int64) {
	slabStats := make(map[string]map[string]int)
	generalStats := make(map[string]int)

	arena.Stats(stats)

	fmt.Println("STATS:")
	for key, value := range stats {
		if strings.HasPrefix(key, "slabClass") {
			parts := strings.SplitN(key, "-", 3)
			slabClass := parts[0] + "-" + parts[1]
			statName := parts[2]

			if _, ok := slabStats[slabClass]; !ok {
				slabStats[slabClass] = make(map[string]int)
			}
			slabStats[slabClass][statName] = int(value)
		} else {
			generalStats[key] = int(value)
		}
	}

	// Ordina le chiavi delle statistiche generali
	var generalKeys []string
	for key := range generalStats {
		generalKeys = append(generalKeys, key)
	}
	sort.Strings(generalKeys)

	// Ordina le chiavi slabClass in ordine decrescente in base all'indice
	var slabClassKeys []string
	for key := range slabStats {
		slabClassKeys = append(slabClassKeys, key)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(slabClassKeys)))

	// Stampa le informazioni ordinate per ogni slabClass
	fmt.Println("STATS:")
	for _, slabClass := range slabClassKeys {
		fmt.Println(slabClass + ":")
		stats := slabStats[slabClass]

		var statNames []string
		for name := range stats {
			statNames = append(statNames, name)
		}
		sort.Strings(statNames)

		for _, statName := range statNames {
			fmt.Printf("    %s: %d\n", statName, stats[statName])
		}
		fmt.Println()
	}

	// Stampa le statistiche generali ordinate
	fmt.Println("General statistics:")
	for _, statName := range generalKeys {
		fmt.Printf("    %s: %d\n", statName, generalStats[statName])
	}
}
