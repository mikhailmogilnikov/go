package ledger

import (
	"context"

	"github.com/mikhailmogilnikov/go/hw8/ledger/internal/app"
	"github.com/mikhailmogilnikov/go/hw8/ledger/internal/service"
)

// Service определяет интерфейс приложения для Gateway
type Service = service.Service

// InitService инициализирует сервис и возвращает его вместе с функцией закрытия ресурсов
func InitService(ctx context.Context) (Service, func() error, error) {
	return app.InitService(ctx)
}

