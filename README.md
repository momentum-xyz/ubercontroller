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

