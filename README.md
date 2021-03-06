![logo](logo.png)
# ps
[![GoDoc](https://godoc.org/github.com/sbrow/ps?status.svg)](https://godoc.org/github.com/sbrow/ps) [![Build Status](https://travis-ci.org/sbrow/ps.svg?branch=master)](https://travis-ci.org/sbrow/ps) [![Coverage Status](https://coveralls.io/repos/github/sbrow/ps/badge.svg?branch=master)](https://coveralls.io/github/sbrow/ps?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/sbrow/ps)](https://goreportcard.com/report/github.com/sbrow/ps)

`import "github.com/sbrow/ps"`

* [Overview](#pkg-overview)
* [Installation](pkg-installation)
* [Subdirectories](#pkg-subdirectories)
* [TODO](#pkg-note-TODO)
* [Documentation](#pkg-doc)

## <a name="pkg-overview">Overview</a>
Package ps is a rudimentary API between Adobe Photoshop CS5 and Golang.
The interaction between the two is implemented using Javascript/VBScript.

Use it to control Photoshop, edit documents, and perform batch operations.

Currently only supports Photoshop CS5 Windows x86_64.





## <a name="pkg-installation">Installation</a>
```sh
$ go get -u github.com/sbrow/ps
```
<!---

#### <a name="pkg-examples">Examples</a>
* [JSLayer](example_JSLayer_test.go)

--->



## <a name="pkg-note-TODO">TODO</a>

`sbrow:` (2) Make TextLayer a subclass of ArtLayer.

`sbrow:` Reduce cyclomatic complexity of ActiveDocument().

`sbrow:` refactor Close to Document.Close

## <a name="pkg-doc">Documentation</a>
For full Documentation please visit https://godoc.org/github.com/sbrow/ps
- - -


Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
