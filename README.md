[![Test(https://github.com/jacobweinstock/rollzap/workflows/Test/badge.svg)](https://github.com/jacobweinstock/rollzap/actions?query=workflow%3A%22Test)
[![Go Report](https://goreportcard.com/badge/github.com/jacobweinstock/rollzap)](https://goreportcard.com/report/github.com/jacobweinstock/rollzap)

# Rollbar Zap

A simple zapcore.Core implementation to integrate with Rollbar. (forked and modified from [https://github.com/bearcherian/rollzap](https://github.com/bearcherian/rollzap))

To use, initialize rollbar like normal, create a new RollbarCore, then wrap with a NewTee. [See the example code](example/main.go) for a detailed example.

## Testing 

To test this code use `RC_TOKEN=MY_ROLLBAR_TOKEN go test`
