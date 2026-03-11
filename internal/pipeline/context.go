package pipeline

import (
	"context"
)

func checkContext(ctx context.Context, path string) error {
	if ctx == nil {
		return nil
	}
	if ctx.Err() != nil {
		return NewExecutionFailure("execution.canceled", path, "command canceled", ErrExecutionCanceled)
	}
	return nil
}
