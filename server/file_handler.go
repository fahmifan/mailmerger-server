package server

import (
	"errors"
	"net/http"

	"github.com/fahmifan/mailmerger-server/service"
	"github.com/labstack/echo/v4"
)

type FileHandler struct {
	*Server
}

func (f FileHandler) Show(ec echo.Context) (err error) {
	fileName := ec.Param("file_name")
	if fileName == "" {
		return notFound(ec)
	}

	rd, err := f.service.FileService.Find(ec.Request().Context(), fileName)
	if errors.Is(err, service.ErrNotFound) {
		return notFound(ec)
	} else if err != nil {
		return systemError(ec, err)
	}
	defer rd.Close()

	err = ec.Stream(http.StatusOK, "application/octet-stream", rd)
	return err
}
