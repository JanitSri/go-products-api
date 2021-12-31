package main

import (
	"log"
	"net/http"
	"os"
	"products-api/config"
	"products-api/data"
	"products-api/handlers"

	"github.com/gorilla/mux"
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

	l := log.New(os.Stdout, "products-api", log.LstdFlags)
	p := handlers.NewProducts(l, mongo)

	r := mux.NewRouter()

	getRouter := r.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", p.GetProductHandler)
	getRouter.HandleFunc("/products/{id:[0-9]+}", p.GetProductByIdHandler)

	postRouter := r.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", p.AddProductHandler)
	postRouter.Use(p.MiddlewareValidateProduct)

	l.Println("Starting server on port 8000")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		l.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
