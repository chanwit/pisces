package swarm

import (
	"os"
	"testing"

	"github.com/chanwit/pisces/conf"
	"github.com/stretchr/testify/assert"
)

// this test must be run inside integration script
func TestBuild(t *testing.T) {
	spec := BuildSpec{
		Info:       conf.Info{Build: "./web"},
		NodeName:   "node-0",
		NodeAddr:   os.Getenv("DOCKER_NODE_0_ADDR"),
		ProjectDir: os.Getenv("TESTDATA"),
		Project:    "project",
		Service:    "web",
		NoCache:    false,
	}

	imageId := Build(spec)
	assert.NotEmpty(t, imageId)
	assert.Equal(t, 12, len(imageId))
}
