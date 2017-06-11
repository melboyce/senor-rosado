# Señor Rosado

[![Build Status](https://travis-ci.org/weirdtales/senor-rosado.svg?branch=master)](https://travis-ci.org/weirdtales/senor-rosado)

## TODO

* plugins are shit
* dockerfile builds plugins manually

## Summary

This is a Slack bot written in Go.

Of particular note:

* `plugin` is really new and a little raw.
* `plugin` does not support reloading plugins.
* `plugin` does not allow `CGO_ENABLED=0`.

This means until the `plugin` library is matured, the plugin architecture
costs a lot in terms of size and only allows new plugins to be added, but
loaded plugins cannot be modified without restarting the bot.

There's also a non-trivial amount of compile-time involved with building
Go plugins at the moment.


## Running It

There are some pre-requisites:

```
1. get a Slack bot token
2. get a DarkSky API token
3. get a Giphy API token
```

Create a file called `secrets.env` and populate it with the following:

```
SLACK_TOKEN=xxx
DARKSKY_TOKEN=yyy
GIPHY_TOKEN=zzz
```

Build the image:

```
docker build -t sr .
```

Execute the Docker image:

```
docker run -it -v --env-file secrets.env sr
```


## Adding Commands

Commands are Go plugins the source for which should be in the
[_cartridges](_cartridges/) dir. The directory name is prefixed with
an underscore so `go build` doesn't compile its contents automatically.

Commands (cartridges) must be `package main` and must implement the following
functions:

* `Register() (string, string)`
* `Respond(slack.Message, slack.Conn, []string)`

These plugins must be built thusly:

```
go build -buildmode=plugin -o plugins/PLUG.so _cartridges/PLUG.go
```

The [Makefile](Makefile) has a target for doing this en masse.


### Register()

The two strings returned by `Register()` are:

* `regexp` used to match incoming command inputs
* `help` text used by the built-in command `help`


### Respond(slack.Message, slack.Conn, []string)

This function forumlates a reply and sends it. See
[hello.go](_cartridges/hello.go) for an example of how to do this. The
`[]string` is the list of matches that the regexp from `Register()` matched.


## Utility Functions

The `slack` package provides some conveniences:

|Func|Utility|
|:---|:------|
|`GetJSON(url string, target interface{}) error`|mutate a &struct with a JSON unmarshal|
