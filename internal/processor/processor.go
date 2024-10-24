package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bitbomdev/scorecard-downloader/bigquery"
	"github.com/bitbomdev/scorecard-downloader/purl"
	"github.com/bitbomdev/scorecard-downloader/scorecard"
)

const scorecardAPI = "https://api.scorecard.dev/projects/github.com/%s/%s"

type Processor struct {
	accessToken     string
	useBigQuery     bool
	credentialsFile string
}

type Result struct {
	PURL      string          `json:"purl"`
	Success   bool            `json:"success"`
	Scorecard json.RawMessage `json:"scorecard,omitempty"`
	Error     string          `json:"error,omitempty"`
	GitHubURL string          `json:"github_url,omitempty"`
	Date      time.Time       `json:"date,omitempty"`
}

func NewProcessor(useBigQuery bool, credentialsFile string) *Processor {
	p := &Processor{
		useBigQuery:     useBigQuery,
		credentialsFile: credentialsFile,
	}
	return p
}

func (p *Processor) Process(purls []string) ([]Result, error) {
	results := make([]Result, 0, len(purls))
	repoInfos := make([]bigquery.RepoInfo, 0, len(purls))

	// Collect RepoInfo for each PURL
	for _, pl := range purls {
		purlMap, err := purl.BatchPurlLookup([]string{pl})
		if err != nil {
			results = append(results, Result{PURL: pl, Error: fmt.Sprintf("Error looking up purl: %v", err)})
			continue
		}
		githubURL, exists := purlMap[pl]
		if !exists || githubURL == "" {
			results = append(results, Result{PURL: pl, Error: "GitHub URL not found for purl"})
			continue
		}

		owner, repo := parseGitHubPURL(githubURL)
		if owner == "" || repo == "" {
			results = append(results, Result{PURL: pl, Error: "Invalid GitHub URL"})
			continue
		}

		repoInfos = append(repoInfos, bigquery.RepoInfo{Org: owner, Repo: repo, PURL: pl})
	}

	// Fetch scorecard data for all valid RepoInfos
	var scorecardDataList []*scorecard.ScorecardData
	var err error
	if p.useBigQuery {
		scorecardDataList, err = bigquery.GetScorecardData(context.Background(), repoInfos, p.credentialsFile)
	} else {
		// If not using BigQuery, handle each repo individually
		for _, repoInfo := range repoInfos {
			data, err := scorecard.GetScorecardDataFromAPI(repoInfo.Org, repoInfo.Repo, scorecardAPI)
			if err != nil {
				results = append(results, Result{PURL: fmt.Sprintf("github.com/%s/%s", repoInfo.Org, repoInfo.Repo), Error: fmt.Sprintf("Error fetching scorecard data: %v", err)})
				continue
			}
			scorecardDataList = append(scorecardDataList, data)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("error fetching scorecard data: %v", err)
	}

	// Map results back to PURLs
	for i, data := range scorecardDataList {
		if data != nil {
			jsonData, err := json.Marshal(data)
			if err != nil {
				results = append(results, Result{PURL: fmt.Sprintf("github.com/%s/%s", repoInfos[i].Org, repoInfos[i].Repo), Error: fmt.Sprintf("Error marshaling scorecard data: %v", err)})
			} else {
				results = append(results, Result{
					PURL:      repoInfos[i].PURL,
					Success:   true,
					Scorecard: jsonData,
					GitHubURL: fmt.Sprintf("https://github.com/%s/%s", repoInfos[i].Org, repoInfos[i].Repo),
					Date:      time.Now().UTC(),
				})
			}
		} else {
			results = append(results, Result{PURL: fmt.Sprintf("github.com/%s/%s", repoInfos[i].Org, repoInfos[i].Repo), Error: "Scorecard data not found"})
		}
	}
	return results, nil
}

func parseGitHubPURL(purl string) (string, string) {
	purl = strings.TrimSpace(purl)
	purl = strings.TrimPrefix(purl, "https://")

	const prefix = "github.com/"
	if !strings.HasPrefix(purl, prefix) {
		return "", ""
	}
	parts := strings.Split(strings.TrimPrefix(purl, prefix), "/")
	if len(parts) < 2 {
		return "", ""
	}
	owner := strings.TrimSpace(parts[0])
	repo := strings.TrimSpace(parts[1])
	return owner, repo
}
