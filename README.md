# DICOM Parser

This is a command-line tool for organizing a directory of unstructured DICOM files into a structured `Patient/Study/Series` directory hierarchy based on the files' DICOM tag metadata. It also generates a CSV report mapping the original file paths to their new locations.

## Features

-   Scans a directory recursively for DICOM files.
-   Parses DICOM tags to identify Patient ID, Study Instance UID, and Series Instance UID.
-   Copies files into a structured output directory: `<output_dir>/<PatientID>/<StudyInstanceUID>/<SeriesInstanceUID>/`.
-   Generates a `report.csv` file in the output directory detailing the file movements.
-   Includes a `--dry-run` mode to generate the report without moving any files.
-   Includes a `--verbose` mode for detailed logging.

## Installation

You can build the binary from the source code. Make sure you have Go installed (version 1.22 or later).

```bash
go build -o dicom-parser ./cmd/dicom-parser
```

This will create a `dicom-parser` binary in your project's root directory.

## Usage

The tool requires an input directory and an output directory.

```bash
./dicom-parser --input <path_to_dicom_files> --output <path_to_output_dir>
```

### Command-Line Flags

-   `--input` (required): Path to the directory containing the unstructured DICOM files.
-   `--output` (required): Path to the directory where the structured files and report will be saved.
-   `--dry-run` (optional): If set to `true`, the tool will not copy any files but will still generate the `report.csv`. Defaults to `false`.
-   `--verbose` (optional): If set to `true`, the tool will print detailed logs of its progress to the terminal. Defaults to `false`.

### Example

```bash
# Perform a dry run with verbose logging
./dicom-parser --input ./my_dicoms --output ./organized_dicoms --dry-run --verbose

# Run the tool to organize the files
./dicom-parser --input ./my_dicoms --output ./organized_dicoms
```

## Output

### Directory Structure

The tool will create a new directory structure inside your specified output directory, like so:

```
<output_dir>/
├── <PatientID_1>/
│   └── <StudyInstanceUID_A>/
│       └── <SeriesInstanceUID_X>/
│           ├── file1.dcm
│           └── file2.dcm
├── <PatientID_2>/
│   └── <StudyInstanceUID_B>/
│       └── <SeriesInstanceUID_Y>/
│           └── file3.dcm
└── report.csv
```

### CSV Report

The `report.csv` file will contain the following columns:

-   `OriginalPath`: The original absolute path to the DICOM file.
-   `NewPath`: The new absolute path where the file has been copied.
-   `PatientID`: The Patient ID extracted from the DICOM header.
-   `StudyInstanceUID`: The Study Instance UID extracted from the DICOM header.
-   `SeriesInstanceUID`: The Series Instance UID extracted from the DICOM header.

## Sample DICOM Files

You can find sample DICOM files for testing at this site: [https://medimodel.com/sample-dicom-files/](https://medimodel.com/sample-dicom-files/)
