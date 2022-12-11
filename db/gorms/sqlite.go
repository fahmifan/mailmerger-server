package gorms

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type SQLite struct {
	DSN              string
	EnableWAL        bool
	EnableForeignKey bool
}

func (s SQLite) MustOpen() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(s.DSN+"?doNotInterpretDatetime=1"), &gorm.Config{})
	panicErr(err)

	if s.EnableWAL {
		err := db.Exec(`PRAGMA journal_mode=WAL`).Error
		panicErr(err)
	}
	if s.EnableForeignKey {
		err := db.Exec(`PRAGMA foreign_keys=ON`).Error
		panicErr(err)
	}

	return db
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
