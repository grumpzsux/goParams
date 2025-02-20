package utils

import (
	"net/url"
	"path"
	"strings"
)

// HardcodedExtensions are the file extensions that should be ignored.
var HardcodedExtensions = []string{
	".jpg", ".jpeg", ".png", ".gif", ".pdf", ".svg", ".json",
	".css", ".js", ".webp", ".woff", ".woff2", ".eot", ".ttf", ".otf", ".mp4", ".txt",
}

// HasExtension checks if the provided URL has a file extension matching any of the given extensions.
func HasExtension(rawURL string, extensions []string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		// If URL parsing fails, assume no valid extension is present.
		return false
	}
	ext := strings.ToLower(path.Ext(u.Path))
	for _, e := range extensions {
		if ext == e {
			return true
		}
	}
	return false
}

// CleanURL cleans a URL by removing redundant port information.
// For example, if the scheme is "http" and the port is 80 (or "https" with port 443), the port is removed.
func CleanURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		// On error, return the original URL.
		return rawURL
	}

	host := u.Hostname()
	port := u.Port()
	if (u.Scheme == "http" && port == "80") || (u.Scheme == "https" && port == "443") {
		u.Host = host
	}
	return u.String()
}

// CleanURLs processes a list of URLs:
//  1. It first cleans each URL (removing redundant port info).
//  2. Then it skips any URLs whose path has one of the excluded file extensions.
//  3. For the remaining URLs, it replaces each query parameter's value with the provided placeholder.
//  4. Finally, it returns a deduplicated list of cleaned URLs.
func CleanURLs(urls []string, extensions []string, placeholder string) []string {
	cleanedSet := make(map[string]struct{})

	for _, rawURL := range urls {
		// Clean the URL.
		cleanedURL := CleanURL(rawURL)
		// Skip URL if it has one of the hardcoded unwanted extensions.
		if HasExtension(cleanedURL, extensions) {
			continue
		}

		u, err := url.Parse(cleanedURL)
		if err != nil {
			// If URL parsing fails, keep the original cleaned URL.
			cleanedSet[cleanedURL] = struct{}{}
			continue
		}

		// Replace each query parameter's value with the placeholder.
		q := u.Query()
		for key := range q {
			q.Set(key, placeholder)
		}
		u.RawQuery = q.Encode()

		cleanedSet[u.String()] = struct{}{}
	}

	// Convert the set to a slice.
	var result []string
	for s := range cleanedSet {
		result = append(result, s)
	}
	return result
}
