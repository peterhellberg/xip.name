[xip.name](http://xip.name/)
========

[![Build Status](https://travis-ci.org/peterhellberg/xip.name.svg?branch=master)](https://travis-ci.org/peterhellberg/xip.name)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/peterhellberg/xip.name)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/peterhellberg/xip.name#license-mit)

A simple wildcard DNS inspired by [xip.io](http://xip.io/)

```bash
        10.0.0.1.xip.name  resolves to  10.0.0.1
    www.10.0.0.2.xip.name  resolves to  10.0.0.2
    foo.10.0.0.3.xip.name  resolves to  10.0.0.3
bar.baz.10.0.0.4.xip.name  resolves to  10.0.0.4
```

### How does it work?

xip.name runs a custom DNS server which extracts any IP address found
in the requested domain name and sends it back in the response.


### Does it cost anything to use xip.name?

No, but you are welcome to donate if you find the service useful.

Bitcoin: **[16f4vuZcpybd7rfprjB6Ki87BNVbtMA1M5](https://blockchain.info/address/16f4vuZcpybd7rfprjB6Ki87BNVbtMA1M5)**

## Credits

xip.name is built on top of [Miek](http://miek.nl)’s lovely [dns package](https://github.com/miekg/dns) for Go.

## License (MIT)

Copyright (c) 2014 [Peter Hellberg](http://c7.se/)

> Permission is hereby granted, free of charge, to any person obtaining
> a copy of this software and associated documentation files (the
> "Software"), to deal in the Software without restriction, including
> without limitation the rights to use, copy, modify, merge, publish,
> distribute, sublicense, and/or sell copies of the Software, and to
> permit persons to whom the Software is furnished to do so, subject to
> the following conditions:

> The above copyright notice and this permission notice shall be
> included in all copies or substantial portions of the Software.

> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
> MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
> NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
> LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
> OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
> WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
