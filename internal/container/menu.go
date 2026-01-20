package container

import (
	menuapp "go-ai/internal/menu/application/menu"
	"go-ai/internal/menu/infrastructure/db"
	menuhttp "go-ai/internal/menu/transport/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type MenuModule struct {
	Handler *menuhttp.MenuHandler
}

func InitMenuModule(pool *pgxpool.Pool, log zerolog.Logger) *MenuModule {
	repo := db.NewMenuListRepo(pool)
	listUseCase := menuapp.NewListMenusUseCase(repo)
	handler := menuhttp.NewMenuHandler(listUseCase, log)
	return &MenuModule{Handler: handler}
}
