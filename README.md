[xip.name](http://xip.name/)
========

[![Build Status](https://travis-ci.org/peterhellberg/xip.name.svg?branch=master)](https://travis-ci.org/peterhellberg/xip.name)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/peterhellberg/xip.name)
[![License BSD](https://img.shields.io/badge/license-BSD-lightgrey.svg?style=flat)](https://github.com/peterhellberg/xip.name/blob/master/LICENSE)

A simple wildcard DNS inspired by [xip.io](http://xip.io/)

```bash
        10.0.0.1.xip.name  resolves to  10.0.0.1
    www.10.0.0.2.xip.name  resolves to  10.0.0.2
    foo.10.0.0.3.xip.name  resolves to  10.0.0.3
bar.baz.10.0.0.4.xip.name  resolves to  10.0.0.4
```

## How does it work?

xip.name runs a custom Domain Name Server which extracts any IP address found
in the requested domain name and sends it back in the response.

## Does it cost anything to use xip.name?

No, but you are welcome to donate if you find the service useful.

Bitcoin: **[16f4vuZcpybd7rfprjB6Ki87BNVbtMA1M5](https://blockchain.info/address/16f4vuZcpybd7rfprjB6Ki87BNVbtMA1M5)**

You can also help out by signing up to DigitalOcean using my [referral link](https://www.digitalocean.com/?refcode=cd245791f86e).

## Credits

xip.name is built on top of [Miek](http://miek.nl)’s lovely [dns package](https://github.com/miekg/dns) for Go.
