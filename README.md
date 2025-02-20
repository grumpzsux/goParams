# goParams

goParams is a fast, robust, and concurrent tool for harvesting parameterized URLs from various data sources—including the Wayback Machine, Common Crawl, VirusTotal, and AlienVault OTX. It is designed for penetration testers, bug bounty hunters, and security researchers who need to quickly collect and filter URL parameters for further analysis.

![goParams Banner](./assets/banner.png)

## Features

- **Multiple Data Sources:**  
  Harvest URLs from the Wayback Machine, Common Crawl, VirusTotal, and AlienVault OTX.
  
- **Concurrent Processing:**  
  Dynamically control the number of concurrent API requests to avoid rate limiting.
  
- **Configurable Output:**  
  Output results in plain text or JSON format. Optionally save results to a file.
  
- **Robust URL Cleaning:**  
  Filters out URLs with unwanted file extensions (e.g., images, fonts, videos, etc.) and cleans URLs by removing redundant port information and replacing query parameter values with a user‑defined placeholder.
  
- **Customizable Placeholder:**  
  Specify a custom placeholder for query parameter values via the `--canary` flag.
  
- **Structured Logging & Verbosity:**  
  Uses [logrus](https://github.com/sirupsen/logrus) for consistent logging with configurable verbosity levels.
  
- **Context-Aware & Timeout Handling:**  
  Implements context-based cancellation and extended timeouts (e.g., for slow responses from the Wayback Machine).

## Installation

### Prerequisites

- [Go](https://golang.org/) 1.18 or later
- Git

### Clone the Repository

```bash
git clone https://github.com/yourusername/goParams.git
cd goParams
```