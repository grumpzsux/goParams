package api

import (
	"context"
	"net/http"
	"time"

	"github.com/grumpzsux/goParams/internal/config"
	"github.com/grumpzsux/goParams/internal/utils"
)

// HTTPClient is a shared HTTP client with a timeout.
var HTTPClient = &http.Client{
	Timeout: 15 * time.Second,
}

// GetWithRandomUA creates an HTTP GET request with a random User-Agent header from the configuration.
func GetWithRandomUA(ctx context.Context, url string, cfg *config.Config) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Pick a random user agent from the configuration.
	ua := utils.RandomStringFromSlice(cfg.UserAgents)
	if ua == "" {
		ua = "Mozilla/5.0 (compatible)"
	}
	req.Header.Set("User-Agent", ua)
	return HTTPClient.Do(req)
}
