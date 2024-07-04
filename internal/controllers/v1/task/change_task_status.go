package taskcontroller

import (
	httputils "time-tracker/internal/controllers/utils/http"
	"time-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
)

type changeTaskStatusRequest struct {
	TaskStatus domain.TaskStatus `json:"status" binding:"required"`
}

// @Summary Set task status
// @Description Set the status of a task for a specific user
// @Accept json
// @Param userID path string true "User ID"
// @Param taskID path string true "Task ID"
// @Param request body changeTaskStatusRequest true "Task status request"
// @Success 200 "Success response"
// @Failure 400 {object} httputils.HTTPError "Bad request"
// @Failure 500 {object} httputils.HTTPError "Internal server error"
// @Router /users/{userID}/tasks/{taskID}/status [put]
func (c *TaskController) setTaskStatus(ctx *gin.Context) {
	userIDStr := ctx.Param("userID")
	userID, err := ksuid.Parse(userIDStr)
	if err != nil {
		httputils.BadRequest(ctx, errors.Wrap(err, "parse user id"))
		return
	}
	taskIDStr := ctx.Param("taskID")
	taskID, err := ksuid.Parse(taskIDStr)
	if err != nil {
		httputils.BadRequest(ctx, errors.Wrap(err, "parse task id"))
		return
	}

	var req changeTaskStatusRequest
	err = ctx.BindJSON(&req)
	if err != nil {
		httputils.BindJSONError(ctx, err)
		return
	}
	if err = req.TaskStatus.Validate(); err != nil {
		httputils.BindJSONError(ctx, err)
		return
	}

	err = c.tasksService.SetTaskStatus(taskID, userID, req.TaskStatus)
	if err != nil {
		httputils.InternalError(ctx)
	}
}
