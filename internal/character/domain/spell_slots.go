package domain

import "strings"

// HalfCasterSpellSlots returns spell slots for half-caster classes (paladin/ranger)
func HalfCasterSpellSlots(level int) map[int]int {
	// D&D 5e half-caster spell slots (paladin/ranger)
	// Source: PHB Table
	slots := map[int][]int{
		1:  {0, 0, 0, 0, 0}, // Level 1
		2:  {2, 0, 0, 0, 0},
		3:  {3, 0, 0, 0, 0},
		4:  {3, 0, 0, 0, 0},
		5:  {4, 2, 0, 0, 0},
		6:  {4, 2, 0, 0, 0},
		7:  {4, 3, 0, 0, 0},
		8:  {4, 3, 0, 0, 0},
		9:  {4, 3, 2, 0, 0},
		10: {4, 3, 2, 0, 0},
		11: {4, 3, 3, 0, 0},
		12: {4, 3, 3, 0, 0},
		13: {4, 3, 3, 1, 0},
		14: {4, 3, 3, 1, 0},
		15: {4, 3, 3, 2, 0},
		16: {4, 3, 3, 2, 0},
		17: {4, 3, 3, 2, 1},
		18: {4, 3, 3, 2, 1},
		19: {4, 3, 3, 2, 2},
		20: {4, 3, 3, 3, 2},
	}
	// Default to highest if level > 20
	if level > 20 {
		level = 20
	}
	arr := slots[level]
	result := map[int]int{}
	for i, v := range arr {
		if v > 0 {
			result[i+1] = v
		}
	}
	return result
}

// FullCasterSpellSlots returns spell slots for full-caster classes (wizard, cleric, etc.)
func FullCasterSpellSlots(level int) map[int]int {
	// D&D 5e full-caster spell slots (wizard, cleric, etc.)
	// Source: PHB Table
	slots := map[int][]int{
		1:  {2, 0, 0, 0, 0, 0, 0, 0, 0},
		2:  {3, 0, 0, 0, 0, 0, 0, 0, 0},
		3:  {4, 2, 0, 0, 0, 0, 0, 0, 0},
		4:  {4, 3, 0, 0, 0, 0, 0, 0, 0},
		5:  {4, 3, 2, 0, 0, 0, 0, 0, 0},
		6:  {4, 3, 3, 0, 0, 0, 0, 0, 0},
		7:  {4, 3, 3, 1, 0, 0, 0, 0, 0},
		8:  {4, 3, 3, 2, 0, 0, 0, 0, 0},
		9:  {4, 3, 3, 3, 1, 0, 0, 0, 0},
		10: {4, 3, 3, 3, 2, 0, 0, 0, 0},
		11: {4, 3, 3, 3, 2, 1, 0, 0, 0},
		12: {4, 3, 3, 3, 2, 1, 0, 0, 0},
		13: {4, 3, 3, 3, 2, 1, 1, 0, 0},
		14: {4, 3, 3, 3, 2, 1, 1, 0, 0},
		15: {4, 3, 3, 3, 2, 1, 1, 1, 0},
		16: {4, 3, 3, 3, 2, 1, 1, 1, 0},
		17: {4, 3, 3, 3, 2, 1, 1, 1, 1},
		18: {4, 3, 3, 3, 3, 1, 1, 1, 1},
		19: {4, 3, 3, 3, 3, 2, 1, 1, 1},
		20: {4, 3, 3, 3, 3, 2, 2, 1, 1},
	}
	if level > 20 {
		level = 20
	}
	arr := slots[level]
	result := map[int]int{}
	for i, v := range arr {
		if v > 0 {
			result[i+1] = v
		}
	}
	return result
}

// MaxSpellSlotLevel returns the maximum spell slot level for a class
func MaxSpellSlotLevel(class string) int {
	switch strings.ToLower(class) {
	case "paladin", "ranger":
		return 5
	case "wizard", "cleric", "druid", "bard", "sorcerer":
		return 9
	default:
		return 0
	}
}

// PactMagicSpellSlots returns Warlock Pact Magic spell slots
func PactMagicSpellSlots(level int) map[int]int {
	// D&D 5e Warlock Pact Magic slots
	// Source: PHB Table
	type pact struct {
		slots int
		level int
	}
	table := map[int]pact{
		1:  {1, 1},
		2:  {2, 1},
		3:  {2, 2},
		4:  {2, 2},
		5:  {2, 3},
		6:  {2, 3},
		7:  {2, 4},
		8:  {2, 4},
		9:  {2, 5},
		10: {2, 5},
		11: {3, 5},
		12: {3, 5},
		13: {3, 5},
		14: {3, 5},
		15: {3, 5},
		16: {4, 5},
		17: {4, 5},
		18: {4, 5},
		19: {4, 5},
		20: {4, 5},
	}
	result := map[int]int{}
	if pact, ok := table[level]; ok {
		result[pact.level] = pact.slots
	}
	// Warlocks also have cantrips (level 0)
	cantrips := map[int]int{
		1: 2, 2: 2, 3: 2, 4: 3, 5: 3, 6: 3, 7: 3, 8: 3, 9: 3, 10: 4,
		11: 4, 12: 4, 13: 4, 14: 4, 15: 4, 16: 4, 17: 4, 18: 4, 19: 4, 20: 4,
	}
	result[0] = cantrips[level]
	return result
}

// FullCasterCantrips returns the number of cantrips for full casters
func FullCasterCantrips(level int) int {
	switch {
	case level >= 17:
		return 5
	case level >= 10:
		return 5
	case level >= 4:
		return 4
	default:
		return 3
	}
}
