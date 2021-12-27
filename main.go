package main

import (
	"fmt"
	"products-api/config"
	"products-api/data"
)

var mongoUsername = config.GetEnvVaraiable("MONGO_USERNAME")
var mongoPassword = config.GetEnvVaraiable("MONGO_PASSWORD")
var mongoDB = config.GetEnvVaraiable("MONGO_DATABASE")

func main() {
	fmt.Println("Products API")

	mongo := data.NewMongoDB(mongoUsername, mongoPassword, mongoDB)
	defer data.CloseDBConnection(mongo)
	data.InitializeDBConnection(mongo)
	data.HealthCheck(mongo)
}
