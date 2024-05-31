package configserver

type HTTPServerConfig struct {
	Host string `yaml:"host,omitempty"`
	Port string `yaml:"port,omitempty"`
}
