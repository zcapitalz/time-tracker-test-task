package usercontroller

import (
	"time-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

type UserController struct {
	userService UserService
}

type UserService interface {
	CreateUser(passportSeriesAndNumber string) (*domain.User, error)
	GetUsersPage(afterUserID *ksuid.KSUID, limit int, filters *domain.UserFilters) ([]domain.User, error)
	UpdateUser(userUpdate *domain.UserUpdate) error
	DeleteUser(userID ksuid.KSUID) error
}

func NewUserController(UserService UserService) *UserController {
	return &UserController{
		userService: UserService,
	}
}

func (c *UserController) RegisterRoutes(engine *gin.Engine) {
	usersGroup := engine.Group("/api/v1/users")
	usersGroup.POST("", c.createUser)
	usersGroup.GET("", c.getUsers)
	usersGroup.DELETE(":userID", c.deleteUser)
	usersGroup.PATCH(":userID", c.updateUser)
}
