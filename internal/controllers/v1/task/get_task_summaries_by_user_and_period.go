package taskcontroller

import (
	"net/http"
	"sort"
	"time"
	httputils "time-tracker/internal/controllers/utils/http"
	"time-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
)

type getUserTaskSummariesByPeriodRequest struct {
	PeriodStart time.Time `form:"periodStart"`
	PeriodEnd   time.Time `form:"periodEnd"`
}

type getUserTaskSummariesByPeriodResponse struct {
	TaskSummaries []taskSummary `json:"taskSummaries"`
}

// @Summary Get task summaries by user and period
// @Description Get the task summaries for a specific user within a given period
// @Produce json
// @Param userID path string true "User ID"
// @Param periodStart query string true "Start of the period" format(date-time)
// @Param periodEnd query string true "End of the period" format(date-time)
// @Success 200 {object} getUserTaskSummariesByPeriodResponse "Success response"
// @Failure 400 {object} httputils.HTTPError "Bad request"
// @Failure 500 {object} httputils.HTTPError "Internal server error"
// @Router /users/{userID}/task-summaries [get]
func (c *TaskController) getTaskSummariesByUserAndPeriod(ctx *gin.Context) {
	userIDStr := ctx.Param("userID")
	userID, err := ksuid.Parse(userIDStr)
	if err != nil {
		httputils.BadRequest(ctx, errors.Wrap(err, "parse user id"))
		return
	}
	var req getUserTaskSummariesByPeriodRequest
	err = ctx.BindQuery(&req)
	if err != nil {
		httputils.BindQueryError(ctx, err)
		return
	}

	taskSummaryEntities, err := c.tasksService.GetTaskSummariesByUserAndPeriod(userID, req.PeriodStart, req.PeriodEnd)
	switch err.(type) {
	case nil:
	case domain.IncorrectPeriodError:
		httputils.BadRequest(ctx, err)
		return
	default:
		httputils.InternalError(ctx)
		return
	}
	sort.Slice(taskSummaryEntities, func(i, j int) bool {
		return taskSummaryEntities[i].TotalInWorkDuration > taskSummaryEntities[j].TotalInWorkDuration
	})

	taskSummaries := make([]taskSummary, len(taskSummaryEntities))
	for i := 0; i < len(taskSummaries); i++ {
		(&taskSummaries[i]).fromEntity(&taskSummaryEntities[i])
	}
	for i := 0; i < len(taskSummaries); i++ {
		if taskSummaries[i].TotalInWorkDuration == "0h 0m" { // task has been in work for less then a minute
			taskSummaries = taskSummaries[:i]
			break
		}
	}

	ctx.JSON(http.StatusOK, getUserTaskSummariesByPeriodResponse{
		TaskSummaries: taskSummaries,
	})
}
