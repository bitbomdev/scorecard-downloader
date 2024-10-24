package utils

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/bitbomdev/scorecard-downloader/internal/processor"
)

func ReadPurlsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var purls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		purl := strings.TrimSpace(scanner.Text())
		if purl != "" {
			purls = append(purls, purl)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return purls, nil
}

func SaveResultsToZip(results []processor.Result, outputZip string) error {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	fileCount := 0
	for _, result := range results {
		if !result.Success || result.Scorecard == nil {
			log.Printf("Skipping result for PURL %s: success=%v, data=%v", result.PURL, result.Success, result.Scorecard != nil)
			continue
		}

		// Use the PURL as the filename, replacing '/' with '_' to avoid directory issues
		fileName := strings.ReplaceAll(result.PURL, "/", "_") + ".json"
		fileWriter, err := zipWriter.Create(fileName)
		if err != nil {
			return fmt.Errorf("error creating file in zip for PURL %s: %v", result.PURL, err)
		}

		var data map[string]interface{}
		if err := json.Unmarshal(result.Scorecard, &data); err != nil {
			return fmt.Errorf("error unmarshaling result data: %v", err)
		}

		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling result data: %v", err)
		}

		if _, err := fileWriter.Write(jsonData); err != nil {
			return fmt.Errorf("error writing data to zip file: %v", err)
		}
		fileCount++
	}

	if fileCount == 0 {
		return fmt.Errorf("no valid results to write to zip file")
	}

	if err := zipWriter.Close(); err != nil {
		return fmt.Errorf("error closing zip writer: %v", err)
	}

	outFile, err := os.Create(outputZip)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, buf); err != nil {
		return fmt.Errorf("error writing zip data to file: %v", err)
	}

	return nil
}

func SaveJSONToFile(data []byte, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error writing JSON data to file: %v", err)
	}

	return nil
}
