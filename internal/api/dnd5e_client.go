package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	BaseURL = "https://www.dnd5eapi.co/api"
	// Rate limit: 5-10 requests per second for development (being nice to volunteers)
	RequestsPerSecond = 8
	RateLimitDelay    = time.Second / RequestsPerSecond
)

// Client represents a D&D 5e API client with rate limiting
type Client struct {
	httpClient  *http.Client
	rateLimiter *time.Ticker
	mu          sync.Mutex
}

// NewClient creates a new D&D 5e API client
func NewClient() *Client {
	return &Client{
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		rateLimiter: time.NewTicker(RateLimitDelay),
	}
}

// Close stops the rate limiter
func (c *Client) Close() {
	c.rateLimiter.Stop()
}

// makeRequest makes a rate-limited HTTP request to the API
func (c *Client) makeRequest(endpoint string) (*http.Response, error) {
	// Wait for rate limiter
	<-c.rateLimiter.C

	url := fmt.Sprintf("%s%s", BaseURL, endpoint)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("API returned status %d for %s", resp.StatusCode, url)
	}

	return resp, nil
}

// SpellDetails represents detailed spell information from the API
type SpellDetails struct {
	Index  string `json:"index"`
	Name   string `json:"name"`
	Level  int    `json:"level"`
	School struct {
		Index string `json:"index"`
		Name  string `json:"name"`
	} `json:"school"`
	Range         string   `json:"range"`
	Components    []string `json:"components"`
	Duration      string   `json:"duration"`
	CastingTime   string   `json:"casting_time"`
	Description   []string `json:"desc"`
	HigherLevel   []string `json:"higher_level,omitempty"`
	Ritual        bool     `json:"ritual"`
	Concentration bool     `json:"concentration"`
	Classes       []struct {
		Index string `json:"index"`
		Name  string `json:"name"`
	} `json:"classes"`
}

// WeaponDetails represents detailed weapon information from the API
type WeaponDetails struct {
	Index             string `json:"index"`
	Name              string `json:"name"`
	EquipmentCategory struct {
		Index string `json:"index"`
		Name  string `json:"name"`
	} `json:"equipment_category"`
	WeaponCategory string `json:"weapon_category"`
	WeaponRange    string `json:"weapon_range"`
	CategoryRange  string `json:"category_range"`
	Range          struct {
		Normal int `json:"normal,omitempty"`
		Long   int `json:"long,omitempty"`
	} `json:"range,omitempty"`
	Damage struct {
		DamageDice string `json:"damage_dice"`
		DamageType struct {
			Index string `json:"index"`
			Name  string `json:"name"`
		} `json:"damage_type"`
	} `json:"damage"`
	Properties []struct {
		Index string `json:"index"`
		Name  string `json:"name"`
	} `json:"properties"`
	TwoHanded bool `json:"two_handed_damage,omitempty"`
}

// ArmorDetails represents detailed armor information from the API
type ArmorDetails struct {
	Index             string `json:"index"`
	Name              string `json:"name"`
	EquipmentCategory struct {
		Index string `json:"index"`
		Name  string `json:"name"`
	} `json:"equipment_category"`
	ArmorCategory string `json:"armor_category"`
	ArmorClass    struct {
		Base     int  `json:"base"`
		DexBonus bool `json:"dex_bonus"`
		MaxBonus int  `json:"max_bonus,omitempty"`
	} `json:"armor_class"`
	StrMinimum          int  `json:"str_minimum,omitempty"`
	StealthDisadvantage bool `json:"stealth_disadvantage,omitempty"`
}

// GetSpell fetches detailed spell information from the API
func (c *Client) GetSpell(spellName string) (*SpellDetails, error) {
	// Convert spell name to API index format (lowercase, hyphens instead of spaces)
	index := strings.ToLower(strings.ReplaceAll(spellName, " ", "-"))
	index = strings.ReplaceAll(index, "'", "")

	endpoint := fmt.Sprintf("/spells/%s", url.PathEscape(index))

	resp, err := c.makeRequest(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var spell SpellDetails
	if err := json.NewDecoder(resp.Body).Decode(&spell); err != nil {
		return nil, fmt.Errorf("failed to decode spell response: %w", err)
	}

	return &spell, nil
}

// GetEquipment fetches detailed equipment information from the API
func (c *Client) GetEquipment(equipmentName string) (interface{}, error) {
	// Convert equipment name to API index format
	index := strings.ToLower(strings.ReplaceAll(equipmentName, " ", "-"))
	index = strings.ReplaceAll(index, "'", "")

	endpoint := fmt.Sprintf("/equipment/%s", url.PathEscape(index))

	resp, err := c.makeRequest(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// First, decode to determine equipment category
	var basicEquipment struct {
		EquipmentCategory struct {
			Index string `json:"index"`
			Name  string `json:"name"`
		} `json:"equipment_category"`
	}

	// Read response body
	var responseData json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return nil, fmt.Errorf("failed to decode equipment response: %w", err)
	}

	// Decode basic info to check category
	if err := json.Unmarshal(responseData, &basicEquipment); err != nil {
		return nil, fmt.Errorf("failed to decode equipment category: %w", err)
	}

	// Decode based on equipment category
	category := basicEquipment.EquipmentCategory.Index
	switch category {
	case "weapon":
		var weapon WeaponDetails
		if err := json.Unmarshal(responseData, &weapon); err != nil {
			return nil, fmt.Errorf("failed to decode weapon: %w", err)
		}
		return &weapon, nil
	case "armor":
		var armor ArmorDetails
		if err := json.Unmarshal(responseData, &armor); err != nil {
			return nil, fmt.Errorf("failed to decode armor: %w", err)
		}
		return &armor, nil
	default:
		// For other equipment types, return basic info
		var equipment map[string]interface{}
		if err := json.Unmarshal(responseData, &equipment); err != nil {
			return nil, fmt.Errorf("failed to decode equipment: %w", err)
		}
		return equipment, nil
	}
}

// BatchResult represents the result of a batch operation
type BatchResult struct {
	Name  string
	Data  interface{}
	Error error
}

// GetSpellsBatch fetches multiple spells concurrently with rate limiting
func (c *Client) GetSpellsBatch(spellNames []string) []BatchResult {
	results := make([]BatchResult, len(spellNames))

	// Use worker pattern for controlled concurrency
	jobsChan := make(chan int, len(spellNames))

	// Start workers (limit concurrent requests)
	numWorkers := 3 // Limit concurrent requests to be nice to the API
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobsChan {
				spell, err := c.GetSpell(spellNames[idx])
				results[idx] = BatchResult{
					Name:  spellNames[idx],
					Data:  spell,
					Error: err,
				}
			}
		}()
	}

	// Send jobs
	for i := range spellNames {
		jobsChan <- i
	}
	close(jobsChan)

	// Wait for completion
	wg.Wait()

	return results
}

// GetEquipmentBatch fetches multiple equipment items concurrently with rate limiting
func (c *Client) GetEquipmentBatch(equipmentNames []string) []BatchResult {
	results := make([]BatchResult, len(equipmentNames))

	// Use worker pattern for controlled concurrency
	jobsChan := make(chan int, len(equipmentNames))

	// Start workers (limit concurrent requests)
	numWorkers := 3 // Limit concurrent requests to be nice to the API
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobsChan {
				equipment, err := c.GetEquipment(equipmentNames[idx])
				results[idx] = BatchResult{
					Name:  equipmentNames[idx],
					Data:  equipment,
					Error: err,
				}
			}
		}()
	}

	// Send jobs
	for i := range equipmentNames {
		jobsChan <- i
	}
	close(jobsChan)

	// Wait for completion
	wg.Wait()

	return results
}
