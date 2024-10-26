package config

type RedisConfig struct {
	HOST     string `env:"HOST" envDefault:"localhost"`
	PORT     string `env:"PORT" envDefault:"6379"`
	PASSWORD string `env:"PASSWORD" envDefault:""`
	DB       int    `env:"DB" envDefault:"0"`
}
