# dead-simple-service-template
WIP

## Components

### Http Server

https://github.com/gin-gonic/gin/tree/master

#### Logs

https://github.com/uber-go/zap

### Worker

#### Kafka (#Kafka)

#### Tasks

aka scheduler

https://github.com/fieldryand/goflow

### Kafka

#### Producer

#### Consumer

### Redis

https://github.com/redis/go-redis

### Http Client

### Postgres
- [x] Pool
- [ ] Single

#### Driver

https://github.com/jackc/pgx

#### Logs
WIP
https://github.com/jackc/pgx-zap/tree/master

#### Migrations

<!-- https://github.com/ariga/atlas <- very complex system -->

https://github.com/pressly/goose

##### Install goose

`go install github.com/pressly/goose/v3/cmd/goose@latest`

##### Check migrations status

Plese consider to be in migrations dir or specify it with flag or use envs
`goose --dir=migrations postgres "$POSTGRES_DSN" status`

### Mongo

#### Driver?

#### Migrations?

## App State

https://github.com/uber-go/fx

## Styleguid

https://github.com/uber-go/guide/blob/master/style.md

### Better async

https://github.com/sourcegraph/conc

## Logs

- https://github.com/uber-go/zap
- https://github.com/jackc/pgx-zap/tree/master
- https://github.com/getsentry/sentry-go/blob/master/gin/README.md
You can throw panics, but should you?

## Tests


## Linting
https://github.com/uber-go/guide/blob/master/style.md#linting

tl;dr
errcheck to ensure that errors are handled
goimports to format code and manage imports
golint to point out common style mistakes
govet to analyze code for common mistakes
staticcheck to do various static analysis checks
