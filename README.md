Public API for People (public-people-api)
=========================================

Installation & running locally
------------------------------

1. Download the source code, dependencies and its test dependencies:

        go get -u github.com/Financial-Times/public-people-api
        cd $GOPATH/src/github.com/Financial-Times/public-people-api
        go get -t

1. Run the tests and install the binary:

        go test ./...
        go install

1. Run the binary (using the `help` flag to see the available optional arguments):

        $GOPATH/bin/public-people-api [--help]

1. Test:

    1. Either using curl:

            curl http://localhost:8080/people/143ba45c-2fb3-35bc-b227-a6ed80b5c517 | json_pp

    1. Or using [httpie](https://github.com/jkbrzt/httpie):

            http GET http://localhost:8080/people/143ba45c-2fb3-35bc-b227-a6ed80b5c517


Endpoints
---------

### GET

Using curl:

    curl http://localhost:8080/people/143ba45c-2fb3-35bc-b227-a6ed80b5c517 | json_pp`

Or using [httpie](https://github.com/jkbrzt/httpie):

    http GET http://localhost:8080/people/143ba45c-2fb3-35bc-b227-a6ed80b5c517

The expected response will contain information about the person, and the organisations they are connected to (via memberships).

Based on the following [google doc](https://docs.google.com/document/d/1SC4Uskl-VD78y0lg5H2Gq56VCmM4OFHofZM-OvpsOFo/edit#heading=h.qjo76xuvpj83).


Health Checks
-------------

Health checks: [http://localhost:8080/__health](http://localhost:8080/__health)


### Logging

* The application uses [logrus](https://github.com/Sirupsen/logrus); the log file is initialised in [app.go](app.go).
* Logging requires an `env` app parameter, for all environments other than `local` logs are written to file.
* When running locally, logs are written to console. If you want to log locally to file, you need to pass in an env parameter that is != `local`.
* NOTE: `/build-info` and `/__gtg` endpoints are not logged as they are called every second from varnish/vulcand and this information is not needed in logs/splunk.
