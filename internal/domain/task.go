package domain

import (
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
)

type TaskStatus string

const (
	TaskStatusIddle  TaskStatus = "iddle"
	TaskStatusInWork TaskStatus = "in-work"
)

var taskStatusesSet = map[TaskStatus]struct{}{
	"iddle":   {},
	"in-work": {},
}

func (t TaskStatus) Validate() error {
	_, exists := taskStatusesSet[t]
	if !exists {
		return fmt.Errorf("task status is invalid")
	}
	return nil
}

type TaskSummary struct {
	TaskID              ksuid.KSUID
	TotalInWorkDuration time.Duration
}

type TaskStatusChange struct {
	TaskID     ksuid.KSUID
	TaskStatus TaskStatus
	Time       time.Time
}
