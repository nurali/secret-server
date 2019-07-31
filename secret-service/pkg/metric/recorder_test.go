package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractEndpoint(t *testing.T) {
	tables := []struct {
		path string
		want string
	}{
		{"/api/secret/123", "/api/secret"},
		{"/api/secret", "/api/secret"},
	}

	for _, table := range tables {
		got := extractEndpoint(table.path)
		assert.Equal(t, table.want, got)
	}
}
