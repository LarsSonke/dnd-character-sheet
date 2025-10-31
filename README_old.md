# D&D Character Sheet Manager

A comprehensive D&D 5e character sheet management system with CLI and web interface, featuring API integration for spell and equipment enrichment.

## Features

### Core Functionality
- **Character Management**: Create, view, list, update, and delete D&D 5e characters
- **Equipment System**: Equip weapons, armor, and shields with automatic AC calculation
- **Spell System**: Learn and prepare spells with spell slot management
- **Web Interface**: Beautiful HTML character sheets with dynamic data
- **CLI Interface**: Full command-line interface for all operations

### API Integration
- **D&D 5e API Integration**: Automatic enrichment from [dnd5eapi.co](https://www.dnd5eapi.co/)
- **Spell Enrichment**: Adds school, range, components, and descriptions
- **Equipment Enrichment**: Weapon categories, damage types, armor properties
- **Rate Limiting**: Respectful 8 req/sec rate limiting (configurable)
- **Batch Processing**: Efficient concurrent API requests

### Advanced Features
- **Weapon Attack Calculations**: Automatic attack bonuses and damage calculations
- **Spellcasting Support**: Spell attack bonuses for all spellcasting classes
- **Smart Weapon Detection**: Finesse, ranged, two-handed weapon handling
- **Ability Score Modifiers**: Automatic calculation of all modifiers
- **Proficiency System**: Class-based skill and weapon proficiencies

## Quick Start

### Prerequisites
- Go 1.19 or later
- Internet connection for API features

### Installation
```bash
git clone <repository-url>
cd DnD-sheet
go build -o dnd-sheet ./cmd/cli/
```

### Basic Usage

#### Create a Character
```bash
./dnd-sheet create -name "Aragorn" -race "human" -class "ranger" -level 5 \
  -str 16 -dex 18 -con 14 -int 12 -wis 15 -cha 13 -background "folk hero"
```

#### Equip Character
```bash
./dnd-sheet equip -name "Aragorn" -weapon "longsword" -armor "studded leather" -shield ""
```

#### Start Web Server
```bash
./dnd-sheet web -port 8080
# Open http://localhost:8080 in your browser
```

#### View Character
```bash
./dnd-sheet view -name "Aragorn"
```

#### Test API Integration
```bash
./dnd-sheet api-test -spells -class "wizard" -limit 5
./dnd-sheet api-test -equipment -type "weapon" -limit 8
```

## Project Structure

```
DnD-sheet/
├── cmd/cli/                    # Main application entry point
├── internal/
│   ├── api/                    # D&D 5e API client
│   │   └── dnd5e_client.go     # Rate-limited API client
│   ├── character/              # Character domain
│   │   ├── domain/             # Character entities
│   │   ├── infrastructure/     # JSON persistence
│   │   └── service/            # Business logic
│   ├── cli/                    # Command-line interface
│   ├── equipment/              # Equipment enrichment
│   │   └── enrichment_service.go
│   ├── spell/                  # Spell enrichment
│   │   └── enrichment_service.go
│   └── web/                    # Web interface
│       ├── server.go           # HTTP server
│       └── template_data.go    # Template data structures
├── web/
│   ├── templates/              # HTML templates
│   │   └── charactersheet.html # Main character sheet
│   └── static/                 # CSS/JS assets
├── storage/                    # CSV data files
│   ├── classes.csv             # Class definitions
│   ├── equipment.csv           # Equipment database
│   ├── races.csv              # Race definitions
│   └── spells.csv             # Spell database
└── data/                      # Character JSON files (created at runtime)
```

## API Integration

### Spell Enrichment
The system automatically enriches spells with additional data from the D&D 5e API:
- **School**: Evocation, Necromancy, etc.
- **Range**: Touch, 120 feet, etc.
- **Components**: Verbal, Somatic, Material
- **Descriptions**: Full spell descriptions

### Equipment Enrichment
Weapons and armor are enriched with:
- **Weapon Properties**: Finesse, Two-handed, Ranged, etc.
- **Damage Information**: Damage dice and damage types
- **Armor Class**: Base AC and Dex modifier limits
- **Categories**: Simple/Martial weapons, Light/Medium/Heavy armor

### Rate Limiting
- Development: 8 requests/second (respectful to volunteer API)
- Production ready: Configurable up to 50 requests/second
- Concurrent processing with worker pools
- Automatic retry logic with exponential backoff

## Web Interface

The web interface provides:
- **Character List**: Overview of all characters
- **Character Sheets**: Full D&D 5e character sheet layout
- **Dynamic Attacks**: Calculated weapon attacks with bonuses
- **Spell Information**: Enriched spell data display
- **Equipment Details**: Enhanced equipment with API properties

### Character Sheet Features
- Automatic ability score modifiers
- Calculated attack bonuses (Str/Dex + proficiency)
- Weapon damage with proper types
- Spell attack calculations for casters
- Armor class calculations with Dex limits
- Proficiency bonus scaling by level

## Development

### Running Tests
```bash
go test ./...
```

### API Testing
```bash
# Test spell enrichment
./dnd-sheet api-test -spells -class "wizard" -limit 5

# Test equipment enrichment  
./dnd-sheet api-test -equipment -type "weapon" -limit 8
```

### Adding New Features
The project follows clean architecture principles:
- **Domain**: Core business entities
- **Infrastructure**: External concerns (file I/O, APIs)
- **Service**: Business logic and use cases
- **Interface**: CLI and web presentation layers

## Character Data

Characters are stored as JSON files in the `data/` directory. The system supports:
- All official D&D 5e races and classes
- Complete ability score tracking
- Spell slot management by level
- Equipment with automatic calculations
- Background and proficiency tracking

## Commands Reference

| Command | Description | Example |
|---------|-------------|---------|
| `create` | Create new character | `create -name "Hero" -race "elf" -class "wizard" ...` |
| `view` | Display character sheet | `view -name "Hero"` |
| `list` | List all characters | `list` |
| `delete` | Remove character | `delete -name "Hero"` |
| `update` | Update character level | `update -name "Hero" -level 5` |
| `equip` | Equip items | `equip -name "Hero" -weapon "staff" -armor "robes"` |
| `prepare-spell` | Prepare spell | `prepare-spell -name "Hero" -spell "fireball"` |
| `learn-spell` | Learn new spell | `learn-spell -name "Hero" -spell "magic missile"` |
| `web` | Start web server | `web -port 8080` |
| `api-test` | Test API integration | `api-test -spells -equipment` |

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is open source. Please respect the D&D 5e API's terms of service when using the API integration features.

## Acknowledgments

- [D&D 5e API](https://www.dnd5eapi.co/) - Excellent free API for D&D data
- Wizards of the Coast - For creating D&D 5e
- Go Community - For excellent tooling and libraries