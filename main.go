package main

import (
	"log"
	"os"

	"github.com/bitbomdev/scorecard-downloader/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "scorecard-downloader",
		Usage: "Download and process scorecard data for GitHub repositories",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "purls",
				Usage:    "PURLs of the repositories to process",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "purls-file",
				Usage:    "File containing PURLs, one per line",
				Required: false,
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "Output file name (default: results.json)",
				Value: "results.json",
			},
			&cli.BoolFlag{
				Name:  "use-bigquery",
				Usage: "Use BigQuery instead of the Scorecard API",
			},
			&cli.StringFlag{
				Name:  "credentials-file",
				Usage: "Path to the BigQuery credentials file",
			},
		},
		Action: cmd.Run,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
