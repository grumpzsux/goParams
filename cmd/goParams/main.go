package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/grumpzsux/goParams/internal/api"
	"github.com/grumpzsux/goParams/internal/config"
	"github.com/grumpzsux/goParams/internal/logger"
	"github.com/grumpzsux/goParams/internal/utils"
)

var (
	cfgFile      string
	verbose      bool
	concurrency  int
	outputFormat string
	domain       string
	domainList   string
	placeholder  string // Canary placeholder for cleaning URLs.
	outputFile   string // New flag for output file.
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "goParams",
		Short: "goParams is a parameterized URL harvester",
		Long:  "goParams is a robust tool for harvesting parameterized URLs from various data sources.",
		Run: func(cmd *cobra.Command, args []string) {
			printBanner()
			runApp(args)
		},
	}

	// Persistent flags.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Path to configuration file (default is config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", 5, "Number of concurrent API requests")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output-format", "f", "plain", "Output format: plain or json")

	// Command-specific flags.
	rootCmd.Flags().StringVarP(&domain, "domain", "d", "", "Target domain (e.g., example.com)")
	rootCmd.Flags().StringVarP(&domainList, "list", "l", "", "File containing a list of domains/subdomains")
	rootCmd.Flags().StringVar(&placeholder, "canary", "PLACEHOLDER", "Custom placeholder for URL query parameters when cleaning URLs")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file for results (if not provided, results are printed to stdout)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runApp(args []string) {
	// Initialize logger with the chosen verbosity level.
	logger.Init(verbose)
	logrus.Info("Starting goParams...")

	// Load configuration.
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}
	if err := config.Validate(cfg); err != nil {
		logrus.Fatalf("Invalid configuration: %v", err)
	}
	// Override concurrency if provided from CLI.
	cfg.Concurrency = concurrency

	// Collect target domains.
	var domains []string
	if domain != "" {
		domains = append(domains, domain)
	}
	if domainList != "" {
		list, err := utils.LoadDomainList(domainList)
		if err != nil {
			logrus.Fatalf("Error reading domain list: %v", err)
		}
		domains = append(domains, list...)
	}
	if len(domains) == 0 {
		logrus.Fatal("No domains provided. Use -d or -l flag to supply target domains.")
	}

	// Create a cancellable context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create a semaphore channel for dynamic concurrency.
	sem := make(chan struct{}, cfg.Concurrency)
	var wg sync.WaitGroup
	results := make(map[string][]string)
	var mu sync.Mutex

	// Process each domain concurrently.
	for _, d := range domains {
		wg.Add(1)
		sem <- struct{}{}
		go func(target string) {
			defer wg.Done()
			logrus.Infof("Processing domain: %s", target)
			// Query all available data sources concurrently.
			urls, err := api.FetchAll(ctx, target, cfg)
			if err != nil {
				logrus.Errorf("Error fetching URLs for %s: %v", target, err)
			}
			// Clean the URLs using the user-defined placeholder.
			cleanedURLs := utils.CleanURLs(urls, utils.HardcodedExtensions, placeholder)
			mu.Lock()
			results[target] = cleanedURLs
			mu.Unlock()
			<-sem
		}(d)
	}
	wg.Wait()

	// Output results.
	if outputFile != "" {
		// Write results to file.
		if err := utils.WriteResultsToFile(outputFile, results, outputFormat); err != nil {
			logrus.Errorf("Failed to write output file: %v", err)
		} else {
			logrus.Infof("Output written to %s", outputFile)
		}
	} else {
		// Print results to stdout.
		if outputFormat == "json" {
			utils.OutputJSON(results)
		} else {
			utils.OutputPlain(results)
		}
	}
}

// printBanner displays an ASCII banner at startup.
func printBanner() {
	banner := `
              __________
   ____   ____\______   \_____ ____________    _____   ______
  / ___\ /  _ \|     ___/\__  \\_  __ \__  \  /     \ /  ___/
 / /_/  >  <_> )    |     / __ \|  | \// __ \|  Y Y  \\___ \
 \___  / \____/|____|    (____  /__|  (____  /__|_|  /____  >
/_____/                       \/           \/      \/     \/ [v1.0]
             goParams - Parameterized URL Harvester.
                         [@GrumpzSux]
`
	fmt.Println(banner)
}
