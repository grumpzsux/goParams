package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/grumpzsux/goParams/internal/config"
)

// VirusTotalResponse represents a simplified structure for the VirusTotal domain report.
type VirusTotalResponse struct {
	DetectedURLs   []struct {
		URL string `json:"url"`
	} `json:"detected_urls"`
	UndetectedURLs [][]interface{} `json:"undetected_urls"`
}

// FetchVirusTotal fetches URLs from VirusTotal for the given domain.
func FetchVirusTotal(ctx context.Context, domain string, cfg *config.Config) ([]string, error) {
	if cfg.VirusTotalAPIKey == "" {
		color.Yellow("No VirusTotal API key provided. Skipping VirusTotal lookup for %s", domain)
		return nil, nil
	}
	apiURL := fmt.Sprintf("https://www.virustotal.com/vtapi/v2/domain/report?apikey=%s&domain=%s", cfg.VirusTotalAPIKey, domain)
	color.Blue("[*] Fetching from VirusTotal for domain: %s", domain)

	resp, err := GetWithRandomUA(ctx, apiURL, cfg)
	if err != nil {
		return nil, fmt.Errorf("error fetching from VirusTotal: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("VirusTotal returned status code %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading VirusTotal response: %w", err)
	}

	var vtResp VirusTotalResponse
	if err := json.Unmarshal(bodyBytes, &vtResp); err != nil {
		return nil, fmt.Errorf("error parsing VirusTotal JSON: %w", err)
	}

	urlSet := make(map[string]struct{})
	for _, entry := range vtResp.DetectedURLs {
		if entry.URL != "" && strings.Contains(entry.URL, "?") {
			urlSet[entry.URL] = struct{}{}
		}
	}
	for _, arr := range vtResp.UndetectedURLs {
		if len(arr) > 0 {
			if urlStr, ok := arr[0].(string); ok && strings.Contains(urlStr, "?") {
				urlSet[urlStr] = struct{}{}
			}
		}
	}

	var results []string
	for u := range urlSet {
		results = append(results, u)
	}
	return results, nil
}
