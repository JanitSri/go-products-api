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
	data.GetProductByProductId(mongo, 2)
}
