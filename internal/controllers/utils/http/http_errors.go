package httputils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type HTTPError struct {
	Message string `json:"error,omitempty"`
}

func Error(ctx *gin.Context, status int, err error) {
	ctx.JSON(status, HTTPError{Message: err.Error()})
}

func BadRequest(ctx *gin.Context, err error) {
	Error(ctx, http.StatusBadRequest, err)
}

func BindJSONError(ctx *gin.Context, err error) {
	BadRequest(ctx, errors.Wrap(err, "parse and validate json body"))
}

func BindQueryError(ctx *gin.Context, err error) {
	BadRequest(ctx, errors.Wrap(err, "parse and validate query"))
}

func Conflict(ctx *gin.Context, err error) {
	Error(ctx, http.StatusConflict, err)
}

func UprocessableContent(ctx *gin.Context, err error) {
	Error(ctx, http.StatusUnprocessableEntity, err)
}

func InternalError(ctx *gin.Context) {
	Error(ctx, http.StatusInternalServerError, errors.New(""))
}
