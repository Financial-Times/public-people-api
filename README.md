# Public API for People (public-people-api-neo4j)
__Provides a public API for People stored in a Neo4J graph database__

## Build & deployment etc:
* General view http://ftjen10085-lvpr-uk-p:8181/view/public-people-api
* Build and publish to forge http://ftjen10085-lvpr-uk-p:8181/job/public-people-api-build
* Deploy to test or production http://ftjen10085-lvpr-uk-p:8181/job/public-people-api-deploy

## Installation & running locally
* `go get -u github.com/Financial-Times/public-people-api`
* `cd $GOPATH/src/github.com/Financial-Times/public-people-api`
* `go install`
* `$GOPATH/bin/public-people-api --neo-url={neo4jUrl} --port={port}`
_Both arguments are optional, they default to a local Neo4j install and port 8080._
* `curl http://localhost:8080/people/143ba45c-2fb3-35bc-b227-a6ed80b5c517 | json_pp`

## API definition
Based on the following [google doc](https://docs.google.com/document/d/1SC4Uskl-VD78y0lg5H2Gq56VCmM4OFHofZM-OvpsOFo/edit#heading=h.qjo76xuvpj83)

## Todo
* Test cases
* Metrics
* Health checks
* Logging levels (ERROR, WARN, INFO & DEBUG)

### Healthchecks (NOT IMPLEMENTED)

Healthchecks: [http://localhost:8080/__health](http://localhost:8080/__health)

Good-to-go: [http://localhost:8080/__gtg](http://localhost:8080/__gtg)
