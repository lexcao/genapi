package printer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// stringSlice implements Item for a slice of strings
type stringSlice []string

func (s stringSlice) Name() string {
	return "[]string"
}

func (s stringSlice) Print(p *Printer) {
	for _, item := range s {
		p.Indent()
		p.Quote(item)
		p.Unquoted(",")
		p.NewLine()
	}
}

func TestStringSlice(t *testing.T) {
	slice := stringSlice{"item1", "item2", "item3"}

	actual := Print(slice)
	expect := `[]string{
	"item1",
	"item2",
	"item3",
}`

	assert.Equal(t, expect, actual)
}

type stringSliceSlice []stringSlice

func (s stringSliceSlice) Name() string {
	return "[][]string"
}

func (s stringSliceSlice) Print(p *Printer) {
	for _, item := range s {
		p.Indent()
		p.Item(item)
		p.Unquoted(",")
		p.NewLine()
	}
}

func TestStringSliceSlice(t *testing.T) {
	slice := stringSliceSlice{stringSlice{"item1", "item2"}, stringSlice{"item3", "item4"}}

	actual := Print(slice)
	expect := `[][]string{
	[]string{
		"item1",
		"item2",
	},
	[]string{
		"item3",
		"item4",
	},
}`

	assert.Equal(t, expect, actual)
}
