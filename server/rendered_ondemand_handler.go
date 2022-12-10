package server

import (
	"errors"
	"net/http"

	"github.com/fahmifan/mailmerger-server/service"
	"github.com/fahmifan/ulids"
	"github.com/labstack/echo/v4"
)

type RenderOnDemandHandler struct {
	*Server
}

func (handler RenderOnDemandHandler) Show(ec echo.Context) error {
	campaignID, err := ulids.Parse(ec.Param("campaign_id"))
	if err != nil {
		return badRequest(ec, "invalid campaign_id")
	}

	templateIdStr := ec.QueryParam("templateID")
	body := ec.QueryParam("body")

	var (
		rendered   []byte
		templateID ulids.ULID
	)
	if templateIdStr != "" {
		templateID, err = ulids.Parse(templateIdStr)
		if err != nil {
			return badRequest(ec, "invalid query param templateID")
		}
		rendered, err = handler.service.CampaignService.RenderByIDAndBodyAndTemplate(ec.Request().Context(), campaignID, templateID, body)
	} else {
		rendered, err = handler.service.CampaignService.RenderByIDAndBody(ec.Request().Context(), campaignID, body)
	}
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			return notFound(ec)
		default:
			return systemError(ec, err)
		}
	}

	return ec.HTML(http.StatusOK, string(rendered))
}
