package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/grumpzsux/goParams/internal/config"
)

// AlienVaultResponse represents the JSON response structure from Alien Vault OTX.
type AlienVaultResponse struct {
	FullSize int `json:"full_size"`
	UrlList  []struct {
		URL      string      `json:"url"`
		HTTPCode interface{} `json:"httpcode"` // May be int or string.
	} `json:"url_list"`
}

// BaseAlienVaultURL is the API endpoint template.
const BaseAlienVaultURL = "https://otx.alienvault.com/api/v1/indicators/{TYPE}/{DOMAIN}/url_list?limit=500"

// getIndicatorType returns "hostname" if the domain appears to be a subdomain.
func getIndicatorType(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) > 2 {
		return "hostname"
	}
	return "domain"
}

// FetchAlienVault retrieves URLs from Alien Vault OTX for the given domain.
func FetchAlienVault(ctx context.Context, domain string, cfg *config.Config) ([]string, error) {
	if cfg.AlienVaultAPIKey == "" {
		color.Yellow("No Alien Vault API key provided. Skipping Alien Vault lookup for %s", domain)
		return nil, nil
	}
	indicatorType := getIndicatorType(domain)
	escapedDomain := url.QueryEscape(domain)
	baseURL := strings.Replace(BaseAlienVaultURL, "{TYPE}", indicatorType, 1)
	baseURL = strings.Replace(baseURL, "{DOMAIN}", escapedDomain, 1)

	// Get total pages by requesting the showNumPages parameter.
	initialURL := baseURL + "&showNumPages=True"
	color.Blue("[*] Fetching Alien Vault page count from: %s", initialURL)

	resp, err := GetWithRandomUA(ctx, initialURL, cfg)
	if err != nil {
		return nil, fmt.Errorf("error fetching Alien Vault initial page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, errors.New("Alien Vault rate limit reached (429)")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Alien Vault returned status code %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var initResp AlienVaultResponse
	if err := json.Unmarshal(bodyBytes, &initResp); err != nil {
		return nil, fmt.Errorf("error parsing Alien Vault initial JSON: %w", err)
	}

	totalURLs := initResp.FullSize
	if totalURLs == 0 {
		color.Yellow("Alien Vault returned zero results for %s", domain)
		return nil, nil
	}

	totalPages := int(math.Ceil(float64(totalURLs) / 500.0))
	color.Blue("[*] Alien Vault reports %d results over %d pages", totalURLs, totalPages)

	var wg sync.WaitGroup
	var mu sync.Mutex
	urlSet := make(map[string]struct{})
	for page := 1; page <= totalPages; page++ {
		wg.Add(1)
		pageURL := baseURL + "&page=" + fmt.Sprintf("%d", page)
		go func(pageURL string) {
			defer wg.Done()
			pageURLs, err := processAlienVaultPage(ctx, pageURL, domain, cfg)
			if err != nil {
				color.Yellow("Error processing Alien Vault page %s: %v", pageURL, err)
				return
			}
			mu.Lock()
			for _, u := range pageURLs {
				urlSet[u] = struct{}{}
			}
			mu.Unlock()
		}(pageURL)
	}
	wg.Wait()

	var results []string
	for u := range urlSet {
		results = append(results, u)
	}
	return results, nil
}

// processAlienVaultPage makes a request to the provided page URL and returns a slice of valid URLs.
func processAlienVaultPage(ctx context.Context, pageURL, targetDomain string, cfg *config.Config) ([]string, error) {
	color.Blue("[*] Processing Alien Vault page: %s", pageURL)
	resp, err := GetWithRandomUA(ctx, pageURL, cfg)
	if err != nil {
		return nil, fmt.Errorf("error requesting Alien Vault page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, errors.New("Alien Vault rate limit reached (429)")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Alien Vault page returned status code %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(bodyBytes) == 0 {
		color.Yellow("Alien Vault page %s returned an empty response", pageURL)
		return nil, nil
	}

	var avResp AlienVaultResponse
	if err := json.Unmarshal(bodyBytes, &avResp); err != nil {
		return nil, fmt.Errorf("error parsing Alien Vault page JSON: %w", err)
	}

	var urlsFound []string
	for _, entry := range avResp.UrlList {
		foundURL := entry.URL
		if foundURL == "" {
			continue
		}
		// Basic filtering: include only URLs with query parameters and matching target domain.
		if strings.Contains(foundURL, "?") && strings.Contains(strings.ToLower(foundURL), strings.ToLower(targetDomain)) {
			urlsFound = append(urlsFound, foundURL)
		}
	}
	return urlsFound, nil
}
