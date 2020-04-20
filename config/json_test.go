package config

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

const testFile = "deploy_test.json"

func TestUnmarshal_ContainsRepos(t *testing.T) {
	deployment := UnmarshalFile(testFile)
	log.Printf("%v\n", deployment)
	assert.Equal(t, "dirTest", deployment.Dir)
	assert.Equal(t, "revTest", deployment.Rev)
}
