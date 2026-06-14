package job

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/ORaneezy/go-runner/internal/domain/entity"
)

type LogsManager interface {
	Insert(ctx context.Context, runID int, stepID int, message string) error
}

type RunManager interface {
	SetRunStatus(ctx context.Context, runID int, status entity.RunStatus) error
}

type PipelineRunJob struct {
	jobs        chan entity.Job
	ctx         context.Context
	cancel      context.CancelFunc
	logsManager LogsManager
	runManager  RunManager
}

func NewPipelineRunJob(
	logsManager LogsManager,
	runManager RunManager,
) *PipelineRunJob {
	ctx, cancel := context.WithCancel(context.Background())
	return &PipelineRunJob{
		jobs:        make(chan entity.Job, 5),
		ctx:         ctx,
		logsManager: logsManager,
		runManager:  runManager,
		cancel:      cancel,
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
	_ = j.runManager.SetRunStatus(ctx, job.RunID, entity.RunStatusRunning)
	ss := make([]bool, len(job.Pipeline.Steps), len(job.Pipeline.Steps))

	for i, s := range job.Pipeline.Steps {
		if i != 0 && !ss[i-1] {
			continue
		}

		cmd := exec.CommandContext(ctx, "sh", "-c", s.Command)
		cmd.Dir = job.Pipeline.WorkDirectory

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			_ = j.logsManager.Insert(
				ctx, job.RunID, s.ID,
				fmt.Sprintf("failed to pipe stdout: %v", err),
			)
			continue
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			_ = j.logsManager.Insert(
				ctx, job.RunID, s.ID, fmt.Sprintf(
					"failed to pipe stderr: %v",
					err,
				),
			)
			continue
		}

		if err = cmd.Start(); err != nil {
			_ = j.logsManager.Insert(
				ctx, job.RunID, s.ID, fmt.Sprintf(
					"failed to start command: %v",
					err,
				),
			)
			continue
		}

		wg := sync.WaitGroup{}
		for _, rc := range []io.ReadCloser{stdout, stderr} {
			wg.Add(1)
			go func(rc io.ReadCloser) {
				defer wg.Done()
				defer rc.Close()
				scanner := bufio.NewScanner(rc)
				for scanner.Scan() {
					line := scanner.Text()
					_ = j.logsManager.Insert(ctx, job.RunID, s.ID, line)
				}
			}(rc)
		}

		if err = cmd.Wait(); err != nil {
			_ = j.logsManager.Insert(
				ctx, job.RunID, s.ID,
				fmt.Sprintf("failed to wait command: %v", err),
			)

			continue
		}

		wg.Wait()

		ss[i] = true
	}

	for i, success := range ss {
		if !success {
			_ = j.runManager.SetRunStatus(ctx, job.RunID, entity.RunStatusFailure)
			return
		}

		if i == len(ss)-1 {
			_ = j.runManager.SetRunStatus(ctx, job.RunID, entity.RunStatusSuccess)
		}
	}
}
