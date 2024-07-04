package storages

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"time-tracker/internal/domain"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
)

type TaskStorage struct {
	SQLStorage
}

func NewTaskStorage(db *sqlx.DB) *TaskStorage {
	storage := new(TaskStorage)
	storage.Init(db)
	return storage
}

func (s *TaskStorage) CreateTask(taskID ksuid.KSUID, userID ksuid.KSUID) error {
	builder := s.builder.
		Insert("tasks").
		Columns("id, user_id").
		Values(taskID, userID)

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "build query")
	}

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "execute query")
	}

	return nil
}

func (s *TaskStorage) GetTaskStatus(taskID ksuid.KSUID) (domain.TaskStatus, error) {
	builder := s.builder.
		Select("task_status").
		From("task_status_changes").
		Where(sq.Eq{"task_id": taskID}).
		OrderBy("time DESC").
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return "", errors.Wrap(err, "build query")
	}

	var status domain.TaskStatus
	err = s.db.QueryRowx(query, args...).Scan(&status)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return "", domain.TaskNotFoundError{}
	default:
		return "", errors.Wrap(err, "execute query")
	}
	if err := status.Validate(); err != nil {
		return "", fmt.Errorf("validate db returned data")
	}

	return status, nil
}

func (s *TaskStorage) SetTaskStatus(taskID ksuid.KSUID, userID ksuid.KSUID, status domain.TaskStatus) error {
	currentStatus, err := s.GetTaskStatus(taskID)
	switch err.(type) {
	case nil:
		if status == currentStatus {
			return nil
		}
	case domain.TaskNotFoundError:
		err := s.CreateTask(taskID, userID)
		if err != nil {
			return errors.Wrap(err, "create task")
		}
	default:
		return errors.Wrap(err, "get current task status")
	}

	builder := s.builder.
		Insert("task_status_changes").
		Columns("task_id, task_status").
		Values(taskID, string(status))

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "build query")
	}

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "execute query")
	}

	return nil
}

func (s *TaskStorage) GetTaskStatusesByUserAndTime(userID ksuid.KSUID, t time.Time) (map[ksuid.KSUID]domain.TaskStatus, error) {
	builder := s.builder.
		Select("DISTINCT ON(task_id) task_id, task_status").
		From("task_status_changes tsc").
		InnerJoin("tasks t ON tsc.task_id = t.id").
		Where(sq.Eq{"user_id": userID}).
		Where(sq.Lt{"tsc.time": t}).
		OrderBy("task_id, time DESC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	rows, err := s.db.Queryx(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "execute query")
	}

	taskStatuses := make(map[ksuid.KSUID]domain.TaskStatus)
	for rows.Next() {
		var taskID ksuid.KSUID
		var taskStatus domain.TaskStatus
		err = rows.Scan(&taskID, &taskStatus)
		if err != nil {
			return nil, errors.Wrap(err, "scan returned rows")
		}
		taskStatuses[taskID] = taskStatus
	}

	return taskStatuses, nil
}

func (s *TaskStorage) GetTaskStatusChangesByUserAndPeriod(
	userID ksuid.KSUID,
	periodStart,
	periodEnd time.Time) ([]domain.TaskStatusChange, error) {

	builder := s.builder.
		Select("task_id, task_status, time").
		From("task_status_changes tsc").
		InnerJoin("tasks t ON tsc.task_id = t.id").
		Where(sq.Eq{"user_id": userID}).
		Where(sq.GtOrEq{"time": periodStart}).
		Where(sq.LtOrEq{"time": periodEnd})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	var taskStatusChanges []domain.TaskStatusChange
	err = sqlscan.Select(context.TODO(), s.db, &taskStatusChanges, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "execute query")
	}

	return taskStatusChanges, nil
}
