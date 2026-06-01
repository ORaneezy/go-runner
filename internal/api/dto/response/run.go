package response

type Run struct {
	ID         int    `json:"id"`
	PipelineID int    `json:"pipeline_id"`
	Status     string `json:"status"`
}

type RunLog struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}
