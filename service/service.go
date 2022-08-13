package service

import (
	"github.com/fahmifan/mailmerger-server/pkg/localfs"
	"go.etcd.io/bbolt"
)

type Config struct {
	db            *bbolt.DB
	localStorage  *localfs.Storage
	blastEmailCfg *BlastEmailConfig
}

type Service struct {
	CampaignService *CampaignService
	FileService     *FileService
}

func NewService(db *bbolt.DB, localStorage *localfs.Storage, blastEmailCfg *BlastEmailConfig) *Service {
	cfg := Config{
		db:            db,
		localStorage:  localStorage,
		blastEmailCfg: blastEmailCfg,
	}
	db.Update(func(tx *bbolt.Tx) (err error) {
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
