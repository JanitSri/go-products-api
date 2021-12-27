package data

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
	err = client.Connect(m.ctx)
	if err != nil {
		log.Fatal(err)
	}

	m.client = client

	databases, err := m.client.ListDatabaseNames(m.ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
}

func (m mongoDB) disconnect() {
	if err := m.client.Disconnect(m.ctx); err != nil {
		panic(err)
	}
}

func (m mongoDB) ping() {
	if err := m.client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")
}
