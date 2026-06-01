package pipeline

import (
	"context"

	"github.com/ORaneezy/go-runner/internal/api/dto/response"
	"github.com/ORaneezy/go-runner/internal/domain/entity"
	"github.com/ORaneezy/go-runner/pkg/mapper"
)

type GetUsecase struct {
	getter PipelineGetter
}

func NewGetUsecase(getter PipelineGetter) *GetUsecase {
	return &GetUsecase{
		getter: getter,
	}
}

func (u *GetUsecase) Execute(ctx context.Context, id int) (*response.Pipeline, error) {
	p, err := u.getter.GetPipelineByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := response.Pipeline{
		ID:            p.ID,
		Name:          p.Name,
		WorkDirectory: p.WorkDirectory,
		Steps: mapper.Map(
			p.Steps, func(o entity.Step) response.Step {
				return response.Step{
					Name: o.Name,
					Run:  o.Run,
				}
			},
		),
	}

	return &resp, nil
}
