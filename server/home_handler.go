package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HomeHandler struct{}

func (h HomeHandler) Index(ec echo.Context) error {
	return ec.Render(http.StatusOK, "pages/home/index.html", echo.Map{})
}
