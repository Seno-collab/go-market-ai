package container

import (
	menuitemapp "go-ai/internal/menu/application/menu_item"
	"go-ai/internal/menu/infrastructure/db"
	menuhttp "go-ai/internal/menu/transport/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type MenuItemModule struct {
	MenuItemHandler *menuhttp.MenuItemHandler
}

func InitMenuItemModule(pool *pgxpool.Pool, log zerolog.Logger) *MenuItemModule {
	repo := db.NewMenuRepo(pool)
	topicRepo := db.NewTopicRepo(pool)
	optionGroupRepo := db.NewOptionGroupRepo(pool)
	optionItemRepo := db.NewOptionItemRepo(pool)

	createUseCase := menuitemapp.NewCreateUseCase(repo, topicRepo)
	getUseCase := menuitemapp.NewGetUseCase(repo)
	getDetailUseCase := menuitemapp.NewDetailUseCase(repo, optionGroupRepo, optionItemRepo)
	updateUseCase := menuitemapp.NewUpdateUseCase(repo)
	deleteUseCase := menuitemapp.NewDeleteUseCase(repo)
	getMenuItems := menuitemapp.NewGetMenuItemsUseCase(repo)
	updateStatusUseCase := menuitemapp.NewUpdateStatusUseCase(repo)
	menuItemHandler := menuhttp.NewMenuItemHandler(
		createUseCase,
		getUseCase,
		getDetailUseCase,
		updateUseCase,
		deleteUseCase,
		getMenuItems,
		updateStatusUseCase,
		log)
	return &MenuItemModule{
		MenuItemHandler: menuItemHandler,
	}
}
