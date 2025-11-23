package ledger

import (
	"context"

	"github.com/mikhailmogilnikov/go/hw9/ledger/internal/app"
	"github.com/mikhailmogilnikov/go/hw9/ledger/internal/service"
)

type Service = service.Service

func InitService(ctx context.Context) (Service, func() error, error) {
	return app.InitService(ctx)
}
