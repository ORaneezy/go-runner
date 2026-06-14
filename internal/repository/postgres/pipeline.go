package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/ORaneezy/go-runner/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PipelineRepository struct {
	pool *pgxpool.Pool
}

func NewPipelineRepository(pool *pgxpool.Pool) *PipelineRepository {
	return &PipelineRepository{
		pool: pool,
	}
}

func (r *PipelineRepository) CreatePipeline(ctx context.Context, pipeline entity.Pipeline) (
	int,
	error,
) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}

	defer tx.Rollback(ctx)

	var id int
	err = tx.QueryRow(
		ctx,
		"INSERT INTO pipelines (name, work_directory) VALUES ($1, $2) RETURNING id",
		pipeline.Name, pipeline.WorkDirectory,
	).
		Scan(&id)

	if err != nil {
		return id, err
	}

	valueStr := make([]string, 0, len(pipeline.Steps))
	args := make([]any, 0, 4*len(pipeline.Steps))
	total := 0
	for i, s := range pipeline.Steps {
		valueStr = append(
			valueStr, fmt.Sprintf("($%d, $%d, $%d, $%d)", total+1, total+2, total+3, total+4),
		)
		args = append(args, s.Name, i+1, s.Command, id)
		total += 4
	}

	sql := fmt.Sprintf(
		`INSERT INTO pipeline_steps (name, sequence_order, command, pipeline_id) VALUES %s`,
		strings.Join(valueStr, ","),
	)

	if _, err = tx.Exec(ctx, sql, args...); err != nil {
		return id, err
	}

	if err = tx.Commit(ctx); err != nil {
		return id, err
	}

	return id, nil
}

func (r *PipelineRepository) GetPipelineByID(ctx context.Context, id int) (
	*entity.Pipeline,
	error,
) {
	var pipeline entity.Pipeline
	err := r.pool.QueryRow(
		ctx,
		"SELECT id, name, work_directory FROM pipelines WHERE id = $1",
		id,
	).
		Scan(
			&pipeline.ID, &pipeline.Name, &pipeline.WorkDirectory,
		)

	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(
		ctx, `SELECT id, name, sequence_order, command FROM pipeline_steps WHERE pipeline_id = $1`,
		pipeline.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []entity.Step
	for rows.Next() {
		var step entity.Step
		err = rows.Scan(&step.ID, &step.Name, &step.SequenceOrder, &step.Command)
		if err != nil {
			return nil, err
		}

		steps = append(steps, step)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	pipeline.Steps = steps

	return &pipeline, nil
}

func (r *PipelineRepository) CreatePipelineRun(ctx context.Context, pipelineID int) (int, error) {
	var runID int
	err := r.pool.QueryRow(
		ctx,
		"INSERT INTO pipeline_runs (status, pipeline_id) VALUES ($1, $2) RETURNING id",
		entity.RunStatusWaiting,
		pipelineID,
	).Scan(&runID)

	if err != nil {
		return 0, err
	}

	return runID, nil
}
