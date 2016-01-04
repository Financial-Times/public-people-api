# Public API for People (public-people-api-neo4j)

__Provides a public API for People stored in a Neo4J graph database__


## Build & deployment etc:

* Build and publish to forge http://ftjen10085-lvpr-uk-p:8181/job/public-people-api-neo4j-build
* Deploy to test *tbd*
* Deploy to production *tbd*

## Installation & running locally

* `go get -u github.com/Financial-Times/public-people-api-neo4j`

* `$GOPATH/bin/public-people-api-neo4j --neo-url={neo4jUrl} --port={port}`
_Both arguments are optional, they default to a local Neo4j install and port 8080._

## API definition
Based on the following [google doc](https://docs.google.com/document/d/1SC4Uskl-VD78y0lg5H2Gq56VCmM4OFHofZM-OvpsOFo/edit#heading=h.qjo76xuvpj83)

## Todo
* Test cases
* Logging levels (ERROR, WARN, INFO & DEBUG)
* General approach to empty fields
  * Null versus "" or other empty value for simple fields
  * Empty list / array / object for compound fields
* Deployment and build

### Try it!

`curl -XPUT -H "X-Request-Id: 123" -H "Content-Type: application/json" localhost:8080/people/3fa70485-3a57-3b9b-9449-774b001cd965 --data '{"uuid":"3fa70485-3a57-3b9b-9449-774b001cd965", "name":"Robert W. Addington", "identifiers":[{ "authority":"http://api.ft.com/system/FACTSET-PPL", "identifierValue":"000BJG-E"}]}'`

`curl -H "X-Request-Id: 123" localhost:8080/people/3fa70485-3a57-3b9b-9449-774b001cd965`

### Healthchecks

Healthchecks: [http://localhost:8080/__health](http://localhost:8080/__health)

Good-to-go: [http://localhost:8080/__gtg](http://localhost:8080/__gtg)
