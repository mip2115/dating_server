package DB

import (
	"context"
	"errors"

	//	"fmt"
	// "go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"log"
	"os"
	"time"
)

type DB_Connection_Struct struct {
	DB_Handle *mongo.Database
	Ctx       *context.Context
	Client    *mongo.Client
}

func init() {

}

var DB_Connection *DB_Connection_Struct

// var DB_Connection *mongo.Database

func SetupDB() (*DB_Connection_Struct, error) {
	DB_Connection = &DB_Connection_Struct{}

	uri := os.Getenv("MONGO_URI_LOCAL")
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.Println("Started DB")

	DB_Connection.Client = client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	DB_Connection.Ctx = &ctx

	DB_Connection.DB_Handle = client.Database("kama")
	return DB_Connection, nil

}

// GetDataBase returns the DB
func GetDatabase() (*DB_Connection_Struct, error) {
	if DB_Connection == nil {
		return nil, errors.New("DB is null")
	}
	return DB_Connection, nil
}

func GetDatabaseHandle() (*mongo.Database, error) {
	if DB_Connection == nil {
		return nil, errors.New("DB is null")
	}
	if DB_Connection.DB_Handle == nil {
		return nil, errors.New("Handle is null")
	}

	return DB_Connection.DB_Handle, nil
}

func (DB_Connection *DB_Connection_Struct) TearDownCollection(collection string) error {

	return nil
}

func GetCollection(c string) (*mongo.Collection, error) {
	DB_Handle, err := GetDatabaseHandle()
	if err != nil {
		return nil, err
	}
	return DB_Handle.Collection(c), nil

}

func Disconnect() error {
	if DB_Connection == nil {
		return errors.New("DB is null")
	}

	DB_Connection.Client.Disconnect(context.Background())
	return nil
}
