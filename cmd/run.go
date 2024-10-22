package cmd

import (
	"fmt"
	"log"

	"github.com/bitbomdev/scorecard-downloader/internal/processor"
	"github.com/bitbomdev/scorecard-downloader/internal/utils"
	"github.com/urfave/cli/v2"
)

func Run(c *cli.Context) error {
	purls := c.StringSlice("purls")
	purlsFile := c.String("purls-file")
	outputZip := c.String("output")

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

	p := processor.NewProcessor()
	results, err := p.Process(purls)
	if err != nil {
		return fmt.Errorf("error processing scorecard: %v", err)
	}

	if err := utils.SaveResultsToZip(results, outputZip); err != nil {
		return fmt.Errorf("error saving results to zip: %v", err)
	}

	fmt.Printf("Results saved to %s\n", outputZip)
	return nil
}
