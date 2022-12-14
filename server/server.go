package server

import (
	"github.com/fahmifan/mailmerger-server/service"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

type Server struct {
	echo       *echo.Echo
	service    *service.Service
	csrfSecret string
}

func NewServer(service *service.Service) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) Run() {
	s.routes()
	s.echo.Start(":8080")
}

func (s *Server) routes() {
	s.echo = echo.New()
	s.echo.Renderer = NewPongoRenderer(&PongoRendererConfig{
		RootDir:      "server/templates",
		DebugEnabled: true,
	})
	s.echo.Debug = true
	setSecureCSRF := csrf.Secure(true)
	csrfOpts := []csrf.Option{csrf.SameSite(csrf.SameSiteLaxMode), setSecureCSRF}
	csrfMiddleware := csrf.Protect([]byte(s.csrfSecret), csrfOpts...)

	home := HomeHandler{}
	campaign := CampaignHandler{s}
	event := EventHandler{s}
	files := FileHandler{s}
	template := TemplateHandler{s}
	renderOnDemand := RenderOnDemandHandler{s}

	s.echo.Use(echo.WrapMiddleware(csrfMiddleware))

	s.echo.GET("/", home.Index).Name = "homes"

	s.echo.GET("/campaigns", campaign.List).Name = "campaigns"
	s.echo.GET("/campaigns/new", campaign.New).Name = "campaigns-new"
	s.echo.POST("/campaigns", campaign.Create).Name = "campaigns-create"
	s.echo.GET("/campaigns/:id", campaign.Show).Name = "campaigns-show"
	s.echo.GET("/campaigns/:id/edit", campaign.Edit).Name = "campaigns-edit"
	s.echo.POST("/campaigns/:id/update", campaign.Update).Name = "campaigns-update"
	s.echo.POST("/campaigns/:id/delete", campaign.Delete).Name = "campaigns-delete"

	s.echo.GET("/campaigns/:campaign_id/render-on-demand", renderOnDemand.Show).Name = "render-ondemand"

	s.echo.POST("/events", event.Create).Name = "events-create"

	s.echo.GET("/files/:file_name", files.Show).Name = "files-show"

	s.echo.GET("/templates", template.List).Name = "templates"
	s.echo.GET("/templates/:id", template.Show).Name = "templates-show"
	s.echo.GET("/templates/new", template.New).Name = "templates-new"
	s.echo.POST("/templates", template.Create).Name = "templates-create"
	s.echo.GET("/templates/:id/edit", template.Edit).Name = "templates-edit"
	s.echo.POST("/templates/:id/update", template.Update).Name = "templates-update"

}
