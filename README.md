Public API for People (public-people-api)
=========================================

Generates JSON representation of a Person in public-friendly format.

People are being migrated to be served from the new [Public Concepts API](https://github.com/Financial-Times/public-concepts-api) and as such this API will eventually be deprecated. From July 2018 requests to this service will be redirected via the concepts api then transformed to match the existing contract and returned.

Build & deployment
------------------

* Built by Docker Hub: [coco/public-people-api](https://hub.docker.com/r/coco/public-people-api/)
* CI provided by CircleCI: [public-people-api](https://circleci.com/gh/Financial-Times/public-people-api)

[![Coverage Status](https://coveralls.io/repos/github/Financial-Times/public-people-api/badge.svg?branch=master)](https://coveralls.io/github/Financial-Times/public-people-api?branch=master)[![CircleCI](https://circleci.com/gh/Financial-Times/public-people-api.svg?style=svg)](https://circleci.com/gh/Financial-Times/public-people-api)

Installation & running locally
------------------------------

1. Run the tests and install the binary:

        go test -race ./...
        go install

2. Run the binary (using the `help` flag to see the available optional arguments):

        $GOPATH/bin/public-people-api [--help]  

Options:

      --app-system-code         System Code of the application (env $APP_SYSTEM_CODE) (default "public-people-api")
      --app-name                Application name (env $APP_NAME) (default "Public People API")
      --log-level               App log level (env $LOG_LEVEL) (default "info")
      --port                    Port to listen on (env $PORT) (default 8080)
      --cache-duration          Duration Get requests should be cached for. e.g. 2h45m would set the max-age value to '7440' seconds (default:30s)
      --requestLoggingEnabled   Whether to log requests (env $REQUEST_LOGGING_ENABLED) (default true)
      --publicConceptsApiURL    Public concepts API endpoint URL. ($CONCEPTS_API) (default: "http://localhost:8080")

            
Test locally
------------------------------
```
go test -v -race ./...
```

Endpoints
---------

* Based on the following [google doc](https://docs.google.com/document/d/1SC4Uskl-VD78y0lg5H2Gq56VCmM4OFHofZM-OvpsOFo/edit#heading=h.qjo76xuvpj83).
* See the [api](_ft/api.yml) for the swagger definitions of the endpoints below.  


