package handler

import (
	"net/http"

	"github.com/ORaneezy/go-runner/internal/usecase/run"
	"github.com/labstack/echo/v5"
)

type Run struct {
	getUsecase     *run.GetUsecase
	getLogsUsecase *run.GetLogsUsecase
}

func NewRunHandler(getUsecase *run.GetUsecase, getLogsUsecase *run.GetLogsUsecase) *Run {
	return &Run{
		getUsecase:     getUsecase,
		getLogsUsecase: getLogsUsecase,
	}
}

func (h *Run) GetRun(c *echo.Context) error {
	var runID int
	err := echo.PathValuesBinder(c).Int("id", &runID).BindError()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := h.getUsecase.Execute(c.Request().Context(), runID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Run) GetRunLogs(c *echo.Context) error {
	var runID int
	err := echo.PathValuesBinder(c).Int("id", &runID).BindError()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := h.getLogsUsecase.Execute(c.Request().Context(), runID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}
