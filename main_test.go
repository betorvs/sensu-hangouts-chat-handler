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

func TestAnnotationsSlice(t *testing.T) {
	expectedTags := []string{"runbook_url"}
	plugin.AnnotationsAsLink = "runbook_url"
	tags := annotationsSlice()
	assert.Equal(t, tags, expectedTags)
	expectedTags2 := []string{"runbook_url", "docs_url"}
	plugin.AnnotationsAsLink = "runbook_url,docs_url"
	tags2 := annotationsSlice()
	assert.Equal(t, tags2, expectedTags2)
}

func TestStringInSlice(t *testing.T) {
	testSlice := []string{"foo", "bar", "test"}
	testString := "test"
	testResult := stringInSlice(testString, testSlice)
	assert.True(t, testResult)
}

func TestValidateDescription(t *testing.T) {
	event := types.FixtureEvent("entity1", "check1")
	event.Check.Annotations["runbook_url"] = "https://play.golang.org"
	event.Check.Annotations["sensu.io/plugins/sensu-hangouts-chat-handler/config/webhook"] = "https://LongWebhookURLHere.com"
	plugin.AnnotationsAsLink = "runbook_url"
	test1 := validateDescription("runbook_url")
	assert.Equal(t, test1, false)
	test2 := validateDescription("sensu.io/plugins/sensu-hangouts-chat-handler/config/webhook")
	assert.Equal(t, test2, false)
}

func TestEventDescription(t *testing.T) {
	event := types.FixtureEvent("entity1", "check1")
	event.Check.Annotations["runbook_url"] = "https://play.golang.org"
	plugin.AnnotationsAsLink = "runbook_url"
	test1 := eventDescription(event)
	assert.Contains(t, test1, "check1")
	assert.NotContains(t, test1, "check_runbook_url")
}

func TestTrim(t *testing.T) {
	testString := "This string is 33 characters long"
	assert.Equal(t, trim(testString, 40), testString)
	assert.Equal(t, trim(testString, 4), "This")
}

func TestAnnotationsExcludeSlice(t *testing.T) {
	expectedTags := []string{"runbook_url"}
	plugin.AnnotationsSuffixExcludeList = "runbook_url"
	tags := annotationsExcludeSlice()
	assert.Equal(t, tags, expectedTags)
	expectedTags2 := []string{"runbook_url", "docs_url"}
	plugin.AnnotationsSuffixExcludeList = "runbook_url,docs_url"
	tags2 := annotationsExcludeSlice()
	assert.Equal(t, tags2, expectedTags2)
}

func TestTitlePrettify(t *testing.T) {
	test1 := "long-check-with-too-many-dashes"
	res1 := titlePrettify(test1)
	val1 := "Long Check With Too Many Dashes"
	assert.Equal(t, val1, res1)
	test2 := "long-check-with-too-many-dashes/and/slashes-and\\others"
	res2 := titlePrettify(test2)
	val2 := "Long Check With Too Many Dashes And Slashes And Others"
	assert.Equal(t, val2, res2)
}
