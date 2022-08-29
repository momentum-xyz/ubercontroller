package config

// MQTT : structure to hold configuration
type MQTT struct {
	HOST     string `yaml:"host" envconfig:"MQTT_BROKER_HOST"`
	PORT     uint   `yaml:"port" envconfig:"MQTT_BROKER_PORT"`
	USER     string `yaml:"user" envconfig:"MQTT_BROKER_USER"`
	PASSWORD string `yaml:"password" envconfig:"MQTT_BROKER_PASSWORD"`
}

func (x *MQTT) Init() {
	x.HOST = "localhost"
	x.PORT = 1883
	x.USER = ""
	x.PASSWORD = ""
}
