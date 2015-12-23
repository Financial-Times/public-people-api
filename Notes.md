
## Access to metadata

Compare json structure in papi.json to the curl request. In the curl request you can get the labels from the metadata. Might be good is we could avoid additional calls.

Wonder if this can be done through some setting like http://neo4j.com/docs/stable/rest-api-transactional.html#rest-api-include-query-statistics

`curl -H 'Content-type:application/json' http://localhost:7474/db/data/node/1185210`
```
{
  "extensions" : { },
  "metadata" : {
    "id" : 1185210,
    "labels" : [ "Thing", "Concept", "Organisation", "PublicCompany" ]
  },
  "paged_traverse" : "http://localhost:7474/db/data/node/1185210/paged/traverse/{returnType}{?pageSize,leaseTime}",
  "outgoing_relationships" : "http://localhost:7474/db/data/node/1185210/relationships/out",
  "outgoing_typed_relationships" : "http://localhost:7474/db/data/node/1185210/relationships/out/{-list|&|types}",
  "create_relationship" : "http://localhost:7474/db/data/node/1185210/relationships",
  "labels" : "http://localhost:7474/db/data/node/1185210/labels",
  "traverse" : "http://localhost:7474/db/data/node/1185210/traverse/{returnType}",
  "all_relationships" : "http://localhost:7474/db/data/node/1185210/relationships/all",
  "all_typed_relationships" : "http://localhost:7474/db/data/node/1185210/relationships/all/{-list|&|types}",
  "property" : "http://localhost:7474/db/data/node/1185210/properties/{key}",
  "self" : "http://localhost:7474/db/data/node/1185210",
  "incoming_relationships" : "http://localhost:7474/db/data/node/1185210/relationships/in",
  "properties" : "http://localhost:7474/db/data/node/1185210/properties",
  "incoming_typed_relationships" : "http://localhost:7474/db/data/node/1185210/relationships/in/{-list|&|types}",
  "data" : {
    "legalName" : "American Electric Power Company, Inc.",
    "hiddenLabel" : "AMERICAN ELECTRIC POWER CO INC",
    "factsetIdentifier" : "000BY1-E",
    "properName" : "American Electric Power Co., Inc.",
    "prefLabel" : "American Electric Power Co., Inc.",
    "shortName" : "American Electric Power",
    "leiIdentifier" : "1B4S6S7G0TW5EE83BO58",
    "uuid" : "af6d5434-8fd9-39e7-a95d-97c7743f4a77"
  }
```
