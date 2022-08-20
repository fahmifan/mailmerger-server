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
	templates, err := m.service.TemplateService.FindAll(ec.Request().Context())
	if err != nil {
		return systemError(ec, err)
	}

	payload := echo.Map{
		"templates": templates,
	}
	return ec.Render(http.StatusOK, "pages/campaigns/new.html", payload)
}

func (m CampaignHandler) Create(ec echo.Context) (err error) {
	req := service.CreateCampaignRequest{}
	if err := ec.Bind(&req); err != nil {
		log.Err(err).Msg("create campaign - bind")
		return ec.Redirect(http.StatusSeeOther, m.echo.Reverse("campaigns-new"))
	}

	parseEmptyNilID(&req.TemplateID, ec.FormValue("template_id"))

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

	campaign, err := m.service.CampaignService.Create(ec.Request().Context(), req)
	if err != nil {
		return systemError(ec, err)
	}

	return ec.Redirect(http.StatusSeeOther, ec.Echo().Reverse("campaigns-show", campaign.ID))
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

	templates, err := m.service.TemplateService.FindAll(ec.Request().Context())
	if err != nil {
		return systemError(ec, err)
	}

	payload := echo.Map{
		"campaign":  campaign,
		"templates": templates,
	}

	return ec.Render(http.StatusOK, "pages/campaigns/edit.html", payload)
}

func (m CampaignHandler) Update(ec echo.Context) (err error) {
	req := service.UpdateCampaignRequest{}
	if err := ec.Bind(&req); err != nil {
		log.Err(err).Msg("update campaign - bind")
		return ec.Redirect(http.StatusSeeOther, m.echo.Reverse("campaigns-new"))
	}
	parseEmptyNilID(&req.TemplateID, ec.FormValue("template_id"))

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

func (m CampaignHandler) Delete(ec echo.Context) (err error) {
	id, err := ulids.Parse(ec.Param("id"))
	if err != nil {
		return notFound(ec)
	}

	_, err = m.service.CampaignService.Delete(ec.Request().Context(), id)
	if err != nil {
		return systemError(ec, err)
	}

	return ec.Redirect(http.StatusSeeOther, m.echo.Reverse("campaigns"))
}

// parseEmptyNilID set id to nil if id is empty
// else parse the id
func parseEmptyNilID(uid **ulids.ULID, id string) error {
	if id == "" {
		uid = nil
		return nil
	}
	templateID, err := ulids.Parse(id)
	if err != nil {
		return err
	}

	if uid == nil {
		*uid = &templateID
		return nil
	}

	*uid = &templateID
	return nil
}
