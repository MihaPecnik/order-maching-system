package config

import "flag"

type Config struct {
	Migrate     bool   `key:"migrate"`
	PostgresURI string `key:"postgres_uri"`
	Populate    bool   `key:"populate"`
}

func GetConfig() *Config {
	migrate := flag.Bool("migrate", false, "database migration")
	postgresUrl := flag.String("postgres_url", "", "postgres url")
	populate := flag.Bool("populate", false, "populate database")
	flag.Parse()

	return &Config{
		Migrate:     *migrate,
		PostgresURI: *postgresUrl,
		Populate:    *populate,
	}
}
