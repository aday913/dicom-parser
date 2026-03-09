package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aday913/dicom-parser/pkg/parser"
)

func main() {
	// Define flags
	inputDir := flag.String("input", "", "Input directory containing DICOM files")
	outputDir := flag.String("output", "", "Output directory for structured DICOM files")
	dryRun := flag.Bool("dry-run", false, "Perform a dry run without moving files")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	flag.Parse()

	// Validate flags
	if *inputDir == "" || *outputDir == "" {
		fmt.Println("Error: --input and --output flags are required.")
		flag.Usage()
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Starting DICOM parser...\n")
		fmt.Printf("Input Directory: %s\n", *inputDir)
		fmt.Printf("Output Directory: %s\n", *outputDir)
		fmt.Printf("Dry Run: %t\n", *dryRun)
	}

	// Create output directory if it doesn't exist
	if !*dryRun {
		if err := os.MkdirAll(*outputDir, 0755); err != nil {
			fmt.Printf("Error creating output directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Parse the DICOM directory
	dicomFiles, err := parser.ParseDicomDir(*inputDir, *verbose)
	if err != nil {
		fmt.Printf("Error parsing DICOM directory: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Found %d DICOM files to process.\n", len(dicomFiles))
	}

	// Prepare CSV report
	csvFile, err := os.Create(filepath.Join(*outputDir, "report.csv"))
	if err != nil {
		fmt.Printf("Error creating CSV report: %v\n", err)
		os.Exit(1)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	// Write CSV header
	csvHeader := []string{"OriginalPath", "NewPath", "PatientID", "StudyInstanceUID", "SeriesInstanceUID"}
	if err := csvWriter.Write(csvHeader); err != nil {
		fmt.Printf("Error writing CSV header: %v\n", err)
	}

	// Process each DICOM file
	for _, df := range dicomFiles {
		newPath := filepath.Join(*outputDir, df.PatientID, df.StudyInstanceUID, df.SeriesInstanceUID, filepath.Base(df.FilePath))

		if *verbose {
			fmt.Printf("Processing %s -> %s\n", df.FilePath, newPath)
		}

		// Write to CSV
		csvRow := []string{df.FilePath, newPath, df.PatientID, df.StudyInstanceUID, df.SeriesInstanceUID}
		if err := csvWriter.Write(csvRow); err != nil {
			fmt.Printf("Error writing row to CSV for file %s: %v\n", df.FilePath, err)
		}

		// Copy file if not a dry run
		if !*dryRun {
			if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
				fmt.Printf("Error creating directory for %s: %v\n", newPath, err)
				continue // Skip to next file
			}
			if err := copyFile(df.FilePath, newPath); err != nil {
				fmt.Printf("Error copying file from %s to %s: %v\n", df.FilePath, newPath, err)
			}
		}
	}

	if *verbose {
		fmt.Println("Processing complete.")
		fmt.Printf("CSV report generated at: %s\n", filepath.Join(*outputDir, "report.csv"))
	}
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}
