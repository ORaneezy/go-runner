package run

import (
	"context"

	"github.com/ORaneezy/go-runner/internal/domain/entity"
	"github.com/ORaneezy/go-runner/internal/usecase"
)

type CreateUsecase struct {
	runCreator     usecase.RunCreator
	pipelineGetter usecase.PipelineGetter
	jobEnqueuer    usecase.JobEnqueuer
}

func NewRunUsecase(
	runCreator usecase.RunCreator,
	pipelineGetter usecase.PipelineGetter,
	jobEnqueuer usecase.JobEnqueuer,
) *CreateUsecase {
	return &CreateUsecase{
		runCreator: runCreator, pipelineGetter: pipelineGetter,
		jobEnqueuer: jobEnqueuer,
	}
}

func (u *CreateUsecase) CreatePipelineRun(ctx context.Context, ID int) (int, error) {
	p, err := u.pipelineGetter.GetPipelineByID(ctx, ID)
	if err != nil {
		return 0, err
	}

	runID, err := u.runCreator.CreatePipelineRun(ctx, p.ID)
	if err != nil {
		return 0, err
	}

	err = u.jobEnqueuer.Enqueue(ctx, entity.Job{RunID: runID, Pipeline: *p})
	if err != nil {
		return 0, err
	}

	return runID, err
}
