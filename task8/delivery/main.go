package main

import (
	"context"
	"log"
	"os"
	"task8/delivery/controllers"
	"task8/delivery/routers"
	"task8/infrastructure"
	"task8/repositories"
	"task8/usecases"

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

	ps := infrastructure.NewPasswordService()
	js := infrastructure.NewJWTService()
	usererpository := repositories.NewUserRepository(db, ps)
	userusecase := usecases.NewUserUsecase(usererpository, js)
	usercontroller := controllers.NewUserController(userusecase)

	router := routers.SetRouter(taskcontroller, usercontroller)
	router.Run(":8080")
}
