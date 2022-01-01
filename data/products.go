package data

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Rating defines the structure for an rating of a Product
// swagger:model
type Rating struct {
	// the rating for the product
	//
	// required: false
	Rate float64 `bson:"rate,omitempty" json:"rate,omitempty"`

	// the number of total ratings for the product
	//
	// required: false
	// min: 1
	Count uint64 `bson:"count,omitempty" json:"count,omitempty"`
}

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the mongo id for the product - generated automaticlly
	//
	// required: false
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`

	// the id for the product
	//
	// required: false
	// min: 1
	ProductId uint32 `bson:"id,omitempty" json:"id,omitempty"`

	// the name of the product
	//
	// required: false
	// max length: 255
	Title string `bson:"title,omitempty" json:"title,omitempty"`

	// the price of the product
	//
	// required: false
	// min: 0.01
	Price float64 `bson:"price,omitempty" json:"price,omitempty"`

	// the description of the product
	//
	// required: false
	// max length: 500
	Description string `bson:"description,omitempty" json:"description,omitempty"`

	// the category that the product belongs to
	//
	// required: false
	// max length: 50
	Category string `bson:"category,omitempty" json:"category,omitempty"`

	// the image url of the product
	//
	// required: false
	Image string `bson:"image,omitempty" json:"image,omitempty"`

	// the rating for the product
	//
	// required: false
	Ratings *Rating `bson:"rating,omitempty" json:"rating,omitempty"`
}

type Products []Product

func (p Product) toJson() []byte {
	b, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return b
}

func (p Product) getProductId() uint32 {
	return p.ProductId
}

func GetAllProducts(d DataStore) Products {
	results := readData(d)

	if results == nil {
		return nil
	}

	return results
}

func GetProductByProductId(d DataStore, productId uint32) Products {
	results := readDataById(d, int(productId))

	if len(results) == 0 {
		return nil
	}

	return results
}

func AddProduct(d DataStore, p Product) string {
	result := insertData(d, p)
	return result
}

func DeleteProduct(d DataStore, productId uint32) int {
	result := deleteData(d, int(productId))
	return result
}

func UpdateProduct(d DataStore, p Product) int {
	result := updateData(d, p)
	return result
}

func SearchProducts(d DataStore, searchTerm string) Products {
	results := searchData(d, searchTerm)
	return results
}
