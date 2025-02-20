package api

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/grumpzsux/goParams/internal/config"
)

// FetchFunc defines the signature for API fetching functions.
type FetchFunc func(ctx context.Context, domain string, cfg *config.Config) ([]string, error)

// FetchAll queries all available data sources concurrently and returns a deduplicated list of URLs.
// If one source fails (for example, Wayback times out), its error is logged as a warning while continuing with the others.
func FetchAll(ctx context.Context, domain string, cfg *config.Config) ([]string, error) {
	apis := []FetchFunc{
		FetchWayback,     // Extended timeout is handled within FetchWayback.
		FetchCommonCrawl, // Likewise, these functions accept context.
		FetchVirusTotal,
		FetchAlienVault,
	}

	var wg sync.WaitGroup
	urlCh := make(chan []string)
	errCh := make(chan error, len(apis))

	for _, apiFunc := range apis {
		wg.Add(1)
		go func(f FetchFunc) {
			defer wg.Done()
			urls, err := f(ctx, domain, cfg)
			if err != nil {
				errCh <- err
				return
			}
			urlCh <- urls
		}(apiFunc)
	}

	// Close channels once all goroutines have finished.
	go func() {
		wg.Wait()
		close(urlCh)
		close(errCh)
	}()

	// Collect results.
	urlSet := make(map[string]struct{})
	for urls := range urlCh {
		for _, u := range urls {
			urlSet[u] = struct{}{}
		}
	}

	// Log any errors encountered.
	for err := range errCh {
		logrus.Warnf("An API error occurred: %v", err)
	}

	// Convert set to slice.
	var results []string
	for u := range urlSet {
		results = append(results, u)
	}
	return results, nil
}
