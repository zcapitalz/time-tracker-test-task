package storages

import (
	"context"
	"log/slog"
	"time-tracker/internal/domain"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
)

type UserStorage struct {
	SQLStorage
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	storage := new(UserStorage)
	storage.Init(db)
	return storage
}

func (s *UserStorage) CreateUser(user *domain.User) (*domain.User, error) {
	builder := s.builder.
		Insert("users").
		Columns(`id, passport_series_and_number, name,
			surname, patronymic, address`).
		Values(user.ID, user.PassportSeriesAndNumber, user.Name,
			user.Surname, user.Patronymic, user.Address)

	query, params, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	_, err = s.db.Query(query, params...)
	if err != nil {
		return nil, errors.Wrap(err, "execute query")
	}

	return user, nil
}

func (s *UserStorage) GetUsersPage(afterUserID *ksuid.KSUID, limit int, filters *domain.UserFilters) ([]domain.User, error) {
	builder := s.builder.
		Select(`id, passport_series_and_number,
			name, surname, patronymic, address`).
		From("users").
		OrderBy("id").
		Limit((uint64)(limit))

	if afterUserID != nil {
		builder = builder.Where(sq.Gt{"id": afterUserID})
	}

	if filters != nil {
		for col, val := range map[string]struct {
			val   any
			isNil bool
		}{
			"id":                       {val: filters.Equal.ID, isNil: filters.Equal.ID == nil},
			"passport_sies_and_number": {val: filters.Equal.PassportSeriesAndNumber, isNil: filters.Equal.PassportSeriesAndNumber == nil},
			"name":                     {val: filters.Equal.Name, isNil: filters.Equal.Name == nil},
			"surname":                  {val: filters.Equal.Surname, isNil: filters.Equal.Surname == nil},
			"patronymic":               {val: filters.Equal.Patronymic, isNil: filters.Equal.Patronymic == nil},
			"address":                  {val: filters.Equal.Address, isNil: filters.Equal.Address == nil},
		} {
			if !val.isNil {
				builder = builder.Where(sq.Eq{col: val.val})
			}
		}
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}
	slog.Debug("query", "query", query)

	var users []domain.User
	err = sqlscan.Select(context.TODO(), s.db, &users, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "execute query")
	}

	return users, nil
}

func (s *UserStorage) UpdateUser(userUpdate *domain.UserUpdate) error {
	builder := s.builder.
		Update("users").
		Where(sq.Eq{"id": userUpdate.ID})

	for col, val := range map[string]struct {
		val   any
		isNil bool
	}{
		"name":                       {val: userUpdate.Name, isNil: userUpdate.Name == nil},
		"surname":                    {val: userUpdate.Surname, isNil: userUpdate.Surname == nil},
		"patronymic":                 {val: userUpdate.Patronymic, isNil: userUpdate.Patronymic == nil},
		"address":                    {val: userUpdate.Address, isNil: userUpdate.Address == nil},
		"passport_series_and_number": {val: userUpdate.PassportSeriesAndNumber, isNil: userUpdate.PassportSeriesAndNumber == nil},
	} {
		if !val.isNil {
			builder = builder.Set(col, val.val)
		}
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "build query")
	}
	slog.Info("query", "query", query)

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "execute query")
	}

	return nil
}

func (s *UserStorage) DeleteUser(userID ksuid.KSUID) error {
	builder := s.builder.
		Delete("users").
		Where(sq.Eq{"id": userID})

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
