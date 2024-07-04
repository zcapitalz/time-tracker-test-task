package usercontroller

import (
	"net/http"
	httputils "time-tracker/internal/controllers/utils/http"
	"time-tracker/internal/domain"

	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	PassportNumberAndSeries string `json:"passportNumber" binding:"required"`
}

type createUserResponse struct {
	User user `json:"user"`
}

// @Summary Create a new user
// @Description Create a new user with the provided passport number and series
// @Accept json
// @Produce json
// @Param request body createUserRequest true "User creation request"
// @Success 200 {object} createUserResponse "User created successfully"
// @Failure 400 {object} httputils.HTTPError
// @Failure 409 {object} httputils.HTTPError
// @Failure 500 {object} httputils.HTTPError
// @Router /users [post]
func (c *UserController) createUser(ctx *gin.Context) {
	var req createUserRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		httputils.BindJSONError(ctx, err)
		return
	}

	userEntity, err := c.userService.CreateUser(req.PassportNumberAndSeries)
	switch err.(type) {
	case nil:
	case domain.UserAlreadyExistsError:
		httputils.Conflict(ctx, err)
		return
	default:
		httputils.InternalError(ctx)
		return
	}

	user := new(user)
	user.fromEntity(userEntity)
	ctx.JSON(http.StatusOK, createUserResponse{
		User: *user,
	})
}
