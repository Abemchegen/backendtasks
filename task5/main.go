package main

import (
	"context"
	"log"
	"os"
	"task5/controllers"
	"task5/data"
	"task5/router"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	uri := os.Getenv("MONGO_URL")

	if uri == "" {
		log.Fatal("MONGO_URI environment variable not set")
	}

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("task_manager")
	taskservice := data.NewTaskService(db)
	taskcontroller := controllers.NewTaskController(taskservice)

	router := router.SetupRouter(taskcontroller)
	e := router.Run(":8080")

	if e != nil {
		log.Fatalf("Error running server: %v", e)
	}

}
