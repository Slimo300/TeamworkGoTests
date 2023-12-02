package customerimporter

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestDomains_Sort(t *testing.T) {
	testCases := []struct {
		desc           string
		input          Domains
		expectedOutput Domains
	}{
		{
			desc:           "test",
			input:          []Domain{{Name: "blogger.com", Occurences: 2}, {Name: "gmail.com", Occurences: 1}, {Name: "mail.com", Occurences: 4}},
			expectedOutput: []Domain{{Name: "mail.com", Occurences: 4}, {Name: "blogger.com", Occurences: 2}, {Name: "gmail.com", Occurences: 1}},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tC.input.Sort()
			if !reflect.DeepEqual(tC.expectedOutput, tC.input) {
				t.Errorf("Output: %v is different than expected: %v", tC.input, tC.expectedOutput)
			}
		})
	}
}

func ExampleIsValidFormatType() {
	fmt.Println(IsValidFormatType("json"))
	fmt.Println(IsValidFormatType("yaml"))
	fmt.Println(IsValidFormatType("csv"))
	// Output:
	// true
	// true
	// false
}

func TestDomains_writeJSON(t *testing.T) {
	domains := Domains{{Name: "mail.com", Occurences: 2}, {Name: "gmail.com", Occurences: 1}}

	buf := bytes.NewBuffer(nil)
	if err := domains.writeJSON(buf); err != nil {
		t.Errorf("Error writing json to buffer: %v", err)
	}

	reader := bufio.NewReader(buf)
	res, err := reader.ReadString('\n')
	if err != nil {
		t.Errorf("Error reading from buffer: %v", err)
	}

	expectedResult := "[{\"Name\":\"mail.com\",\"Occurences\":2},{\"Name\":\"gmail.com\",\"Occurences\":1}]"

	if strings.TrimSpace(res) != expectedResult {
		t.Errorf("Result different than expected. \nResult: %v\nExpected: %v", res, expectedResult)
	}

}

func TestDomains_writeYAML(t *testing.T) {
	domains := Domains{{Name: "mail.com", Occurences: 2}, {Name: "gmail.com", Occurences: 1}}

	buf := bytes.NewBuffer(nil)
	if err := domains.writeYAML(buf); err != nil {
		t.Errorf("Error writing json to buffer: %v", err)
	}

	reader := bufio.NewReader(buf)
	res, err := io.ReadAll(reader)
	t.Log(res)
	if err != nil {
		t.Errorf("Error reading from buffer: %v", err)
	}

	expectedResult := `- name: mail.com
  occurences: 2
- name: gmail.com
  occurences: 1
`

	if string(res) != expectedResult {
		t.Errorf("Result different than expected. \nResult: %v\nExpected: %v", string(res), expectedResult)
	}

}

func TestDomains_WriteTo(t *testing.T) {
	testCases := []struct {
		desc        string
		writer      io.Writer
		format      string
		expectedErr error
	}{
		{
			desc:        "writer is nil",
			writer:      nil,
			format:      "json",
			expectedErr: ErrWriterIsNil{},
		},
		{
			desc:        "invalid format",
			writer:      bytes.NewBuffer(nil),
			format:      "csv",
			expectedErr: ErrUnsupportedFormatType{format: "csv"},
		},
		{
			desc:        "format json",
			writer:      bytes.NewBuffer(nil),
			format:      "json",
			expectedErr: nil,
		},
		{
			desc:        "format yaml",
			writer:      bytes.NewBuffer(nil),
			format:      "yaml",
			expectedErr: nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			domains := Domains{{Name: "mail.com", Occurences: 2}}

			err := domains.WriteTo(tC.writer, tC.format)
			if !errors.Is(err, tC.expectedErr) {
				t.Errorf("Received error: %s is different than expected: %s", err, tC.expectedErr)
			}
		})

	}
}
