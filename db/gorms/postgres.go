package gorms

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DSN string
}

func (p Postgres) MustOpen() *gorm.DB {
	db, err := gorm.Open(postgres.Open(p.DSN), &gorm.Config{})
	panicErr(err)
	return db
}
