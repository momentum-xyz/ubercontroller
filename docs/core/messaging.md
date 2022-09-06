# Message broker
This page describes how the MQTT messaging protocol is being utilized in the application. 

### Configuration
The following environment variables are required in order to connect to the broker instance, these variables need to be set in the config.yaml file.

Yaml Module: `mqtt`

| Variable  | Description                |
|-----------|----------------------------|
| host      | Host address of the broker |
| port      | Port number                |
| username  | Broker username            |
| password  | Broker password            |