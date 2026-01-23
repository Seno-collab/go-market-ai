package app

import (
	healthapp "go-ai/internal/health/application/health"
	"go-ai/internal/transport/response"
)

type HealthSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data healthapp.HealthResponse `json:"data"`
}

type HealthFailureResponseDoc struct {
	response.SuccessBaseDoc
	Data healthapp.HealthResponse `json:"data"`
}
