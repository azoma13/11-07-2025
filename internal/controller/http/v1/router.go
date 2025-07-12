package v1

import (
	"github.com/azoma13/archiving-service/internal/service"
	"github.com/labstack/echo/v4"
)

func NewRouter(handler *echo.Echo, services *service.Services) {

	handler.GET("/health", func(c echo.Context) error { return c.NoContent(200) })

	v1 := handler.Group("/api/v1")
	{
		newTaskRoutes(v1.Group("/task"), services.Task)
	}
}
