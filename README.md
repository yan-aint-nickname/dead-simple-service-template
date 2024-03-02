# dead-simple-service-template

WIP

## DB

### Migrations

#### Install goose

`go install github.com/pressly/goose/v3/cmd/goose@latest`

#### Check migrations status
Plese consider to be in migrations dir or specify it with flag
`goose --dir=migrations postgres "$POSTGRES_DSN" status`
