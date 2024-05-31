package test

import (
	"flag"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/flag/flagenum"
)

func Test_Multiple_String(t *testing.T) {

	type testCase struct {
		name          string
		defaultValues []string
		allowedValues []string
		arguments     []string
		expected      []string
		initErr       error
		parseErr      error
	}

	tests := []testCase{
		//positive scenarios
		{
			name:      "no default, no allowed, with arguments",
			arguments: []string{"first", "second"},
			expected:  []string{"first", "second"},
		},
		{
			name:          "with default, no allowed, with argumens",
			arguments:     []string{"first", "second"},
			defaultValues: []string{"third"},
			expected:      []string{"first", "second"},
		},
		{
			name:          "with default, no allowed, no arguments",
			defaultValues: []string{"third"},
			expected:      []string{"third"},
		},
		{
			name:          "with default, with allowed, with arguments",
			defaultValues: []string{"third"},
			allowedValues: []string{"first", "second", "third"},
			arguments:     []string{"third"},
			expected:      []string{"third"},
		},
		{
			name:          "with default, with allowed, no arguments",
			defaultValues: []string{"third"},
			allowedValues: []string{"third"},
			expected:      []string{"third"},
		},
		{
			name:          "no default, with allowed, no arguments",
			allowedValues: []string{"third"},
		},
		{
			name:      "no default, no allowed, with arguments",
			arguments: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29"},
			expected:  []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29"},
		},
		//negative scenarios
		{
			name:          "no default, with allowed, bad arguments",
			allowedValues: []string{"second", "third"},
			arguments:     []string{"first"},
			parseErr:      fmt.Errorf("invalid value \"first\" for flag -val: must be one of second,third"),
		},
		{
			name:          "no default, bad allowed, no arguments",
			allowedValues: []string{"second", "second"},
			arguments:     []string{"first"},
			initErr:       fmt.Errorf("duplicated allowed value \"second\" for flag -val"),
		},
		{
			name:          "bad default, with allowed, with arguments",
			defaultValues: []string{"fifth"},
			allowedValues: []string{"second", "third"},
			arguments:     []string{"second"},
			initErr:       fmt.Errorf("unexpected default value \"fifth\" for flag -val: must be one of second,third"),
		},
	}

	assertCase := func(test testCase, flag *flag.FlagSet, selected *[]string, err error) {
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
	}

	for _, test := range tests {
		t.Run("MultipleVar:"+test.name, func(t *testing.T) {
			flag := flag.NewFlagSet("test", flag.ContinueOnError)
			var selected []string
			err := flagenum.MultipleVar(flag, &selected, "val", test.defaultValues, test.allowedValues, strAsIs, strAsIs, "enumerated parameter")
			assertCase(test, flag, &selected, err)
		})

		t.Run("Multiple:"+test.name, func(t *testing.T) {
			flag := flag.NewFlagSet("test", flag.ContinueOnError)

			selected, err := flagenum.Multiple(flag, "val", test.defaultValues, test.allowedValues, strAsIs, strAsIs, "enumerated parameter")

			assertCase(test, flag, selected, err)
		})
	}
}

func Test_Multiple_String_Usage(t *testing.T) {
	out := &strings.Builder{}
	flag := flag.NewFlagSet("test", flag.ContinueOnError)
	flag.SetOutput(out)
	_, _ = flagenum.Multiple(flag, "val", []string{"v1", "v3"}, []string{"v1", "v2", "v3"}, strAsIs, strAsIs, "enumerated parameter")

	flag.Usage()

	assert.Equal(t, "Usage of test:\n  -val value\n    \tenumerated parameter (allowed any of v1,v2,v3) (default v1,v3)\n", out.String())
}

func strAsIs(s string) string {
	return s
}
