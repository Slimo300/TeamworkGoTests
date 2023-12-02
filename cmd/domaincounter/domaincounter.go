package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/Slimo300/TeamworkGoTests/pkg/customerimporter"
)

func main() {

	outputFile := flag.String("output-file", "", "file in which the result should be saved")
	columnName := flag.String("column", "email", "column name for reader to search in csv file")
	format := flag.String("format", "yaml", "yaml or json, format in which result should be displayed, defaults to yaml")

	flag.Parse()

	// We check if customerimporter library supports given format
	if !customerimporter.IsValidFormatType(*format) {
		log.Fatalf("%s is an invalid format type. Use 'yaml' or 'json'", *format)
	}

	filename := flag.Args()[0]

	// opening file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error openning file: %v", err)
	}
	defer file.Close()

	emailReader, err := customerimporter.NewCSVEmailReader(file, *columnName)
	if err != nil {
		log.Fatalf("Error creating new csv email reader: %v", err)
	}

	emailCounter, err := customerimporter.NewDomainCounter(emailReader)
	if err != nil {
		log.Fatalf("Error creating new domain counter: %v", err)
	}

	res, err := emailCounter.CountEmailDomains()
	if err != nil {
		log.Fatalf("Error counting email domains: %v", err)
	}

	var writer io.Writer

	if *outputFile == "" {
		writer = os.Stdout
	} else {
		writer, err = os.OpenFile(*outputFile, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			log.Fatalf("Couldn't open file: %s. Error: %v", *outputFile, err)
		}
	}

	if err := res.WriteTo(writer, *format); err != nil {
		log.Fatalf("Error writing results: %v", err)
	}
}
