package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Server       ServerConfig
	Secret       string      `env:"SECRET"`
	DbSecretPath string      `env:"DB_SECRET_PATH"`
	Vault        VaultConfig `envPrefix:"VAULT_"`
	Mongo        MongoConfig `envPrefix:"MONGO_"`
	Redis        RedisConfig `envPrefix:"REDIS"`
}

const envPrefix = "INVENTORY_SERVICE_"

func ReadEnv() (Config, error) {
	opts := env.Options{
		RequiredIfNoDef: true,
		Prefix:          envPrefix,
	}

	cfg := new(Config)
	err := env.ParseWithOptions(cfg, opts)
	if err != nil {
		return Config{}, err
	}

	return *cfg, nil
}
