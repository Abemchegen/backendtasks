package main

import (
	"context"
	"log"
	"os"
	"task7/delivery/controllers"
	"task7/delivery/routers"
	"task7/repositories"
	"task7/usecases"

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
	taskrepository := repositories.NewTaskRepository(db)
	taskusecase := usecases.NewTaskUsecase(taskrepository)
	taskcontroller := controllers.NewTaskController(taskusecase)

	usererpository := repositories.NewUserRepository(db)
	userusecase := usecases.NewUserUsecase(usererpository)
	usercontroller := controllers.NewUserController(userusecase)

	router := routers.SetRouter(taskcontroller, usercontroller)
	router.Run(":8080")
}
