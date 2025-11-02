package database

import (
	"context"
	"time"
)

func mongoCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
