package api

import (
	"context"
	"net/http"

	"github.com/ORaneezy/go-runner/internal/api/transport/http/handler"
	"github.com/ORaneezy/go-runner/internal/infra/job"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type API struct {
	httpSrv          *http.Server
	pipelineQueueJob *job.PipelineRunJob
}

func New(
	pipelineQueueJob *job.PipelineRunJob, pipelineHandler *handler.Pipeline,
	runHandler *handler.Run,
) *API {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	api := e.Group("/api")
	v1 := api.Group("/v1")

	pipelineGroup := v1.Group("/pipeline")

	pipelineGroup.POST("", pipelineHandler.CreatePipeline)
	pipelineGroup.GET("/:id", pipelineHandler.GetPipeline)
	pipelineGroup.POST("/:id/run", pipelineHandler.RunPipeline)

	runGroup := v1.Group("/run")
	runGroup.GET("/:id", runHandler.GetRun)
	runGroup.GET("/:id/logs", runHandler.GetRunLogs)

	srv := http.Server{
		Addr:    ":8080",
		Handler: e,
	}

	return &API{
		httpSrv:          &srv,
		pipelineQueueJob: pipelineQueueJob,
	}
}

func (api *API) Start() error {
	go api.pipelineQueueJob.Start()
	return api.httpSrv.ListenAndServe()
}

func (api *API) Stop(ctx context.Context) error {
	api.pipelineQueueJob.Stop()
	return api.httpSrv.Shutdown(ctx)
}
