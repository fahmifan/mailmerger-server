package server

import (
	"errors"
	"net/http"

	"github.com/fahmifan/mailmerger-server/service"
	"github.com/fahmifan/ulids"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type CampaignHandler struct {
	*Server
}

func systemError(ec echo.Context, err error) error {
	log.Err(err).Msg("")
	return ec.Render(http.StatusOK, "pages/system_error.html", echo.Map{})
}

func notFound(ec echo.Context) error {
	return ec.Render(http.StatusNotFound, "pages/not_found_error.html", echo.Map{})
}

func (m CampaignHandler) List(ec echo.Context) (err error) {
	campaigns, err := m.service.CampaignService.List(ec.Request().Context())
	if err != nil {
		log.Err(err).Msg("")
		return systemError(ec, err)
	}

	return ec.Render(http.StatusOK, "pages/campaigns/index.html", echo.Map{
		"campaigns": campaigns,
	})
}

// render new form
func (m CampaignHandler) New(ec echo.Context) (err error) {
	return ec.Render(http.StatusOK, "pages/campaigns/new.html", echo.Map{})
}

func (m CampaignHandler) Create(ec echo.Context) (err error) {
	req := service.CreateCampaignRequest{}
	if err := ec.Bind(&req); err != nil {
		log.Err(err).Msg("create campaign - bind")
		return ec.Redirect(http.StatusSeeOther, m.echo.Reverse("campaigns-new"))
	}

	const mb = 1024 * 1024
	const maxMem = 2 * mb
	if err := ec.Request().ParseMultipartForm(maxMem); err != nil {
		return systemError(ec, err)
	}

	csvFile, _, err := ec.Request().FormFile("csv")
	if err != nil {
		return systemError(ec, err)
	}
	req.CSV = csvFile
	defer csvFile.Close()

	_, err = m.service.CampaignService.Create(ec.Request().Context(), req)
	if err != nil {
		return systemError(ec, err)
	}

	return ec.Redirect(http.StatusSeeOther, m.echo.Reverse("campaigns"))
}

func (m CampaignHandler) Show(ec echo.Context) (err error) {
	id, err := ulids.Parse(ec.Param("id"))
	if err != nil {
		return notFound(ec)
	}

	campaign, err := m.service.CampaignService.Find(ec.Request().Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		return notFound(ec)
	} else if err != nil {
		return systemError(ec, err)
	}

	return ec.Render(http.StatusOK, "pages/campaigns/show.html", echo.Map{"campaign": campaign})
}

func (m CampaignHandler) Edit(ec echo.Context) (err error) {
	id, err := ulids.Parse(ec.Param("id"))
	if err != nil {
		return notFound(ec)
	}

	campaign, err := m.service.CampaignService.Find(ec.Request().Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		return notFound(ec)
	} else if err != nil {
		return systemError(ec, err)
	}

	return ec.Render(http.StatusOK, "pages/campaigns/edit.html", echo.Map{"campaign": campaign})
}

func (m CampaignHandler) Update(ec echo.Context) (err error) {
	req := service.UpdateCampaignRequest{}
	if err := ec.Bind(&req); err != nil {
		log.Err(err).Msg("create campaign - bind")
		return ec.Redirect(http.StatusSeeOther, m.echo.Reverse("campaigns-new"))
	}

	const mb = 1024 * 1024
	const maxMem = 2 * mb
	if err := ec.Request().ParseMultipartForm(maxMem); err != nil {
		return systemError(ec, err)
	}

	csvFile, header, err := ec.Request().FormFile("csv")
	if err == nil && header.Size > 0 {
		req.CSV = csvFile
		defer csvFile.Close()
	}

	_, err = m.service.CampaignService.Update(ec.Request().Context(), req)
	if errors.Is(err, service.ErrNotFound) {
		return notFound(ec)
	} else if err != nil {
		return systemError(ec, err)
	}

	return ec.Redirect(http.StatusSeeOther, m.echo.Reverse("campaigns-show", req.ID))
}
