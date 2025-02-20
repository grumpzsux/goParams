package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// LoadDomainList reads a file line-by-line and returns a slice of non-empty, trimmed domain strings.
func LoadDomainList(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var domains []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			domains = append(domains, line)
		}
	}
	return domains, scanner.Err()
}

// WriteOutput writes the provided URLs to the specified file, one URL per line.
func WriteOutput(filename string, urls []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, url := range urls {
		_, err := writer.WriteString(url + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

// WriteResultsToFile writes the aggregated results to a file in the specified format.
// 'results' is a map from domain to a slice of URLs.
// 'format' can be "json" or "plain".
func WriteResultsToFile(filename string, results map[string][]string, format string) error {
	var content string
	if format == "json" {
		b, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return fmt.Errorf("error formatting JSON output: %w", err)
		}
		content = string(b)
	} else {
		var sb strings.Builder
		for domain, urls := range results {
			sb.WriteString(fmt.Sprintf("Domain: %s\n", domain))
			for _, url := range urls {
				sb.WriteString(url + "\n")
			}
			sb.WriteString("\n")
		}
		content = sb.String()
	}
	return os.WriteFile(filename, []byte(content), 0644)
}

// HumanReadableSize converts a byte count into a human-readable string format.
func HumanReadableSize(bytes uint64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	exp := int(math.Log(float64(bytes)) / math.Log(1024))
	pre := "KMGTPE"[exp-1]
	value := float64(bytes) / math.Pow(1024, float64(exp))
	return fmt.Sprintf("%.1f %cB", value, pre)
}

// GetMemoryUsage returns the current heap allocation (in bytes) by reading runtime statistics.
func GetMemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

// PrintProgressBar displays a simple progress bar in the terminal.
func PrintProgressBar(current, total int, prefix, suffix string, length int, fill string) {
	percent := float64(current) / float64(total)
	filledLength := int(math.Round(float64(length) * percent))
	bar := strings.Repeat(fill, filledLength) + strings.Repeat("-", length-filledLength)
	fmt.Printf("\r%s |%s| %d/%d %s", prefix, bar, current, total, suffix)
	if current >= total {
		fmt.Println()
	}
}

// RandomInt returns a random integer between min and max (inclusive).
func RandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

// ParseInt safely converts a string to an integer.
func ParseInt(s string, defaultVal int) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}

// RandomStringFromSlice returns a random element from a slice of strings.
func RandomStringFromSlice(slice []string) string {
	if len(slice) == 0 {
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	return slice[rand.Intn(len(slice))]
}

// OutputJSON prints the results in JSON format.
func OutputJSON(results map[string][]string) {
	b, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Printf("Error formatting JSON output: %v\n", err)
		return
	}
	fmt.Println(string(b))
}

// OutputPlain prints the results in plain text format.
func OutputPlain(results map[string][]string) {
	for domain, urls := range results {
		fmt.Printf("Domain: %s\n", domain)
		for _, url := range urls {
			fmt.Println(url)
		}
	}
}
