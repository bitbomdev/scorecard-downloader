package bigquery

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/bitbomdev/scorecard-downloader/scorecard"
	"google.golang.org/api/iterator"
)

const (
	projectID = "openssf"
	datasetID = "scorecardcron"
	tableID   = "scorecard-v2_latest"
)

func GetScorecardData(ctx context.Context, owner, repo string) (*scorecard.ScorecardData, error) {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	query := client.Query(fmt.Sprintf(`
		SELECT
			date,
			repo.name,
			repo.commit,
			scorecard.version,
			scorecard.commit,
			score,
			ARRAY_AGG(STRUCT(
				checks.name,
				checks.score,
				checks.reason,
				checks.documentation.short AS doc_short,
				checks.documentation.url AS doc_url
			)) AS checks
		FROM %s.%s.%s
		WHERE repo.name = "github.com/%s/%s"
		GROUP BY date, repo.name, repo.commit, scorecard.version, scorecard.commit, score
		ORDER BY date DESC
		LIMIT 1`, projectID, datasetID, tableID, owner, repo))

	it, err := query.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("query.Read: %v", err)
	}

	var scorecardData scorecard.ScorecardData
	for {
		var row struct {
			Date             string
			RepoName         string
			RepoCommit       string
			ScorecardVersion string
			ScorecardCommit  string
			Score            float64
			Checks           []struct {
				Name     string
				Score    int
				Reason   string
				DocShort string
				DocURL   string
			}
		}
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating over results: %v", err)
		}

		scorecardData.Date = row.Date
		scorecardData.Repo.Name = row.RepoName
		scorecardData.Repo.Commit = row.RepoCommit
		scorecardData.Scorecard.Version = row.ScorecardVersion
		scorecardData.Scorecard.Commit = row.ScorecardCommit
		scorecardData.Score = row.Score

		scorecardData.Checks = make([]scorecard.Check, len(row.Checks))
		for i, check := range row.Checks {
			scorecardData.Checks[i] = scorecard.Check{
				Name:     check.Name,
				Score:    check.Score,
				Reason:   check.Reason,
				DocShort: check.DocShort,
				DocURL:   check.DocURL,
			}
		}
	}

	if scorecardData.Repo.Name == "" {
		return nil, nil // No data found
	}

	return &scorecardData, nil
}
