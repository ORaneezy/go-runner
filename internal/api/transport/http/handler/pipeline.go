package handler

import (
	"io"
	"net/http"

	"github.com/ORaneezy/go-runner/internal/usecase/pipeline"
	"github.com/labstack/echo/v5"
)

type Pipeline struct {
	createUsecase *pipeline.CreateUsecase
	getUsecase    *pipeline.GetUsecase
	runUsecase    *pipeline.RunUsecase
}

func NewPipelineHandler(
	createUsecase *pipeline.CreateUsecase,
	getUsecase *pipeline.GetUsecase,
	runUsecase *pipeline.RunUsecase,
) *Pipeline {
	return &Pipeline{
		createUsecase: createUsecase,
		getUsecase:    getUsecase,
		runUsecase:    runUsecase,
	}
}

func (h *Pipeline) CreatePipeline(c *echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	defer c.Request().Body.Close()

	id, err := h.createUsecase.Execute(c.Request().Context(), body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(
		http.StatusCreated, map[string]interface{}{
			"id": id,
		},
	)
}

func (h *Pipeline) GetPipeline(c *echo.Context) error {
	var pipelineID int
	err := echo.PathValuesBinder(c).
		Int("id", &pipelineID).
		BindError()

	if err != nil {
		return err
	}

	resp, err := h.getUsecase.Execute(c.Request().Context(), pipelineID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Pipeline) RunPipeline(c *echo.Context) error {
	var pipelineID int
	err := echo.PathValuesBinder(c).
		Int("id", &pipelineID).
		BindError()

	if err != nil {
		return err
	}

	id, err := h.runUsecase.CreatePipelineRun(c.Request().Context(), pipelineID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(
		http.StatusCreated, map[string]interface{}{
			"id": id,
		},
	)
}
