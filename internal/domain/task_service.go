package domain

import (
	"log/slog"
	"sort"
	"time"
	"time-tracker/internal/utils/slogutils"

	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
)

type TaskService struct {
	taskStorage TaskStorage
}

type TaskStorage interface {
	SetTaskStatus(taskID ksuid.KSUID, userID ksuid.KSUID, status TaskStatus) error
	GetTaskStatusesByUserAndTime(userID ksuid.KSUID, t time.Time) (map[ksuid.KSUID]TaskStatus, error)
	GetTaskStatusChangesByUserAndPeriod(userID ksuid.KSUID, periodStart, periodEnd time.Time) ([]TaskStatusChange, error)
}

func NewTaskService(taskStorage TaskStorage) *TaskService {
	return &TaskService{
		taskStorage: taskStorage,
	}
}

func (t *TaskService) SetTaskStatus(taskID ksuid.KSUID, userID ksuid.KSUID, status TaskStatus) error {
	err := t.taskStorage.SetTaskStatus(taskID, userID, status)
	if err != nil {
		err = errors.Wrap(err, "set task status")
		slog.Error("", slogutils.ErrorAttr(err))
	}
	return err
}

func (t *TaskService) GetTaskSummariesByUserAndPeriod(
	userID ksuid.KSUID,
	periodStart time.Time,
	periodEnd time.Time) ([]TaskSummary, error) {

	if !periodEnd.After(periodStart) {
		return nil, IncorrectPeriodError{"period end should be after period start"}
	}
	if periodEnd.After(time.Now()) {
		periodEnd = time.Now()
	}

	startTaskStatuses, err := t.taskStorage.
		GetTaskStatusesByUserAndTime(userID, periodStart)
	if err != nil {
		err = errors.Wrap(err, "get task statuses by user")
		slog.Error("", slogutils.ErrorAttr(err))
		return nil, err
	}

	taskStatusChangesSlice, err := t.taskStorage.
		GetTaskStatusChangesByUserAndPeriod(userID, periodStart, periodEnd)
	if err != nil {
		err = errors.Wrap(err, "get task status changes by user and period")
		slog.Error("", slogutils.ErrorAttr(err))
		return nil, err
	}
	for taskID, startTaskStatus := range startTaskStatuses {
		taskStatusChangesSlice = append(taskStatusChangesSlice, TaskStatusChange{
			TaskID:     taskID,
			TaskStatus: startTaskStatus,
			Time:       periodStart,
		})
	}

	taskStatusChangesMap := make(map[ksuid.KSUID][]TaskStatusChange)
	for _, taskStatusChange := range taskStatusChangesSlice {
		taskStatusChangesMap[taskStatusChange.TaskID] = append(
			taskStatusChangesMap[taskStatusChange.TaskID],
			taskStatusChange)
	}

	taskSummaries := make([]TaskSummary, 0)
	for taskID, taskStatusChanges := range taskStatusChangesMap {
		sort.Slice(taskStatusChanges, func(i, j int) bool { // sort descending by id, time
			cmp := ksuid.Compare(taskStatusChanges[i].TaskID, taskStatusChanges[j].TaskID)
			if cmp == 0 {
				return taskStatusChanges[i].Time.Before(taskStatusChanges[j].Time)
			}
			return cmp < 0
		})

		var totalInWorkDuration time.Duration
		for i := 1; i < len(taskStatusChanges); i++ {
			if taskStatusChanges[i].TaskStatus == TaskStatusIddle {
				totalInWorkDuration += taskStatusChanges[i].Time.Sub(taskStatusChanges[i-1].Time)
			}
		}
		if taskStatusChanges[len(taskStatusChanges)-1].TaskStatus == TaskStatusInWork {
			totalInWorkDuration += periodEnd.Sub(taskStatusChanges[len(taskStatusChanges)-1].Time)
		}

		taskSummaries = append(taskSummaries, TaskSummary{
			TaskID:              taskID,
			TotalInWorkDuration: totalInWorkDuration,
		})
	}

	return taskSummaries, nil
}
