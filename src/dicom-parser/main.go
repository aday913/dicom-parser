package main

import (
    "fmt"
    "os"

    "github.com/suyashkumar/dicom"
    "github.com/suyashkumar/dicom/pkg/tag"
)

func main() {
    fmt.Println("Hello, World!")

    dicomFilePath := os.Args[1]
    dataset, err := dicom.ParseFile(dicomFilePath, nil)
    if err != nil {
        fmt.Printf("Error parsing DICOM file: %v\n", err)
        return
    }

    // fmt.Printf("Parsed DICOM dataset: %v\n", dataset)

    element, err := dataset.FindElementByTag(tag.SeriesInstanceUID) // Patient's Name
    if err != nil {
        fmt.Printf("Error finding element: %v\n", err)
        return
    }

    fmt.Printf("Patient's Name: %v\n", element.Value)
}
