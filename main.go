package main

import (
	"./routes"

	//"context"
	//"fmt"
	"./DB/aws"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"./DB"
	//"./aws"

	"log"
	"net/http"
	"os"
	//"time"
)

var (
	PORT = "PORT"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

}
func main() {
	var BaseRouter = mux.NewRouter()
	/*
		err := aws.SetupAWS()
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	*/

	routes.CreateRoutes(BaseRouter)

	db, err := DB.SetupDB()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Client.Disconnect(*db.Ctx)

	_, err := aws.SetAWSConnection()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = http.ListenAndServe(os.Getenv(PORT), BaseRouter)
	if err != nil {
		log.Fatal(err.Error())
	}

}
