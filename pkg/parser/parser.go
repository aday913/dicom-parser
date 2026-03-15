package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

// DicomInfo holds the extracted metadata from a DICOM file.
type DicomInfo struct {
	FilePath          string
	PatientID         string
	StudyInstanceUID  string
	SeriesInstanceUID string
}

// ParseDicomDir walks a directory and parses DICOM files.
func ParseDicomDir(dirPath string, verbose bool) ([]DicomInfo, error) {
	var dicomFiles []DicomInfo

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if verbose {
				fmt.Printf("Processing file: %s\n", path)
			}
			dataset, err := dicom.ParseFile(path, nil)
			if err != nil {
				if verbose {
					fmt.Printf("Could not parse %s as DICOM: %v\n", path, err)
				}
				return nil // Continue walking even if a file is not DICOM
			}

			patientID, err := getElementValue(dataset, tag.PatientID)
			if err != nil {
				if verbose {
					fmt.Printf("Skipping file %s: could not get PatientID: %v\n", path, err)
				}
				return nil
			}
			studyUID, err := getElementValue(dataset, tag.StudyInstanceUID)
			if err != nil {
				if verbose {
					fmt.Printf("Skipping file %s: could not get StudyInstanceUID: %v\n", path, err)
				}
				return nil
			}
			seriesUID, err := getElementValue(dataset, tag.SeriesInstanceUID)
			if err != nil {
				if verbose {
					fmt.Printf("Skipping file %s: could not get SeriesInstanceUID: %v\n", path, err)
				}
				return nil
			}

			dicomFiles = append(dicomFiles, DicomInfo{
				FilePath:          path,
				PatientID:         patientID,
				StudyInstanceUID:  studyUID,
				SeriesInstanceUID: seriesUID,
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return dicomFiles, nil
}

func getElementValue(dataset dicom.Dataset, tag tag.Tag) (string, error) {
	element, err := dataset.FindElementByTag(tag)
	if err != nil {
		return "", err
	}
	if element.Value == nil {
		return "", fmt.Errorf("value is nil for tag %s", tag)
	}
	// The value is returned as a slice of strings, we take the first one.
	return element.Value.GetValue().([]string)[0], nil
}
