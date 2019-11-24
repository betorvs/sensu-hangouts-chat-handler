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
	webhook     string
	annotations string
	stdin       *os.File
)

// payload struct for post in sensu
type payload struct {
	Text string `json:"text"`
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

// parseAnnotations func try to find a predeterminated keys
func parseAnnotations(event *types.Event) string {
	var output string
	// localannotations := make(map[string]string)
	tags := strings.Split(annotations, ",")
	if event.Check.Annotations != nil {
		for key, value := range event.Check.Annotations {
			if stringInSlice(key, tags) {
				output += fmt.Sprintf("  %s: %s ,\n", key, value)
			}
		}
	}
	if event.Entity.Annotations != nil {
		for key, value := range event.Check.Annotations {
			if stringInSlice(key, tags) {
				output += fmt.Sprintf("  %s: %s ,\n", key, value)
			}
		}
	}
	output += fmt.Sprintf("Check output: %s", event.Check.Output)

	return output
}

// eventDescription func return an message to use it
func eventDescription(event *types.Event) string {
	return fmt.Sprintf("*%s*\nServer: %s, \nCheck: %s, \nMore Information:\n%s", formattedEventAction(event), event.Entity.Name, event.Check.Name, parseAnnotations(event))
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
	formPost := payload{
		Text: eventDescription(event),
	}
	bodymarshal, err := json.Marshal(&formPost)
	if err != nil {
		fmt.Printf("[ERROR] %s", err)
	}
	Post(webhook, bodymarshal)
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
