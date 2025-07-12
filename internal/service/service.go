package service

import (
	"context"

	"github.com/azoma13/archiving-service/internal/entity"
)

type TaskAddFileInput struct {
	TaskId  int
	UrlFile string
}

type TaskGetStatusInput struct {
	TaskId int
}

type Task interface {
	CreateTask(ctx context.Context) (int, error)
	AddFile(ctx context.Context, input TaskAddFileInput) error
	GetStatusTask(ctx context.Context, taskId int) (string, string, error)
	ArchivingFiles(ctx context.Context, task entity.Task) (string, error)
}

type Services struct {
	Task
}

func NewServices() *Services {
	return &Services{
		Task: NewTaskService(),
	}
}
