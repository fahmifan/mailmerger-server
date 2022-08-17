package service

import (
	"context"

	"github.com/fahmifan/ulids"
)

type Template struct {
	ID   ulids.ULID `gorm:"primary_key"`
	Name string
	HTML string
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
