package memory

import (
	"context"
	"errors"

	"github.com/ORaneezy/go-runner/internal/domain/entity"
	"github.com/ORaneezy/go-runner/internal/infra/database"
)

type PipelineRepository struct {
	db *database.MemoryDB
}

func NewPipelineRepository(db *database.MemoryDB) *PipelineRepository {
	return &PipelineRepository{
		db: db,
	}
}

func (r *PipelineRepository) CreatePipeline(ctx context.Context, pipeline entity.Pipeline) (
	int,
	error,
) {
	id := len(r.db.Pipelines) + 1
	pipeline.ID = id

	r.db.Pipelines = append(r.db.Pipelines, pipeline)
	return id, nil
}

func (r *PipelineRepository) GetPipelineByID(ctx context.Context, id int) (
	*entity.Pipeline,
	error,
) {
	for _, pipeline := range r.db.Pipelines {
		if pipeline.ID == id {
			return &pipeline, nil
		}
	}

	return nil, errors.New("not found")
}

func (r *PipelineRepository) CreatePipelineRun(ctx context.Context, pipelineID int) (int, error) {
	id := len(r.db.Runs) + 1
	r.db.Runs = append(r.db.Runs, entity.Run{ID: id, PipelineID: pipelineID})

	return id, nil
}
