# goParams

goParams is a fast, robust, and concurrent tool for harvesting parameterized URLs from various data sources—including the Wayback Machine, Common Crawl, VirusTotal, and AlienVault OTX. It is designed for penetration testers, bug bounty hunters, and security researchers who need to quickly collect and filter URL parameters for further analysis.
<p align="center">
<img src="https://github.com/user-attachments/assets/455b3ef2-d35e-4277-809b-78958f45225a" width=500 height=300>
</p>

## Features

- **Multiple Data Sources:** Harvest URLs from the Wayback Machine, Common Crawl, VirusTotal, and AlienVault OTX.
- **Concurrent Processing:** Dynamically control the number of concurrent API requests to avoid rate limiting.
- **Configurable Output:** Output results in plain text or JSON format. Optionally save results to a file.
- **Robust URL Cleaning:** Filters out URLs with unwanted file extensions (e.g., images, fonts, videos, etc.) and cleans URLs by removing redundant port information and replacing query parameter values with a user‑defined placeholder.
- **Customizable Canaries:** Specify a custom placeholder for query parameter values via the `--canary` flag.
- **Structured Logging & Verbosity:** Uses [logrus](https://github.com/sirupsen/logrus) for consistent logging with configurable verbosity levels.
- **Context-Aware & Timeout Handling:** Implements context-based cancellation and extended timeouts (e.g., for slow responses from the Wayback Machine).

## Installation

### Prerequisites

- [Go](https://golang.org/) 1.18 or later
- Git

### Install with Go
Run the following command to install goParams:
```bash
go install github.com/grumpzsux/goParams/cmd/goParams@latest
```
Then ensure that the installed binary is in your PATH (for example, by adding $(go env GOPATH)/bin to your PATH). Once that’s done, you can run the command below to verify the installation and see the help output:
```bash
goParams --help
```

### Manually Build
This creates an executable named goParams in your project root.
```bash
git clone https://github.com/yourusername/goParams.git
cd goParams
go build -o goParams ./cmd/goParams
```

## Usage

### Command-Line Flags
```yaml
Usage:
  goParams [flags]

Flags:
  -c, --concurrency int        Number of concurrent API requests (default 5)
  -d, --domain string          Target domain (e.g., example.com)
  -f, --output-format string   Output format: plain or json (default "plain")
  -l, --list string            File containing a list of domains/subdomains
  -o, --output string          Output file for results (if not provided, prints to stdout)
      --canary string          Custom placeholder for URL query parameters when cleaning URLs (default "PLACEHOLDER")
  -v, --verbose                Enable verbose logging
  --config string               Path to configuration file (default "config.yaml")
  -h, --help                   help for goParams
```

![image](https://github.com/user-attachments/assets/50b38bed-5a5f-4751-b967-461824e38ad7)

### Examples
- **Harvest URLs for a Single Domain**
```bash
./goParams -d example.com
```
- **Harvest URLs for a List of Domains**
```bash
./goParams -l domains.txt
```
- **Use Verbose Logging and Output JSON**
```bash
./goParams -d example.com -v -f json
```
- **Save Results to a File with a Custom Canary**
```bash
./goParams -d example.com --canary "MYCANARY" -o results.txt
```
## Configuration
goParams uses a YAML configuration file for API keys and other settings. By default, it looks for `config.yaml` in the project root.

**Example `config.yaml`**
```yaml
virustotal_api_key: "YOUR_VIRUSTOTAL_API_KEY"
alienvault_api_key: "YOUR_ALIENVAULT_API_KEY"
concurrency: 5
user_agents:
  - "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)"
  - "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko)"
  - "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko)"
rate_limit: 60
```
- **virustotal_api_key:** Your VirusTotal API key.
- **alienvault_api_key:** Your AlienVault OTX API key.
- **concurrency:** Default number of concurrent API requests.
- **user_agents:** Custom list of user agent strings for rotating requests.
- **rate_limit:** (Optional) API rate limit settings.

## Contributing
Contributions are welcome! Please follow these steps:

- Fork the repository.
- Create a feature branch: `git checkout -b my-feature`.
- Commit your changes: `git commit -m "Add some feature"`.
- Push your branch: `git push origin my-feature`.
- Open a pull request describing your changes.

Please ensure your code follows Go conventions and passes all tests.

## Contributers
- Shout out to `xnl-h4ck3r` the <a href="https://github.com/xnl-h4ck3r/waymore">waymore</a> creator for the inspiration.
- Shout out to `Devansh Batham` the <a href="https://github.com/devanshbatham/ParamSpider">ParamSpider</a> creator for the inspiration.




