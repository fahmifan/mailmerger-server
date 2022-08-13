package server

import (
	"net/http"

	"github.com/fahmifan/mailmerger-server/service"
	"github.com/labstack/echo/v4"
)

type EventHandler struct {
	*Server
}

func (e EventHandler) Create(ec echo.Context) (err error) {
	req := service.CreateBlastEmailEventRequest{}
	if err = ec.Bind(&req); err != nil {
		return systemError(ec, err)
	}

	_, err = e.service.CampaignService.CreateBlastEmailEvent(ec.Request().Context(), req)
	if err != nil {
		return systemError(ec, err)
	}

	return ec.Redirect(http.StatusSeeOther, e.echo.Reverse("campaigns-show", req.CampaignID))
}
