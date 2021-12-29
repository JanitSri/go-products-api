package main

import (
	"products-api/config"
	"products-api/data"
)

var mongoUsername = config.GetEnvVaraiable("MONGO_USERNAME")
var mongoPassword = config.GetEnvVaraiable("MONGO_PASSWORD")
var mongoDB = config.GetEnvVaraiable("MONGO_DATABASE")

func main() {
	mongo := data.NewMongoDB(mongoUsername, mongoPassword, mongoDB)
	data.InitializeDBConnection(mongo)
	defer data.CloseDBConnection(mongo)
	//data.GetAllProducts(mongo)
	//data.GetProductByProductId(mongo, 2)

	/*product := data.Product{
		Title:       "Test Item",
		Price:       99.99,
		Description: "Test Description",
		Category:    "Test Category",
		Image:       "www.testimage.com",
		Ratings: &data.Rating{
			Rate:  2.5,
			Count: 565,
		},
	}

	data.AddProduct(mongo, product)*/

	//data.DeleteProduct(mongo, 90)

	/*updateProduct := data.Product{
		ProductId: 90,
		Title:     "Test Item123",
		Category:  "Test Category123",
		Image:     "www.testimage123.com",
		Ratings: &data.Rating{
			Rate:  4.5,
			Count: 123,
		},
	}

	data.UpdateProduct(mongo, updateProduct)*/

	//data.SearchProducts(mongo, "jacket")
}
