package parse

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed sample.json
var sampleJSON []byte

func TestParseSample(t *testing.T) {
	want := []Collection{
		{
			[]string{"Brown", "Blue"},
			"Tuesday, May 14, 2024",
		},
		{
			[]string{"Grey"},
			"Tuesday, May 21, 2024",
		},
		{
			[]string{"Brown", "Blue"},
			"Tuesday, May 28, 2024",
		},
		{
			[]string{"Green"},
			"Tuesday, June 4, 2024",
		},
		{
			[]string{"Brown", "Blue"},
			"Tuesday, June 11, 2024",
		},
		{
			[]string{"Grey"},
			"Tuesday, June 18, 2024",
		},
		{
			[]string{"Brown", "Blue"},
			"Tuesday, June 25, 2024",
		},
		{
			[]string{"Green"},
			"Tuesday, July 2, 2024",
		},
		{
			[]string{"Brown", "Blue"},
			"Tuesday, July 9, 2024",
		},
	}

	collections, err := Parse(sampleJSON)

	assert.Nil(t, err)
	assert.Equal(t, want, collections)
}
