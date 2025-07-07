package utils

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey string

const TraceCtxKey ctxKey = "trace_id"

func GetTraceID(ctx context.Context) uuid.UUID {
	traceID, ok := ctx.Value(TraceCtxKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return traceID
}
