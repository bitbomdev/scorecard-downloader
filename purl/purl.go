package purl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// PurlLookupBatchRequest represents the request body for PurlLookupBatch API.
type PurlLookupBatchRequest struct {
	Requests []struct {
		Purl string `json:"purl"`
	} `json:"requests"`
}

// PurlLookupBatchResponse represents the response from PurlLookupBatch API.
type PurlLookupBatchResponse struct {
	Responses []struct {
		Request struct {
			Purl string `json:"purl"`
		} `json:"request"`
		Result struct {
			Version struct {
				Links []struct {
					Label string `json:"label"`
					URL   string `json:"url"`
				} `json:"links"`
			} `json:"version"`
			Package struct {
				Links []struct {
					Label string `json:"label"`
					URL   string `json:"url"`
				} `json:"links"`
			} `json:"package"`
		} `json:"result"`
	} `json:"responses"`
}

// BatchPurlLookup performs batch lookup of pURLs and returns a map of pURL to GitHub repository URL.
func BatchPurlLookup(purls []string) (map[string]string, error) {
	const (
		apiURL       = "https://api.deps.dev/v3alpha/purlbatch"
		maxBatchSize = 100
	)

	result := make(map[string]string)

	// Chunk pURLs into batches of maxBatchSize
	for i := 0; i < len(purls); i += maxBatchSize {
		end := i + maxBatchSize
		if end > len(purls) {
			end = len(purls)
		}
		batch := purls[i:end]

		// Create request body
		reqBody := PurlLookupBatchRequest{}
		for _, p := range batch {
			reqBody.Requests = append(reqBody.Requests, struct {
				Purl string `json:"purl"`
			}{Purl: p})
		}

		// Marshal request to JSON
		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}

		// Create HTTP POST request
		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, fmt.Errorf("failed to create HTTP request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Send HTTP request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send HTTP request: %v", err)
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}

		// Check for non-200 status codes
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API returned non-OK status: %s, body: %s", resp.Status, string(body))
		}

		// Unmarshal response JSON
		var apiResp PurlLookupBatchResponse
		err = json.Unmarshal(body, &apiResp)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response JSON: %v", err)
		}

		// Process each response
		for _, r := range apiResp.Responses {
			githubURL := ""
			// Check version links
			for _, link := range r.Result.Version.Links {
				if link.Label == "SOURCE_REPO" {
					if isGitHubURL(link.URL) {
						githubURL = convertGitURLToHTTPS(link.URL)
						break
					}
				}
			}
			// If not found in version, check package links
			if githubURL == "" {
				for _, link := range r.Result.Package.Links {
					if link.Label == "SOURCE_REPO" {
						if isGitHubURL(link.URL) {
							githubURL = convertGitURLToHTTPS(link.URL)
							break
						}
					}
				}
			}
			result[r.Request.Purl] = githubURL
		}
	}

	return result, nil
}

// isGitHubURL checks if the URL is a GitHub repository.
func isGitHubURL(repoURL string) bool {
	parsedURL, err := url.Parse(repoURL)
	if err != nil {
		return false
	}
	return parsedURL.Host == "github.com"
}

// convertGitURLToHTTPS converts a git+ssh URL to HTTPS URL.
func convertGitURLToHTTPS(gitURL string) string {
	// Example: git+ssh://git@github.com/DABH/colors.js.git -> https://github.com/DABH/colors.js
	parsedURL, err := url.Parse(gitURL)
	if err != nil {
		return ""
	}
	path := parsedURL.Path
	if len(path) > 0 && path[len(path)-4:] == ".git" {
		path = path[:len(path)-4]
	}
	return fmt.Sprintf("https://github.com%s", path)
}
