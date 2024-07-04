package usercontroller

import (
	"fmt"
	"net/http"
	httputils "time-tracker/internal/controllers/utils/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
)

// @Summary Delete a user
// @Description Delete a user with the provided user ID
// @Param userID path string true "User ID to delete"
// @Success 200 {object} nil "User deleted successfully"
// @Failure 400 {object} httputils.HTTPError "Bad request"
// @Failure 500 {object} httputils.HTTPError "Internal server error"
// @Router /users/{userID} [delete]
func (c *UserController) deleteUser(ctx *gin.Context) {
	if ctx.Param("userID") == "" {
		httputils.BadRequest(ctx, fmt.Errorf("user id not provided"))
		return
	}
	userID, err := ksuid.Parse(ctx.Param("userID"))
	if err != nil {
		httputils.BadRequest(ctx, errors.Wrap(err, "parse user id"))
		return
	}

	err = c.userService.DeleteUser(userID)
	if err != nil {
		httputils.InternalError(ctx)
		return
	}

	ctx.Status(http.StatusOK)
}
