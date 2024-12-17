package middlewares

import (
	"context"
	"net/http"

	"github.com/danushk97/image-analyzer/pkg/errors"
	pkgLogger "github.com/danushk97/image-analyzer/pkg/logger"
	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, appCtx context.Context, err error) {
	_, log := pkgLogger.Ctx(appCtx)
	log.WithError(err).Error(err.Error())

	var problemDetail = gin.H{
		"code":        errors.INTERNAL_SERVER_ERROR,
		"description": "Something went wrong, please try again in some time.",
	}

	iErr, ok := err.(errors.IError)
	if !ok {
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	if iErr.IsOfType(errors.BAD_REQUEST_ERROR) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":        errors.BAD_REQUEST_ERROR,
			"description": iErr.Error(),
		})

		return
	} else if iErr.IsOfType(errors.AUTHORIZATION_ERROR) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":        errors.AUTHORIZATION_ERROR,
			"description": iErr.Error(),
		})

		return
	} else {
		c.JSON(http.StatusInternalServerError, problemDetail)
	}
}
