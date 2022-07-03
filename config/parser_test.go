package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var configFilepath string = "config.yml"
var invalidFilepath string = "gibberish"

func TestParseConfigValidFilepath(t *testing.T) {
	_, err := ParseConfig(configFilepath)
	assert.Equal(t, nil, err)
}
