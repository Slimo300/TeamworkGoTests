// package customerimporter reads from the given customers.csv file and returns a
// sorted (data structure of your choice) of email domains along with the number
// of customers with e-mail addresses for each domain.  Any errors should be
// logged (or handled). Performance matters (this is only ~3k lines, but *could*
// be 1m lines or run on a small machine).
package customerimporter

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Slimo300/TeamworkGoTests/pkg/emailvalidator"
)

type ErrColumnNotFound struct {
	column string
}

func (e ErrColumnNotFound) Error() string {
	return fmt.Sprintf("There is no column named %s in given file", e.column)
}

type ErrReaderNil struct{}

func (e ErrReaderNil) Error() string {
	return "Reader cannot be nil"
}

// EmailReader is an interface describing ReadEmail functionality
type EmailReader interface {
	ReadEmail() (string, error)
}

type mockEmailReader struct {
	current int
	emails  []string
}

// ReadEmail with every call iterates over its internal slice and returns next item in it until it is out of slice's range
// Then it returns io.EOF
func (r *mockEmailReader) ReadEmail() (string, error) {
	if r.current >= len(r.emails) {
		return "", io.EOF
	}

	// we first save email to be returned then increment current.
	email := r.emails[r.current]
	r.current += 1

	return email, nil
}

type csvEmailReader struct {
	reader      *csv.Reader
	emailColumn int
}

// NewCSVEmailReader creates a new email reader
func NewCSVEmailReader(reader io.Reader, columnName string) (EmailReader, error) {

	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = 0

	// Reading first line of csv containing column names
	record, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	emailColumn := -1
	// Finding a column containing email addresses
	for i, field := range record {
		if field == columnName {
			emailColumn = i
		}
	}
	// if emailColumn remains -1 it means we haven't found given column name
	if emailColumn == -1 {
		return nil, ErrColumnNotFound{column: columnName}
	}

	return &csvEmailReader{
		reader:      csvReader,
		emailColumn: emailColumn,
	}, nil
}

// ReadEmail method returns next email in csv file and error if one occurs.
// It returns io.EOF when at end of the file
func (r *csvEmailReader) ReadEmail() (string, error) {
	record, err := r.reader.Read()
	if err != nil {
		return "", err
	}

	email := record[r.emailColumn]

	return email, nil

}

// DomainCounter is responsible for counting domains
type DomainCounter interface {
	CountEmailDomains() (Domains, error)
}

type domainCounter struct {
	reader EmailReader
}

// NewDomainCounter takes EmailReader as argument and returns instance of DomainCounter interface
func NewDomainCounter(emailReader EmailReader) (DomainCounter, error) {
	if emailReader == nil {
		return nil, ErrReaderNil{}
	}

	return &domainCounter{
		reader: emailReader,
	}, nil
}

// CountEmailDomains is a method that reads from EmailReader interface until EOF, validates
// received emails, extracts domains from them and counts their occurences. It returns a sorted slice
// of found domains ordered decrementally by their number of occurences.
func (c *domainCounter) CountEmailDomains() (Domains, error) {
	domainCounter := make(map[string]int)

	for i := 2; ; i++ {
		email, err := c.reader.ReadEmail()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if !emailvalidator.IsValidEmail(email) {
			// we are omitting invalid emails and logging iformation about them to terminal
			log.Printf("Error validating email at line %d. %s is not a valid email address", i, email)
			continue
		}

		// since email passed validation we don't need to worry about panic
		domain := strings.Split(email, "@")[1]

		domainCounter[domain]++
	}

	sortedDomains := MapToDomainSlice(domainCounter)
	sortedDomains.Sort()

	return sortedDomains, nil
}
