# spamd-client

Golang SpamAssassin Client Library and Commandline tool

[![MPLv2 License](https://img.shields.io/badge/license-MPLv2-blue.svg?style=flat-square)](https://www.mozilla.org/MPL/2.0/)

## Description

spamd-client is a Golang library and cmdline tool that implements the
SPAMD client protocol used by SpamAssassin.

## Requirements

* Golang 1.20.x or higher
* Pflag - github.com/spf13/pflag

## Getting started

### spamd-client client

The spamd-client client can be installed as follows

```console
$ go get github.com/jniltinho/spamd-client/cmd/spamd-client
```

Or by cloning the repo and then running

```console
$ make build
$ ./bin/spamd-client
```

### spamd-client library

You can import the library in your code

```golang
import "github.com/jniltinho/spamd-client/pkg"
```

### Testing

``make test``

## License

MPL-2.0
