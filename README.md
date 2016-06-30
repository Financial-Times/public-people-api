# Public API for People (public-people-api)
__Provides a public API for People stored in a Neo4J graph database__

## Installation & running locally
* `go get -u github.com/Financial-Times/public-people-api`
* `cd $GOPATH/src/github.com/Financial-Times/public-people-api`
* `go test ./...`
* `go install`
* `$GOPATH/bin/public-people-api --neo-url={neo4jUrl} --port={port} --log-level={DEBUG|INFO|WARN|ERROR}--cache-duration{e.g. 22h10m3s}`
_Optional arguments are:
--neo-url defaults to http://localhost:7474/db/data, which is the out of box url for a local neo4j instance.
--port defaults to 8080.
--cache-duration defaults to 1 hour._
* `curl http://localhost:8080/people/143ba45c-2fb3-35bc-b227-a6ed80b5c517 | json_pp`
Or using [httpie](https://github.com/jkbrzt/httpie)
* `http GET http://localhost:8080/people/143ba45c-2fb3-35bc-b227-a6ed80b5c517`

## Endpoints
### GET

* `curl http://localhost:8080/people/143ba45c-2fb3-35bc-b227-a6ed80b5c517 | json_pp`
Or using [httpie](https://github.com/jkbrzt/httpie)
* `http GET http://localhost:8080/people/143ba45c-2fb3-35bc-b227-a6ed80b5c517`

The expected response will contain information about the person, and the organisations they are connected to (via memberships).

## Healthchecks
Healthchecks: [http://localhost:8080/__health](http://localhost:8080/__health)

### Logging
the application uses logrus, the logfile is initilaised in main.go.
 logging requires an env app parameter, for all enviromets  other than local logs are written to file
 when running locally logging is written to console (if you want to log locally to file you need to pass in an env parameter that is != local)
 NOTE: build-info and gtg end points are not logged as they are called every second from varnish/vulcand and this information is not needed in  logs/splunk
