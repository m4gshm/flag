package test

import (
	"flag"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/flagenum/flagenum"
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
		//negative scenarios
		{
			name:          "no default, with allowed, bad arguments",
			allowedValues: []string{"second", "third"},
			arguments:     []string{"first"},
			parseErr:      fmt.Errorf("invalid value \"first\" for flag -val: must be one of [second third]"),
		},
		{
			name:          "bad default, with allowed, with arguments",
			defaultValues: []string{"fifth"},
			allowedValues: []string{"second", "third"},
			arguments:     []string{"second"},
			initErr:       fmt.Errorf("unexpected default value \"fifth\" for flag -val: must be one of [second third]"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flag := flag.NewFlagSet("test", flag.ContinueOnError)
			selected, err := flagenum.Multiple(flag, "val", "enumerated parameter", func(s string) string { return s }, test.defaultValues, test.allowedValues)

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
