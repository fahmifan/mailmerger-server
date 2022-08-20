package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func systemError(ec echo.Context, err error) error {
	log.Err(err).Msg("")
	return ec.Render(http.StatusOK, "pages/system_error.html", echo.Map{})
}

func notFound(ec echo.Context) error {
	return ec.Render(http.StatusNotFound, "pages/not_found_error.html", echo.Map{})
}

func badRequest(ec echo.Context) error {
	return ec.Render(http.StatusBadRequest, "pages/not_found_error.html", echo.Map{})
}
