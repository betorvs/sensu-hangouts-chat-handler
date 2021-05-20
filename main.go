package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu-community/sensu-plugin-sdk/templates"
	"github.com/sensu/sensu-go/types"
)

// Config represents the handler plugin config.
type Config struct {
	sensu.PluginConfig
	Webhook                      string
	SensuDashboard               string
	WithAnnotations              bool
	WithLabels                   bool
	MessageTemplate              string
	MessageLimit                 int
	DescriptionTemplate          string
	DescriptionLimit             int
	TitlePrettify                bool
	AnnotationsAsLink            string
	AnnotationsSuffixAsLink      string
	AnnotationsSuffixExcludeList string
	EnableThreadKey              bool
	ThreadKeyValue               string
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-hangouts-chat-handler",
			Short:    "The Sensu Go Google Hangsout handler for alerting",
			Keyspace: "sensu.io/plugins/sensu-hangouts-chat-handler/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		{
			Path:      "webhook",
			Env:       "HANGOUTSCHAT_WEBHOOK",
			Argument:  "webhook",
			Shorthand: "w",
			Default:   "",
			Usage:     "The Webhook URL, use default from HANGOUTSCHAT_WEBHOOK env var",
			Value:     &plugin.Webhook,
		},
		{
			Path:      "sensuDashboard",
			Env:       "HANGOUTSCHAT_SENSU_DASHBOARD",
			Argument:  "sensuDashboard",
			Shorthand: "s",
			Default:   "disabled",
			Usage:     "The HANGOUTS Chat Handler will use it to create a source Sensu Dashboard URL. Use HANGOUTSCHAT_SENSU_DASHBOARD. Example: http://sensu-dashboard.example.local/c/~/n",
			Value:     &plugin.SensuDashboard,
		},
		{
			Path:      "withAnnotations",
			Env:       "",
			Argument:  "withAnnotations",
			Shorthand: "a",
			Default:   false,
			Usage:     "Include the event.metadata.Annotations in details to send to Hangouts Chat",
			Value:     &plugin.WithAnnotations,
		},
		{
			Path:      "withLabels",
			Env:       "",
			Argument:  "withLabels",
			Shorthand: "W",
			Default:   false,
			Usage:     "Include the event.metadata.Labels in details to send to Hangouts Chat",
			Value:     &plugin.WithLabels,
		},
		{
			Path:      "messageTemplate",
			Env:       "HANGOUTSCHAT_MESSAGE_TEMPLATE",
			Argument:  "messageTemplate",
			Shorthand: "m",
			Default:   "{{.Entity.Name}}/{{.Check.Name}}",
			Usage:     "The template for the message to be sent",
			Value:     &plugin.MessageTemplate,
		},
		{
			Path:      "messageLimit",
			Env:       "HANGOUTSCHAT_MESSAGE_LIMIT",
			Argument:  "messageLimit",
			Shorthand: "l",
			Default:   130,
			Usage:     "The maximum length of the message field",
			Value:     &plugin.MessageLimit,
		},
		{
			Path:      "descriptionTemplate",
			Env:       "HANGOUTSCHAT_DESCRIPTION_TEMPLATE",
			Argument:  "descriptionTemplate",
			Shorthand: "d",
			Default:   "{{.Check.Output}}",
			Usage:     "The template for the description to be sent",
			Value:     &plugin.DescriptionTemplate,
		},
		{
			Path:      "descriptionLimit",
			Env:       "HANGOUTSCHAT_DESCRIPTION_LIMIT",
			Argument:  "descriptionLimit",
			Shorthand: "L",
			Default:   1500,
			Usage:     "The maximum length of the description field",
			Value:     &plugin.DescriptionLimit,
		},
		{
			Path:      "titlePrettify",
			Env:       "",
			Argument:  "titlePrettify",
			Shorthand: "T",
			Default:   false,
			Usage:     "Remove all -, /, \\ and apply strings.Title in message title",
			Value:     &plugin.TitlePrettify,
		},
		{
			Path:      "annotations-as-link",
			Env:       "HANGOUTSCHAT_ANNOTATIONS_LINK",
			Argument:  "annotations-as-link",
			Shorthand: "A",
			Default:   "",
			Usage:     "Parse Check.metadata.annotations as link to post in Hangouts Chat. e. prometheus_url",
			Value:     &plugin.AnnotationsAsLink,
		},
		{
			Path:      "annotations-suffix-link",
			Env:       "HANGOUTSCHAT_ANNOTATIONS_SUFFIX_LINK",
			Argument:  "annotations-suffix-link",
			Shorthand: "S",
			Default:   "",
			Usage:     "Parse Check.metadata.annotations as link to post in Hangouts Chat. e. prometheus_url",
			Value:     &plugin.AnnotationsSuffixAsLink,
		},
		{
			Path:      "annotations-suffix-exclude",
			Env:       "HANGOUTSCHAT_ANNOTATIONS_SUFFIX_EXCLUDE",
			Argument:  "annotations-suffix-exclude",
			Shorthand: "E",
			Default:   "",
			Usage:     "Parse Check.metadata.annotations as link to post in Hangouts Chat. e. prometheus_url",
			Value:     &plugin.AnnotationsSuffixExcludeList,
		},
		{
			Path:      "threadKey",
			Env:       "",
			Argument:  "threadKey",
			Shorthand: "",
			Default:   false,
			Usage:     "Send message with &threadKey=event.id",
			Value:     &plugin.EnableThreadKey,
		},
		{
			Path:      "threadKey-value",
			Env:       "",
			Argument:  "threadKey-value",
			Shorthand: "",
			Default:   "id",
			Usage:     "Send message with &threadKey=event.id. Options: id or alias (entity.name-check.name)",
			Value:     &plugin.ThreadKeyValue,
		},
	}
)

// SliceCard struct
type SliceCard struct {
	Cards []Cards `json:"cards"`
}

// Header struct
type Header struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	ImageURL string `json:"imageUrl"`
}

// KeyValue struct
type KeyValue struct {
	TopLabel         string `json:"topLabel"`
	Content          string `json:"content"`
	ContentMultiline bool   `json:"contentMultiline"`
}

// OpenLink struct
type OpenLink struct {
	URL string `json:"url"`
}

// OnClick struct
type OnClick struct {
	OpenLink OpenLink `json:"openLink"`
}

// TextButton struct
type TextButton struct {
	Text    string  `json:"text"`
	OnClick OnClick `json:"onClick"`
}

// Buttons struct
type Buttons struct {
	TextButton TextButton `json:"textButton"`
}

// Widgets struct
type Widgets struct {
	KeyValue *KeyValue `json:"keyValue,omitempty"`
	Buttons  []Buttons `json:"buttons,omitempty"`
}

// Sections struct
type Sections struct {
	Widgets []Widgets `json:"widgets"`
}

// Cards struct
type Cards struct {
	Header   Header     `json:"header"`
	Sections []Sections `json:"sections"`
}

// parseDescription func returns string with custom template string to use in description
func parseDescription(event *types.Event) (description string) {
	description, err := templates.EvalTemplate("description", plugin.DescriptionTemplate, event)
	if err != nil {
		return ""
	}
	// allow newlines to get expanded
	description = strings.Replace(description, `\n`, "\n", -1)
	return trim(description, plugin.DescriptionLimit)
}

// parseEventTitle func returns string
func parseEventTitle(event *types.Event) (title string) {
	title, err := templates.EvalTemplate("title", plugin.MessageTemplate, event)
	if err != nil {
		return ""
	}
	if plugin.TitlePrettify {
		newTitle := titlePrettify(title)
		return trim(newTitle, plugin.MessageLimit)
	}
	return trim(title, plugin.MessageLimit)
}

func main() {
	handler := sensu.NewGoHandler(&plugin.PluginConfig, options, checkArgs, executeHandler)
	handler.Execute()
}

func checkArgs(_ *types.Event) error {
	if len(plugin.Webhook) == 0 {
		return fmt.Errorf("webhook url for Hangsout Chat is empty")
	}
	return nil
}

// formattedEventAction func
func formattedEventAction(event *types.Event) string {
	switch event.Check.Status {
	case 0:
		return "RESOLVED"
	default:
		return "ALERT"
	}
}

// parseAnnotationsToButton func
func parseAnnotationsToButton(event *types.Event) []Buttons {
	var button []Buttons

	if plugin.SensuDashboard != "disabled" && checkURL(plugin.SensuDashboard) {
		newbutton := Buttons{}
		newbutton.TextButton.Text = "Sensu Source"
		newbutton.TextButton.OnClick.OpenLink.URL = fmt.Sprintf("%s/%s/events/%s/%s", plugin.SensuDashboard, event.Entity.Namespace, event.Entity.Name, event.Check.Name)
		button = append(button, newbutton)
	}
	if plugin.AnnotationsAsLink != "" {
		tags := annotationsSlice()
		if event.Check.Annotations != nil {
			for key, value := range event.Check.Annotations {
				if stringInSlice(key, tags) && checkURL(value) {
					newbutton := Buttons{}
					newbutton.TextButton.Text = key
					newbutton.TextButton.OnClick.OpenLink.URL = value
					button = append(button, newbutton)
				}
			}
		}
	}
	if plugin.AnnotationsSuffixAsLink != "" {
		// tags := annotationsSlice()
		if event.Check.Annotations != nil {
			for key, value := range event.Check.Annotations {
				if plugin.AnnotationsSuffixExcludeList != "" {
					tags := annotationsExcludeSlice()
					if stringInSlice(key, tags) {
						continue
					}
				}
				// && strings.Contains(key, )
				if strings.HasSuffix(key, plugin.AnnotationsSuffixAsLink) && checkURL(value) {

					newbutton := Buttons{}
					newbutton.TextButton.Text = key
					newbutton.TextButton.OnClick.OpenLink.URL = value
					button = append(button, newbutton)
				}

			}
		}
	}
	if len(button) == 0 {
		newbutton := Buttons{}
		newbutton.TextButton.Text = "Sensu Documentation"
		newbutton.TextButton.OnClick.OpenLink.URL = "https://docs.sensu.io/sensu-go/latest"
		button = append(button, newbutton)
	}
	return button
}

// // eventDescription func return an message to use it
func eventDescription(event *types.Event) string {
	var (
		annotations string
		labels      string
		message     string
	)
	message = fmt.Sprintf("Entity: %s, \nCheck: %s, \nCommand: %s, ", event.Entity.Name, event.Check.Name, event.Check.Command)
	if event.Check.ProxyEntityName != "" {
		message += fmt.Sprintf("\nProxy_Entity: %s, \n", event.Check.ProxyEntityName)
	}
	if plugin.WithAnnotations {
		if event.Annotations != nil {
			for key, value := range event.Annotations {
				if validateDescription(key) {
					annotations += fmt.Sprintf("%s_%s: %s, \n", "event", key, value)
				}
			}
		}
		if event.Check.Annotations != nil {
			for key, value := range event.Check.Annotations {
				if validateDescription(key) {
					annotations += fmt.Sprintf("%s_%s: %s, \n", "check", key, value)
				}
			}
		}
		if event.Entity.Annotations != nil {
			for key, value := range event.Check.Annotations {
				if validateDescription(key) {
					annotations += fmt.Sprintf("%s_%s: %s, \n", "entity", key, value)
				}
			}
		}
		if annotations != "" {
			message += fmt.Sprintf("\n Annotations: \n%s", annotations)
		}
	}
	if plugin.WithLabels {
		if event.Labels != nil {
			for key, value := range event.Labels {
				if validateDescription(key) {
					labels += fmt.Sprintf("%s_%s: %s, \n", "event", key, value)
				}
			}
		}
		if event.Check.Labels != nil {
			for key, value := range event.Check.Labels {
				if validateDescription(key) {
					labels += fmt.Sprintf("%s_%s: %s, \n", "check", key, value)
				}
			}
		}
		if event.Entity.Labels != nil {
			for key, value := range event.Check.Labels {
				if validateDescription(key) {
					labels += fmt.Sprintf("%s_%s: %s, \n", "entity", key, value)
				}
			}
		}
		if labels != "" {
			message += fmt.Sprintf("\n Labels: \n%s", labels)
		}
	}
	return message
}

func validateDescription(key string) bool {
	// tags := []string{plugin.Keyspace}
	if plugin.AnnotationsAsLink != "" {
		tags := annotationsSlice()
		if !stringInSlice(key, tags) && !strings.Contains(key, plugin.Keyspace) {
			return true
		}
	}

	return false
}

// run func do everything
func executeHandler(event *types.Event) error {

	keyvalue1 := KeyValue{
		TopLabel:         formattedEventAction(event),
		Content:          eventDescription(event),
		ContentMultiline: true,
	}
	// keyvalue2 := KeyValue{
	// 	TopLabel: "More Information",
	// 	Content:  parseAnnotations(event),
	// }
	keyvalue3 := KeyValue{
		TopLabel:         "Check Output",
		Content:          parseDescription(event),
		ContentMultiline: true,
	}
	widget1 := Widgets{
		KeyValue: &keyvalue1,
	}
	widget2 := Widgets{
		// KeyValue: &keyvalue2,
		Buttons: parseAnnotationsToButton(event),
	}
	widget3 := Widgets{
		KeyValue: &keyvalue3,
	}
	header := Header{
		Title:    "Sensu Event",
		Subtitle: parseEventTitle(event),
	}
	section1 := Sections{
		Widgets: []Widgets{widget1},
	}
	section2 := Sections{
		Widgets: []Widgets{widget2},
	}
	section3 := Sections{
		Widgets: []Widgets{widget3},
	}
	card := Cards{
		Header:   header,
		Sections: []Sections{section1, section2, section3},
	}
	formPost := SliceCard{
		Cards: []Cards{card},
	}
	// prettyJSON, err := json.MarshalIndent(formPost, "", "  ")
	// if err != nil {
	// 	fmt.Printf("Failed to generate json %v", err)
	// }
	// fmt.Printf("%s\n", string(prettyJSON))
	bodymarshal, err := json.Marshal(&formPost)
	if err != nil {
		fmt.Printf("[ERROR] %s", err)
	}
	webhook := plugin.Webhook
	if plugin.EnableThreadKey {
		switch plugin.ThreadKeyValue {
		case "id":
			webhook = fmt.Sprintf("%s&threadKey=%s", plugin.Webhook, event.GetUUID())
		case "alias":
			webhook = fmt.Sprintf("%s&threadKey=%s", plugin.Webhook, parseAlias(event))
		default:
			webhook = fmt.Sprintf("%s&threadKey=%s", plugin.Webhook, event.GetUUID())
		}
	}
	err = Post(webhook, bodymarshal)
	if err != nil {
		fmt.Printf("[ERROR] %s", err)
	}
	return nil
}

//Post func to send the json to hangouts chat
func Post(url string, body []byte) error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("[ERROR] %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[ERROR] %s", err)
	}
	if resp.StatusCode != 200 {
		bodyText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("[ERROR] %s", err)
		}
		s := string(bodyText)
		fmt.Printf("[LOG]: %s ; %s", resp.Status, s)
	}
	defer resp.Body.Close()
	return nil
}

func annotationsSlice() []string {
	tags := []string{}
	if strings.Contains(plugin.AnnotationsAsLink, ",") {
		tags = strings.Split(plugin.AnnotationsAsLink, ",")
	}
	if !strings.Contains(plugin.AnnotationsAsLink, ",") {
		tags = append(tags, plugin.AnnotationsAsLink)
	}

	return tags
}

func annotationsExcludeSlice() []string {
	tags := []string{}
	if strings.Contains(plugin.AnnotationsSuffixExcludeList, ",") {
		tags = strings.Split(plugin.AnnotationsSuffixExcludeList, ",")
	}
	if !strings.Contains(plugin.AnnotationsSuffixExcludeList, ",") {
		tags = append(tags, plugin.AnnotationsSuffixExcludeList)
	}

	return tags
}

// use to parse annotations to send as link
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// time func returns only the first n bytes of a string
func trim(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

func titlePrettify(s string) string {
	var title string
	title = strings.ReplaceAll(s, "-", " ")
	title = strings.ReplaceAll(title, "\\", " ")
	title = strings.ReplaceAll(title, "/", " ")
	title = strings.Title(title)
	return title
}

func parseAlias(event *types.Event) string {
	return fmt.Sprintf("%s-%s", event.Entity.Name, event.Check.Name)
}

func checkURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
