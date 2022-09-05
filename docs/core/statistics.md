# Statistics
This page describes how statistics are being tracked inside the application.

### Prerequisites
- InfluxDB

### Configuration
The following environment variables are required in order to connect to the influx instance, these variables need to be set in the config.yaml file.

| Variable           | Description  |
|--------------------|--------------|
| INFLUXDB_URL       | Instance url |
| INFLUXDB_ORG       | Organisation |
| INFLUXDB_BUCKET    | Bucket name  |
| INFLUXDB_TOKEN     | Influx token |

### Grafana
At the moment of writing, statistics collected by the uber controller are being deployed on a Grafana instance. You can reach this instance by browsing: https://stats.dev.odyssey.ninja/grafana