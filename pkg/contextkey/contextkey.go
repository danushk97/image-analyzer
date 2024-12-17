package contextkey

import (
	"context"
)

type key string

// String converts key
func (c key) String() string {
	return string(c)
}

var (
	// UserID is the authenticated user making the request
	UserID      = key("userID")
	RequestID   = key("requestID")
	RequestPath = key("requestPath")
	AppCtx      = key("appCtx")
)

// GetUserIDFromRequest ...
func GetUserIDFromCtx(ctx context.Context) string {
	userVal := ctx.Value(UserID)
	if userVal == nil {
		return ""
	}

	strVal, ok := userVal.(string)
	if !ok {
		return ""
	}
	return strVal

}

// SetUserIDFromRequest ... sets input userID in context userID header
func SetUserIDFromRequest(ctx context.Context, userID string) context.Context {
	userVal := ctx.Value(UserID)
	if userVal == nil {
		ctx = context.WithValue(ctx, UserID, userID)
	}
	return ctx
}
