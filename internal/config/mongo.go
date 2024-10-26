package config

type MongoConfig struct {
	Host           string `env:"HOST" envDefault:"localhost"`
	Port           string `env:"PORT" envDefault:"27017"`
	User           string `env:"USER" envDefault:"root"`
	Pass           string `env:"PASS" envDefault:"root"`
	DbName         string `env:"DB_NAME" envDefault:"test"`
	CollectionName string `env:"COLLECTION_NAME" envDefault:"test"`
}
