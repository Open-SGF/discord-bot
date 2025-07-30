package resource

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNamer(t *testing.T) {
	tests := []struct {
		testName string
		prefix   string
		name     string
		output   string
	}{
		{
			testName: "no prefix",
			prefix:   "",
			name:     "somename",
			output:   "somename",
		},
		{
			testName: "with prefix",
			prefix:   "prefix",
			name:     "somename",
			output:   "prefix-somename",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			namer := NewNamer(test.prefix, test.name)

			assert.Equal(t, test.name, namer.Name())
			assert.Equal(t, test.output, namer.FullName())
		})
	}
}
