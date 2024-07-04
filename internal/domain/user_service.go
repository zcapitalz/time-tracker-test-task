package domain

import (
	"log/slog"
	"time-tracker/internal/utils/slogutils"

	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
)

type UserService struct {
	userStorage       UserStorage
	userInfoConnector UserInfoConnector
}

func NewUserService(
	userStorage UserStorage,
	userInfoConnector UserInfoConnector) *UserService {
	return &UserService{
		userStorage:       userStorage,
		userInfoConnector: userInfoConnector,
	}
}

type UserStorage interface {
	CreateUser(user *User) (*User, error)
	GetUsersPage(afterUserID *ksuid.KSUID, limit int, filters *UserFilters) ([]User, error)
	UpdateUser(userUpdate *UserUpdate) error
	DeleteUser(userID ksuid.KSUID) error
}

type UserInfoConnector interface {
	GetUserInfo(passportNumberAndSeries string) (*UserInfo, error)
}

func (s *UserService) CreateUser(passportSeriesAndNumber string) (*User, error) {
	userInfo, err := s.userInfoConnector.GetUserInfo(passportSeriesAndNumber)
	if err != nil {
		err = errors.Wrap(err, "get user info")
		slog.Error("", slogutils.ErrorAttr(err))
		return nil, err
	}

	user := &User{
		ID:                      ksuid.New(),
		PassportSeriesAndNumber: passportSeriesAndNumber,
		Name:                    userInfo.Name,
		Surname:                 userInfo.Surname,
		Patronymic:              userInfo.Patronymic,
		Address:                 userInfo.Address,
	}
	user, err = s.userStorage.CreateUser(user)
	if err != nil {
		err = errors.Wrap(err, "save users to storage")
		slog.Error("", slogutils.ErrorAttr(err))
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUsersPage(afterUserID *ksuid.KSUID, limit int, filters *UserFilters) ([]User, error) {
	users, err := s.userStorage.GetUsersPage(afterUserID, limit, filters)
	if err != nil {
		err = errors.Wrap(err, "get users from storage")
		slog.Error("", slogutils.ErrorAttr(err))
	}
	return users, err
}

func (s *UserService) UpdateUser(userUpdate *UserUpdate) error {
	err := s.userStorage.UpdateUser(userUpdate)
	if err != nil {
		err = errors.Wrap(err, "update user in storage")
		slog.Error("", slogutils.ErrorAttr(err))
	}
	return err
}

func (s *UserService) DeleteUser(userID ksuid.KSUID) error {
	err := s.userStorage.DeleteUser(userID)
	if err != nil {
		err = errors.Wrap(err, "delete user from storage")
		slog.Error("", slogutils.ErrorAttr(err))
	}
	return err
}
