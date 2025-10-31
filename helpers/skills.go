package helpers

var ClassSkillProficiencies = map[string][]string{
	"barbarian": {"animal handling", "athletics", "intimidation", "nature", "perception", "survival"},
	"bard":      {"acrobatics", "animal handling", "arcana", "athletics", "deception", "history", "insight", "intimidation", "investigation", "medicine", "nature", "perception", "performance", "persuasion", "religion", "sleight of hand", "stealth", "survival"},
	"cleric":    {"history", "insight", "medicine", "persuasion", "religion"},
	"druid":     {"arcana", "animal handling", "insight", "medicine", "nature", "perception", "religion", "survival"},
	"fighter":   {"acrobatics", "animal handling", "athletics", "history", "insight", "intimidation", "perception", "survival"},
	"monk":      {"acrobatics", "athletics", "history", "insight", "religion", "stealth"},
	"paladin":   {"athletics", "insight", "intimidation", "medicine", "persuasion", "religion"},
	"ranger":    {"animal handling", "athletics", "insight", "investigation", "nature", "perception", "stealth", "survival"},
	"rogue":     {"acrobatics", "athletics", "deception", "insight", "intimidation", "investigation", "perception", "performance", "persuasion", "sleight of hand", "stealth"},
	"sorcerer":  {"arcana", "deception", "insight", "intimidation", "persuasion", "religion"},
	"warlock":   {"arcana", "deception", "history", "intimidation", "investigation", "nature", "religion"},
	"wizard":    {"arcana", "history", "insight", "investigation", "medicine", "religion"},
}

var BackgroundSkillProficiencies = map[string][]string{
	"acolyte":       {"insight", "religion"},
	"charlatan":     {"deception", "sleight of hand"},
	"criminal":      {"deception", "stealth"},
	"entertainer":   {"acrobatics", "performance"},
	"folk hero":     {"animal handling", "survival"},
	"guild artisan": {"insight", "persuasion"},
	"hermit":        {"medicine", "religion"},
	"noble":         {"history", "persuasion"},
	"outlander":     {"athletics", "survival"},
	"sage":          {"arcana", "history"},
	"sailor":        {"athletics", "perception"},
	"soldier":       {"athletics", "intimidation"},
	"urchin":        {"sleight of hand", "stealth"},
}

var ClassSkillCount = map[string]int{
	"barbarian": 2,
	"bard":      3,
	"cleric":    2,
	"druid":     2,
	"fighter":   2,
	"monk":      2,
	"paladin":   2,
	"ranger":    3,
	"rogue":     4,
	"sorcerer":  2,
	"warlock":   2,
	"wizard":    2,
}

var CasterClasses = map[string]bool{
	"bard":      true,
	"cleric":    true,
	"druid":     true,
	"paladin":   true,
	"ranger":    true,
	"sorcerer":  true,
	"warlock":   true,
	"wizard":    true,
}
