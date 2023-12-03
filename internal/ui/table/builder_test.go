package table

import (
	"testing"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/assert"
)

func TestColumnsWithFixedWidth(t *testing.T) {
	longestResourceType := 0

	for i := 0; i < int(resource.NumberOfResourceTypes); i++ {
		typeName := resource.Type(i).String()
		if runewidth.StringWidth(typeName) > longestResourceType {
			longestResourceType = runewidth.StringWidth(typeName)
		}
	}

	assert.Equal(t, longestResourceType, columnsWithFixedWidth["Type"])
}
