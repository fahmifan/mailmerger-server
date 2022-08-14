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
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found error")

type Audit struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Campaign struct {
	ID     ulids.ULID `gorm:"primary_key"`
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
	ID       ulids.ULID `gorm:"primary_key"`
	Folder   string
	FileName string
	Audit
}

type Template struct {
	ID         ulids.ULID `gorm:"primary_key"`
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

	if err = tx.Create(&newCampaign).Omit("Events", "Template", "File").Error; err != nil {
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

	tx := c.cfg.db.WithContext(ctx)

	if req.CSV != nil {
		campaign.File, err = c.createFile(ctx, req.CSV)
		if err != nil {
			return Campaign{}, err
		}

		if err = tx.Create(&campaign.File).Error; err != nil {
			return Campaign{}, err
		}
	}

	if err = c.cfg.db.Updates(&campaign).Omit("Events", "Template", "File").Error; err != nil {
		return Campaign{}, err
	}

	tpl := Template{
		ID:         ulids.New(),
		CampaignID: campaign.ID,
		Body:       req.BodyTemplate,
		Subject:    req.SubjectTemplate,
	}
	if err = c.replaceTemplate(ctx, &campaign, &tpl); err != nil {
		return
	}

	return campaign, nil
}

func (c *CampaignService) replaceTemplate(ctx context.Context, campaign *Campaign, tpl *Template) (err error) {
	return c.cfg.db.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Delete(&campaign.Template).Error; err != nil {
			return
		}

		if err = tx.Create(tpl).Error; err != nil {
			return
		}

		campaign.Template = *tpl
		return
	})
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
		ID:         ulids.New(),
		CampaignID: campaign.ID,
		Status:     EventStatusSuccess,
	}
	if err = mailer.SendAll(ctx); err != nil {
		log.Err(err).Msg("sendAll")
		event.Status = EventStatusFailed
		event.Detail = err.Error()
	}

	err = c.cfg.db.Create(&event).Error
	if err != nil {
		return Event{}, err
	}

	return event, nil
}

func MarshalJson(i interface{}) []byte {
	bt, _ := json.Marshal(i)
	return bt
}
