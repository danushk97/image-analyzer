package middlewares

import (
	"context"
	"net/http"

	"github.com/danushk97/image-analyzer/internal/constants"
	internaErr "github.com/danushk97/image-analyzer/internal/errors"
	"github.com/danushk97/image-analyzer/pkg/contextkey"
	"github.com/danushk97/image-analyzer/pkg/errors"
	pkgLogger "github.com/danushk97/image-analyzer/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ErrorResponse(c *gin.Context, ctx context.Context, err error) {
	log := pkgLogger.Ctx(ctx)
	log.WithError(err).Error(err.Error())

	var problemDetail = gin.H{
		"code":        internaErr.ServerError,
		"description": "Something went wrong, please try again in some time.",
	}

	iErr, ok := err.(errors.IError)
	if !ok {
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	if iErr.IsOfType(errors.BAD_REQUEST_ERROR) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":        internaErr.BadRequesterror,
			"description": iErr.Error(),
		})

		return
	} else if iErr.IsOfType(errors.AUTHORIZATION_ERROR) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":        internaErr.Unauthorized,
			"description": iErr.Error(),
		})

		return
	} else {
		c.JSON(http.StatusInternalServerError, problemDetail)
	}
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := AppCtx(c)
		logger := pkgLogger.Ctx(ctx)

		logger.Info("AUTH_VALIDATION")

		// Extract user_id from header
		userID := c.GetHeader(constants.HeaderUserId)
		if userID == "" {
			logger.Warn("USER_ID_NOT_FOUND")
			ErrorResponse(
				c,
				ctx,
				errors.NewAuthorizationError(internaErr.Unauthorized),
			)
			c.Abort()
		}

		// Set user_id in context
		c.Set("user_id", userID)

		logger.Info("AUTH_VALIDATION_SUCCESS")

		c.Next()
	}
}

func AppCtx(gc *gin.Context) context.Context {

	val, ok := gc.Get(contextkey.AppCtx.String())
	if ok {
		appCtx := val.(context.Context)
		return appCtx
	}

	requestId, ok := gc.Keys[constants.HeaderRequestId]
	if !ok {
		requestId = uuid.NewString()
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, contextkey.RequestID, requestId)
	ctx = context.WithValue(ctx, contextkey.RequestPath, gc.FullPath())

	gc.Set(contextkey.AppCtx.String(), ctx)

	return ctx
}
