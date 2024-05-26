package test

import (
	"flag"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/flagenum/flagenum"
)

func Test_Single_String(t *testing.T) {

	type testCase struct {
		name          string
		allowedValues []string
		arguments     []string
		expected      string
		initErr       error
		parseErr      error
	}

	tests := []testCase{
		//positive scenarios
		{
			name:      "no allowed, with arguments",
			arguments: []string{"first"},
			expected:  "first",
		},
		{
			name:          "with allowed, with arguments",
			allowedValues: []string{"first", "second", "third"},
			arguments:     []string{"third"},
			expected:      "third",
		},
		{
			name:          "with allowed, no arguments",
			allowedValues: []string{"third"},
		},
		{
			name:          "with allowed, no arguments",
			allowedValues: []string{"third"},
		},
		//negative scenarios
		{
			name:          "no default, with allowed, bad arguments",
			allowedValues: []string{"second", "third"},
			arguments:     []string{"first"},
			parseErr:      fmt.Errorf("invalid value \"first\" for flag -val: must be one of [second third]"),
		},
		{
			name:          "bad allowed, with arguments",
			allowedValues: []string{"second", "second"},
			arguments:     []string{"second"},
			initErr:       fmt.Errorf("duplicated allowed value \"second\" for flag -val"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flag := flag.NewFlagSet("test", flag.ContinueOnError)
			selected, err := flagenum.Single(flag, "val", "enumerated parameter", func(s string) string { return s }, test.allowedValues)

			if test.initErr != nil {
				assert.EqualError(t, err, test.initErr.Error())
			} else {

				a := []string{}
				for _, arg := range test.arguments {
					a = append(a, "--val", arg)
				}

				err = flag.Parse(a)
				if test.parseErr != nil {
					assert.EqualError(t, err, test.parseErr.Error())
				} else {
					assert.NoError(t, err)
					assert.Equal(t, test.expected, *selected)
				}
			}
		})
	}
}
