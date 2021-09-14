package impl

import (
	"context"
	"github.com/benc-uk/dapr-store/cmd/products/spec"
	dapr "github.com/dapr/go-sdk/client"
)

type FrontendService struct {
	dapr.Client
	ctx context.Context
	serviceName string
}


func NewService(serviceName string,ctx context.Context, dsn string) *FrontendService {
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	return &FrontendService{
		client,
		ctx,
		serviceName,
	}
}

func(f FrontendService) AllProducts() ([]spec.Product,error){
	result = f.Client.InvokeMethod(f.ctx,)
}