# runner
[![GoDoc](https://godoc.org/github.com/sbrow/runner?status.svg)](https://godoc.org/github.com/sbrow/runner) [![Build Status](https://travis-ci.org/sbrow/runner.svg?branch=master)](https://travis-ci.org/sbrow/runner) [![Coverage Status](https://coveralls.io/repos/github/sbrow/runner/badge.svg?branch=master)](https://coveralls.io/github/sbrow/runner?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/sbrow/runner)](https://goreportcard.com/report/github.com/sbrow/runner)

Package runner runs the non-go code that Photoshop understands, and passes it to
back to the go program. Currently, this is primarily implemented through Adobe
Extendscript, but hopefully in the future it will be upgraded to a C++ plugin.

## Installation
```bash
$ go get -u github.com/sbrow/ps/runner
```
