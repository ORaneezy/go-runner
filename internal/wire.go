//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/ORaneezy/go-runner/internal/api"
	"github.com/ORaneezy/go-runner/internal/api/transport/http/handler"
	"github.com/ORaneezy/go-runner/internal/config"
	"github.com/ORaneezy/go-runner/internal/infra/database"
	"github.com/ORaneezy/go-runner/internal/infra/job"
	"github.com/ORaneezy/go-runner/internal/repository/postgres"
	"github.com/ORaneezy/go-runner/internal/usecase"
	"github.com/ORaneezy/go-runner/internal/usecase/pipeline"
	"github.com/ORaneezy/go-runner/internal/usecase/run"
	"github.com/google/wire"
)

var PipelineRepositorySet = wire.NewSet(
	postgres.NewPipelineRepository,
	wire.Bind(new(pipeline.PipelineCreator), new(*postgres.PipelineRepository)),
	wire.Bind(new(usecase.PipelineGetter), new(*postgres.PipelineRepository)),
	wire.Bind(new(usecase.RunCreator), new(*postgres.PipelineRepository)),
)

var RunRepositorySet = wire.NewSet(
	postgres.NewRunRepository,
	wire.Bind(new(run.RunGetter), new(*postgres.RunRepository)),
	wire.Bind(new(run.LogsGetter), new(*postgres.RunRepository)),
	wire.Bind(new(job.LogsManager), new(*postgres.RunRepository)),
	wire.Bind(new(job.RunManager), new(*postgres.RunRepository)),
)

var PipelineUsecaseSet = wire.NewSet(
	pipeline.NewCreateUsecase,
	pipeline.NewGetUsecase,
	run.NewRunUsecase,
)

var RunUsecaseSet = wire.NewSet(
	run.NewGetUsecase,
	run.NewGetLogsUsecase,
)

var PipelineRunJobSet = wire.NewSet(
	job.NewPipelineRunJob,
	wire.Bind(new(usecase.JobEnqueuer), new(*job.PipelineRunJob)),
)

func InitAPI(isDebug bool) (*api.API, error) {
	wire.Build(
		config.NewConfig,
		database.NewPostgres,
		PipelineRepositorySet,
		RunRepositorySet,
		PipelineRunJobSet,
		PipelineUsecaseSet,
		handler.NewPipelineHandler,
		RunUsecaseSet,
		handler.NewRunHandler,
		api.New,
	)

	return nil, nil
}
