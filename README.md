# Napping: HTTP for Gophers

Package `napping` is a [Go][] client library for interacting with
[RESTful APIs][].  Napping was inspired  by Python's excellent [Requests][]
library.


## Status

[![Drone Build Status](https://drone.io/github.com/jmcvetta/napping/status.png)](https://drone.io/github.com/jmcvetta/napping/latest)
[![Travis Build Status](https://travis-ci.org/jmcvetta/napping.png)](https://travis-ci.org/jmcvetta/napping)
[![Coverage Status](https://coveralls.io/repos/jmcvetta/napping/badge.png?branch=napping)](https://coveralls.io/r/jmcvetta/napping)

API is fairly stable, but there may be additions and small changes from time to
time.  All API changes will be made via Pull Request, so it's recommended you
Watch the repo Issues if using `napping` in production.

Used by, and developed in conjunction with, package [neoism][].


## Documentation

See GoDoc for [automatically generated API documentation][godoc].

Check out [examples/github_auth_token.go][auth-token] for a working example
showing how to retrieve an auth token from the Github API.


## License

This is Free Software, released under the terms of the [GPL v3][].


[Go]:           http://golang.org
[RESTful APIs]: http://en.wikipedia.org/wiki/Representational_state_transfer#RESTful_web_APIs
[Requests]:     http://python-requests.org
[GPL v3]:       http://www.gnu.org/copyleft/gpl.html
[auth-token]:   https://github.com/jmcvetta/napping/blob/master/examples/github_auth_token.go
[godoc]:        http://godoc.org/github.com/jmcvetta/napping
[neoism]:       https://github.com/jmcvetta/neoism
