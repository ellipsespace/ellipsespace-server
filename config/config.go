package config

import "os"

type Config struct {
	AppIP              string
	AppPort            string
	JwtSecrets         string
	DBConnectionString string
}

func New() Config {
	return Config{
		AppIP:              getValue("APP_IP", "localhost"),
		AppPort:            getValue("APP_PORT", "8888"),
		JwtSecrets:         getValue("JWT_SECRETS", "CHANGE_ME"),
		DBConnectionString: getValue("CONN_STRING", "CHANGE_ME"),
	}
}

func getValue(key, def string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return def
}
