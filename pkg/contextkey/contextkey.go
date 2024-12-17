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
func GetFromFromCtx(ctx context.Context, key key) string {
	val := ctx.Value(key)
	if val == nil {
		return ""
	}

	strVal, ok := val.(string)
	if !ok {
		return ""
	}
	return strVal

}

// SetUserIDFromRequest ... sets input userID in context userID header
func SetInContext(
	ctx context.Context,
	key key,
	value string,
) context.Context {
	val := ctx.Value(key)
	if val == nil {
		ctx = context.WithValue(ctx, key, value)
	}
	return ctx
}
