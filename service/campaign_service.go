package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"github.com/fahmifan/mailmerger"
	"github.com/fahmifan/ulids"
	"github.com/flosch/pongo2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
)

var ErrNotFound = errors.New("not found error")

type Audit struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Campaign struct {
	ID         ulids.ULID `gorm:"primary_key"`
	FileID     *ulids.ULID
	Name       string
	Body       string
	Subject    string
	TemplateID *ulids.ULID
	Audit

	File     File
	Events   []Event
	Template *Template `gorm:"->;foreignKey:TemplateID"`
}

func (c Campaign) HasNoTemplate() bool {
	return c.TemplateID == nil || c.Template == nil
}

func (c *Campaign) BeforeCreate(tx *gorm.DB) error {
	omitFields := []string{"Events"}
	gormOmit(tx, omitFields...)
	return nil
}

func (c *Campaign) BeforeUpdate(tx *gorm.DB) error {
	omitFields := []string{"Events"}
	gormOmit(tx, omitFields...)
	return nil
}

func gormOmit(tx *gorm.DB, columns ...string) {
	if len(columns) == 1 && strings.ContainsRune(columns[0], ',') {
		tx.Statement.Omits = strings.FieldsFunc(columns[0], utils.IsValidDBNameChar)
	} else {
		tx.Statement.Omits = columns
	}
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
	Name            string      `form:"name"`
	BodyTemplate    string      `form:"body"`
	SubjectTemplate string      `form:"subject"`
	CSV             io.Reader   `form:"-"`
	TemplateID      *ulids.ULID `form:"-"`
}

func (c *CampaignService) Create(ctx context.Context, req CreateCampaignRequest) (campaign Campaign, err error) {
	tx := c.cfg.db.WithContext(ctx)

	campaign = Campaign{
		ID:         ulids.New(),
		Name:       req.Name,
		TemplateID: req.TemplateID,
	}

	if req.CSV != nil {
		file, err := c.createFileIfAny(ctx, req.CSV)
		if err != nil {
			return Campaign{}, err
		}
		campaign.FileID = &file.ID
	}
	campaign.Body = req.BodyTemplate
	campaign.Subject = req.SubjectTemplate

	if err = tx.Create(&campaign).Error; err != nil {
		return Campaign{}, err
	}

	return campaign, nil
}

const csvFolder = "csvs"

func (c *CampaignService) createFileIfAny(ctx context.Context, csvFile io.Reader) (_ File, err error) {
	id := ulids.New()

	fileName := id.String() + ".csv"
	filePath := path.Join(csvFolder, fileName)

	err = c.cfg.localStorage.Save(ctx, filePath, csvFile)
	if err != nil {
		return File{}, err
	}

	file := File{
		ID:       id,
		Folder:   filePath,
		FileName: fileName,
	}

	if err = c.cfg.db.Create(&file).Error; err != nil {
		return File{}, err
	}

	return file, nil
}

func (c *CampaignService) List(ctx context.Context) (campaigns []Campaign, err error) {
	if err = c.cfg.db.Model(&Campaign{}).
		Preload("Events").
		Preload("File").
		Order("created_at desc").
		Find(&campaigns).
		Error; err != nil {
		return
	}

	return
}

func (c *CampaignService) Find(ctx context.Context, id ulids.ULID) (campaign Campaign, err error) {
	if err = c.cfg.db.
		Preload("Template").
		Preload("File").
		Preload("Events").
		Take(&campaign, "id = ?", id).
		Error; err != nil {
		return
	}
	return campaign, nil
}

type UpdateCampaignRequest struct {
	ID         ulids.ULID  `form:"id"`
	Name       string      `form:"name"`
	Body       string      `form:"body"`
	Subject    string      `form:"subject"`
	CSV        io.Reader   `form:"-"`
	TemplateID *ulids.ULID `form:"-"`
}

func (c *CampaignService) Update(ctx context.Context, req UpdateCampaignRequest) (_ Campaign, err error) {
	campaign, err := c.Find(ctx, req.ID)
	if err != nil {
		return
	}

	campaign.Name = req.Name
	campaign.Body = req.Body
	campaign.Subject = req.Subject
	campaign.TemplateID = req.TemplateID

	if req.CSV != nil {
		campaign.File, err = c.createFileIfAny(ctx, req.CSV)
		if err != nil {
			return Campaign{}, err
		}
		campaign.FileID = &campaign.File.ID
	}

	payload := map[string]interface{}{
		"name":        campaign.Name,
		"body":        campaign.Body,
		"subject":     campaign.Subject,
		"template_id": campaign.TemplateID,
		"updated_at":  "now()",
	}
	if err = c.cfg.db.Model(&campaign).Updates(payload).Error; err != nil {
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

	body := bytes.NewBuffer(nil)
	if campaign.TemplateID != nil {
		tpl, err := c.findTemplate(ctx, *campaign.TemplateID)
		if err != nil {
			return Event{}, err
		}
		pongoTpl, err := pongo2.FromString(tpl.HTML)
		if err != nil {
			return Event{}, err
		}

		err = pongoTpl.ExecuteWriter(pongo2.Context{"body": campaign.Body}, body)
		if err != nil {
			return Event{}, err
		}
	} else {
		body.WriteString(campaign.Body)
	}

	csvFile, err := c.cfg.localStorage.Seek(ctx, campaign.File.Folder)
	if err != nil {
		return
	}
	defer csvFile.Close()

	mailer := mailmerger.NewMailer(&mailmerger.MailerConfig{
		SenderEmail:     c.cfg.blastEmailCfg.Sender,
		CsvSrc:          csvFile,
		BodyTemplate:    body,
		SubjectTemplate: strings.NewReader(campaign.Subject),
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

func (c *CampaignService) Delete(ctx context.Context, id ulids.ULID) (campaign Campaign, err error) {
	if campaign, err = c.Find(ctx, id); err != nil {
		return Campaign{}, fmt.Errorf("find old campaign: %w", err)
	}

	if err = c.cfg.db.WithContext(ctx).Delete(&campaign).Error; err != nil {
		return Campaign{}, fmt.Errorf("delete: %w", err)
	}

	return campaign, nil
}

func (c *CampaignService) findTemplate(ctx context.Context, id ulids.ULID) (tpl Template, err error) {
	err = c.cfg.db.Take(&tpl, "id = ?", id).Error
	return tpl, unwrapErr(err)
}
