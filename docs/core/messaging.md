# Message broker
This page describes how the MQTT messaging protocol is being utilized in the application. 

### Configuration
The following environment variables are required in order to connect to the broker instance, these variables need to be set in the config.yaml file.

| Variable     | Description                   |
|--------------|-------------------------------|
| DB_DATABASE  | Name of the database schema   |
| PGDB_HOST    | Instance hostname             |
| DB_PORT      | Instance port number          |
| DB_USERNAME  | Username                      |
| DB_PASSWORD  | Password (if set)             |
| DB_MAX_CONNS | Maximum amount of connections |