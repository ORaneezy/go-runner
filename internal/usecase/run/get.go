package run

import (
	"context"

	"github.com/ORaneezy/go-runner/internal/api/dto/response"
	"github.com/ORaneezy/go-runner/internal/domain/entity"
)

type RunGetter interface {
	GetRunByID(ctx context.Context, runID int) (*entity.Run, error)
}

type GetUsecase struct {
	runGetter RunGetter
}

func NewGetUsecase(getter RunGetter) *GetUsecase {
	return &GetUsecase{
		runGetter: getter,
	}
}

func (u *GetUsecase) Execute(ctx context.Context, runID int) (*response.Run, error) {
	r, err := u.runGetter.GetRunByID(ctx, runID)
	if err != nil {
		return nil, err
	}

	resp := &response.Run{
		ID:         r.ID,
		PipelineID: r.PipelineID,
		Status:     string(r.Status),
	}

	return resp, nil
}
