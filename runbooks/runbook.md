# UPP - Public People Api

The Public People API is a micro-service which aims to provide Person and
related data given a Person identifier.

## Code

public-people-api

## Primary URL

<https://github.com/Financial-Times/public-people-api>

## Service Tier

Bronze

## Lifecycle Stage

Production

## Delivered By

content

## Supported By

content

## Known About By

- elitsa.pavlova
- kalin.arsov
- ivan.nikolov
- miroslav.gatsanoga

## Host Platform

AWS

## Architecture

The Public People API is a micro-service which aims to provide Person and
related data given a Person identifier.

## Contains Personal Data

No

## Contains Sensitive Data

No

## Dependencies

- public-concepts-api

## Failover Architecture Type

ActiveActive

## Failover Process Type

FullyAutomated

## Failback Process Type

PartiallyAutomated

## Failover Details

The service is deployed in the delivery clusters as a deployment.

## Data Recovery Process Type

NotApplicable

## Data Recovery Details

The service does not store data, so it does not require any data recovery steps.

## Release Process Type

PartiallyAutomated

## Rollback Process Type

Manual

## Release Details

The release is triggered by making a Github release which is then picked up by a Jenkins multibranch pipeline. The Jenkins pipeline should be manually started in order for it to deploy the helm package to the Kubernetes clusters.

## Key Management Process Type

Manual

## Key Management Details

To access the service clients need to provide basic auth credentials.
To rotate credentials you need to login to a particular cluster and update varnish-auth secrets.

## Monitoring

Pod health:

- <https://upp-prod-delivery-eu.upp.ft.com/__health/__pods-health?service-name=public-people-api>
- <https://upp-prod-delivery-us.upp.ft.com/__health/__pods-health?service-name=public-people-api>

## First Line Troubleshooting

<https://github.com/Financial-Times/upp-docs/tree/master/guides/ops/first-line-troubleshooting>

## Second Line Troubleshooting

Please refer to the GitHub repository README for troubleshooting information.
