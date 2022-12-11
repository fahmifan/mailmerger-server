package gorms

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func Close(db *gorm.DB) {
	if db == nil {
		return
	}
	dbConn, err := db.DB()
	if err != nil {
		return
	}
	if err = dbConn.Close(); err != nil {
		log.Err(err).Msg("close db")
	}
}
