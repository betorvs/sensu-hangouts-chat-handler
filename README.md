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
sensuctl asset create sensu-hangouts-chat-handler --url "https://assets.bonsai.sensu.io/1daec49623e9384d5374f7e11f12a343cf374e5f/sensu-hangouts-chat-handler_0.0.1_linux_amd64.tar.gz" --sha512 "59fd8fd9909819ad9eb1897814e38903634e6ed38cadc50a2ad75a069b466bc9b4eb23bc33ce51fc58bb8cff47b931256b0e2e8cfe60a89c4ebfafca097e8c45"
```


## Contributing

See https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md

[1]: https://github.com/sensu/sensu-go
[2]: https://developers.google.com/hangouts/chat
[3]: https://docs.sensu.io/sensu-go/5.0/reference/handlers/#how-do-sensu-handlers-work
[4]: https://github.com/betorvs/sensu-hangouts-chat-handler/releases
