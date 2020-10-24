package DB

import (
	"context"
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var lock = &sync.Mutex{}

// GetTestDB -
func GetTestDB() (*DB_Connection_Struct, error) {
	lock.Lock()
	defer lock.Unlock()

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err)
	}

	var DB_Connection *DB_Connection_Struct
	if DatabaseConnection != nil {
		return DatabaseConnection, nil
	}

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
	log.Println("Started test DB")

	DB_Connection.Client = client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	DB_Connection.Ctx = &ctx

	DB_Connection.DB_Handle = client.Database("testing_db")
	DatabaseConnection = DB_Connection

	collections, err := DB_Connection.DB_Handle.ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	for _, c := range collections {
		if err := DB_Connection.DB_Handle.Collection(c).Drop(context.Background()); err != nil {
			return nil, err
		}
	}

	return DB_Connection, nil
}

// DisconnectTestDB -
func DisconnectTestDB() error {
	if DatabaseConnection == nil {
		return errors.New("DB is null")
	}

	DatabaseConnection.Client.Disconnect(context.Background())
	return nil
}
