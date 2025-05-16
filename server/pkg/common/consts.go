package common

import "context"

type CtxKey string

const (
	CtxKeyLogger    CtxKey = "logger"
	CtxKeyProfileID CtxKey = "profile_id"
)

func GetProfileIDFromContext(ctx context.Context) string {
	return ctx.Value(CtxKeyProfileID).(string)
}
