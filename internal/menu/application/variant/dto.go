package variantapp

import "go-ai/pkg/utils"

type CreateVariantRequest struct {
	Name       string      `json:"name"`
	PriceDelta utils.Money `json:"price_delta"`
}
