package main

import (
	"context"
	"log"
	"os"
	"task6/controllers"
	"task6/data"
	"task6/routers"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env file")
	}

	uri := os.Getenv("MONGO_URL")
	if uri == "" {
		log.Fatal("MONGO_URI not set")
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

	db := client.Database("taskmanager")
	taskservice := data.NewTaskService(db)
	taskcontroller := controllers.NewTaskController(taskservice)

	userservice := data.NewUserService(db)
	usercontroller := controllers.NewUserController(*userservice)

	router := routers.SetRouter(taskcontroller, usercontroller)
	router.Run(":8080")
}
