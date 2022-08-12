package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"path"

	"github.com/fahmifan/mailmerger-server/localfs"
	"github.com/fahmifan/ulids"
	"github.com/oklog/ulid/v2"
	"go.etcd.io/bbolt"
)

var ErrNotFound = errors.New("not found error")

type Campaign struct {
	ID       ulids.ULID
	Name     string
	CSV      CSV
	Template Template
}

type CSV struct {
	ID       ulids.ULID
	Path     string
	FileName string
}

type Template struct {
	Body    string
	Subject string
}

type Config struct {
	db           *bbolt.DB
	localStorage *localfs.Storage
}

type Service struct {
	CampaignService *CampaignService
}

func NewService(db *bbolt.DB, localStorage *localfs.Storage) *Service {
	cfg := Config{
		db:           db,
		localStorage: localStorage,
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
	}
}

type CampaignService struct {
	cfg *Config
}

const CampaignBucket = "campaigns"

type CreateCampaignRequest struct {
	Name            string    `form:"name"`
	BodyTemplate    string    `form:"body_template"`
	SubjectTemplate string    `form:"subject_template"`
	CSV             io.Reader `form:"-"`
}

func (c *CampaignService) Create(ctx context.Context, req CreateCampaignRequest) (campaign Campaign, err error) {
	campaign = Campaign{
		ID:   ulids.New(),
		Name: req.Name,
		Template: Template{
			Body:    req.BodyTemplate,
			Subject: req.SubjectTemplate,
		},
	}

	if req.CSV != nil {
		campaign.CSV, err = c.createFile(req.CSV)
		if err != nil {
			return Campaign{}, err
		}
	}

	err = c.cfg.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(CampaignBucket))
		err = bucket.Put([]byte(campaign.ID.ULID.String()), MarshalJson(campaign))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return Campaign{}, err
	}
	return
}

func (c *CampaignService) createFile(csvFile io.Reader) (_ CSV, err error) {
	id := ulids.New()

	const folder = "csvs"
	fileName := id.String() + ".csv"
	filePath := path.Join(folder, fileName)

	err = c.cfg.localStorage.Save(filePath, csvFile)
	if err != nil {
		return CSV{}, err
	}

	return CSV{
		ID:       id,
		Path:     filePath,
		FileName: fileName,
	}, nil
}

func (c *CampaignService) List(ctx context.Context) (campaigns []Campaign, err error) {
	err = c.cfg.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(CampaignBucket))
		return bucket.ForEach(func(k, v []byte) error {
			campaign := Campaign{}
			err := json.Unmarshal(v, &campaign)
			if err != nil {
				return err
			}
			campaigns = append(campaigns, campaign)
			return nil
		})
	})
	return campaigns, err
}

func (c *CampaignService) Show(ctx context.Context, id ulid.ULID) (campaign Campaign, err error) {
	err = c.cfg.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(CampaignBucket))
		v := bucket.Get([]byte(id.String()))
		err := json.Unmarshal(v, &campaign)
		if err != nil {
			return err
		}
		if v == nil {
			return ErrNotFound
		}
		return json.Unmarshal(v, &campaign)
	})
	if err != nil {
		return Campaign{}, err
	}
	return campaign, nil
}

func MarshalJson(i interface{}) []byte {
	bt, _ := json.Marshal(i)
	return bt
}
