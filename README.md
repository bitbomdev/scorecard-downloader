# Scorecard Downloader

A tool to download and process OpenSSF Scorecard data for GitHub repositories.

## Overview

Scorecard Downloader is a command-line tool that allows you to download and process OpenSSF Scorecard data for specified GitHub repositories. It can handle multiple repositories at once and saves the results in a JSON file. The tool uses the deps.dev API to convert Package URLs (pURLs) to GitHub URLs before downloading the Scorecard data. It also supports using BigQuery to download the Scorecard data.

## Installation

```bash
go get github.com/bitbomdev/scorecard-downloader
```

## Usage

```bash
scorecard-downloader [global options] command [command options] [arguments...]
```

### Global Options

- `--purls value`: PURLs of the repositories to process (can be specified multiple times)
- `--purls-file value`: File containing PURLs, one per line
- `--output value`: Output file name (default: "results.json")
- `--use-bigquery`: Use BigQuery instead of the Scorecard API
- `--credentials-file value`: Path to the BigQuery credentials file
- `--help, -h`: Show help

### Examples

Process repositories specified directly:

```bash
scorecard-downloader --purls pkg:github/kubernetes/kubernetes --purls pkg:github/golang/go
```

Process repositories from a file:

```bash
scorecard-downloader --purls-file repos.txt
```

Specify custom output file:

```bash
scorecard-downloader --purls pkg:github/kubernetes/kubernetes --output custom_results.json
```

Use BigQuery with a credentials file:

```bash
scorecard-downloader --purls pkg:github/kubernetes/kubernetes --use-bigquery --credentials-file path/to/credentials.json
```

## Output

The tool saves the processed Scorecard data in a JSON file. By default, the output file is named `results.json`, but you can specify a custom name using the `--output` option.
