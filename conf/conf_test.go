package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// this test must be run inside integration script
func TestParseConfig(t *testing.T) {
	content := `
web:
  build: ./web
`

	config, err := parseConfig([]byte(content))

	assert.NoError(t, err)
	assert.Nil(t, config.PodSpec)
	assert.Equal(t, "./web", config.Services["web"].Build)
}
