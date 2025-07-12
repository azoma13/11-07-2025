package service

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/azoma13/archiving-service/config"
	"github.com/azoma13/archiving-service/internal/entity"
)

type TaskService struct {
	tasks  map[int]entity.Task
	taskID int
}

func NewTaskService() *TaskService {
	return &TaskService{
		tasks:  make(map[int]entity.Task),
		taskID: 1,
	}
}

func (s *TaskService) CreateTask(ctx context.Context) (int, error) {

	if len(s.tasks) >= config.Cfg.App.MaxNumTasks {
		return 0, fmt.Errorf("error maximum number of tasks is %d", config.Cfg.App.MaxNumTasks)
	}

	dir := fmt.Sprintf("content/task_%d/", s.taskID)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return 0, fmt.Errorf("not create dir with named %s: %v", dir, err)
	}

	newTask := entity.Task{
		Id:     s.taskID,
		Status: StatusContentEmpty,
		Files:  []entity.File{},
	}

	s.tasks[s.taskID] = newTask

	id := s.taskID
	s.taskID++

	return id, nil
}

func (s *TaskService) AddFile(ctx context.Context, input TaskAddFileInput) error {
	task, ok := s.tasks[input.TaskId]
	if !ok {
		return fmt.Errorf("error task not fount")
	}

	countFile := len(task.Files)
	if countFile >= config.Cfg.App.MaxNumFiles {
		return fmt.Errorf("error maximum number of files is %d", config.Cfg.App.MaxNumFiles)
	}

	err := dowlandFile(input, countFile+1)
	if err != nil {
		return fmt.Errorf("error dowland file: %v", err)
	}

	task.Files = append(task.Files, entity.File{
		IdFile:  countFile + 1,
		UrlFile: input.UrlFile,
	})

	statusContent := StatusTask(countFile + 1)
	task.Status = statusContent

	s.tasks[input.TaskId] = task

	return nil
}

func (s *TaskService) GetStatusTask(ctx context.Context, taskId int) (string, string, error) {
	task, ok := s.tasks[taskId]
	if !ok {
		return "", "", fmt.Errorf("error task not fount")
	}

	countFile := len(task.Files)
	statusContent := StatusTask(countFile)
	if statusContent == StatusContentComplete {
		url, err := s.ArchivingFiles(ctx, task)
		if err != nil {
			return "", "", fmt.Errorf("error archiving file: %v", err)
		}

		err = os.RemoveAll(fmt.Sprintf("./content/task_%d/", taskId))
		if err != nil {
			return "", "", fmt.Errorf("error remove task content dir")
		}

		delete(s.tasks, taskId)

		return statusContent, url, nil
	}

	return statusContent, "", nil
}

func (s *TaskService) ArchivingFiles(ctx context.Context, task entity.Task) (string, error) {

	linkFile := fmt.Sprintf("archives/task_%d.zip", task.Id)
	zipFile, err := os.Create(linkFile)
	if err != nil {
		return "", fmt.Errorf("error create zip file: %v", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	root := fmt.Sprintf("content/task_%d", task.Id)
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)

		return err
	})
	if err != nil {
		return "", fmt.Errorf("error walk in %s dir: %v", root, err)
	}

	absolutePath, err := filepath.Abs(linkFile)
	if err != nil {
		return "", err
	}
	return absolutePath, nil
}

func dowlandFile(input TaskAddFileInput, countFile int) error {
	fileName := path.Base(input.UrlFile)
	if idx := strings.Index(fileName, "?"); idx != -1 {
		fileName = fileName[:idx]
	}

	ext := path.Ext(fileName)
	ok := slices.Contains(config.Cfg.AllowedFileExtensions, ext)
	if !ok {
		return fmt.Errorf("extension file not allowed")
	}

	dst := fmt.Sprintf("content/task_%d/file_%d_%s", input.TaskId, countFile, fileName)
	file, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create file error: %v", err)
	}
	defer file.Close()

	response, err := http.Get(input.UrlFile)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("file was not found link: %s", input.UrlFile)
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func StatusTask(countFile int) string {
	switch countFile {
	case 0:
		return StatusContentEmpty
	case config.Cfg.MaxNumFiles:
		return StatusContentComplete
	default:
		return StatusContentInProgress
	}
}
