package config

type ServerConfig struct {
	Host string `env:"HOST"`
	Port string `env:"PORT,required"`
}
