package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bitbomdev/scorecard-downloader/internal/processor"
	"github.com/bitbomdev/scorecard-downloader/internal/utils"
	"github.com/urfave/cli/v2"
)

func Run(c *cli.Context) error {
	purls := c.StringSlice("purls")
	purlsFile := c.String("purls-file")
	outputFile := c.String("output")

	if len(purls) == 0 && purlsFile == "" {
		return fmt.Errorf("either --purls or --purls-file must be provided")
	}

	if purlsFile != "" {
		filePurls, err := utils.ReadPurlsFromFile(purlsFile)
		if err != nil {
			return fmt.Errorf("error reading PURLs from file: %v", err)
		}
		purls = append(purls, filePurls...)
	}

	log.Printf("Input PURLs: %v", purls)

	useBigQuery := c.Bool("use-bigquery")
	credentialsFile := c.String("credentials-file")

	p := processor.NewProcessor(useBigQuery, credentialsFile)
	results, err := p.Process(purls)
	if err != nil {
		return fmt.Errorf("error processing scorecard: %v", err)
	}

	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling results to JSON: %v", err)
	}
	if outputFile != "" {
		if err := utils.SaveJSONToFile(jsonData, outputFile); err != nil {
			return fmt.Errorf("error saving results to file: %v", err)
		}
		fmt.Printf("Results saved to %s\n", outputFile)
	} else {
		fmt.Println(string(jsonData))
	}

	return nil
}
