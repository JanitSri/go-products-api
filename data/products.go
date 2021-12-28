package data

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rating struct {
	Rate  float64 `bson:"rate,omitempty" json:"rate,omitempty"`
	Count uint64  `bson:"count,omitempty" json:"count,omitempty"`
}

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ProductId   uint32             `bson:"id,omitempty" json:"productId,omitempty"`
	Title       string             `bson:"title,omitempty" json:"title,omitempty"`
	Price       float64            `bson:"price,omitempty" json:"price,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Category    string             `bson:"category,omitempty" json:"category,omitempty"`
	Image       string             `bson:"image,omitempty" json:"image,omitempty"`
	Ratings     Rating             `bson:"rating" json:"rating,omitempty"`
}

type Products []Product

func (p Product) toJson() string {
	b, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}

func GetAllProducts(d dataStore) Products {
	results := readData(d, bson.D{})
	for _, result := range results {
		fmt.Println(result.toJson())
	}
	return results
}

func GetProductByProductId(d dataStore, productId uint32) {
	filter := bson.D{{"id", productId}}
	results := readData(d, filter)

	if len(results) == 0 {
		notFound := `{"Error":"Resource Not Found"}`
		rawNotFound := json.RawMessage(notFound)
		bytes, _ := rawNotFound.MarshalJSON()
		fmt.Println(string(bytes))
		return
	}

	for _, result := range results {
		fmt.Println(result.toJson())
	}
}
