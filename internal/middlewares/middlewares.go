package middlewares

import (
	"net/http"

	"github.com/danushk97/image-analyzer/internal/constants"
	internaErr "github.com/danushk97/image-analyzer/internal/errors"
	"github.com/danushk97/image-analyzer/pkg/contextkey"
	"github.com/danushk97/image-analyzer/pkg/errors"
	pkgLogger "github.com/danushk97/image-analyzer/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ErrorResponse(ctx *gin.Context, err error) {
	log := pkgLogger.Ctx(ctx)
	log.WithError(err).Error(err.Error())

	var problemDetail = gin.H{
		"code":        internaErr.ServerError,
		"description": "Something went wrong, please try again in some time.",
	}

	iErr, ok := err.(errors.IError)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	if iErr.IsOfType(errors.BAD_REQUEST_ERROR) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":        internaErr.BadRequesterror,
			"description": iErr.Error(),
		})

		return
	} else if iErr.IsOfType(errors.AUTHORIZATION_ERROR) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":        internaErr.Unauthorized,
			"description": iErr.Error(),
		})

		return
	} else {
		ctx.JSON(http.StatusInternalServerError, problemDetail)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(gc *gin.Context) {
		logger := pkgLogger.Ctx(gc.Request.Context())

		logger.Info("AUTH_VALIDATION")

		// Extract user_id from header
		userID := gc.GetHeader(constants.HeaderUserId)
		if userID == "" {
			logger.Warn("USER_ID_NOT_FOUND")
			ErrorResponse(
				gc,
				errors.NewAuthorizationError(internaErr.Unauthorized),
			)
			gc.Abort()
		}

		ctx := contextkey.SetInContext(
			gc.Request.Context(),
			contextkey.UserID,
			userID,
		)
		gc.Request = gc.Request.WithContext(ctx)

		logger.Info("AUTH_VALIDATION_SUCCESS")

		gc.Next()
	}
}

func CtxMiddleware() gin.HandlerFunc {
	return func(gc *gin.Context) {
		// Convert the Gin context to a Go context (using gin.Context's context method)
		ctx := gc.Request.Context()

		// Retrieve the request ID from gin context keys, if available
		var requestIDStr string
		if requestID, ok := gc.Keys[constants.HeaderRequestId].(string); ok && requestID != "" {
			requestIDStr = requestID
		} else {
			// Generate a new UUID if no request ID is found
			requestIDStr = uuid.NewString()
		}

		// Set the request ID in the context (using the proper key)
		ctx = contextkey.SetInContext(ctx, contextkey.RequestID, requestIDStr)

		// Set the request path in the context
		ctx = contextkey.SetInContext(ctx, contextkey.RequestPath, gc.FullPath())

		// Now, reattach the updated context to the Gin context
		gc.Request = gc.Request.WithContext(ctx)

		// Continue processing the request
		gc.Next()
	}
}
