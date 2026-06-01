package run

import (
	"context"

	"github.com/ORaneezy/go-runner/internal/api/dto/response"
	"github.com/ORaneezy/go-runner/internal/domain/entity"
	"github.com/ORaneezy/go-runner/pkg/mapper"
)

type LogsGetter interface {
	GetLogsByRunID(ctx context.Context, runID int) ([]entity.RunLog, error)
}

type GetLogsUsecase struct {
	logsGetter LogsGetter
}

func NewGetLogsUsecase(getter LogsGetter) *GetLogsUsecase {
	return &GetLogsUsecase{logsGetter: getter}
}

func (u *GetLogsUsecase) Execute(ctx context.Context, runID int) ([]response.RunLog, error) {
	logs, err := u.logsGetter.GetLogsByRunID(ctx, runID)
	if err != nil {
		return nil, err
	}

	resp := mapper.Map(
		logs, func(l entity.RunLog) response.RunLog {
			return response.RunLog{
				Id:      l.ID,
				Message: l.Message,
			}
		},
	)

	return resp, nil
}
