package memory

import (
	"context"
	"errors"

	"github.com/ORaneezy/go-runner/internal/domain/entity"
	"github.com/ORaneezy/go-runner/internal/infra/database"
)

type RunRepository struct {
	db *database.MemoryDB
}

func NewRunRepository(db *database.MemoryDB) *RunRepository {
	return &RunRepository{
		db: db,
	}
}

func (r *RunRepository) GetRunByID(ctx context.Context, runID int) (*entity.Run, error) {
	for _, run := range r.db.Runs {
		if run.ID == runID {
			return &run, nil
		}
	}

	return nil, errors.New("not found")
}

func (r *RunRepository) GetLogsByRunID(ctx context.Context, runID int) ([]entity.RunLog, error) {
	var logs []entity.RunLog
	for _, log := range r.db.RunLogs {
		if log.RunID == runID {
			logs = append(logs, log)
		}
	}

	return logs, nil
}

func (r *RunRepository) InsertBulk(ctx context.Context, runID int, logs []string) error {
	for _, log := range logs {
		id := len(r.db.RunLogs) + 1
		r.db.RunLogs = append(r.db.RunLogs, entity.RunLog{ID: id, RunID: runID, Message: log})
	}

	return nil
}

func (r *RunRepository) Insert(ctx context.Context, runID int, message string) error {

	id := len(r.db.RunLogs)
	r.db.RunLogs = append(r.db.RunLogs, entity.RunLog{ID: id, RunID: runID, Message: message})

	return nil
}
