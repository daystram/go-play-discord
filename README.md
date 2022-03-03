# go-play-discord

[![Github Actions](https://github.com/daystram/go-play-discord/actions/workflows/ci.yml/badge.svg)](https://github.com/daystram/go-play-discord/actions/workflows/ci.yml)
[![Docker Pulls](https://img.shields.io/docker/pulls/daystram/go-play-discord)](https://hub.docker.com/r/daystram/go-play-discord)
[![MIT License](https://img.shields.io/github/license/daystram/go-play-discord)](https://github.com/daystram/go-play-discord/blob/master/LICENSE)

A Discord bot to run and format Go code via the [Go Playground](https://go.dev/play).

## Installation

### Go version < 1.16

```shell
$ go get -u github.com/daystram/go-play-discord/cmd/go-play-discord
```

### Go 1.16+

```shell
$ go install github.com/daystram/go-play-discord/cmd/go-play-discord@latest
```

## Usage

After providing the required configuration, the bot can simply be run as follows:

```shell
$ go-play-discord
```

### Docker

Instead of installing the command itself, you can run the bot via Docker:

```shell
$ docker run --name go-play-discord --env-file ./.env -d daystram/go-play-discord
```

## Configuration

The bot could be configured by setting the following environment variables.

| Name        | Description       | Default | Required |
| ----------- | ----------------- | ------- | -------- |
| `BOT_TOKEN` | Discord Bot token | `""`    | ✅       |

## License

This project is licensed under the [MIT License](./LICENSE).
