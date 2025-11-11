package service

import (
	"DnD-sheet/internal/character/domain"
	"DnD-sheet/internal/spell"
	"strings"
	"testing"
)

func TestMarkdownFormatter_FormatCharacter(t *testing.T) {
	formatter := NewMarkdownFormatter()

	tests := []struct {
		name     string
		char     *domain.Character
		expected []string // strings that should be present in output
	}{
		{
			name: "Basic Fighter Character",
			char: &domain.Character{
				Name:               "Test Fighter",
				Class:              "fighter",
				Race:               "human",
				Background:         "soldier",
				Level:              3,
				Str:                16,
				Dex:                14,
				Con:                15,
				Int:                10,
				Wis:                12,
				Cha:                10,
				ProficiencyBonus:   2,
				SkillProficiencies: []string{"athletics", "intimidation"},
				Weapon:             "longsword",
				Armor:              "chain mail",
				Shield:             "shield",
			},
			expected: []string{
				"# Test Fighter",
				"Class: fighter",
				"STR: 16 (+3)",
				"[x] Athletics (Str)",
				"[x] Intimidation (Cha)",
				"[] Acrobatics (Dex)",
				"Main hand: longsword",
				"Armor: chain mail",
				"Shield: shield",
				"Armor class: 18",
			},
		},
		{
			name: "Spellcaster with Spells",
			char: &domain.Character{
				Name:               "Test Wizard",
				Class:              "wizard",
				Race:               "elf",
				Background:         "sage",
				Level:              5,
				Str:                8,
				Dex:                14,
				Con:                12,
				Int:                16,
				Wis:                13,
				Cha:                10,
				ProficiencyBonus:   3,
				SkillProficiencies: []string{"arcana", "history"},
				SpellSlots:         map[int]int{0: 4, 1: 4, 2: 3, 3: 2},
				PreparedSpells:     []string{"magic missile", "fireball"},
			},
			expected: []string{
				"# Test Wizard",
				"Class: wizard",
				"## Spell slots",
				"Level 0: 4",
				"Level 1: 4",
				"Level 2: 3",
				"Level 3: 2",
				"## Spellcasting",
				"Spellcasting ability: INT",
				"Spell save DC: 14",      // 8 + 3 prof + 3 int mod
				"Spell attack bonus: +6", // 3 prof + 3 int mod
				"## Spells",
			},
		},
		{
			name: "Non-Spellcaster No Spell Sections",
			char: &domain.Character{
				Name:             "Test Barbarian",
				Class:            "barbarian",
				Race:             "half-orc",
				Background:       "outlander",
				Level:            2,
				Str:              16,
				Dex:              13,
				Con:              15,
				Int:              8,
				Wis:              12,
				Cha:              9,
				ProficiencyBonus: 2,
			},
			expected: []string{
				"# Test Barbarian",
				"Class: barbarian",
				"## Combat stats",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := formatter.FormatCharacter(tt.char)

			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain %q, but it didn't.\nOutput:\n%s", expected, output)
				}
			}

			// Test that non-casters don't have spell sections
			if !tt.char.IsSpellcaster() {
				forbiddenSections := []string{"## Spell slots", "## Spellcasting", "## Spells"}
				for _, section := range forbiddenSections {
					if strings.Contains(output, section) {
						t.Errorf("Non-spellcaster should not have %q section", section)
					}
				}
			}
		})
	}
}

func TestMarkdownFormatter_AbilityModifiers(t *testing.T) {
	formatter := NewMarkdownFormatter()

	tests := []struct {
		score    int
		expected string
	}{
		{score: 8, expected: "(-1)"},
		{score: 10, expected: "(+0)"},
		{score: 12, expected: "(+1)"},
		{score: 16, expected: "(+3)"},
		{score: 20, expected: "(+5)"},
	}

	for _, tt := range tests {
		char := &domain.Character{Str: tt.score}
		output := formatter.FormatCharacter(char)

		if !strings.Contains(output, tt.expected) {
			t.Errorf("Expected STR line to contain %q for score %d", tt.expected, tt.score)
		}
	}
}

func TestMarkdownFormatter_ArmorClass(t *testing.T) {
	tests := []struct {
		name     string
		armor    string
		shield   string
		dex      int
		expected int
	}{
		{name: "No Armor", armor: "", shield: "", dex: 14, expected: 12},                            // 10 + 2 dex
		{name: "Leather + Shield", armor: "leather armor", shield: "shield", dex: 14, expected: 15}, // 11 + 2 dex + 2 shield
		{name: "Chain Mail", armor: "chain mail", shield: "", dex: 14, expected: 16},                // 16 base (dex ignored)
		{name: "Chain Shirt High Dex", armor: "chain shirt", shield: "", dex: 18, expected: 15},     // 13 + 2 (max dex)
		{name: "Plate + Shield", armor: "plate", shield: "shield", dex: 10, expected: 20},           // 18 + 2 shield
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := &domain.Character{
				Armor:  tt.armor,
				Shield: tt.shield,
				Dex:    tt.dex,
			}

			ac := char.ArmorClass()
			if ac != tt.expected {
				t.Errorf("Expected AC %d, got %d for %s", tt.expected, ac, tt.name)
			}
		})
	}
}

func TestMarkdownFormatter_SpellLevels(t *testing.T) {
	tests := []struct {
		spell    string
		expected int
	}{
		{"command", 1},
		{"beacon of hope", 3},
		{"unknown spell", 1}, // default
	}

	for _, tt := range tests {
		level := spell.GetSpellLevel(tt.spell)
		if level != tt.expected {
			t.Errorf("Expected spell %q to be level %d, got %d", tt.spell, tt.expected, level)
		}
	}
}

func TestMarkdownFormatter_EdgeCases(t *testing.T) {
	formatter := NewMarkdownFormatter()

	t.Run("Empty Character", func(t *testing.T) {
		char := &domain.Character{}
		output := formatter.FormatCharacter(char)

		// Should not panic and should produce valid markdown
		if !strings.HasPrefix(output, "# ") {
			t.Error("Should start with markdown title")
		}
	})

	t.Run("Character with No Equipment", func(t *testing.T) {
		char := &domain.Character{
			Name:  "Naked Fighter",
			Class: "fighter",
			Level: 1,
			Str:   10, Dex: 10, Con: 10, Int: 10, Wis: 10, Cha: 10,
		}
		output := formatter.FormatCharacter(char)

		// Should have equipment section but empty
		if !strings.Contains(output, "## Equipment\n\n") {
			t.Error("Should have empty equipment section")
		}
	})

	t.Run("Character with Empty Spell List", func(t *testing.T) {
		char := &domain.Character{
			Name:           "Spellcaster No Spells",
			Class:          "wizard",
			Level:          1,
			PreparedSpells: []string{},
			SpellSlots:     map[int]int{0: 3, 1: 2},
		}
		output := formatter.FormatCharacter(char)

		// Should have spell slots but no spells section
		if !strings.Contains(output, "## Spell slots") {
			t.Error("Should have spell slots section")
		}
		if strings.Contains(output, "## Spells") {
			t.Error("Should not have spells section when no spells prepared")
		}
	})
}
