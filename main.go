package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"products-api/config"
	"products-api/data"
	"products-api/handlers"
	"time"

	"github.com/gorilla/mux"
)

var mongoUsername = config.GetEnvVaraiable("MONGO_USERNAME")
var mongoPassword = config.GetEnvVaraiable("MONGO_PASSWORD")
var mongoDB = config.GetEnvVaraiable("MONGO_DATABASE")

const bindAddr = ":8000"

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

	l := log.New(os.Stdout, "products-api: ", log.LstdFlags)
	p := handlers.NewProducts(l, mongo)

	r := mux.NewRouter()

	getRouter := r.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", p.GetProductHandler)
	getRouter.HandleFunc("/products/{id:[0-9]+}", p.GetProductByIdHandler)

	postRouter := r.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", p.AddProductHandler)
	postRouter.Use(p.MiddlewareValidateProduct)

	deleteRouter := r.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", p.DeleteProductHandler)

	putRouter := r.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products", p.UpdateProductHandler)
	putRouter.Use(p.MiddlewareValidateProduct)

	s := &http.Server{
		Addr:         bindAddr,
		Handler:      r,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		l.Println("Starting server on port 8000")
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Recieved terminate", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
