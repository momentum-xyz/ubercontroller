package config

type Influx struct {
	URL    string `yaml:"url" envconfig:"INFLUXDB_URL"`
	ORG    string `yaml:"org" envconfig:"INFLUXDB_ORG"`
	BUCKET string `yaml:"bucket" envconfig:"INFLUXDB_BUCKET"`
	TOKEN  string `yaml:"token" envconfig:"INFLUXDB_TOKEN"`
}

func (x *Influx) Init() {

}
