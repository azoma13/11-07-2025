package v1

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/azoma13/archiving-service/internal/service"
	"github.com/labstack/echo/v4"
)

type addFileInput struct {
	TaskId  int    `json:"task_id"`
	UrlFile string `json:"url_file"`
}

type getStatusInput struct {
	TaskId int `json:"task_id"`
}

type taskRoutes struct {
	taskService service.Task
}

func newTaskRoutes(g *echo.Group, taskService service.Task) {
	r := &taskRoutes{
		taskService: taskService,
	}

	g.POST("/create", r.createTask)
	g.POST("/add-file", r.addFile)
	g.GET("/status", r.getStatus)
}

func (r *taskRoutes) createTask(c echo.Context) error {
	id, err := r.taskService.CreateTask(c.Request().Context())
	if err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusInternalServerError, "The server is busy, please try again later")
		return err
	}

	type response struct {
		Id int `json:"task_id"`
	}

	return c.JSON(http.StatusCreated, response{
		Id: id,
	})
}

func (r *taskRoutes) addFile(c echo.Context) error {
	var input addFileInput

	if err := c.Bind(&input); err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	err := r.taskService.AddFile(c.Request().Context(), service.TaskAddFileInput{
		TaskId:  input.TaskId,
		UrlFile: input.UrlFile,
	})
	if err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (r *taskRoutes) getStatus(c echo.Context) error {
	taskIdStr := c.Request().URL.Query().Get("task_id")
	if taskIdStr == "" {
		log.Println("invalid request body")
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return fmt.Errorf("invalid request body")
	}

	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	status, url, err := r.taskService.GetStatusTask(c.Request().Context(), taskId)
	if err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if url == "" {
		type response struct {
			StatusContent string `json:"status_content"`
		}
		return c.JSON(http.StatusOK, response{
			StatusContent: status,
		})
	}

	type response struct {
		StatusContent string `json:"status_content"`
		ArchivingUrl  string `json:"archiving_url"`
	}

	return c.JSON(http.StatusOK, response{
		StatusContent: status,
		ArchivingUrl:  url,
	})
}
