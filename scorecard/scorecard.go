package scorecard

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ScorecardData struct {
	Date string `json:"date"`
	Repo struct {
		Name   string `json:"name"`
		Commit string `json:"commit"`
	} `json:"repo"`
	Scorecard struct {
		Version string `json:"version"`
		Commit  string `json:"commit"`
	} `json:"scorecard"`
	Score  float64 `json:"score"`
	Checks []Check
}

type Check struct {
	Name     string
	Score    int
	Reason   string
	DocShort string
	DocURL   string
}

func GetScorecardDataFromAPI(owner, repo, scorecardAPI string) (*ScorecardData, error) {
	url := fmt.Sprintf(scorecardAPI, owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code")
	}

	var scorecardData ScorecardData
	if err := json.NewDecoder(resp.Body).Decode(&scorecardData); err != nil {
		return nil, err
	}

	return &scorecardData, nil
}