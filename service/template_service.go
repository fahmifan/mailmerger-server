package service

import (
	"context"
	"errors"

	"github.com/fahmifan/ulids"
	"gorm.io/gorm"
)

type Template struct {
	ID   ulids.ULID `json:"id" gorm:"primary_key"`
	Name string     `json:"name"`
	HTML string     `json:"html"`
	Audit
}

type TemplateService struct {
	cfg *Config
}

type CreateTemplateRequest struct {
	Name string `form:"name"`
	HTML string `form:"html"`
}

func (t *TemplateService) FindAll(ctx context.Context) (tmplts []Template, err error) {
	err = t.cfg.db.WithContext(ctx).Order("created_at desc").Find(&tmplts).Error
	return
}

func unwrapErr(err error) error {
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return ErrNotFound
}

func (t *TemplateService) FindByID(ctx context.Context, id ulids.ULID) (tpl Template, err error) {
	err = t.cfg.db.Take(&tpl, "id = ?", id).Error
	if err != nil {
		return Template{}, unwrapErr(err)
	}
	return tpl, nil
}

func (t *TemplateService) Create(ctx context.Context, req CreateTemplateRequest) (tpl Template, err error) {
	tpl = Template{
		ID:   ulids.New(),
		HTML: req.HTML,
		Name: req.Name,
	}
	if err = t.cfg.db.WithContext(ctx).Create(&tpl).Error; err != nil {
		return Template{}, err
	}

	return
}

type UpdateTemplateRequest struct {
	ID   ulids.ULID `form:"id"`
	HTML string     `form:"html"`
	Name string     `form:"name"`
}

func (t *TemplateService) Update(ctx context.Context, req UpdateTemplateRequest) (tpl Template, err error) {
	err = t.cfg.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Take(&tpl, "id = ?", req.ID).Error; err != nil {
			return
		}
		tpl.HTML = req.HTML
		tpl.Name = req.Name
		if err = tx.Updates(&tpl).Error; err != nil {
			return
		}

		return
	})
	if err != nil {
		return Template{}, unwrapErr(err)
	}

	return
}
