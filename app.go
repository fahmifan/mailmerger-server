package main

import (
	"context"

	"github.com/fahmifan/mailmerger-server/service"
	"github.com/fahmifan/ulids"
)

// App struct
type App struct {
	ctx context.Context
	*service.Service
}

// NewApp creates a new App application struct
func NewApp(svc *service.Service) *App {
	return &App{Service: svc}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ShowCampaign(idStr string) (service.Campaign, error) {
	id, err := ulids.Parse(idStr)
	if err != nil {
		return service.Campaign{}, err
	}

	campaign, err := a.CampaignService.Find(a.ctx, id)
	if err != nil {
		return service.Campaign{}, err
	}

	return campaign, nil
}

func (app *App) ListCampaigns() ([]service.Campaign, error) {
	return app.CampaignService.List(app.ctx)
}

func (app *App) CreateRenderedTemplate(templateIdStr, body string) (string, error) {
	templateID, err := ulids.Parse(templateIdStr)
	if err != nil {
		return "", err
	}

	result, err := app.CampaignService.RenderByBodyAndTemplate(app.ctx, templateID, body)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (app *App) ListTemplates() ([]service.Template, error) {
	return app.TemplateService.FindAll(app.ctx)
}

func (app *App) ShowTemplate(idStr string) (service.Template, error) {
	id, err := ulids.Parse(idStr)
	if err != nil {
		return service.Template{}, err
	}
	return app.TemplateService.FindByID(app.ctx, id)
}
