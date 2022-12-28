package config

import "os"

type Config struct {
	Port               string
	JwtSecrets         string
	DBConnectionString string
	AppMode            string
}

func New() Config {
	return Config{
		Port:               getValue("PORT", "8080"),
		JwtSecrets:         getValue("JWT_SECRETS", "CHANGE_ME"),
		DBConnectionString: getValue("CONN_STRING", "CHANGE_ME"),
		AppMode:            getValue("APP_MODE", "develop"),
	}
}

func getValue(key, def string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return def
}
