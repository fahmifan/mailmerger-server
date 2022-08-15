package service

import (
	"github.com/fahmifan/mailmerger-server/pkg/localfs"
	"gorm.io/gorm"
)

type Config struct {
	db            *gorm.DB
	localStorage  *localfs.Storage
	blastEmailCfg *BlastEmailConfig
}

type Service struct {
	CampaignService *CampaignService
	FileService     *FileService
}

func NewService(db *gorm.DB, localStorage *localfs.Storage, blastEmailCfg *BlastEmailConfig) *Service {
	cfg := Config{
		db:            db,
		localStorage:  localStorage,
		blastEmailCfg: blastEmailCfg,
	}
	return &Service{
		CampaignService: &CampaignService{&cfg},
		FileService:     &FileService{&cfg},
	}
}
