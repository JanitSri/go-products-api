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
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(d)
}

func (p *Products) GetProductByIdHandler(w http.ResponseWriter, r *http.Request) {
	p.l.Println("GET - Product By ID")

	id := getId(w, r)

	lp := data.GetProductByProductId(p.db, uint32(id))
	d, err := json.Marshal(lp)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}

	if string(d) == "null" {
		http.Error(w, "Product not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(d)
}

func (p *Products) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	p.l.Println("POST - Add Product")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	result := data.AddProduct(p.db, prod)
	res := fmt.Sprintf(`{"Product Added With Mongo ID":"%s"}`, result)
	rawNotFound := json.RawMessage(res)
	bytes, _ := rawNotFound.MarshalJSON()

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func (p *Products) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	p.l.Println("DELETE - Delete Product")

	id := getId(w, r)
	result := data.DeleteProduct(p.db, uint32(id))

	res := fmt.Sprintf(`{"Number of Products Deleted":"%d"}`, result)
	rawNotFound := json.RawMessage(res)
	bytes, _ := rawNotFound.MarshalJSON()

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func (p *Products) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	p.l.Println("PUT - Update Product")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	result := data.UpdateProduct(p.db, prod)
	res := fmt.Sprintf(`{"Number of Product Updated":"%d"}`, result)
	rawNotFound := json.RawMessage(res)
	bytes, _ := rawNotFound.MarshalJSON()

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func (p *Products) SearchProductHandler(w http.ResponseWriter, r *http.Request) {
	p.l.Println("GET - SEARCH PRODUCTS")

	searchKey := r.FormValue("searchKey")

	sp := data.SearchProducts(p.db, searchKey)
	ps, err := json.Marshal(sp)

	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}

	if string(ps) == "null" {
		res := fmt.Sprintf(`{"Result":"%s"}`, string(ps))
		rawNotFound := json.RawMessage(res)
		ps, _ = rawNotFound.MarshalJSON()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(ps)
}

type KeyProduct struct{}

func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		var prod data.Product
		err := json.NewDecoder(r.Body).Decode(&prod)

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

func getId(w http.ResponseWriter, r *http.Request) int {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 0 {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
		return -1
	}
	return id
}
