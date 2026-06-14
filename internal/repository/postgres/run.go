package postgres

import (
	"context"

	"github.com/ORaneezy/go-runner/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RunRepository struct {
	pool *pgxpool.Pool
}

func NewRunRepository(pool *pgxpool.Pool) *RunRepository {
	return &RunRepository{
		pool: pool,
	}
}

func (r *RunRepository) GetRunByID(ctx context.Context, runID int) (*entity.Run, error) {
	var run entity.Run
	err := r.pool.QueryRow(ctx, "SELECT id, status FROM pipeline_runs WHERE id = $1", runID).Scan(
		&run.ID, &run.Status,
	)

	if err != nil {
		return nil, err
	}

	return &run, nil
}

func (r *RunRepository) SetRunStatus(
	ctx context.Context, runID int, status entity.RunStatus,
) error {
	_, err := r.pool.Exec(ctx, "UPDATE pipeline_runs SET status = $1 WHERE id = $2", status, runID)
	return err
}

func (r *RunRepository) GetLogsByRunID(ctx context.Context, runID int) ([]entity.RunLog, error) {
	rows, err := r.pool.Query(
		ctx, "SELECT id, message, created_at FROM run_logs WHERE run_id = $1",
		runID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []entity.RunLog
	for rows.Next() {
		var log entity.RunLog
		err = rows.Scan(&log.ID, &log.Message, &log.CreatedAt)
		if err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *RunRepository) Insert(ctx context.Context, runID int, stepID int, message string) error {
	_, err := r.pool.Exec(
		ctx,
		"INSERT INTO run_logs (run_id, step_id, message) VALUES ($1, $2, $3)",
		runID, stepID, message,
	)

	if err != nil {
		return err
	}

	return nil
}
