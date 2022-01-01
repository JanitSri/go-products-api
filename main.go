// Package classification Products API.
//
// The products API will allow for reading, creating, updating,
// deleting, & searching for products. It was built using GO, gorilla/mux,
// and mongo-go-driver. Go-swagger was used for API documentation.
//
//     Schemes: http
//     Host: localhost
//     BasePath: /
//     Version: 1.0.0
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Janit Sri<janits_27@hotmail.com> https://janit.dev
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
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

	"github.com/go-openapi/runtime/middleware"
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

	l := log.New(os.Stdout, "products-api: ", log.LstdFlags)
	p := handlers.NewProducts(l, mongo)

	r := mux.NewRouter()

	getRouter := r.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", p.GetProductHandler)
	getRouter.HandleFunc("/products/{id:[0-9]+}", p.GetProductByIdHandler)
	getRouter.HandleFunc("/products/search", p.SearchProductHandler)

	opts := middleware.RedocOpts{SpecURL: "./swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

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
