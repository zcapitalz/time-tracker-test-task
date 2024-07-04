package domain

import "github.com/segmentio/ksuid"

type User struct {
	ID                      ksuid.KSUID
	PassportSeriesAndNumber string
	Name                    string
	Surname                 string
	Patronymic              string
	Address                 string
}

type UserInfo struct {
	Name       string
	Surname    string
	Patronymic string
	Address    string
}

type UserUpdate struct {
	ID                      ksuid.KSUID
	Name                    *string
	Surname                 *string
	Patronymic              *string
	Address                 *string
	PassportSeriesAndNumber *string
}

type UserFilters struct {
	Equal struct {
		ID                      *ksuid.KSUID
		Name                    *string
		Surname                 *string
		Patronymic              *string
		Address                 *string
		PassportSeriesAndNumber *string
	}
}
