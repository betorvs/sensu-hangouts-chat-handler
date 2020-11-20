package main

import (
	"encoding/json"
	"testing"

	"github.com/sensu/sensu-go/types"
	"github.com/stretchr/testify/assert"
)

func TestParseEventKeyTags(t *testing.T) {
	event := types.FixtureEvent("foo", "bar")
	_, err := json.Marshal(event)
	assert.NoError(t, err)
	plugin.MessageTemplate = "{{.Entity.Name}}/{{.Check.Name}}"
	plugin.MessageLimit = 100
	title := parseEventTitle(event)
	assert.Contains(t, title, "foo")
}

func TestParseDescription(t *testing.T) {
	event := types.FixtureEvent("foo", "bar")
	event.Check.Output = "Check OK"
	_, err := json.Marshal(event)
	assert.NoError(t, err)
	plugin.DescriptionTemplate = "{{.Check.Output}}"
	plugin.DescriptionLimit = 100
	description := parseDescription(event)
	assert.Equal(t, description, "Check OK")
}

func TestCheckArgs(t *testing.T) {
	assert := assert.New(t)
	event := types.FixtureEvent("entity1", "check1")
	assert.Error(checkArgs(event))
	plugin.Webhook = "https://longwebhook.com"
	assert.NoError(checkArgs(event))
}

func TestTrim(t *testing.T) {
	testString := "This string is 33 characters long"
	assert.Equal(t, trim(testString, 40), testString)
	assert.Equal(t, trim(testString, 4), "This")
}
