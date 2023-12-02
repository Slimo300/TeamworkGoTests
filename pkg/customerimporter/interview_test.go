package customerimporter

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"reflect"
	"testing"
)

func TestNewCSVEmailReader(t *testing.T) {
	testCases := []struct {
		desc           string
		column         string
		fileData       string
		returnReader   bool
		expectedReader EmailReader
		expectedError  error
	}{
		{
			desc:           "empty file",
			fileData:       "",
			column:         "email",
			returnReader:   false,
			expectedReader: nil,
			expectedError:  io.EOF,
		},
		{
			desc:           "file without specified column",
			fileData:       "name,mail\nJohn,john.doe@mail.com",
			column:         "email",
			returnReader:   false,
			expectedReader: nil,
			expectedError:  ErrColumnNotFound{column: "email"},
		},
		{
			desc:         "csvEmailReader created",
			fileData:     "name,email\nJohn,john.doe@mail.com",
			returnReader: true,
			column:       "email",
			// We don't expect first row as it is read at constructor
			expectedReader: &csvEmailReader{reader: csv.NewReader(bytes.NewReader([]byte("John,john.doe@mail.com"))), emailColumn: 1},
			expectedError:  nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// For test purposes we use bytes.Reader instead of a file
			reader, err := NewCSVEmailReader(bytes.NewReader([]byte(tC.fileData)), tC.column)

			if err != tC.expectedError {
				t.Errorf("Returned error: %v is different than expected: %v", err, tC.expectedError)
			}

			if tC.returnReader {
				readerAsCSVReader, ok := reader.(*csvEmailReader)
				if !ok {
					t.Error("Returned reader is not of type csvEmailReader")
				}
				expectedReaderAsCSVReader, _ := tC.expectedReader.(*csvEmailReader)

				// checking if emailColumn parameters match
				if readerAsCSVReader.emailColumn != expectedReaderAsCSVReader.emailColumn {
					t.Errorf("csvEmailReader emailColumn (%d) don't match expected value: %d", readerAsCSVReader.emailColumn, expectedReaderAsCSVReader.emailColumn)
				}

				// we compare data that can be read from readers we use ReadAll method because we know size of our test data
				readerData, err := readerAsCSVReader.reader.ReadAll()
				if err != nil {
					t.Errorf("Error reading data from reader: %v", err)
				}
				expectedReaderData, err := expectedReaderAsCSVReader.reader.ReadAll()
				if err != nil {
					t.Errorf("Error reading data from (expected) reader: %v", err)
				}

				if !reflect.DeepEqual(readerData, expectedReaderData) {
					t.Errorf("Values in reader: %v values don't match expected: %v", readerData, expectedReaderData)
				}

			}
		})
	}
}

func TestCSVEmailReader_ReadEmail(t *testing.T) {
	testCases := []struct {
		desc             string
		fileData         string
		expectedResponse string
		expectedError    error
	}{
		{
			desc:             "empty file",
			fileData:         "name,email", // we pass first row as its going to be read by EmailReader constructor to determine email column
			expectedResponse: "",
			expectedError:    io.EOF,
		},
		{
			desc:             "successful read",
			fileData:         "name,email\nJonh,john.doe@mail.com",
			expectedResponse: "john.doe@mail.com",
			expectedError:    nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			emailReader, err := NewCSVEmailReader(bytes.NewReader([]byte(tC.fileData)), "email")
			if err != nil {
				t.Errorf("Error creating csvEmailReader: %v", err)
			}

			response, err := emailReader.ReadEmail()
			if response != tC.expectedResponse {
				t.Errorf("Response from ReadEmail: '%s' is different than expected: '%s'", response, tC.expectedResponse)
			}
			if err != tC.expectedError {
				t.Errorf("Error from ReadEmail: '%s' is different than expected: '%s'", err, tC.expectedError)
			}

		})
	}
}

func TestNewDomainCounter(t *testing.T) {
	testCases := []struct {
		desc           string
		reader         EmailReader
		expectedResult DomainCounter
		expectedErr    error
	}{
		{
			desc:           "reader is nil",
			reader:         nil,
			expectedResult: nil,
			expectedErr:    ErrReaderNil{},
		},
		{
			desc:           "success",
			reader:         &mockEmailReader{},
			expectedResult: &domainCounter{reader: &mockEmailReader{}},
			expectedErr:    nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			counter, err := NewDomainCounter(tC.reader)
			if !errors.Is(err, tC.expectedErr) {
				t.Errorf("Received error is different than expected. Received: %v, expected: %v", err, tC.expectedErr)
			}
			if !reflect.DeepEqual(counter, tC.expectedResult) {
				t.Errorf("Received counter is different than expected. Received: %v, expected: %v", counter, tC.expectedResult)
			}
		})
	}
}

func TestDomainCounter_CountEmailDomains(t *testing.T) {
	testCases := []struct {
		desc           string
		emails         []string
		expectedOutput Domains
	}{
		{
			desc:           "count_domains",
			emails:         []string{"host@gmail.com", "host2@gmail.com", "host@gmail.com", "host@mail.com"},
			expectedOutput: []Domain{{Name: "gmail.com", Occurences: 3}, {Name: "mail.com", Occurences: 1}},
		},
		{
			desc:           "count_domains_with_invalid",
			emails:         []string{"host@gmail.com", "host2@gmail.com", "host@gmail.com", "host@mail.com", "email", "mail@.com"},
			expectedOutput: []Domain{{Name: "gmail.com", Occurences: 3}, {Name: "mail.com", Occurences: 1}},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			domainCounter, err := NewDomainCounter(&mockEmailReader{emails: tC.emails})
			if err != nil {
				t.Errorf("Error when creating domain counter: %v", err)
			}

			result, err := domainCounter.CountEmailDomains()
			if err != nil {
				t.Errorf("Error when trying to count email domains")
			}
			if !reflect.DeepEqual(result, tC.expectedOutput) {
				t.Errorf("Output: %v is different than expected: %v", result, tC.expectedOutput)
			}
		})
	}
}
