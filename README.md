# Ubercontroller
New controller for Odyssey Momentum platform.


## Development

### Prerequisites

- [Go environment](https://go.dev/doc/install)
- Connection information for a running Postgres database

### TL;DR (quick start)

For local development:

 - copy config.example.yaml to config.yaml
 - configure the database connection in config.yaml
 - make run

### Configuration

Configuration can be done through environment variables or a YAML file. A (minimal) configuration file example is included at `config.example.yaml`. Copy this file to `config.yaml` and it will be used by default. The minimal required configuration are the PostgreSQL database connection fields. All others have default values setup for a development environment.

The configuration is implemented with Go struct tags.
For the available configuration options see config/config.go


## API documentation

Documentation for the JSON API is automatically generated.
For the develop branch this is done for every change and [deployed](https://momentum-xyz.github.io/ubercontroller/api).

To generate it locally NodeJS is required (`npx` should be available).

Run `make docs-html`.
Then open `build/docs/api.html` in a web browser.


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
