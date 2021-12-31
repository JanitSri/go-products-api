package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"products-api/data"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l  *log.Logger
	db data.DataStore
}

func NewProducts(l *log.Logger, db data.DataStore) *Products {
	return &Products{l, db}
}

func (p *Products) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	p.l.Println("GET - Products")

	lp := data.GetAllProducts(p.db)
	d, err := json.Marshal(lp)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}

	w.Write(d)
}

func (p *Products) GetProductByIdHandler(w http.ResponseWriter, r *http.Request) {
	p.l.Println("GET - Product By ID")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 0 {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
		return
	}

	lp := data.GetProductByProductId(p.db, uint32(id))
	d, err := json.Marshal(lp)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}

	w.Write(d)
}

func (p *Products) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	p.l.Println("POST - Add Product")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	result := data.AddProduct(p.db, prod)
	res := fmt.Sprintf(`{"Result":"%s"}`, result)
	rawNotFound := json.RawMessage(res)
	bytes, _ := rawNotFound.MarshalJSON()
	w.Write(bytes)
}

type KeyProduct struct{}

func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		var prod data.Product
		err := json.NewDecoder(r.Body).Decode(&prod)

		fmt.Println("adding", p)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
