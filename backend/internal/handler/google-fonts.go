package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const GOOGLE_FONTS_API_URL = "https://www.googleapis.com/webfonts/v1/webfonts"

// FontInfo represents a font with its metadata
type FontInfo struct {
	Family   string   `json:"family"`
	Category string   `json:"category"`
	Variants []string `json:"variants"`
	Variable bool     `json:"variable"`
}

// GoogleFont represents the structure from Google Fonts API
type GoogleFont struct {
	Family   string   `json:"family"`
	Category string   `json:"category"`
	Variants []string `json:"variants"`
}

// GoogleFontsAPIResponse represents the response from Google Fonts API
type GoogleFontsAPIResponse struct {
	Items []GoogleFont `json:"items"`
}

// PaginatedFontsResponse represents the paginated response
type PaginatedFontsResponse struct {
	Fonts   []FontInfo `json:"fonts"`
	Total   int        `json:"total"`
	Offset  int        `json:"offset"`
	Limit   int        `json:"limit"`
	HasMore bool       `json:"hasMore"`
}

// Cache structure
var (
	fontsCache      []FontInfo
	fontsCacheMutex sync.RWMutex
	fontsCacheTime  time.Time
	cacheDuration   = 24 * time.Hour // Cache for 24 hours
)

// Fallback fonts if API fails
var FALLBACK_FONTS = []FontInfo{
	// Sans Serif fonts
	{Family: "Inter", Category: "sans-serif", Variants: []string{"400", "600", "700"}, Variable: true},
	{Family: "Roboto", Category: "sans-serif", Variants: []string{"400", "600", "700"}, Variable: false},
	{Family: "Open Sans", Category: "sans-serif", Variants: []string{"400", "600", "700"}, Variable: true},
	{Family: "Poppins", Category: "sans-serif", Variants: []string{"400", "600", "700"}, Variable: false},
	{Family: "Montserrat", Category: "sans-serif", Variants: []string{"400", "600", "700"}, Variable: false},
	{Family: "Outfit", Category: "sans-serif", Variants: []string{"400", "600", "700"}, Variable: false},
	{Family: "Plus Jakarta Sans", Category: "sans-serif", Variants: []string{"400", "600", "700"}, Variable: false},
	{Family: "DM Sans", Category: "sans-serif", Variants: []string{"400", "600", "700"}, Variable: true},
	{Family: "Geist", Category: "sans-serif", Variants: []string{"400", "600", "700"}, Variable: true},
	{Family: "Oxanium", Category: "sans-serif", Variants: []string{"400", "600", "700"}, Variable: false},
	{Family: "Architects Daughter", Category: "handwriting", Variants: []string{"400", "600", "700"}, Variable: false},
	// Serif fonts
	{Family: "Merriweather", Category: "serif", Variants: []string{"400", "600", "700"}, Variable: false},
	{Family: "Playfair Display", Category: "serif", Variants: []string{"400", "600", "700"}, Variable: true},
	{Family: "Lora", Category: "serif", Variants: []string{"400", "600", "700"}, Variable: true},
	{Family: "Source Serif Pro", Category: "serif", Variants: []string{"400", "600", "700"}, Variable: false},
	{Family: "Libre Baskerville", Category: "serif", Variants: []string{"400", "600", "700"}, Variable: false},
	{Family: "Space Grotesk", Category: "serif", Variants: []string{"400", "600", "700"}, Variable: false},
	// Monospace fonts
	{Family: "JetBrains Mono", Category: "monospace", Variants: []string{"400", "600", "700"}, Variable: true},
	{Family: "Fira Code", Category: "monospace", Variants: []string{"400", "600", "700"}, Variable: true},
	{Family: "Source Code Pro", Category: "monospace", Variants: []string{"400", "600", "700"}, Variable: false},
	{Family: "IBM Plex Mono", Category: "monospace", Variants: []string{"400", "600", "700"}, Variable: false},
	{Family: "Roboto Mono", Category: "monospace", Variants: []string{"400", "600", "700"}, Variable: false},
	{Family: "Space Mono", Category: "monospace", Variants: []string{"400", "600", "700"}, Variable: false},
	{Family: "Geist Mono", Category: "monospace", Variants: []string{"400", "600", "700"}, Variable: true},
}

// fetchGoogleFonts fetches fonts from Google Fonts API
func fetchGoogleFonts(apiKey string) ([]FontInfo, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Google Fonts API key is required")
	}

	url := fmt.Sprintf("%s?key=%s", GOOGLE_FONTS_API_URL, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Google Fonts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google Fonts API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var apiResponse GoogleFontsAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Transform to our format
	fonts := make([]FontInfo, len(apiResponse.Items))
	for i, font := range apiResponse.Items {
		variable := false
		for _, variant := range font.Variants {
			if strings.Contains(variant, "wght") || strings.Contains(variant, "ital,wght") {
				variable = true
				break
			}
		}

		fonts[i] = FontInfo{
			Family:   font.Family,
			Category: font.Category,
			Variants: font.Variants,
			Variable: variable,
		}
	}

	fmt.Printf("✅ Fetched %d fonts from Google Fonts API\n", len(fonts))
	return fonts, nil
}

// getCachedGoogleFonts returns cached fonts or fetches new ones
func getCachedGoogleFonts(apiKey string) []FontInfo {
	fontsCacheMutex.RLock()
	if len(fontsCache) > 0 && time.Since(fontsCacheTime) < cacheDuration {
		fonts := fontsCache
		fontsCacheMutex.RUnlock()
		return fonts
	}
	fontsCacheMutex.RUnlock()

	// Try to fetch from API
	fonts, err := fetchGoogleFonts(apiKey)
	if err != nil {
		fmt.Printf("Error fetching Google Fonts: %v\n", err)
		fmt.Println("Using fallback fonts")
		return FALLBACK_FONTS
	}

	// Update cache
	fontsCacheMutex.Lock()
	fontsCache = fonts
	fontsCacheTime = time.Now()
	fontsCacheMutex.Unlock()

	return fonts
}

// GetGoogleFonts returns the handler for Google Fonts API
func GetGoogleFonts() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get query parameters
		query := strings.ToLower(c.DefaultQuery("q", ""))
		category := strings.ToLower(c.DefaultQuery("category", ""))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

		// Limit max to 100
		if limit > 100 {
			limit = 100
		}

		// Get API key from environment
		apiKey := os.Getenv("GOOGLE_FONTS_API_KEY")

		// Get fonts (cached or fetched)
		googleFonts := getCachedGoogleFonts(apiKey)

		// Filter fonts based on query and category
		filteredFonts := make([]FontInfo, 0)
		for _, font := range googleFonts {
			// Filter by query
			if query != "" && !strings.Contains(strings.ToLower(font.Family), query) {
				continue
			}

			// Filter by category
			if category != "" && category != "all" && strings.ToLower(font.Category) != category {
				continue
			}

			filteredFonts = append(filteredFonts, font)
		}

		// Paginate
		total := len(filteredFonts)
		end := offset + limit
		if end > total {
			end = total
		}

		paginatedFonts := []FontInfo{}
		if offset < total {
			paginatedFonts = filteredFonts[offset:end]
		}

		// Build response
		response := PaginatedFontsResponse{
			Fonts:   paginatedFonts,
			Total:   total,
			Offset:  offset,
			Limit:   limit,
			HasMore: offset+limit < total,
		}

		c.JSON(http.StatusOK, response)
	}
}
