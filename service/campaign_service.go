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
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found error")

type Audit struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Campaign struct {
	ID     ulids.ULID
	FileID *ulids.ULID
	Name   string
	Audit

	File     File
	Template Template
	Events   []Event
}

func (c Campaign) IsNoEvent() bool {
	return len(c.Events) == 0
}

type File struct {
	ID       ulids.ULID
	Folder   string
	FileName string
	Audit
}

type Template struct {
	ID         ulids.ULID
	CampaignID ulids.ULID
	Body       string
	Subject    string
	Audit
}

type EventStatus string

const (
	EventStatusSuccess EventStatus = "success"
	EventStatusFailed  EventStatus = "failed"
)

type Event struct {
	ID         ulids.ULID
	CampaignID ulids.ULID `gorm:"references:CampaignID"`
	Detail     string
	CreatedAt  time.Time
	Status     EventStatus
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
	tx := c.cfg.db.WithContext(ctx)

	newCampaign := Campaign{
		ID:   ulids.New(),
		Name: req.Name,
	}

	if req.CSV != nil {
		file, err := c.createFile(ctx, req.CSV)
		if err != nil {
			return Campaign{}, err
		}

		if err = tx.Create(&file).Error; err != nil {
			return Campaign{}, err
		}
		newCampaign.FileID = &file.ID
	}

	template := Template{
		ID:      ulids.New(),
		Body:    req.BodyTemplate,
		Subject: req.SubjectTemplate,
	}
	newCampaign.Template = template

	if err = tx.Create(&newCampaign).Error; err != nil {
		return
	}

	template.CampaignID = newCampaign.ID
	if err = tx.Create(&template).Error; err != nil {
		return
	}

	campaign = newCampaign
	return
}

const csvFolder = "csvs"

func (c *CampaignService) createFile(ctx context.Context, csvFile io.Reader) (_ File, err error) {
	id := ulids.New()

	fileName := id.String() + ".csv"
	filePath := path.Join(csvFolder, fileName)

	err = c.cfg.localStorage.Save(ctx, filePath, csvFile)
	if err != nil {
		return File{}, err
	}

	return File{
		ID:       id,
		Folder:   filePath,
		FileName: fileName,
	}, nil
}

func (c *CampaignService) List(ctx context.Context) (campaigns []Campaign, err error) {
	if err = c.cfg.db.Model(&Campaign{}).
		Preload("Events").
		Preload("Template").
		Preload("File").
		Find(&campaigns).
		Error; err != nil {
		return
	}

	return
}

func (c *CampaignService) Find(ctx context.Context, id ulids.ULID) (campaign Campaign, err error) {
	if err = c.cfg.db.
		Preload("File").
		Preload("Template").
		Preload("Events").
		Take(&campaign, "id = ?", id).
		Error; err != nil {
		return
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
	campaign, err := c.Find(ctx, req.ID)
	if err != nil {
		return
	}

	campaign.Name = req.Name
	campaign.Template = Template{
		Body:    req.BodyTemplate,
		Subject: req.SubjectTemplate,
	}

	if req.CSV != nil {
		campaign.File, err = c.createFile(ctx, req.CSV)
		if err != nil {
			return Campaign{}, err
		}
	}

	err = c.cfg.db.Updates(&campaign).Error
	if err != nil {
		return Campaign{}, err
	}
	return campaign, nil
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

	csvFile, err := c.cfg.localStorage.Seek(ctx, campaign.File.Folder)
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
	err = c.cfg.boltDB.Update(func(tx *bbolt.Tx) error {
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
