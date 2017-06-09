# Señor Rosado

## Summary

This is a Slack bot written in Go.


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

Commands are Go plugins the source for which should be in the
[_cartridges](_cartridges/) dir. The directory name is prefixed with
an underscore so `go build` doesn't compile its contents automatically.

Commands (cartridges) must be `package main` and must implement the following
functions:

* `Register() string`
* `Respond(slack.Message, slack.Conn, []string)`

These plugins must be built thusly:

```
go build -buildmode=plugin -o plugins/PLUG.so _cartridges/PLUG.go
```

The [Makefile](Makefile) has a target for doing this en masse.


### Register()

This function returns a string used as a regexp which is matched against
the input text sans the bot's UID. If there is a match, the `Respond()`
function is called.


### Respond(slack.Message, slack.Conn, []string)

This function forumlates a reply and sends it. See
[hello.go](_cartridges/hello.go) for an example of how to do this. The
`[]string` is the list of matches that the regexp from `Register()` matched.
