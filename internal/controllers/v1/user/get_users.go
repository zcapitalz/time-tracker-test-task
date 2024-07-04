package usercontroller

import (
	"encoding/json"
	"net/http"
	httputils "time-tracker/internal/controllers/utils/http"
	"time-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
)

const (
	defaultUserPaginationLimit = 10
)

type getUsersRequest struct {
	UserFilters    *string `form:"filters"`
	UserPagination *string `form:"pagination"`
}

type userPagination struct {
	AfterUserID *ksuid.KSUID `json:"afterUserID"`
	Limit       int          `json:"limit"`
}

type getUsersResponse struct {
	Users []user `json:"users"`
}

// @Summary Get users
// @Description Get a list of users based on filters and pagination
// @Produce json
// @Param filters query string false "User filters"
// @Param pagination query string false "User pagination"
// @Success 200 {object} getUsersResponse "List of users"
// @Failure 400 {object} httputils.HTTPError "Bad request"
// @Failure 500 {object} httputils.HTTPError "Internal server error"
// @Router /users [get]
func (c *UserController) getUsers(ctx *gin.Context) {
	var req getUsersRequest
	err := ctx.BindQuery(&req)
	if err != nil {
		httputils.BadRequest(ctx, errors.Wrap(err, "parse query"))
		return
	}

	var userFilters *domain.UserFilters
	if req.UserFilters != nil {
		userFilters = new(domain.UserFilters)
		err = json.Unmarshal([]byte(*req.UserFilters), userFilters)
		if err != nil {
			httputils.BindQueryError(ctx, errors.Wrap(err, "parse filters"))
			return
		}
	}
	userPagination := userPagination{
		Limit: defaultUserPaginationLimit,
	}
	if req.UserPagination != nil {
		err = json.Unmarshal([]byte(*req.UserPagination), &userPagination)
		if err != nil {
			httputils.BindQueryError(ctx, errors.Wrap(err, "parse pagination"))
			return
		}
	}

	userEntities, err := c.userService.GetUsersPage(
		userPagination.AfterUserID,
		userPagination.Limit,
		userFilters)
	if err != nil {
		httputils.InternalError(ctx)
		return
	}

	users := make([]user, len(userEntities))
	for i := 0; i < len(users); i++ {
		(&users[i]).fromEntity(&userEntities[i])
	}
	ctx.JSON(http.StatusOK, getUsersResponse{
		Users: users,
	})
}
