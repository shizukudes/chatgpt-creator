package register

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

// CookieExport represents a cookie in JSON format
type CookieExport struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Domain   string `json:"domain"`
	Path     string `json:"path"`
	Expires  string `json:"expires,omitempty"`
	Secure   bool   `json:"secure"`
	HttpOnly bool   `json:"httpOnly"`
	SameSite string `json:"sameSite,omitempty"`
}

// SaveCookiesToFile saves cookies to a JSON file
func SaveCookiesToFile(cookies []*http.Cookie, email, outputDir string) error {
	if outputDir == "" {
		outputDir = "cookies"
	}
	
	// Create cookies directory if not exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create cookies directory: %w", err)
	}
	
	// Convert cookies to export format
	var exports []CookieExport
	for _, cookie := range cookies {
		// Skip sameSite field for better compatibility with browser extensions
		// (especially mobile versions like Firefox Android Cookie-Editor)
		
		expires := ""
		if !cookie.Expires.IsZero() {
			expires = cookie.Expires.Format(time.RFC3339)
		}
		
		exports = append(exports, CookieExport{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Expires:  expires,
			Secure:   cookie.Secure,
			HttpOnly: cookie.HttpOnly,
			// SameSite omitted for compatibility
		})
	}
	
	// Generate filename based on email
	filename := filepath.Join(outputDir, fmt.Sprintf("%s.json", email))
	
	// Marshal to JSON
	data, err := json.MarshalIndent(exports, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cookies: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write cookie file: %w", err)
	}
	
	return nil
}

// SaveCookiesToFileBrowserFriendly saves only chatgpt.com cookies (for browser import without cross-domain errors)
func SaveCookiesToFileBrowserFriendly(cookies []*http.Cookie, email, outputDir string) error {
	if outputDir == "" {
		outputDir = "cookies"
	}
	
	// Filter only chatgpt.com cookies
	var chatgptCookies []*http.Cookie
	for _, cookie := range cookies {
		if cookie.Domain == "chatgpt.com" || cookie.Domain == "" {
			chatgptCookies = append(chatgptCookies, cookie)
		}
	}
	
	// Create cookies directory if not exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create cookies directory: %w", err)
	}
	
	// Convert cookies to export format (without sameSite for compatibility)
	var exports []CookieExport
	for _, cookie := range chatgptCookies {
		expires := ""
		if !cookie.Expires.IsZero() {
			expires = cookie.Expires.Format(time.RFC3339)
		}
		
		exports = append(exports, CookieExport{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Expires:  expires,
			Secure:   cookie.Secure,
			HttpOnly: cookie.HttpOnly,
		})
	}
	
	// Generate filename
	filename := filepath.Join(outputDir, fmt.Sprintf("%s-browser.json", email))
	
	// Marshal to JSON
	data, err := json.MarshalIndent(exports, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cookies: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write cookie file: %w", err)
	}
	
	return nil
}
