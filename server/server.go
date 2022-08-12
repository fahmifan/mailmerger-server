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
	s.echo.Renderer = &PongoRenderer{
		PongoRendererConfig: &PongoRendererConfig{
			RootDir:      "server/templates",
			DebugEnabled: true,
		},
	}
	setSecureCSRF := csrf.Secure(true)
	csrfOpts := []csrf.Option{csrf.SameSite(csrf.SameSiteLaxMode), setSecureCSRF}
	csrfMiddleware := csrf.Protect([]byte(s.csrfSecret), csrfOpts...)

	campaign := CampaignHandler{s}

	s.echo.Use(echo.WrapMiddleware(csrfMiddleware))
	s.echo.GET("/campaigns", campaign.List).Name = "campaigns"
	s.echo.GET("/campaigns/new", campaign.New).Name = "campaigns-new"
	s.echo.POST("/campaigns", campaign.Create).Name = "campaigns-create"
	s.echo.GET("/campaigns/:id", campaign.Show).Name = "campaigns-show"
	s.echo.GET("/campaigns/:id/edit", campaign.Edit).Name = "campaigns-edit"
	s.echo.POST("/campaigns/:id/update", campaign.Update).Name = "campaigns-update"
}
