package data

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rating struct {
	Rate  float64 `bson:"rate,omitempty" json:"rate,omitempty"`
	Count uint64  `bson:"count,omitempty" json:"count,omitempty"`
}

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ProductId   uint32             `bson:"id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title,omitempty" json:"title,omitempty"`
	Price       float64            `bson:"price,omitempty" json:"price,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Category    string             `bson:"category,omitempty" json:"category,omitempty"`
	Image       string             `bson:"image,omitempty" json:"image,omitempty"`
	Ratings     *Rating            `bson:"rating,omitempty" json:"rating,omitempty"`
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

func SearchProducts(d DataStore, searchTerm string) {
	results := searchData(d, searchTerm)

	for _, result := range results {
		fmt.Println(string(result.toJson()))
	}
}
