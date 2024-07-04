package connectors

import "time-tracker/internal/domain"

type TestUserInfoConnector struct{}

func NewTestUserInfoConnector() *TestUserInfoConnector {
	return &TestUserInfoConnector{}
}

func (c *TestUserInfoConnector) GetUserInfo(passportNumberAndSeries string) (*domain.UserInfo, error) {
	return &domain.UserInfo{
		Name:       "Nikolay",
		Surname:    "Nikolayev",
		Patronymic: "Nikolayevich",
		Address:    "ул. Пушкина, д. Колотушкина",
	}, nil
}
