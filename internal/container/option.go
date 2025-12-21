package container

import (
	optiongroupapp "go-ai/internal/menu/application/option-group"
	optionitemapp "go-ai/internal/menu/application/option_item"
	"go-ai/internal/menu/infrastructure/db"
	menuhttp "go-ai/internal/menu/transport/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type OptionModule struct {
	OptionGroupHandler *menuhttp.OptionGroupHandler
	OptionItemHandler  *menuhttp.OptionItemHandler
}

func InitOptionModule(pool *pgxpool.Pool, log zerolog.Logger) *OptionModule {
	optionGroupRepo := db.NewOptionGroupRepo(pool)
	optionItemRepo := db.NewOptionItemRepo(pool)
	menuItemRepo := db.NewMenuRepo(pool)

	createGroup := optiongroupapp.NewCreateUseCase(optionGroupRepo, menuItemRepo)
	getGroup := optiongroupapp.NewGetUseCase(optionGroupRepo)
	getGroupsByMenuItem := optiongroupapp.NewGetByMenuItemUseCase(optionGroupRepo, menuItemRepo)
	updateGroup := optiongroupapp.NewUpdateUseCase(optionGroupRepo)
	deleteGroup := optiongroupapp.NewDeleteUseCase(optionGroupRepo)
	optionGroupHandler := menuhttp.NewOptionGroupHandler(createGroup, getGroup, getGroupsByMenuItem, updateGroup, deleteGroup, log)

	createItem := optionitemapp.NewCreateUseCase(optionItemRepo, optionGroupRepo)
	getItem := optionitemapp.NewGetUseCase(optionItemRepo)
	getItemsByGroup := optionitemapp.NewGetByGroupUseCase(optionItemRepo, optionGroupRepo)
	updateItem := optionitemapp.NewUpdateUseCase(optionItemRepo)
	deleteItem := optionitemapp.NewDeleteUseCase(optionItemRepo)
	optionItemHandler := menuhttp.NewOptionItemHandler(createItem, getItem, getItemsByGroup, updateItem, deleteItem, log)

	return &OptionModule{
		OptionGroupHandler: optionGroupHandler,
		OptionItemHandler:  optionItemHandler,
	}
}
