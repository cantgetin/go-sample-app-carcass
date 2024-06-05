package config

import "time"

type Config struct {
	AppName  string `env:"APP_NAME" envDefault:"test"`
	HTTPAddr string `env:"HTTP_ADDR" envDefault:":8080"`

	DbHost     string `env:"DB_HOST" envDefault:"db"`
	DbPort     string `env:"DB_PORT" envDefault:"5432"`
	DbUsername string `env:"DB_USERNAME" envDefault:"db"`
	DbPassword string `env:"DB_PASSWORD" envDefault:"db"`
	DbDatabase string `env:"DB_DATABASE" envDefault:"db"`

	PgIdleConn     int           `env:"PG_IDLE_CONN" envDefault:"10"`
	PgMaxOpenConn  int           `env:"PG_MAX_OPEN_CONN" envDefault:"100"`
	PgPingInterval time.Duration `env:"PG_PING_INTERVAL" envDefault:"10s"`
}
