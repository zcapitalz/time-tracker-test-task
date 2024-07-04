package usercontroller

import (
	"fmt"
	httputils "time-tracker/internal/controllers/utils/http"
	"time-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
)

type updateUsersRequest struct {
	User struct {
		Name                    *string `json:"name"`
		Surname                 *string `json:"surname"`
		Patronymic              *string `json:"patronymic"`
		Address                 *string `json:"address"`
		PassportSeriesAndNumber *string `json:"passportSeriesAndNumber"`
	} `json:"user"`
}

type updateUserResponse struct {
	User user `json:"user"`
}

// @Summary Update a user
// @Description Update a user with the provided user ID
// @Accept json
// @Produce json
// @Param userID path string true "User ID to update"
// @Param updateUserRequest body updateUsersRequest true "Update user request"
// @Success 200 {object} updateUserResponse "User updated successfully"
// @Failure 400 {object} httputils.HTTPError "Bad request"
// @Failure 422 {object} httputils.HTTPError "Unprocessable entity"
// @Failure 500 {object} httputils.HTTPError "Internal server error"
// @Router /users/{userID} [patch]
func (c *UserController) updateUser(ctx *gin.Context) {
	if ctx.Param("userID") == "" {
		httputils.BadRequest(ctx, fmt.Errorf("user id not provided"))
		return
	}
	userID, err := ksuid.Parse(ctx.Param("userID"))
	if err != nil {
		httputils.BadRequest(ctx, errors.Wrap(err, "parse user id"))
		return
	}

	var req updateUsersRequest
	err = ctx.BindJSON(&req)
	if err != nil {
		httputils.BindJSONError(ctx, err)
		return
	}

	err = c.userService.UpdateUser(&domain.UserUpdate{
		ID:                      userID,
		Name:                    req.User.Name,
		Surname:                 req.User.Surname,
		Patronymic:              req.User.Patronymic,
		Address:                 req.User.Address,
		PassportSeriesAndNumber: req.User.PassportSeriesAndNumber,
	})
	switch err.(type) {
	case nil:
	case domain.UserNotFoundError:
		httputils.UprocessableContent(ctx, err)
		return
	default:
		httputils.InternalError(ctx)
		return
	}
}
