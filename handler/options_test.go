package handler

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
)

func resetEnvAndArgs() {
	os.Args = make([]string, 0)
	os.Clearenv()
}

func TestEnvVarOverwritesDefault(t *testing.T) {
	defer resetEnvAndArgs()

	defaultOptions, err := DefaultOptions()
	assert.Nil(t, err)
	os.Setenv("SPAS_PORT", "8090")

	options, err := GetOptions()
	assert.Nil(t, err)

	assert.NotEqual(t, defaultOptions.Port, "8090")
	assert.Equal(t, "8090", options.Port)
}

func TestCommandLineArgOverwritesDefaultAndEnvVar(t *testing.T) {
	defer resetEnvAndArgs()

	defaultOptions, err := DefaultOptions()
	assert.Nil(t, err)
	os.Setenv("SPAS_PORT", "8090")
	os.Args = []string { "--port", "8091"}

	options, err := GetOptions()
	assert.Nil(t, err)

	assert.NotEqual(t, defaultOptions.Port, "8091")
	assert.Equal(t, "8091", options.Port)
}

func TestConfigFileOverwritesDefault(t *testing.T) {
	defer resetEnvAndArgs()

	defaultOptions, err := DefaultOptions()
	assert.Nil(t, err)
	os.Setenv("SPAS_CONFIGFILE", "../test_resources/spas.config.json")

	options, err := GetOptions()
	assert.Nil(t, err)

	assert.NotEqual(t, defaultOptions.Port, "8089")
	assert.Equal(t, "8089", options.Port)
}