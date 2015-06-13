package swarm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// this test must be run inside integration script
func TestNodes(t *testing.T) {
	nodes := Nodes()

	assert.Equal(t, 1, len(nodes))
	assert.Equal(t, []string{"node-0"}, nodes)
}
