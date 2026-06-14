package entity

import "time"

type RunStatus string

const (
	RunStatusWaiting RunStatus = "waiting"
	RunStatusRunning RunStatus = "running"
	RunStatusSuccess RunStatus = "success"
	RunStatusFailure RunStatus = "failure"
)

type Run struct {
	ID         int
	PipelineID int
	Status     RunStatus
}

type RunLog struct {
	ID        int
	RunID     int
	CreatedAt time.Time
	Message   string
}
