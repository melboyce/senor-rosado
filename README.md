# Señor Rosado

[![Build Status](https://travis-ci.org/weirdtales/senor-rosado.svg?branch=master)](https://travis-ci.org/weirdtales/senor-rosado)


## Summary

This is a Slack bot written in Go.

I'd prefer to be using `plugin`, but it's too raw at the moment.


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
docker run --env-file secrets.env sr
```


## Adding Commands

All commands are namespaced in the `cmds/` directory. Take a gander at
[hello.go](cmds/cmdhello/hello.go) for a simple overview.
Each command must implement two functions, described below.


### Register()

`Register()` returns two strings:

* `regexp` used to match incoming command inputs
* `help` text used by the built-in command `help`


### Respond()

`Respond(reply slack.Reply, c chan slack.Reply)` should populate the supplied
`slack.Reply` with an appropriate payload and then push it onto the supplied
channel.


## Utility Functions

The `slack` package provides some conveniences:

|Func|Utility|
|:---|:------|
|`GetJSON(url string, target interface{}) error`|mutate a &struct with a JSON unmarshal|
