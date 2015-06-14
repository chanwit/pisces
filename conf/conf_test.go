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

// this test must be run inside integration script
func TestParseConfig2(t *testing.T) {
    content := `
front:
  build: ./front
  links:
  - web

web:
  build: ./web
  links:
  - db

db:
  image: busybox
`
    config, err := parseConfig([]byte(content))

    assert.NoError(t, err)
    assert.Nil(t, config.PodSpec)
    assert.Equal(t, "./web", config.Services["web"].Build)
    assert.Equal(t, "./front", config.Services["front"].Build)
    assert.Equal(t, "busybox", config.Services["db"].Image)
}

func TestTopoSort(t *testing.T) {
    g := make(map[string][]string)
    g["front"] = []string{"web"}
    g["web"] = []string{"db"}
    g["db"] = []string{}
    order, cyclics := topoSortDFS(g)

    assert.Equal(t, []string{"db", "web", "front"}, order)
    assert.Empty(t, cyclics)
}