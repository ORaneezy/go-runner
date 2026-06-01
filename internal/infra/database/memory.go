package database

import "github.com/ORaneezy/go-runner/internal/domain/entity"

type MemoryDB struct {
	Pipelines []entity.Pipeline
	Runs      []entity.Run
	RunLogs   []entity.RunLog
}

func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		Pipelines: []entity.Pipeline{},
	}
}
