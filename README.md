# Sensu Go Hangouts Chat Handler
TravisCI: [![Build Status](https://travis-ci.org/betorvs/sensu-hangouts-chat-handler.svg?branch=master)](https://travis-ci.org/betorvs/sensu-hangouts-chat-handler)

The Sensu Go Hangouts Chat Handler is a [Sensu Event Handler][3] which manages
[Hangouts Chat][2] for alerting purposes. With this handler,
[Sensu][1] can alert systems administrators in Hangouts Chats.

## Installation

Download the latest version of the sensu-hangouts-chat-handler from [releases][4],
or create an executable script from this source.

From the local path of the sensu-hangouts-chat-handler repository:
```
go build -o /usr/local/bin/sensu-hangouts-chat-handler main.go
```

## Configuration

Example Sensu Go handler definition:

```json
{
    "api_version": "core/v2",
    "type": "Handler",
    "metadata": {
        "namespace": "default",
        "name": "hangouts-chat"
    },
    "spec": {
        "type": "pipe",
        "command": "sensu-hangouts-chat-handler",
        "env_vars": [
          "WEBHOOK_HANGOUTSCHAT=https://...."
        ],
        "timeout": 10,
        "filters": [
            "is_incident"
        ]
    }
}
```

Example Sensu Go check definition:

```json
{
    "api_version": "core/v2",
    "type": "CheckConfig",
    "metadata": {
        "namespace": "default",
        "name": "dummy-app-healthz"
    },
    "spec": {
        "command": "check-http -u http://localhost:8080/healthz",
        "subscriptions":[
            "dummy"
        ],
        "publish": true,
        "interval": 10,
        "handlers": [
            "hangouts-chat"
        ]
    }
}
```

## Usage Examples

Help:
```
Usage:
  sensu-hangouts-chat-handler [flags]

Flags:
  -w, --webhook string   The Webhook from Hangouts Chat, use default from WEBHOOK_HANGOUTSCHAT env var
  -h, --help          help for sensu-opsgenie-handler

```

**Note:** Make sure to set the `WEBHOOK_HANGOUTSCHAT` environment variable for sensitive credentials in production to prevent leaking into system process table. Please remember command arguments can be viewed by unprivileged users using commands such as `ps` or `top`. The `--auth` argument is provided as an override primarily for testing purposes. 

### Asset creation

Example: 

```sh
sensuctl asset create sensu-hangouts-chat-handler --url "" --sha512 ""
```


## Contributing

See https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md

[1]: https://github.com/sensu/sensu-go
[2]: https://developers.google.com/hangouts/chat
[3]: https://docs.sensu.io/sensu-go/5.0/reference/handlers/#how-do-sensu-handlers-work
[4]: https://github.com/betorvs/sensu-hangouts-chat-handler/releases
