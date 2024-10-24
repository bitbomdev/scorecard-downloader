package bigquery

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/bitbomdev/scorecard-downloader/scorecard"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	projectID = "openssf"
	datasetID = "scorecardcron"
	tableID   = "scorecard-v2_latest"
)

type RepoInfo struct {
	Org  string
	Repo string
	PURL string
}

func GetScorecardData(ctx context.Context, repos []RepoInfo, credentialsFile string) ([]*scorecard.ScorecardData, error) {
	if credentialsFile == "" {
		return nil, fmt.Errorf("credentials file is required")
	}
	if len(repos) == 0 {
		return nil, fmt.Errorf("at least one repo is required")
	}
	client, err := bigquery.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	var whereClause string
	if len(repos) == 1 {
		// Handle the edge case where only one repository is passed
		whereClause = fmt.Sprintf(`repo.name = "github.com/%s/%s"`, repos[0].Org, repos[0].Repo)
	} else {
		// Construct the WHERE clause for multiple repositories
		repoConditions := make([]string, len(repos))
		for i, repo := range repos {
			repoConditions[i] = fmt.Sprintf(`repo.name = "github.com/%s/%s"`, repo.Org, repo.Repo)
		}
		whereClause = strings.Join(repoConditions, " OR ")
	}

	from := fmt.Sprintf("`%s.%s.%s`", projectID, datasetID, tableID)
	query := client.Query(fmt.Sprintf(`
		SELECT
			FORMAT_DATE("%%Y-%%m-%%d", date) AS date,
			repo.name AS repo_name,
			repo.commit AS repo_commit,
			scorecard.version AS scorecard_version,
			scorecard.commit AS scorecard_commit,
			a.score,
			ARRAY_AGG(
				STRUCT(
					c.name,
					c.score,
					c.reason
				)
			) AS checks
		FROM %s a
		LEFT JOIN UNNEST(checks) AS c
		WHERE %s
		GROUP BY date, repo.name, repo.commit, scorecard.version, scorecard.commit, a.score
		ORDER BY date DESC
	`, from, whereClause))

	it, err := query.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("query.Read: %v", err)
	}
	var scorecardDataList []*scorecard.ScorecardData
	for {
		var row struct {
			Date             string
			RepoName         string `bigquery:"repo_name"`
			RepoCommit       string `bigquery:"repo_commit"`
			ScorecardVersion string `bigquery:"scorecard_version"`
			ScorecardCommit  string `bigquery:"scorecard_commit"`
			Score            float64
			Checks           []struct {
				Name   string
				Score  int
				Reason string
			}
		}
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating over results: %v", err)
		}

		// Find the corresponding RepoInfo to get the PURL
		var purl string
		for _, repo := range repos {
			if "github.com/"+repo.Org+"/"+repo.Repo == row.RepoName {
				purl = repo.PURL
				break
			}
		}

		scorecardData := &scorecard.ScorecardData{
			Date: row.Date,
			Repo: scorecard.Repo{
				Name:   row.RepoName,
				Commit: row.RepoCommit,
			},
			Scorecard: scorecard.Scorecard{
				Version: row.ScorecardVersion,
				Commit:  row.ScorecardCommit,
			},
			Score:  row.Score,
			Checks: make([]scorecard.Check, len(row.Checks)),
			PURL:   purl, // Set the PURL here
		}

		for i, check := range row.Checks {
			scorecardData.Checks[i] = scorecard.Check{
				Name:   check.Name,
				Score:  check.Score,
				Reason: check.Reason,
			}
		}

		scorecardDataList = append(scorecardDataList, scorecardData)
	}
	return scorecardDataList, nil
}
