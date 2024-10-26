package config

type VaultConfig struct {
	Addr  string `env:"ADDRESS,required"`
	Token string `env:"TOKEN,required"`
}
