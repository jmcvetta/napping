# Napping: HTTP for Gophers

Package `napping` is a [Go][] client library for interacting with
[RESTful APIs][].  Napping was inspired  by Python's excellent [Requests][]
library.


## Status

[![Drone Build Status](https://drone.io/github.com/jmcvetta/napping/status.png)](https://drone.io/github.com/jmcvetta/napping/latest)
[![Travis Build Status](https://travis-ci.org/jmcvetta/napping.png)](https://travis-ci.org/jmcvetta/napping)
[![Coverage Status](https://coveralls.io/repos/jmcvetta/restclient/badge.png)](https://coveralls.io/r/jmcvetta/napping)
[![xrefs](https://sourcegraph.com/api/repos/github.com/jmcvetta/napping/badges/xrefs.png)](https://sourcegraph.com/github.com/jmcvetta/napping)
[![funcs](https://sourcegraph.com/api/repos/github.com/jmcvetta/napping/badges/funcs.png)](https://sourcegraph.com/github.com/jmcvetta/napping)
[![top func](https://sourcegraph.com/api/repos/github.com/jmcvetta/napping/badges/top-func.png)](https://sourcegraph.com/github.com/jmcvetta/napping)
[![library users](https://sourcegraph.com/api/repos/github.com/jmcvetta/napping/badges/library-users.png)](https://sourcegraph.com/github.com/jmcvetta/napping)
[![status](https://sourcegraph.com/api/repos/github.com/jmcvetta/napping/badges/status.png)](https://sourcegraph.com/github.com/jmcvetta/napping)

All API changes will be made via Pull Request, so it's recommended you Watch
the repo Issues if using `napping` in production.

Used by, and developed in conjunction with, [Neoism][].


## Documentation

See GoDoc for [automatically generated API documentation][godoc].

Check out [github_auth_token][auth-token] for a working example
showing how to retrieve an auth token from the Github API.

Check out [docker_unixsocket][unix-socket] to see how to make REST
requests to Unix domain socket endpoints.


# Contributing

Contributions, in the form of Pull Requests or Issues, are gladly accepted.
Before submitting a Pull Request, please ensure your code passes all tests, and
that your changes do not decrease test coverage.  I.e. if you add new features,
also add corresponding new tests.


## License

This is Free Software, released under the terms of the [GPL v3][].


[Go]:           http://golang.org
[RESTful APIs]: http://en.wikipedia.org/wiki/Representational_state_transfer#RESTful_web_APIs
[Requests]:     http://python-requests.org
[GPL v3]:       http://www.gnu.org/copyleft/gpl.html
[auth-token]:   https://github.com/jmcvetta/napping/blob/master/examples/github_auth_token/github_auth_token.go
[unix-socket]:  https://github.com/jmcvetta/napping/blob/master/examples/docker_unixsocket/docker_unixsocket.go
[godoc]:        http://godoc.org/github.com/jmcvetta/napping
[Neoism]:       https://github.com/jmcvetta/neoism
