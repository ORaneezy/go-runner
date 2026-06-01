package pipeline

import (
	"context"

	"github.com/ORaneezy/go-runner/internal/domain/entity"
	"github.com/yaml/go-yaml"
)

type PipelineCreator interface {
	CreatePipeline(ctx context.Context, pipeline entity.Pipeline) (int, error)
}

type CreateUsecase struct {
	creator PipelineCreator
}

func NewCreateUsecase(creator PipelineCreator) *CreateUsecase {
	return &CreateUsecase{
		creator: creator,
	}
}

func (u *CreateUsecase) Execute(ctx context.Context, content []byte) (int, error) {
	var pipeline entity.Pipeline

	err := yaml.Unmarshal(content, &pipeline)
	if err != nil {
		return -1, err
	}

	id, err := u.creator.CreatePipeline(ctx, pipeline)
	if err != nil {
		return -1, err
	}

	return id, nil
}
