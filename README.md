# uberubercontroller
New ubercontroller


## Development

### Prerequisites

- [Go environment](https://go.dev/doc/install)
- Connection information for a pre-seeded Postgres database

### TL;DR (quick start)

For local development:

 - copy config.example.yaml to config.yaml
 - configure the database connection in config.yaml
 - make run

### Local development with remote Media Manager

You need to have kubectl configured with dev env and do port-forwarding to dev instance:

```
kubectl port-forward deploy/media-manager-deployment 4002:4000
```

In `config.yaml` add the following to **common** section:

```
render_internal_url: 'http://localhost:4002'
```

## Swagger Documentation
1. Install [swag](https://github.com/swaggo/swag) cli tool
2. Run `swag init -g universe/node/api.go` to generate documentation
3. Open in browser http://localhost:4000/swagger/index.html

## Database migrations

Changes to the database are managed as SQL scripts in `database/migrations`
with [go-migrate](https://github.com/golang-migrate/migrate). See its documentation for detailed instructions. The short version:

```
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate create -ext sql -dir database/migrations/sql/ -seq <SHORT_NAME_FOR_MIGRATION>
```
Two empty .sql files will be created in `database/migrations/sql`.
Implemented them and to test it:
```
go run ./cmd/service/ -migrate
```

And to revert the change:
```
go run ./cmd/service/ -migrate -steps -1
```

If something goes wrong, it will leave the database in a 'dirty' state.
To resolve you have to manually bring back the database to a known state (so use a single transaction in the SQL scripts to make this easier).
After that edit the `schema_migration` table and set the version number back and flip the dirty boolean.
