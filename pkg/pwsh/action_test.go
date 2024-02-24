package pwsh

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

func TestMixin_UnmarshalStep_InlineScript(t *testing.T) {
	b, err := os.ReadFile("testdata/step-input.yaml")
	require.NoError(t, err)

	var action Action
	err = yaml.Unmarshal(b, &action)
	require.NoError(t, err)
	assert.Equal(t, "install", action.Name)
	require.Len(t, action.Steps, 1)

	step := action.Steps[0]
	assert.Equal(t, "Summon Minion", step.Description)
	assert.NotEmpty(t, step.Outputs)
	assert.Equal(t, Output{Name: "VICTORY", JsonPath: "$Id"}, step.Outputs[0])

	assert.Equal(t, "Write-Host \"VICTORY\"", step.InlineScript)

	require.Len(t, step.Arguments, 2)
	assert.Equal(t, "value1", step.Arguments[0])
	assert.Equal(t, "value2", step.Arguments[1])
}

func TestMixin_UnmarshalStep_File(t *testing.T) {
	b, err := os.ReadFile("testdata/step-file.yaml")
	require.NoError(t, err)

	var action Action
	err = yaml.Unmarshal(b, &action)
	require.NoError(t, err)
	assert.Equal(t, "install", action.Name)
	require.Len(t, action.Steps, 1)

	step := action.Steps[0]
	assert.Equal(t, "Summon Minion", step.Description)
	assert.NotEmpty(t, step.Outputs)
	assert.Equal(t, Output{Name: "VICTORY", JsonPath: "$Id"}, step.Outputs[0])

	assert.Equal(t, "./helper.ps1", step.File)

	require.Len(t, step.Arguments, 2)
	assert.Equal(t, "value1", step.Arguments[0])
	assert.Equal(t, "value2", step.Arguments[1])
}
