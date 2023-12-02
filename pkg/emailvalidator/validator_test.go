package emailvalidator_test

import (
	"testing"

	"github.com/Slimo300/TeamworkGoTests/pkg/emailvalidator"
)

func TestIsValidEmail(t *testing.T) {

	testCases := []struct {
		desc           string
		email          string
		expectedResult bool
	}{
		{
			desc:           "email without '@'",
			email:          "email",
			expectedResult: false,
		},
		{
			desc:           "email without '.' in domain",
			email:          "host@com",
			expectedResult: false,
		},
		{
			desc:           "email containing more than one '@'",
			email:          "host@net@net.com",
			expectedResult: false,
		},
		{
			desc:           "valid email",
			email:          "host@net.com",
			expectedResult: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			if emailvalidator.IsValidEmail(tC.email) != tC.expectedResult {
				t.Errorf("%s is expected to be %t, but is %t", tC.email, tC.expectedResult, !tC.expectedResult)
			}
		})
	}
}
