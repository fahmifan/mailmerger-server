package server

import (
	"errors"
	"net/http"

	"github.com/fahmifan/mailmerger-server/service"
	"github.com/fahmifan/ulids"
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

func (t TemplateHandler) Show(ec echo.Context) (err error) {
	id, err := ulids.Parse(ec.Param("id"))
	if err != nil {
		return notFound(ec)
	}

	template, err := t.service.TemplateService.FindByID(ec.Request().Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		return notFound(ec)
	}

	return ec.Render(http.StatusOK, "pages/templates/show.html", echo.Map{"template": template})
}

// render Edit page
func (t TemplateHandler) Edit(ec echo.Context) (err error) {
	id, err := ulids.Parse(ec.Param("id"))
	if err != nil {
		return notFound(ec)
	}

	template, err := t.service.TemplateService.FindByID(ec.Request().Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		return notFound(ec)
	}
	if err != nil {
		return systemError(ec, err)
	}

	return ec.Render(http.StatusOK, "pages/templates/edit.html", echo.Map{"template": template})
}

func (t TemplateHandler) Update(ec echo.Context) (err error) {
	req := service.UpdateTemplateRequest{}
	if ec.Bind(&req) != nil {
		return badRequest(ec)
	}

	template, err := t.service.TemplateService.Update(ec.Request().Context(), req)
	if errors.Is(err, service.ErrNotFound) {
		return notFound(ec)
	}
	if err != nil {
		return systemError(ec, err)
	}

	return ec.Redirect(http.StatusSeeOther, ec.Echo().Reverse("templates-show", template.ID))
}
