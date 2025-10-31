package web

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"DnD-sheet/internal/character/domain"
)

// Server represents the web server for serving character sheets
type Server struct {
	repository domain.CharacterRepository
	templates  *template.Template
}

// NewServer creates a new web server instance
func NewServer(repository domain.CharacterRepository) *Server {
	return &Server{
		repository: repository,
	}
}

// LoadTemplates loads all HTML templates
func (s *Server) LoadTemplates(templateDir string) error {
	templatePath := filepath.Join(templateDir, "*.html")
	templates, err := template.ParseGlob(templatePath)
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}
	s.templates = templates
	return nil
}

// SetupRoutes sets up HTTP routes
func (s *Server) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Serve static files (CSS, JS, etc.)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	// Character routes
	mux.HandleFunc("/", s.handleHome)
	mux.HandleFunc("/character/", s.handleCharacterSheet)

	return mux
}

// handleHome displays a list of all characters
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	characterNames, err := s.repository.List()
	if err != nil {
		http.Error(w, "Failed to load characters", http.StatusInternalServerError)
		return
	}

	// Load character details for each name
	var characters []*domain.Character
	for _, name := range characterNames {
		char, err := s.repository.Load(name)
		if err == nil {
			characters = append(characters, char)
		}
	}

	// For now, create a simple HTML response
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>D&D Character Sheets</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .character-list { margin: 20px 0; }
        .character-item { 
            display: block; 
            padding: 10px; 
            margin: 5px 0; 
            background: #f5f5f5; 
            text-decoration: none; 
            color: #333; 
            border-radius: 5px;
        }
        .character-item:hover { background: #e0e0e0; }
        h1 { color: #8B4513; }
    </style>
</head>
<body>
    <h1>D&D Character Sheets</h1>
    <div class="character-list">
`)

	if len(characters) == 0 {
		fmt.Fprintf(w, "<p>No characters found. Create some characters using the CLI first!</p>")
	} else {
		for _, char := range characters {
			fmt.Fprintf(w, `<a href="/character/%s" class="character-item">%s - Level %d %s %s</a>`,
				char.Name, char.Name, char.Level, char.Race, char.Class)
		}
	}

	fmt.Fprintf(w, `
    </div>
    <p><a href="/">Refresh</a></p>
</body>
</html>`)
}

// handleCharacterSheet displays a character sheet
func (s *Server) handleCharacterSheet(w http.ResponseWriter, r *http.Request) {
	// Extract character name from URL path
	path := strings.TrimPrefix(r.URL.Path, "/character/")
	if path == "" {
		http.Error(w, "Character name required", http.StatusBadRequest)
		return
	}

	// URL decode the character name
	characterName := path

	// Get character from repository
	character, err := s.repository.Load(characterName)
	if err != nil {
		// Check if it's a file not found error (character doesn't exist)
		if os.IsNotExist(err) {
			http.Error(w, fmt.Sprintf("Character '%s' not found", characterName), http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to load character", http.StatusInternalServerError)
		return
	}

	// Convert to template data
	templateData := NewCharacterTemplateData(character)

	// Render the character sheet template
	w.Header().Set("Content-Type", "text/html")
	if err := s.templates.ExecuteTemplate(w, "charactersheet.html", templateData); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		fmt.Printf("Template error: %v\n", err)
		return
	}
}

// Start starts the web server on the specified port
func (s *Server) Start(port int) error {
	mux := s.SetupRoutes()
	addr := fmt.Sprintf(":%d", port)

	fmt.Printf("Starting web server on http://localhost%s\n", addr)
	fmt.Printf("Character sheets available at http://localhost%s\n", addr)

	return http.ListenAndServe(addr, mux)
}
