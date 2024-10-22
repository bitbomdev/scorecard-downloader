package processor

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bitbomdev/scorecard-downloader/scorecard"
)

const scorecardAPI = "https://api.scorecard.dev/projects/github.com/%s/%s"

type Processor struct {
	accessToken string
}

type Result struct {
	PURL    string          `json:"purl"`
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

func NewProcessor() *Processor {
	p := &Processor{}
	return p
}

func (p *Processor) Process(purls []string) ([]Result, error) {
	results := make([]Result, 0, len(purls))

	for _, purl := range purls {
		result := p.processPURL(purl)
		results = append(results, result)
	}

	return results, nil
}

func (p *Processor) processPURL(purl string) Result {
	result := Result{PURL: purl}

	owner, repo, err := parsePURL(purl)
	if err != nil {
		result.Error = fmt.Sprintf("Error parsing purl: %v", err)
		return result
	}

	scorecardData, err := scorecard.GetScorecardDataFromAPI(owner, repo, scorecardAPI)
	if err != nil {
		result.Error = fmt.Sprintf("Error fetching scorecard data from API: %v", err)
		return result
	}

	if scorecardData != nil {
		jsonData, err := json.Marshal(scorecardData)
		if err != nil {
			result.Error = fmt.Sprintf("Error marshaling scorecard data: %v", err)
		} else {
			result.Success = true
			result.Data = jsonData
		}
	} else {
		result.Error = "Scorecard data not found"
	}

	return result
}

// Start of Selection
func parsePURL(purl string) (string, string, error) {
	// Remove "pkg:" prefix if present
	purl = strings.TrimPrefix(purl, "pkg:")

	const prefix = "github.com/"
	if !strings.HasPrefix(purl, prefix) {
		return "", "", fmt.Errorf("invalid purl: %s", purl)
	}
	parts := strings.Split(strings.TrimPrefix(purl, prefix), "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid purl format, expected github.com/owner/repo")
	}
	owner := parts[0]
	repo := parts[1]
	return owner, repo, nil
}

func parseGitHubPURL(purl string) (string, string) {
	// Remove "pkg:" prefix if present
	purl = strings.TrimPrefix(purl, "pkg:")

	const prefix = "github.com/"
	if !strings.HasPrefix(purl, prefix) {
		return "", ""
	}
	parts := strings.Split(strings.TrimPrefix(purl, prefix), "/")
	if len(parts) < 2 {
		return "", ""
	}
	owner := parts[0]
	repo := parts[1]
	return owner, repo
}
