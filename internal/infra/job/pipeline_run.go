package job

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/ORaneezy/go-runner/internal/domain/entity"
	"github.com/ORaneezy/go-runner/internal/usecase/pipeline"
)

type LogsInserter interface {
	InsertBulk(ctx context.Context, runID int, messages []string) error
}

type PipelineRunJob struct {
	jobs           chan entity.Job
	ctx            context.Context
	cancel         context.CancelFunc
	pipelineGetter pipeline.PipelineGetter
	logsInserter   LogsInserter
}

func NewPipelineRunJob(
	pipelineGetter pipeline.PipelineGetter,
	logsInserter LogsInserter,
) *PipelineRunJob {
	ctx, cancel := context.WithCancel(context.Background())
	return &PipelineRunJob{
		jobs:           make(chan entity.Job, 5),
		ctx:            ctx,
		pipelineGetter: pipelineGetter,
		logsInserter:   logsInserter,
		cancel:         cancel,
	}
}

func (j *PipelineRunJob) Enqueue(ctx context.Context, job entity.Job) error {
	select {
	case j.jobs <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-j.ctx.Done():
		return j.ctx.Err()
	}
}

func (j *PipelineRunJob) Start() {
	for {
		select {
		case job := <-j.jobs:
			ctx, cancel := context.WithTimeout(j.ctx, time.Minute*10)

			j.processJob(ctx, job)
			cancel()
		case <-j.ctx.Done():
			return
		}
	}
}

func (j *PipelineRunJob) Stop() {
	j.cancel()
	close(j.jobs)
}

func (j *PipelineRunJob) processJob(ctx context.Context, job entity.Job) {
	p, err := j.pipelineGetter.GetPipelineByID(ctx, job.PipelineID)
	if err != nil {
		return
	}

	for _, s := range p.Steps {
		var logs []string
		logs = append(logs, fmt.Sprintf("running step: %v", s.Name))
		cmd := exec.CommandContext(ctx, "sh", "-c", s.Run)
		cmd.Dir = p.WorkDirectory

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err = cmd.Run(); err != nil {
			logs = append(
				logs, fmt.Sprintf(
					"error: %v, output: %v",
					err,
					stderr.String(),
				),
			)

		}

		logs = append(logs, stdout.String())
		_ = j.logsInserter.InsertBulk(ctx, job.RunID, logs)
	}
}
