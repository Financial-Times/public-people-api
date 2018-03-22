Public API for People (public-people-api)
=========================================

Generates JSON representation of a Person in public-friendly format.

Build & deployment
------------------

* Built by Docker Hub: [coco/public-people-api](https://hub.docker.com/r/coco/public-people-api/)
* CI provided by CircleCI: [public-people-api](https://circleci.com/gh/Financial-Times/public-people-api)

[![Coverage Status](https://coveralls.io/repos/github/Financial-Times/public-people-api/badge.svg?branch=circle-coveralls)](https://coveralls.io/github/Financial-Times/public-people-api?branch=circle-coveralls)[![CircleCI](https://circleci.com/gh/Financial-Times/public-people-api.svg?style=svg)](https://circleci.com/gh/Financial-Times/public-people-api)

Installation & running locally
------------------------------

1. Run the tests and install the binary:

        govendor sync
        govendor test -v -race
        go install

2. Run the binary (using the `help` flag to see the available optional arguments):

        $GOPATH/bin/public-people-api [--help]  

Options:

      --app-system-code         System Code of the application (env $APP_SYSTEM_CODE) (default "public-people-api")
      --app-name                Application name (env $APP_NAME) (default "Public People API")
      --port                    Port to listen on (env $PORT) (default 8080)
      --neoURL                  Connection string for NEO4J (env $NEO4J_CONNECTION) (default "bolt://localhost:7474")
      --requestLoggingEnabled   Whether to log requests (env $REQUEST_LOGGING_ENABLED) (default true)
      --logLevel                App log level (env $LOG_LEVEL) (default "info")
      --graphiteTCPAddress      Graphite TCP address (default: "")
      --graphitePrefix          Prefix to use. Should start with content, include the environment, and the host name. e.g. content.test.public.people.api. (default: "")
      --logMetrics              Whether to log metrics. Set to true if running locally and you want metrics output (default: false)
      --env                     environment this app is running in (default: local) - this is for setting apiUrl
      --cache-duration          Duration Get requests should be cached for. e.g. 2h45m would set the max-age value to '7440' seconds (default:30s)
      --requestLoggingEnabled   Whether to log requests (default: true)

            
Test locally
------------------------------

Tests in neo4j package rely on a running instance of Neo4j installed locally.  

```
docker run \
    --rm \
    --publish=7474:7474 \
    --publish=7687:7687 \
    --env=NEO4J_ACCEPT_LICENSE_AGREEMENT=yes \
    --env=NEO4J_AUTH=none \
    neo4j:3.3.3-enterprise

govendor test -v -race -cover +local
```

Endpoints
---------

* Based on the following [google doc](https://docs.google.com/document/d/1SC4Uskl-VD78y0lg5H2Gq56VCmM4OFHofZM-OvpsOFo/edit#heading=h.qjo76xuvpj83).
* See the [api](_ft/api.yml) for the swagger definitions of the endpoints below.  


