[![Go Report Card](https://goreportcard.com/badge/github.com/lobz1g/page-monitoring)](https://goreportcard.com/report/github.com/lobz1g/page-monitoring)
[![Go](https://github.com/lobz1g/page-monitoring/actions/workflows/go.yml/badge.svg)](https://github.com/lobz1g/page-monitoring/actions/workflows/go.yml)
[![Build](https://github.com/lobz1g/page-monitoring/actions/workflows/release.yaml/badge.svg)](https://github.com/lobz1g/page-monitoring/actions/workflows/release.yaml)

# Page monitoring

It is a small utility for monitoring pages. If something was changed in pages, you can receive notification in Telegram.

## How it works

This utility checks given pages and if some information was changed or the page doesn't answer correctly, the utility
sends notification to your telegram channel.

## How to run

You should:

1. Create a telegram bot. [Here](https://www.google.com) is an instruction.
2. Create a telegram channel and set the bot (from a previous step) as administrator of the channel.
3. Run utility and wait.

## Configuration

You need to change the configuration file `config.json`. By the way, you can see example data in the config file.

### Fields

- `debug` - is boolean value. If you need more logging set it `true`. By default, value is `false`.
- `token` - Your telegram bot token.
- `channel` - Name of channel where bot will send messages about status of pages.
- `timeout` - Timeout before next check the page. By default, the value is `30m`. `timeout` string is a possibly signed
  sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms" or "2h45m". Valid time
  units are "ns", "us" (or "?s"), "ms", "s", "m", "h".
- `url` - Array of pages. You need to paste the full url of the page.

## Run local

You should change config.json file before start. If you want to stop the app, just type `exit` in the console.

### Source

```shell
go run main.go
```

### Binary

Open console in the folder where the binary is and type this command

#### For Linux

```shell
./page-monitoring
```

#### For Windows

```shell
page-monitoring.exe
```

## Run on Docker

```shell
docker build -t monitoring .
docker run --name monitoring monitoring
```
