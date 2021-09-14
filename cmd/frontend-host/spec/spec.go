package spec

import "github.com/benc-uk/dapr-store/cmd/products/spec"

type Frontend struct {
	Products map[string]int `json:"products"`
}

type FrontendService interface {
	AllProducts() ([]spec.Product, error)
}