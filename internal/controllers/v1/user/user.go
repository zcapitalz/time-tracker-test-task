package usercontroller

import (
	"time-tracker/internal/domain"
)

type user struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	Surname                 string `json:"surname"`
	Patronymic              string `json:"patronymic"`
	Address                 string `json:"address"`
	PassportSeriesAndNumber string `json:"passportNumber"`
}

func (u *user) fromEntity(e *domain.User) {
	u.ID = e.ID.String()
	u.Name = e.Name
	u.Surname = e.Surname
	u.Patronymic = e.Patronymic
	u.Address = e.Address
	u.PassportSeriesAndNumber = e.PassportSeriesAndNumber
}
