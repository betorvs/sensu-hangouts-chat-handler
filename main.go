package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sensu/sensu-go/types"
	"github.com/spf13/cobra"
)

var (
	webhook        string
	annotations    string
	sensuDashboard string
	stdin          *os.File
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

func main() {
	rootCmd := configureRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func configureRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sensu-hangouts-chat-handler",
		Short: "The Sensu Go Hangouts Chat handler for incident alerting",
		RunE:  run,
	}
	cmd.Flags().StringVarP(&webhook,
		"webhook",
		"w",
		os.Getenv("WEBHOOK_HANGOUTSCHAT"),
		"The Webhook URL, use default from WEBHOOK_HANGOUTSCHAT env var")

	cmd.Flags().StringVarP(&annotations,
		"withAnnotations",
		"a",
		os.Getenv("HANGOUTSCHAT_ANNOTATIONS"),
		"The Hangouts Chat Handler will parse check and entity annotations with these values. Use HANGOUTSCHAT_ANNOTATIONS env var with commas, like: documentation,playbook")

	cmd.Flags().StringVar(&sensuDashboard,
		"sensuDashboard",
		os.Getenv("HANGOUTSCHAT_SENSU_DASHBOARD"),
		"The HANGOUTS Chat Handler will use it to create a source Sensu Dashboard URL. Use HANGOUTSCHAT_SENSU_DASHBOARD. Example: http://sensu-dashboard.example.local/c/~/n")

	return cmd
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

// stringInSlice checks if a slice contains a specific string
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// parseAnnotationsToButton func
func parseAnnotationsToButton(event *types.Event) []Buttons {
	var button []Buttons
	tags := strings.Split(annotations, ",")
	if event.Check.Annotations != nil {
		for key, value := range event.Check.Annotations {
			if stringInSlice(key, tags) {
				newbutton := Buttons{}
				newbutton.TextButton.Text = fmt.Sprintf("Check %s", key)
				newbutton.TextButton.OnClick.OpenLink.URL = value
				button = append(button, newbutton)
			}
		}
	}
	if event.Entity.Annotations != nil {
		for key, value := range event.Check.Annotations {
			if stringInSlice(key, tags) {
				newbutton := Buttons{}
				newbutton.TextButton.Text = fmt.Sprintf("Entity %s", key)
				newbutton.TextButton.OnClick.OpenLink.URL = value
				button = append(button, newbutton)
			}
		}
	}
	if sensuDashboard != "disabled" {
		newbutton := Buttons{}
		newbutton.TextButton.Text = "Link Event Sensu Source"
		newbutton.TextButton.OnClick.OpenLink.URL = fmt.Sprintf("%s/%s/events/%s/%s", sensuDashboard, event.Entity.Namespace, event.Entity.Name, event.Check.Name)
		button = append(button, newbutton)
	}
	if len(button) == 0 {
		newbutton := Buttons{}
		newbutton.TextButton.Text = "Link Sensu Documentation"
		newbutton.TextButton.OnClick.OpenLink.URL = "https://docs.sensu.io/sensu-go/latest"
		button = append(button, newbutton)
	}
	return button
}

// eventDescription func return an message to use it
func eventDescription(event *types.Event) string {
	return fmt.Sprintf("Entity: %s, \nCheck: %s, \n", event.Entity.Name, event.Check.Name)
}

// run func do everything
func run(cmd *cobra.Command, args []string) error {
	if webhook == "" {
		_ = cmd.Help()
		return fmt.Errorf("webhook is empty")

	}
	if stdin == nil {
		stdin = os.Stdin
	}
	eventJSON, err := ioutil.ReadAll(stdin)
	if err != nil {
		return fmt.Errorf("failed to read stdin: %s", err)
	}

	if annotations == "" {
		annotations = "documentation,playbook"
	}

	if sensuDashboard == "" {
		sensuDashboard = "disabled"
	}

	event := &types.Event{}
	err = json.Unmarshal(eventJSON, event)
	if err != nil {
		return fmt.Errorf("failed to unmarshal stdin data: %s", err)
	}

	if err = event.Validate(); err != nil {
		return fmt.Errorf("failed to validate event: %s", err)
	}

	if !event.HasCheck() {
		return fmt.Errorf("event does not contain check")
	}
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
		Content:          event.Check.Output,
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
		Title:    "Sensu Event (Entity/Check)",
		Subtitle: fmt.Sprintf("%s/%s", event.Entity.Name, event.Check.Name),
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
	// 	log.Fatal("Failed to generate json", err)
	// }
	// fmt.Printf("%s\n", string(prettyJSON))
	bodymarshal, err := json.Marshal(&formPost)
	if err != nil {
		fmt.Printf("[ERROR] %s", err)
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
