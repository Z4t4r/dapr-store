// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// SQLite implementation of the ProductService
// ----------------------------------------------------------------------------

package impl

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/benc-uk/dapr-store/cmd/products/spec"
	"github.com/benc-uk/dapr-store/pkg/env"
	"github.com/benc-uk/dapr-store/pkg/problem"
	"github.com/dapr/go-sdk/service/common"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// ProductService is a Dapr based implementation of ProductService interface
type ProductService struct {
	*gorm.DB
	pubSubName string
	topicName  string
	storeName  string
	serviceName string
}

// NewService creates a new ProductService
func NewService(serviceName string, dsn string) (out *ProductService) {

	topicName := env.GetEnvString("DAPR_ORDERS_TOPIC", "orders-queue")
	storeName := env.GetEnvString("DAPR_STORE_NAME", "statestore")
	pubSubName := env.GetEnvString("DAPR_PUBSUB_NAME", "pubsub")


	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&spec.Product{})
	var products []spec.Product
	result := db.Find(&products)
	if result.RowsAffected <=0 {
		db.Create(&spec.Product{
			Name: "Product1",
			Cost: 1,
			Description: "Description1",
			OnOffer: true,
		})
		db.Create(&spec.Product{
			Name: "Product2",
			Cost: 2,
			Description: "Description2",
			OnOffer: true,
		})

	}


	if err != nil {
		log.Panicf("### Failed to open database %s %+v\n", dsn, err)
		return nil
	}

	log.Printf("### Database %s opened OK\n", dsn)

	out = &ProductService{
		db,
		pubSubName,
		topicName,
		storeName,
		serviceName,
	}
	return
}

// QueryProducts is a simple SQL WHERE query on a single column
func (s ProductService) QueryProducts(column, term string) ([]spec.Product, error) {

	var products []spec.Product
	result := s.First(&products,"name LIKE ?",term)
	//s.DB.Where()
	if result.Error != nil {
		prob := problem.New("err://products-db", "Database query error", 500, result.Error.Error(), s.serviceName)
		return nil, prob
	}

	return products,result.Error
}

// AllProducts returns all products from the DB, yeah this is pretty dumb
func (s ProductService) AllProducts(ctx context.Context,in *common.InvocationEvent)  (out *common.Content, err error) {
	var products []spec.Product
	result := s.First(&products)
	if result.Error != nil {
		prob := problem.New("err://products-db", "Database query error", 500, result.Error.Error(), s.serviceName)
		return nil, prob
	}
	js,err := json.Marshal(products)
	out = &common.Content{
		Data: js,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return
}

// SearchProducts is a text search in name or  product description
func (s ProductService) SearchProducts(query string) ([]spec.Product, error) {
	//rows, err := s.Query("SELECT * FROM products WHERE (description LIKE ? OR name LIKE ?)", "%"+query+"%", "%"+query+"%")
	var products []spec.Product
	result := s.Where("name LIKE ?",query).Or("description LIKE ?", query).Find(&products)

	if result.Error != nil {
		prob := problem.New("err://products-db", "Database query error", 500, result.Error.Error(), s.serviceName)
		return nil, prob
	}

	return products,result.Error
}

// Helper function to take a bunch of rows and return as a slice of Products
func (s ProductService) processRows(rows *sql.Rows) ([]spec.Product, error) {
	products := []spec.Product{}
	defer rows.Close()
	for rows.Next() {
		p := spec.Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Cost, &p.Image, &p.OnOffer)
		if err != nil {
			prob := problem.New("err://products-db", "Error reading row", 500, err.Error(), s.serviceName)
			return nil, prob
		}
		products = append(products, p)
	}

	return products, nil
}
