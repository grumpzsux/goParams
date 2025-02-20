package api

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/fatih/color"
	"github.com/grumpzsux/goParams/internal/config"
)

// CommonCrawlEntry represents one record from the Common Crawl index.
type CommonCrawlEntry struct {
	Timestamp string `json:"timestamp"`
	URL       string `json:"url"`
	Mime      string `json:"mime"`
	Status    string `json:"status"`
	Digest    string `json:"digest"`
}

// BaseIndexURL is the Common Crawl index URL. (This may be updated periodically.)
const BaseIndexURL = "http://index.commoncrawl.org/CC-MAIN-2019-51-index"

// FetchCommonCrawl queries the Common Crawl index for the given domain and returns URLs with parameters.
func FetchCommonCrawl(ctx context.Context, domain string, cfg *config.Config) ([]string, error) {
	// Define filters (exclude "warc/revisit" and status 404).
	filterMIME := "&filter=!~mime:(warc/revisit)"
	filterCode := "&filter=!~status:(404)"
	filterKeywords := ""

	escapedDomain := url.QueryEscape(domain + "/*")
	queryParams := fmt.Sprintf("?output=json&fl=timestamp,url,mime,status,digest&url=%s", escapedDomain)
	fullURL := BaseIndexURL + queryParams + filterMIME + filterCode + filterKeywords
	color.Blue("[*] Fetching from Common Crawl: %s", fullURL)

	resp, err := GetWithRandomUA(ctx, fullURL, cfg)
	if err != nil {
		return nil, fmt.Errorf("error fetching from Common Crawl: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return nil, errors.New("Common Crawl rate limit reached (429)")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Common Crawl returned status code: %d", resp.StatusCode)
	}

	reader := bufio.NewReader(resp.Body)
	var urlsFound []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			color.Yellow("Error reading a line: %v", err)
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			if err == io.EOF {
				break
			}
			continue
		}
		// If the line indicates no captures, exit early.
		if strings.Contains(strings.ToLower(line), "no captures found") {
			color.Yellow("No captures found for %s", domain)
			break
		}
		var entry CommonCrawlEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			color.Yellow("Failed to parse JSON from line: %v", err)
			if err == io.EOF {
				break
			}
			continue
		}
		if strings.Contains(entry.URL, "?") {
			urlsFound = append(urlsFound, entry.URL)
		}
		if err == io.EOF {
			break
		}
	}
	return urlsFound, nil
}
