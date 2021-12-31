package data

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoDB struct {
	username string
	password string
	database string
	client   *mongo.Client
	ctx      context.Context
}

func NewMongoDB(username string, password string, database string) *mongoDB {
	var client *mongo.Client
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return &mongoDB{
		username,
		password,
		database,
		client,
		ctx,
	}
}

func (m *mongoDB) connect() {
	connURI := fmt.Sprintf("mongodb+srv://%s:%s@go-cluster.1uijt.mongodb.net/%s?retryWrit    es=true&w=majority", m.username, m.password, m.database)

	client, err := mongo.NewClient(options.Client().ApplyURI(connURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	m.client = client
}

func (m mongoDB) disconnect() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err := m.client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func (m mongoDB) ping() {
	if err := m.client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")
}

func (m mongoDB) read(productId int) Products {
	filter := bson.D{{"id", productId}}
	if productId == -1 {
		filter = bson.D{{}}
	}

	return find(m, filter)
}

func (m mongoDB) search(searchTerm string) Products {
	filter := bson.D{{"$text", bson.D{{"$search", searchTerm}}}}

	return find(m, filter)
}

func find(m mongoDB, filter bson.D) Products {
	productsCollection := m.client.Database(m.database).Collection("Products")

	cursor, err := productsCollection.Find(context.TODO(), filter)
	defer cursor.Close(m.ctx)
	if err != nil {
		fmt.Println("Find Error: ", err)
		return nil
	}

	var results Products
	if err = cursor.All(m.ctx, &results); err != nil {
		fmt.Println("Find Error: ", err)
		return nil
	}

	return results
}

func (m mongoDB) create(p Product) string {
	productsCollection := m.client.Database(m.database).Collection("Products")
	p.ProductId = getNewProductId(productsCollection, context.TODO())
	result, err := productsCollection.InsertOne(context.TODO(), p)

	if err != nil {
		panic(err)
	}

	return result.InsertedID.(primitive.ObjectID).Hex()
}

func getNewProductId(coll *mongo.Collection, ctx context.Context) uint32 {

	filter := bson.D{}
	sort := bson.D{{"id", -1}}
	opts := options.FindOne().SetSort(sort)

	var p Product
	err := coll.FindOne(ctx, filter, opts).Decode(&p)

	if err != nil {
		panic(err)
	}

	return p.ProductId + 1
}
func (m mongoDB) delete(productId int) int {
	productsCollection := m.client.Database(m.database).Collection("Products")
	filter := bson.D{{"id", productId}}
	result, err := productsCollection.DeleteOne(m.ctx, filter)

	if err != nil {
		panic(err)
	}

	return int(result.DeletedCount)
}

func (m mongoDB) update(p Product) int {
	productsCollection := m.client.Database(m.database).Collection("Products")

	id := p.getProductId()
	interfaceMap := productStructToMap(p)
	update := buildUpdate(interfaceMap)

	filter := bson.D{{"id", id}}
	result, err := productsCollection.UpdateOne(m.ctx, filter, update)

	if err != nil {
		panic(err)
	}

	return int(result.ModifiedCount)
}

func productStructToMap(p Product) map[string]interface{} {
	var productStructMap map[string]interface{}
	product_json := p.toJson()
	json.Unmarshal(product_json, &productStructMap)
	return productStructMap
}

func buildUpdate(m map[string]interface{}) bson.M {
	update := buildUpdateHelper(m, "")
	return bson.M{"$set": update}
}

func buildUpdateHelper(m map[string]interface{}, prefix string) bson.M {
	update := bson.M{}
	for key, value := range m {
		switch valueTypeConv := value.(type) {
		default:
			//fmt.Printf("%v, %v, %T\n", key, value, valueTypeConv)

			if key == "_id" {
				continue
			}

			if reflect.ValueOf(valueTypeConv).Kind() == reflect.Map {
				innerM := valueTypeConv.(map[string]interface{})
				innerUpdate := buildUpdateHelper(innerM, fmt.Sprintf("%s.", key))

				// flatten the map
				for innerKey, innerVal := range innerUpdate {
					update[key+"."+innerKey] = innerVal
				}
			} else {
				update[key] = valueTypeConv
			}
		}
	}
	return update
}
