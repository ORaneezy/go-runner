package entity

type Run struct {
	ID         int
	PipelineID int
	Status     string
}

type RunLog struct {
	ID      int
	RunID   int
	Message string
}
