# Señor Rosado

## Summary

This is a Slack bot written in Go. Two things of note:

* this is my first crack at Go, so it's inelegant (aka shitty)
* ???


## Running It

There are some pre-requisites:

```
1. get a Slack bot token
2. get a DarkSky API token
```

Create a file called `secrets.env` and populate it with the following:

```
SLACK_TOKEN=xxx
DARKSKY_TOKEN=yyy
```

Build the image:

`docker build -t senor-rosado .`

Execute the Docker image:

`docker run -it --env-file secrets.env senor-rosado`


## Adding Commands

`main.go` has an example, but in short:

* commands are go functions
* commands are executed from a map
* the map is keyed with the name of the command
* the map values are functions that execute some go

There's a type used for the map:

```go
map[string]slack.ChatFn
```

Both `ChatLoop` and `ChatFn` are in `slack/chat.go`.
