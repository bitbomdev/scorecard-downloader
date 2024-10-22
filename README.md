# Scorecard Downloader

# Scorecard Downloader

A tool to download and process OpenSSF Scorecard data for GitHub repositories.

## Overview

Scorecard Downloader is a command-line tool that allows you to download and process OpenSSF Scorecard data for specified GitHub repositories. It can handle multiple repositories at once and saves the results in a zip file.

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
- `--output value`: Output zip file name (default: "results.zip")
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
scorecard-downloader --purls pkg:github/kubernetes/kubernetes --output custom_results.zip
```

## Input Format

The tool accepts Package URLs (PURLs) as input. PURLs should be in the format:

```
pkg:github/<owner>/<repo>
```

When using the `--purls-file` option, each PURL should be on a separate line in the file.

## Output

The tool saves the processed Scorecard data in a zip file. By default, the output file is named `results.zip`, but you can specify a custom name using the `--output` option.
