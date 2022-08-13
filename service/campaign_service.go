package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"path"
	"strings"
	"time"

	"github.com/fahmifan/mailmerger"
	"github.com/fahmifan/ulids"
	"github.com/rs/zerolog/log"
	"go.etcd.io/bbolt"
)

var ErrNotFound = errors.New("not found error")

type Campaign struct {
	ID       ulids.ULID
	Name     string
	CSV      CSV
	Template Template
	Events   []Event
}

func (c Campaign) IsNoEvent() bool {
	return len(c.Events) == 0
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

type EventStatus string

const (
	EventStatusSuccess EventStatus = "success"
	EventStatusFailed  EventStatus = "failed"
)

type Event struct {
	CreatedAt time.Time
	Status    EventStatus
	Detail    string
}

type BlastEmailConfig struct {
	Sender      string
	Concurrency uint
	Transporter mailmerger.Transporter
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
		campaign.CSV, err = c.createFile(ctx, req.CSV)
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

const csvFolder = "csvs"

func (c *CampaignService) createFile(ctx context.Context, csvFile io.Reader) (_ CSV, err error) {
	id := ulids.New()

	fileName := id.String() + ".csv"
	filePath := path.Join(csvFolder, fileName)

	err = c.cfg.localStorage.Save(ctx, filePath, csvFile)
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

func (c *CampaignService) Find(ctx context.Context, id ulids.ULID) (campaign Campaign, err error) {
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

type UpdateCampaignRequest struct {
	ID              ulids.ULID `form:"id"`
	Name            string     `form:"name"`
	BodyTemplate    string     `form:"body_template"`
	SubjectTemplate string     `form:"subject_template"`
	CSV             io.Reader  `form:"-"`
}

func (c *CampaignService) Update(ctx context.Context, req UpdateCampaignRequest) (_ Campaign, err error) {
	oldCampaign, err := c.Find(ctx, req.ID)
	if err != nil {
		return
	}

	oldCampaign.Name = req.Name
	oldCampaign.Template = Template{
		Body:    req.BodyTemplate,
		Subject: req.SubjectTemplate,
	}

	if req.CSV != nil {
		oldCampaign.CSV, err = c.createFile(ctx, req.CSV)
		if err != nil {
			return Campaign{}, err
		}
	}

	err = c.cfg.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(CampaignBucket))

		key := []byte(oldCampaign.ID.ULID.String())
		if old := bucket.Get(key); old == nil {
			return ErrNotFound
		}

		err = bucket.Put(key, MarshalJson(oldCampaign))
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

type CreateBlastEmailEventRequest struct {
	CampaignID ulids.ULID `form:"campaign_id"`
}

// CreateBlastEmailEvent create an event
func (c *CampaignService) CreateBlastEmailEvent(ctx context.Context, req CreateBlastEmailEventRequest) (event Event, err error) {
	campaign, err := c.Find(ctx, req.CampaignID)
	if err != nil {
		return
	}

	csvFile, err := c.cfg.localStorage.Seek(ctx, campaign.CSV.Path)
	if err != nil {
		return
	}
	defer csvFile.Close()

	mailer := mailmerger.NewMailer(&mailmerger.MailerConfig{
		SenderEmail:     c.cfg.blastEmailCfg.Sender,
		CsvSrc:          csvFile,
		BodyTemplate:    strings.NewReader(campaign.Template.Body),
		SubjectTemplate: strings.NewReader(campaign.Template.Subject),
		Concurrency:     2,
		Transporter:     c.cfg.blastEmailCfg.Transporter,
	})
	if err = mailer.Parse(); err != nil {
		return
	}

	event = Event{
		CreatedAt: time.Now(),
		Status:    EventStatusSuccess,
	}
	if err = mailer.SendAll(ctx); err != nil {
		log.Err(err).Msg("sendAll")
		event.Status = EventStatusFailed
		event.Detail = err.Error()
	}

	campaign.Events = append(campaign.Events, event)
	err = c.cfg.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(CampaignBucket))
		return bucket.Put([]byte(campaign.ID.String()), MarshalJson(campaign))
	})
	if err != nil {
		return Event{}, err
	}

	return event, nil
}

func MarshalJson(i interface{}) []byte {
	bt, _ := json.Marshal(i)
	return bt
}
