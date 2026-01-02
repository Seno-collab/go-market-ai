package container

import (
	topicapp "go-ai/internal/menu/application/topic"
	"go-ai/internal/menu/infrastructure/db"
	menuhttp "go-ai/internal/menu/transport/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type TopicModule struct {
	TopicHandler *menuhttp.TopicHandler
}

func InitTopicModule(pool *pgxpool.Pool, log zerolog.Logger) *TopicModule {
	repo := db.NewTopicRepo(pool)
	createUseCase := topicapp.NewCreateUseCase(repo)
	getUseCase := topicapp.NewGetUseCase(repo)
	getTopicsUseCase := topicapp.NewGetTopicsUseCase(repo)
	updateUseCase := topicapp.NewUpdateUseCase(repo)
	deleteUseCase := topicapp.NewDeleteUseCase(repo)
	getComboboxUseCase := topicapp.NewComboboxUseCase(repo)
	topicHandler := menuhttp.NewTopicHandler(createUseCase, getUseCase, getTopicsUseCase, updateUseCase, deleteUseCase, getComboboxUseCase, log)
	return &TopicModule{
		TopicHandler: topicHandler,
	}
}
