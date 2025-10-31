package domain

import (
	"math/rand"
	"time"
)

// asiLevels defines which levels grant Ability Score Improvements
var asiLevels = map[int]bool{4: true, 8: true, 12: true, 16: true, 19: true}

// ApplySRDAbilityScoreImprovements applies ASI to a character for levels gained
func (c *Character) ApplySRDAbilityScoreImprovements(oldLevel, newLevel int) {
	rand.Seed(time.Now().UnixNano())
	for lvl := oldLevel + 1; lvl <= newLevel; lvl++ {
		if asiLevels[lvl] {
			// Pick two different stats to increase by 1
			stats := []string{"Str", "Dex", "Con", "Int", "Wis", "Cha"}
			rand.Shuffle(len(stats), func(i, j int) { stats[i], stats[j] = stats[j], stats[i] })
			for i := 0; i < 2; i++ {
				switch stats[i] {
				case "Str":
					c.Str++
				case "Dex":
					c.Dex++
				case "Con":
					c.Con++
				case "Int":
					c.Int++
				case "Wis":
					c.Wis++
				case "Cha":
					c.Cha++
				}
			}
		}
	}
}
