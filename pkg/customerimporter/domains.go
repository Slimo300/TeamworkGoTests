package customerimporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"gopkg.in/yaml.v3"
)

// Domain holds information about domain which is its name and number of times it occures in a file
type Domain struct {
	Name       string
	Occurences int
}

// ErrWriterIsNil is an error returned when nil writer is passed to WriteTo function
type ErrWriterIsNil struct{}

func (e ErrWriterIsNil) Error() string { return "Writer cannot be nil" }

// Error returned when unsupported format is passed
type ErrUnsupportedFormatType struct{ format string }

func (e ErrUnsupportedFormatType) Error() string {
	return fmt.Sprintf("Unsupported format type: %s", e.format)
}

type Domains []Domain

type formatType string

const (
	YAML formatType = "yaml"
	JSON formatType = "json"
)

// IsValidFormatType checks whether given string stands for any of supported formats.
// Format type is also checked in WriteTo method. This method is provided if user would like to validate its method
// before making any actions (counting, sorting, etc.)
func IsValidFormatType(format string) bool {
	if formatType(format) == YAML || formatType(format) == JSON {
		return true
	}
	return false
}

// MapToDomainSlice takes map[string]int and creates a slice of domains from it. Map is a good data structure
// for counting elements but it can't be sorted
func MapToDomainSlice(domainCounter map[string]int) Domains {
	domainSlice := make(Domains, 0, len(domainCounter))

	for domainName, domainOccurences := range domainCounter {
		domainSlice = append(domainSlice, Domain{Name: domainName, Occurences: domainOccurences})
	}

	return domainSlice
}

// Sorting domains based on their occurence counter
func (d Domains) Sort() {
	sort.SliceStable(d, func(i, j int) bool {
		return d[i].Occurences > d[j].Occurences
	})
}

// Write to takes io.Writer and format as its arguments based on given format chooses method to
// apply for writer. It returns error if writer is nil or if format type is unsupported
func (d *Domains) WriteTo(writer io.Writer, format string) error {
	if writer == nil {
		return ErrWriterIsNil{}
	}

	switch formatType(format) {
	case YAML:
		return d.writeYAML(writer)
	case JSON:
		return d.writeJSON(writer)
	default:
		return ErrUnsupportedFormatType{format: format}
	}
}

func (d *Domains) writeJSON(writer io.Writer) error {
	if err := json.NewEncoder(writer).Encode(d); err != nil {
		return err
	}
	return nil
}

func (d *Domains) writeYAML(writer io.Writer) error {
	if err := yaml.NewEncoder(writer).Encode(d); err != nil {
		return err
	}
	return nil
}
