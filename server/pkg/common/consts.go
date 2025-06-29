package common

import "context"

const (
	BytesPerKB = 1 << 10
	BytesPerMB = 1 << 20
	BytesPerGB = 1 << 30
)

type CtxKey string

const (
	CtxKeyLogger    CtxKey = "logger"
	CtxKeyProfileID CtxKey = "profile_id"
)

func GetProfileIDFromContext(ctx context.Context) string {
	return ctx.Value(CtxKeyProfileID).(string)
}
