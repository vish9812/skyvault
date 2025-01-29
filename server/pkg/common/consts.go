package common

type CtxKey string

const (
	CtxKeyApp    CtxKey = "app"
	CtxKeyClaims CtxKey = "claims"
	CtxKeyLogger CtxKey = "logger"
)
