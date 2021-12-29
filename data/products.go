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

func GetAllProducts(d dataStore) Products {
	results := readData(d)
	for _, result := range results {
		fmt.Println(string(result.toJson()))
	}
	return results
}

func GetProductByProductId(d dataStore, productId uint32) {
	results := readDataById(d, int(productId))

	if len(results) == 0 {
		notFound := `{"Error":"Resource Not Found"}`
		rawNotFound := json.RawMessage(notFound)
		bytes, _ := rawNotFound.MarshalJSON()
		fmt.Println(string(bytes))
		return
	}

	for _, result := range results {
		fmt.Println(string(result.toJson()))
	}
}

func AddProduct(d dataStore, p Product) {
	result := insertData(d, p)
	fmt.Printf("Added Product with ID %s\n", result)
}

func DeleteProduct(d dataStore, productId uint32) {
	result := deleteData(d, int(productId))
	fmt.Printf("Number of Products Deleted: %d\n", result)
}

func UpdateProduct(d dataStore, p Product) {
	result := updateData(d, p)
	fmt.Println("Number of Products Updated", result)
}

func SearchProducts(d dataStore, searchTerm string) {
	results := searchData(d, searchTerm)

	for _, result := range results {
		fmt.Println(string(result.toJson()))
	}
}
