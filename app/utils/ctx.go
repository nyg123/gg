package utils

import "context"

var RequestIDKey = "x-request-id"

func GetRequestId(ctx context.Context) string {
	return ctx.Value(RequestIDKey).(string)
}
