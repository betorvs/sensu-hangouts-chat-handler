# Sensu Go Hangouts Chat Handler
![Go Test](https://github.com/betorvs/sensu-hangouts-chat-handler/workflows/Go%20Test/badge.svg)
[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/betorvs/sensu-hangouts-chat-handler)

The Sensu Go Hangouts Chat Handler is a [Sensu Event Handler][3] which manages
[Hangouts Chat][2] for alerting purposes. With this handler,
[Sensu][1] can alert systems administrators in Hangouts Chats.

This handler was inspired by [hangouts ruby plugin][5].

## Installation

Download the latest version of the sensu-hangouts-chat-handler from [releases][4],
or create an executable script from this source.

From the local path of the sensu-hangouts-chat-handler repository:
```
go build -o /usr/local/bin/sensu-hangouts-chat-handler main.go
```

## Configuration

Example Sensu Go handler definition:

```yml
type: Handler
api_version: core/v2
metadata:
  name: hangouts-chat
  namespace: default
spec:
  type: pipe
  command: sensu-hangouts-chat-handler
  env_vars:
  - HANGOUTSCHAT_WEBHOOK="https://...."
  timeout: 10
  runtime_assets:
  - betorvs/sensu-hangouts-chat-handler
  filters:
  - is_incident
```

Example Sensu Go check definition:

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: dummy-app-healthz
  namespace: default
  annotations:
    documentation: https://docs.sensu.io/sensu-go/latest
spec:
  command: check-http -u http://localhost:8080/healthz
  subscriptions:
  - dummy
  handlers:
  - hangouts-chat
  interval: 60
  publish: true
```


## Usage Examples

Help:
```

The Sensu Go Google Hangsout handler for alerting

Usage:
  sensu-hangouts-chat-handler [flags]
  sensu-hangouts-chat-handler [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -A, --annotations-as-link string          Parse Check.metadata.annotations as link to post in Hangouts Chat. e. prometheus_url
  -E, --annotations-suffix-exclude string   Parse Check.metadata.annotations as link to post in Hangouts Chat. e. prometheus_url
  -S, --annotations-suffix-link string      Parse Check.metadata.annotations as link to post in Hangouts Chat. e. prometheus_url
  -L, --descriptionLimit int                The maximum length of the description field (default 1500)
  -d, --descriptionTemplate string          The template for the description to be sent (default "{{.Check.Output}}")
  -h, --help                                help for sensu-hangouts-chat-handler
  -l, --messageLimit int                    The maximum length of the message field (default 130)
  -m, --messageTemplate string              The template for the message to be sent (default "{{.Entity.Name}}/{{.Check.Name}}")
  -s, --sensuDashboard string               The HANGOUTS Chat Handler will use it to create a source Sensu Dashboard URL. Use HANGOUTSCHAT_SENSU_DASHBOARD. Example: http://sensu-dashboard.example.local/c/~/n (default "disabled")
  -T, --titlePrettify                       Remove all -, /, \ and apply strings.Title in message title
  -w, --webhook string                      The Webhook URL, use default from HANGOUTSCHAT_WEBHOOK env var
  -a, --withAnnotations                     Include the event.metadata.Annotations in details to send to Hangouts Chat
  -W, --withLabels                          Include the event.metadata.Labels in details to send to Hangouts Chat

Use "sensu-hangouts-chat-handler [command] --help" for more information about a command.


```

**Note:** Make sure to set the `HANGOUTSCHAT_WEBHOOK` environment variable for sensitive credentials in production to prevent leaking into system process table. Please remember command arguments can be viewed by unprivileged users using commands such as `ps` or `top`. The `--webhook` argument is provided as an override primarily for testing purposes. 

### Argument Annotations

All arguments for this handler are tunable on a per entity or check basis based on annotations.  The
annotations keyspace for this handler is `sensu.io/plugins/sensu-hangouts-chat-handler/config`. 

#### Examples

To change the team argument for a particular check, for that checks's metadata add the following:

```yml
type: CheckConfig
api_version: core/v2
metadata:
  annotations:
    sensu.io/plugins/sensu-hangouts-chat-handler/config/webhook: "https://LongWebhookURLHere"
[...]
```


### Asset creation

The easiest way to get this handler added to your Sensu environment, is to add it as an asset from Bonsai:

```sh
sensuctl asset add betorvs/sensu-hangouts-chat-handler --rename sensu-hangouts-chat-handler
```

See `sensuctl asset --help` for details on how to specify version.

## Contributing

See https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md

[1]: https://github.com/sensu/sensu-go
[2]: https://developers.google.com/hangouts/chat
[3]: https://docs.sensu.io/sensu-go/5.0/reference/handlers/#how-do-sensu-handlers-work
[4]: https://github.com/betorvs/sensu-hangouts-chat-handler/releases
[5]: https://github.com/clevertoday/sensu-plugins-hangouts-chat