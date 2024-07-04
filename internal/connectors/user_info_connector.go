package connectors

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time-tracker/internal/domain"

	"github.com/pkg/errors"
)

const (
	userInfoEndpointURL = "http://localhost:8080/user/info"
)

type UserInfoConnector struct{}

func NewUserInfoConnector() *UserInfoConnector {
	return &UserInfoConnector{}
}

type getUserInfoResponse struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Address    string `json:"address"`
}

func (c *UserInfoConnector) GetUserInfo(passportNumberAndSeries string) (*domain.UserInfo, error) {
	params := url.Values{}
	params.Add("passportNumber", passportNumberAndSeries)
	reqURL, err := url.Parse(userInfoEndpointURL)
	if err != nil {
		return nil, errors.Wrap(err, "parse request url")
	}
	reqURL.RawQuery = params.Encode()

	r, err := http.Get(reqURL.String())
	if err != nil {
		return nil, errors.Wrap(err, "make request")
	}
	defer r.Body.Close()

	var resp getUserInfoResponse
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return nil, errors.Wrap(err, "unarshal json body")
	}

	return &domain.UserInfo{
		Name:       resp.Name,
		Surname:    resp.Surname,
		Patronymic: resp.Patronymic,
		Address:    resp.Address,
	}, nil
}
