package taskcontroller

import (
	"time"
	"time-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

type TaskController struct {
	tasksService TasksService
}

type TasksService interface {
	SetTaskStatus(taskID ksuid.KSUID, userID ksuid.KSUID, taskStatus domain.TaskStatus) error
	GetTaskSummariesByUserAndPeriod(userID ksuid.KSUID, periodStart, priodEnd time.Time) ([]domain.TaskSummary, error)
}

func NewTasksController(tasksService TasksService) *TaskController {
	return &TaskController{
		tasksService: tasksService,
	}
}

func (c *TaskController) RegisterRoutes(engine *gin.Engine) {
	routerGroup := engine.Group("/api/v1")

	routerGroup.PUT("/users/:userID/tasks/:taskID/status", c.setTaskStatus)
	routerGroup.GET("/users/:userID/task-summaries", c.getTaskSummariesByUserAndPeriod)
}
