# Getting Started
Before you are able to start interacting with the controller, you first need to make sure you have everything you need.


### Prerequisites

- GO 1.19
- PostgreSQL
- MQTT
- InfluxDB

This service has been built with GO version 1.19 in mind, other versions might work but are not officially supported.

## Building the service
You can use an IDE to build this service (such as GoLand).
Or you can use the following CLI command:

```bash
make build
```

## Running the service
You can use an IDE to run this service (such as GoLand).
Or you can use the following CLI command:

```bash
make run
```

### Configuration
This application uses a `config.yaml` file to set environment variables for its dependencies. If this file does not exist in the root directory, you should create one yourself.

An example of the config.yaml structure can be seen below:
```yaml
mqtt:
  host: 
  port: 
  user: 
  password: 
postgres:
  database:
  host: 
  port:
  username:
  password:
  max_conns:
influx:
  url: 
  org: 
  bucket:
  token: 
keycloak:
  server: 
  secret: 
  client: 
  realm: 
settings:
  bind_address:
  bind_port:
  loglevel:
common:
  introspect_url:
```

Environment variables are being loaded from inside the config directory. A structure of the config directory can be found below:
<br/>

```
├── config
│   ├── common.go
│   ├── config.go
│   ├── influx.go
│   ├── local.go
│   ├── mqtt.go
│   ├── postgres.go
│   └── uiclient.go
```

Existing environment variables can be found on the tables below:

#### Common variables

Yaml Module: `common`

| Variable        | Description                           |
|-----------------|---------------------------------------|
| introspect_url  | Authentication introspection Endpoint |

#### UI variables

Yaml Module: `ui_client`

| Variable                          | Description                |
|-----------------------------------|----------------------------|
| frontend_url                      | UI Client Url              |
| keycloak_open_id_connect_url      | Keycloak OIDC Url          |
| keycloak_open_id_client_id        | Keycloak OIDC Client ID    |
| keycloak_open_id_scope            | Keycloak OIDC Scope        |
| hydra_open_id_connect_url         | Hydra OIDC Url             |
| hydra_open_id_client_id           | Hydra OIDC Client ID       |
| hydra_open_id_guest_client_id     | Hydra OIDC Guest Client ID |
| hydra_open_id_scope               | Hydra OIDC Scope           |
| web_3_identity_provider_url       | Web3 IDP Url               |
| sentry_dsn                        | Sentry DSN                 |
| agora_app_id                      | Agora App ID               |
| auth_service_url                  | Authentication Service Url |
| google_api_client_id              | Google API Client ID       |
| google_api_developer_key          | Google API Developer Key   |
| react_app_youtube_key             | React Youtube Key          |
| unity_client_streaming_assets_url | Unity asset Url            |
| unity_client_company_name         | Unity company name         |
| unity_client_product_name         | Unity product name         |
| unity_client_product_version      | Unity product version      |

#### MQTT Broker

Yaml Module: `mqtt`

| Variable  | Description                |
|-----------|----------------------------|
| host      | Host address of the broker |
| port      | Port number                |
| username  | Broker username            |
| password  | Broker password            |

#### PostgreSQL

Yaml Module: `postgres`

| Variable    | Description                   |
|-------------|-------------------------------|
| database    | Name of the database schema   |
| host        | Instance hostname             |
| port        | Instance port number          |
| username    | Username                      |
| password    | Password (if set)             |
| max_conns   | Maximum amount of connections |

#### InfluxDB

Yaml Module: `influx`

| Variable | Description  |
|----------|--------------|
| url      | Instance url |
| org      | Organisation |
| bucket   | Bucket name  |
| token    | Influx token |

### Core
<mark>TODO Summarize -> redirect to core/general.md for more information</mark>

### Database
<mark>TODO Summarize -> redirect to database/general.md for more information</mark>

### Plugin infrastructure
<mark>TODO Summarize -> redirect to plugins/general.md for more information</mark>