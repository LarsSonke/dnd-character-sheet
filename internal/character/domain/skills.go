package domain

import "math"

// SkillAbility maps each skill to its governing ability
var SkillAbility = map[string]string{
	"Acrobatics":      "Dex",
	"Animal Handling": "Wis",
	"Arcana":          "Int",
	"Athletics":       "Str",
	"Deception":       "Cha",
	"History":         "Int",
	"Insight":         "Wis",
	"Intimidation":    "Cha",
	"Investigation":   "Int",
	"Medicine":        "Wis",
	"Nature":          "Int",
	"Perception":      "Wis",
	"Performance":     "Cha",
	"Persuasion":      "Cha",
	"Religion":        "Int",
	"Sleight of Hand": "Dex",
	"Stealth":         "Dex",
	"Survival":        "Wis",
}

// Modifier returns the ability modifier for a score
func Modifier(score int) int {
	return int(math.Floor(float64(score-10) / 2.0))
}

// SkillModifiers returns a map of skill name to modifier for the character
func (c *Character) SkillModifiers() map[string]int {
	prof := make(map[string]bool)
	for _, s := range c.SkillProficiencies {
		prof[s] = true
	}
	skillMods := make(map[string]int)
	for skill, ability := range SkillAbility {
		var score int
		switch ability {
		case "Str":
			score = c.Str
		case "Dex":
			score = c.Dex
		case "Con":
			score = c.Con
		case "Int":
			score = c.Int
		case "Wis":
			score = c.Wis
		case "Cha":
			score = c.Cha
		}
		mod := Modifier(score)
		if prof[skill] {
			mod += c.ProficiencyBonus
		}
		skillMods[skill] = mod
	}
	return skillMods
}
