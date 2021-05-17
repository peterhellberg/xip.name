xip.name
========

> The service was shut down on 2021-05-17 after Google marked the entire domain as “Social engineering content” since
> (working as intended) you could link to for example http://1.1.1.1.xip.name/ and getting redirected to <http://1.1.1.1/>

[![Build Status](https://travis-ci.org/peterhellberg/xip.name.svg?branch=master)](https://travis-ci.org/peterhellberg/xip.name)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/peterhellberg/xip.name)
[![License BSD](https://img.shields.io/badge/license-BSD-lightgrey.svg?style=flat)](https://github.com/peterhellberg/xip.name/blob/master/LICENSE)

A simple wildcard DNS inspired by xip.io (which seems to have been shut down now)

```bash
        10.0.0.1.xip.name  resolves to  10.0.0.1
    www.10.0.0.2.xip.name  resolves to  10.0.0.2
    foo.10.0.0.3.xip.name  resolves to  10.0.0.3
bar.baz.10.0.0.4.xip.name  resolves to  10.0.0.4
```

## How does it work?

xip.name runs a custom Domain Name Server which extracts any IP address found
in the requested domain name and sends it back in the response.

## Credits

xip.name is built on top of [Miek](http://miek.nl)’s lovely [dns package](https://github.com/miekg/dns) for Go.
