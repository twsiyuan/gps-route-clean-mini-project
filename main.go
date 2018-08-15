package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var input, output, outputFormat string
	flag.StringVar(&input, "in", "points.csv", "input file")
	flag.StringVar(&output, "out", "", "output file")
	flag.StringVar(&outputFormat, "out_format", "csv", "output format, 'csv' or 'kml'")
	flag.Parse()

	exportFunc := ExportAsCSV
	switch strings.TrimSpace(outputFormat) {
	case "csv":
		if len(output) <= 0 {
			output = "points_output.csv"
		}
		break
	case "kml":
		if len(output) <= 0 {
			output = "points_output.kml"
		}
		exportFunc = ExportAsKML
		break
	default:
		fmt.Fprintln(os.Stderr, "Unexpected output foramt")
		os.Exit(1)
		break
	}

	inputF, err := os.Open(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Open file failed, %v\n", err)
		os.Exit(1)
	}
	defer inputF.Close()

	outputF, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Open file failed, %v\n", err)
		os.Exit(1)
	}
	defer outputF.Close()

	reader := NewGPSReader(inputF)
	reader = NewSpeedFilter(reader, 80)
	reader = NewNoiseFilter(reader)

	if err := exportFunc(reader, outputF); err != nil {
		fmt.Fprintf(os.Stderr, "Export failed, %v\n", err)
		os.Exit(1)
	}
}
