package taskcontroller

import (
	"fmt"
	"time-tracker/internal/domain"
)

type taskSummary struct {
	TaskID              string `json:"id"`
	TotalInWorkDuration string `json:"totalInWorkDuration"`
}

func (t *taskSummary) fromEntity(e *domain.TaskSummary) {
	t.TaskID = e.TaskID.String()
	t.TotalInWorkDuration = fmt.Sprintf(
		"%dh %dm",
		int(e.TotalInWorkDuration.Hours()),
		int(e.TotalInWorkDuration.Minutes())%60)
}
