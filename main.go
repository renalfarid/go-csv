package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

func evaluateFile(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
	}

	return lineCount - 1, scanner.Err()
}

func readCsvFile(filePath, outputPath string, estimatedTotalLines int) error {
	startTime := time.Now()

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	bar := progressbar.Default(int64(estimatedTotalLines))

	var outputRows []map[string]interface{}
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("could not read headers: %w", err)
	}

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		row := make(map[string]interface{})
		for i, value := range record {
			row[strings.ToLower(headers[i])] = value
		}

		outputRows = append(outputRows, row)
		bar.Add(1)
	}

	bar.Finish()

	if outputPath != "" {
		fmt.Println("Writing json File...")
		fmt.Println("=====================")
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("could not create output file: %w", err)
		}
		defer outputFile.Close()

		jsonData, err := json.MarshalIndent(outputRows, "", "  ")
		if err != nil {
			return fmt.Errorf("could not marshal JSON: %w", err)
		}

		_, err = outputFile.Write(jsonData)
		if err != nil {
			return fmt.Errorf("could not write JSON data: %w", err)
		}
	}

	endTime := time.Now()
	readingTime := endTime.Sub(startTime).Seconds()

	fmt.Printf("File name: %s\n", filePath)
	fmt.Printf("Processing time: %.2f seconds\n", readingTime)

	return nil
}

func main() {
	args := os.Args
	fileIndex := -1
	outputIndex := -1

	for i, arg := range args {
		if arg == "--file" && i+1 < len(args) {
			fileIndex = i + 1
		} else if arg == "--output" && i+1 < len(args) {
			outputIndex = i + 1
		}
	}

	if fileIndex != -1 {
		filePath := args[fileIndex]
		outputPath := ""
		if outputIndex != -1 {
			outputPath = args[outputIndex]
		}

		fmt.Println("Reading file...")
		fmt.Println("=================")

		estimatedTotalLines, err := evaluateFile(filePath)
		if err != nil {
			fmt.Println("An error occurred during file evaluation:", err)
			return
		}

		fmt.Printf("Estimated total lines: %d\n", estimatedTotalLines)
		err = readCsvFile(filePath, outputPath, estimatedTotalLines)
		if err != nil {
			fmt.Println("An error occurred:", err)
		}
	} else {
		fmt.Println("Please provide a file path using the --file argument.")
	}
}
