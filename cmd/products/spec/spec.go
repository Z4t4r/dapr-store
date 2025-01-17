// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Specification of the Product entity and service
// ----------------------------------------------------------------------------

package spec

import (
	"context"
	"github.com/dapr/go-sdk/service/common"
	"gorm.io/gorm"
)

// Product holds product data
type Product struct {
	gorm.Model
	Name        string  `json:"name"`
	Cost        float32 `json:"cost"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	OnOffer     bool    `json:"onOffer"`
}

// ProductService defines core CRUD methods a products service should have
type ProductService interface {
	SearchProducts(string) ([]Product, error)
	QueryProducts(string, string) ([]Product, error)
	AllProducts(ctx context.Context,in *common.InvocationEvent)  (out *common.Content, err error)
}
