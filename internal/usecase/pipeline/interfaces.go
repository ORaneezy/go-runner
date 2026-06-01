package pipeline

import (
	"context"

	"github.com/ORaneezy/go-runner/internal/domain/entity"
)

type PipelineGetter interface {
	GetPipelineByID(ctx context.Context, id int) (*entity.Pipeline, error)
}

type RunCreator interface {
	CreatePipelineRun(ctx context.Context, id int) (int, error)
}

type JobEnqueuer interface {
	Enqueue(ctx context.Context, job entity.Job) error
}
