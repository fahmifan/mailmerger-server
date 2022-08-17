package server

import (
	"net/http"

	"github.com/fahmifan/mailmerger-server/service"
	"github.com/labstack/echo/v4"
)

type TemplateHandler struct {
	*Server
}

// render new form
func (t TemplateHandler) New(ec echo.Context) (err error) {
	return ec.Render(http.StatusOK, "pages/templates/new.html", nil)
}

func (t TemplateHandler) List(ec echo.Context) (err error) {
	templates, err := t.service.TemplateService.FindAll(ec.Request().Context())
	if err != nil {
		return systemError(ec, err)
	}
	return ec.Render(http.StatusOK, "pages/templates/index.html", echo.Map{"templates": templates})
}

func (t TemplateHandler) Create(ec echo.Context) (err error) {
	req := service.CreateTemplateRequest{}
	if err = ec.Bind(&req); err != nil {
		return badRequest(ec)
	}
	_, err = t.service.TemplateService.Create(ec.Request().Context(), req)
	if err != nil {
		return systemError(ec, err)
	}

	return ec.Redirect(http.StatusSeeOther, ec.Echo().Reverse("templates"))
}
