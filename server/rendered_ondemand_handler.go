package server

import (
	"errors"
	"net/http"

	"github.com/fahmifan/mailmerger-server/service"
	"github.com/fahmifan/ulids"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type RenderOnDemandHandler struct {
	*Server
}

func (handler RenderOnDemandHandler) Show(ec echo.Context) (err error) {
	templateIdStr := ec.QueryParam("templateID")
	body := ec.QueryParam("body")

	var (
		rendered   []byte
		templateID ulids.ULID
	)
	if templateIdStr == "" {
		return ec.HTML(http.StatusOK, body)
	}

	templateID, err = ulids.Parse(templateIdStr)
	if err != nil {
		return ec.HTML(http.StatusBadRequest, "")
	}

	rendered, err = handler.service.CampaignService.RenderByBodyAndTemplate(ec.Request().Context(), templateID, body)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			return ec.HTML(http.StatusNotFound, "")
		default:
			log.Err(err).Msg("RenderOnDemandHandler-RenderByBodyAndTemplate")
			return ec.HTML(http.StatusInternalServerError, "")
		}
	}

	return ec.HTML(http.StatusOK, string(rendered))
}
