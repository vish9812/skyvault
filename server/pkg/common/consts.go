package common

import "context"

type CtxKey string

const (
	CtxKeyLogger    CtxKey = "logger"
	CtxKeyProfileID CtxKey = "profile_id"
)

func GetProfileIDFromContext(ctx context.Context) int64 {
	return ctx.Value(CtxKeyProfileID).(int64)
}
