# Public API for People (public-people-api-neo4j)
__Provides a public API for People stored in a Neo4J graph database__

## Build & deployment etc:
_NB You will need to tag a commit in order to build, since the UI asks for a tag to build / deploy_
* [Jenkins view](http://ftjen10085-lvpr-uk-p:8181/view/public-people-api)
* [Build and publish to forge](http://ftjen10085-lvpr-uk-p:8181/job/public-people-api-build)
* [Deploy to test or production](http://ftjen10085-lvpr-uk-p:8181/job/public-people-api-deploy)


## Installation & running locally
* `go get -u github.com/Financial-Times/public-people-api`
* `cd $GOPATH/src/github.com/Financial-Times/public-people-api`
* `go test ./...`
* `go install`
* `$GOPATH/bin/public-people-api --neo-url={neo4jUrl} --port={port} --log-level={DEBUG|INFO|WARN|ERROR}`
_Both arguments are optional, they default to a local Neo4j install and port 8080._
* `curl http://localhost:8080/people/143ba45c-2fb3-35bc-b227-a6ed80b5c517 | json_pp`

## API definition
Based on the following [google doc](https://docs.google.com/document/d/1SC4Uskl-VD78y0lg5H2Gq56VCmM4OFHofZM-OvpsOFo/edit#heading=h.qjo76xuvpj83)

## Healthchecks
Healthchecks: [http://localhost:8080/__health](http://localhost:8080/__health)

## Todo
### For parity with existing API
* Add in TMELabels as part of labels (uniq)
* Use annotations for ordering memberships

### API specific
* Complete Test cases
* Runbook

### Cross cutting concerns
* Allow service to start if neo4j is unavailable at startup time
* Add Metrics
* Rework build / deploy (low priority)
  * Suggested flow:
    1. Build & Tests
    1. Publish Release (using konstructor to generate vrm)
    1. Deploy vrm/hash to test/prod
