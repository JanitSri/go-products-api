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

func GetAllProducts(m mongoDB) Products {
	productsCollection := m.client.Database(m.database).Collection("Products")
	cursor, err := productsCollection.Find(m.ctx, bson.D{})
	defer cursor.Close(m.ctx)
	if err != nil {
		panic(err)
	}

	var results Products
	if err = cursor.All(m.ctx, &results); err != nil {
		panic(err)
	}

	for _, result := range results {
		fmt.Println(result.toJson())
	}

	return results
}

func GetProductByProductId(m mongoDB, productId uint32) {
	productsCollection := m.client.Database(m.database).Collection("Products")
	filter := bson.D{{"id", productId}}

	cursor, err := productsCollection.Find(m.ctx, filter)
	defer cursor.Close(m.ctx)
	if err != nil {
		panic(err)
	}
	var results Products
	if err = cursor.All(m.ctx, &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Println(result.toJson())
	}
}
