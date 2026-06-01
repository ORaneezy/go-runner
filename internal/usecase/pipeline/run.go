package pipeline

import (
	"context"

	"github.com/ORaneezy/go-runner/internal/domain/entity"
)

type RunUsecase struct {
	runCreator     RunCreator
	pipelineGetter PipelineGetter
	jobEnqueuer    JobEnqueuer
}

func NewRunUsecase(
	runCreator RunCreator,
	pipelineGetter PipelineGetter,
	jobEnqueuer JobEnqueuer,
) *RunUsecase {
	return &RunUsecase{
		runCreator: runCreator, pipelineGetter: pipelineGetter,
		jobEnqueuer: jobEnqueuer,
	}
}

func (u *RunUsecase) CreatePipelineRun(ctx context.Context, ID int) (int, error) {
	p, err := u.pipelineGetter.GetPipelineByID(ctx, ID)
	if err != nil {
		return 0, err
	}

	runID, err := u.runCreator.CreatePipelineRun(ctx, p.ID)
	if err != nil {
		return 0, err
	}

	err = u.jobEnqueuer.Enqueue(ctx, entity.Job{RunID: runID, PipelineID: p.ID})
	if err != nil {
		return 0, err
	}

	return runID, err
}
