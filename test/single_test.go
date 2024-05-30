package test

import (
	"flag"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/flag/flagenum"
)

func Test_Single_String(t *testing.T) {

	type testCase struct {
		name          string
		defaultValue  string
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
			name:          "with allowed, with default, no arguments",
			allowedValues: []string{"third"},
			defaultValue:  "third",
			expected:      "third",
		},
		//negative scenarios
		{
			name:          "with allowed, no default, bad arguments",
			allowedValues: []string{"second", "third"},
			arguments:     []string{"first"},
			parseErr:      fmt.Errorf("invalid value \"first\" for flag -val: must be one of second,third"),
		},
		{
			name:          "with allowed, bad default, with arguments",
			allowedValues: []string{"first", "second"},
			defaultValue:  "third",
			initErr:       fmt.Errorf("unexpected default value \"third\" for flag -val: must be one of first,second"),
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
			selected, err := flagenum.Single(flag, "val", test.defaultValue, test.allowedValues, strAsIs, strAsIs, "enumerated parameter")

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

func Test_Single_String_Usage(t *testing.T) {
	out := &strings.Builder{}
	flag := flag.NewFlagSet("test", flag.ContinueOnError)
	flag.SetOutput(out)
	_, _ = flagenum.Single(flag, "val", "v1", []string{"v1", "v2"}, strAsIs, strAsIs, "enumerated parameter")

	flag.Usage()

	assert.Equal(t, "Usage of test:\n  -val value\n    \tenumerated parameter (allowed one of v1,v2) (default v1)\n", out.String())
}
