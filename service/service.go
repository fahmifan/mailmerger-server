package service

import (
	"github.com/fahmifan/mailmerger-server/pkg/localfs"
	"go.etcd.io/bbolt"
	"gorm.io/gorm"
)

type Config struct {
	db            *gorm.DB
	boltDB        *bbolt.DB
	localStorage  *localfs.Storage
	blastEmailCfg *BlastEmailConfig
}

type Service struct {
	CampaignService *CampaignService
	FileService     *FileService
}

func NewService(db *gorm.DB, boltDB *bbolt.DB, localStorage *localfs.Storage, blastEmailCfg *BlastEmailConfig) *Service {
	cfg := Config{
		db:            db,
		boltDB:        boltDB,
		localStorage:  localStorage,
		blastEmailCfg: blastEmailCfg,
	}
	boltDB.Update(func(tx *bbolt.Tx) (err error) {
		_, err = tx.CreateBucketIfNotExists([]byte(CampaignBucket))
		if err != nil {
			return
		}
		return err
	})
	return &Service{
		CampaignService: &CampaignService{&cfg},
		FileService:     &FileService{&cfg},
	}
}
