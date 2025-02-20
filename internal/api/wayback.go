package api

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/grumpzsux/goParams/internal/config"
)

// WayBackException is returned when the Wayback Machine response indicates an error.
type WayBackException struct {
	Message string
}

func (w *WayBackException) Error() string {
	return w.Message
}

// fixArchiveOrgUrl removes any trailing "%0A" or "%0a" from the provided URL.
func fixArchiveOrgUrl(urlStr string) string {
	lower := strings.ToLower(urlStr)
	if idx := strings.Index(lower, "%0a"); idx > 0 {
		return urlStr[:idx]
	}
	return urlStr
}

// FetchWayback queries the Wayback Machine CDX API for archived URLs of the given domain.
// It uses an extended timeout (e.g. 2 minutes) so that large datasets can load.
// Returns a deduplicated slice of original URLs that include query parameters.
func FetchWayback(ctx context.Context, domain string, cfg *config.Config) ([]string, error) {
	// For simplicity, we use a default collapse value.
	collapse := "/*"
	apiURL := fmt.Sprintf("https://web.archive.org/cdx/search/cdx?url=%s%s&fl=timestamp,original,mimetype,statuscode,digest", domain, collapse)
	color.Blue("[*] Fetching from Wayback Machine: %s", apiURL)

	// Create a new context with an extended timeout for Wayback.
	waybackTimeout := 2 * time.Minute
	waybackCtx, cancel := context.WithTimeout(ctx, waybackTimeout)
	defer cancel()

	// Use the shared HTTP client with the extended context.
	resp, err := GetWithRandomUA(waybackCtx, apiURL, cfg)
	if err != nil {
		return nil, fmt.Errorf("error fetching from Wayback: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Wayback Machine returned status code %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading Wayback response: %w", err)
	}
	bodyStr := string(bodyBytes)
	lowerBody := strings.ToLower(bodyStr)
	if strings.Contains(lowerBody, "wayback machine has not archived that url") ||
		strings.Contains(lowerBody, "snapshot cannot be displayed due to an internal error") {
		return nil, &WayBackException{Message: "Wayback Machine returned an error response"}
	}

	scanner := bufio.NewScanner(strings.NewReader(bodyStr))
	urlSet := make(map[string]struct{})
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		originalURL := fixArchiveOrgUrl(fields[1])
		if strings.Contains(originalURL, "?") {
			urlSet[originalURL] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning Wayback response: %w", err)
	}

	var results []string
	for url := range urlSet {
		results = append(results, url)
	}
	return results, nil
}
