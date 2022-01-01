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

// List of products returned in the response
// swagger:response productsResponse
type ProductsResponse struct {
	// All the products available
	// in: body
	Body []data.Product
}

// swagger:parameters listProduct deleteProduct
type ProductIDParam struct {
	// The ID of the product
	// in: path
	// required: true
	ID int `json:"id"`
}

// swagger:parameters searchProduct
type SearchKeyQueryParam struct {
	// The search term to search products
	// in: query
	SearchKey string `json:"searchKey"`
}

type SuccessfulResult struct {
	Result string `json:"result"`
}

// Generic successful response returned
// swagger:response genericResult
type GenericResponse struct {
	// Result of action
	// in: body
	Body SuccessfulResult
}

// swagger:route GET /products products listProducts
//
// Lists all products
//
// This will show all available products by default.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Responses:
//       200: productsResponse
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

// swagger:route GET /products/{id} products listProduct
//
// List product by product ID
//
// This will show no products by default.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Responses:
//       200: productsResponse
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

// swagger:route POST /products products addProduct
//
// Add a product
//
// This will add a product to the system.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Responses:
//       200: genericResult
func (p *Products) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	p.l.Println("POST - Add Product")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	result := data.AddProduct(p.db, prod)
	res := fmt.Sprintf(`{"Result":"Product Added With Mongo ID: %s"}`, result)
	rawNotFound := json.RawMessage(res)
	bytes, _ := rawNotFound.MarshalJSON()

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

// swagger:route DELETE /products/{id} products deleteProduct
//
// Delete product by product ID
//
// This will delete a product by product ID.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Responses:
//       200: genericResult
func (p *Products) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	p.l.Println("DELETE - Delete Product")

	id := getId(w, r)
	result := data.DeleteProduct(p.db, uint32(id))

	res := fmt.Sprintf(`{"Result":"Number of Products Deleted: %d"}`, result)
	rawNotFound := json.RawMessage(res)
	bytes, _ := rawNotFound.MarshalJSON()

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

// swagger:route PUT /products/ products updateProduct
//
// Update product
//
// This will update a product in the system.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Responses:
//       200: genericResult
func (p *Products) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	p.l.Println("PUT - Update Product")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	result := data.UpdateProduct(p.db, prod)
	res := fmt.Sprintf(`{"Result":"Number of Product Updated: %d"}`, result)
	rawNotFound := json.RawMessage(res)
	bytes, _ := rawNotFound.MarshalJSON()

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

// swagger:route GET /products/search?searchKey={searchKey} products searchProduct
//
// Search products
//
// This will allow to search the products in the system.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Responses:
//       200: productsResponse
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
