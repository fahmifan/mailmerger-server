package config

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

var cfg Config

var once sync.Once

type DbDriver string

const (
	DbDriverSQLite   = "sqlite"
	DbDriverPostgres = "postgres"

	DefaultDBName = "mailmerger.db"
)

type Config struct {
	Postgres Postgres
	SQLite   SQLite
	DbDriver DbDriver `env:"DB_DRIVER"`
	Mailer   Mailer
}

func init() {
	once.Do(func() {
		loadEnv()
		err := env.Parse(&cfg)
		panicErr(err)
	})
}

type Mailer struct {
	SenderAddress string `env:"MAILER_SENDER_ADDRESS"`
	Concurrency   uint   `env:"MAILER_CONCURRENCY"`
}

type SQLite struct {
	DBName           string `env:"SQLITE_DB_NAME"`
	EnableWAL        bool   `env:"SQLITE_ENABLE_WAL"`
	EnableForeignKey bool   `env:"SQLITE_ENABLE_FOREIGN_KEY"`
}

func (s SQLite) GetDBName() string {
	if s.DBName == "" {
		return DefaultDBName
	}
	return s.DBName
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DBName   string `env:"POSTGRES_DB_NAME"`
	SSLMode  string `env:"POSTGRES_SSL_MODE"`
}

func (p Postgres) DSN() string {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		p.Host, p.User, p.Password, p.DBName, p.Port, p.SSLMode,
	)
	return dsn
}

func GetConfig() Config {
	return cfg
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Info().Msgf("no env file, use os env instead: %w", err)
		return
	}
	log.Info().Msg("using .env file")
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
