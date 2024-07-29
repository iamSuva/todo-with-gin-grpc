package main

import (
	"fmt"
	"os"
	"todowithgin/database"
	grpclient "todowithgin/grpcClient"
	"todowithgin/routes"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("failed to load env")
	}
	postgresURL := os.Getenv("POSTGRES_URL")
	//connect to db
	fmt.Println(postgresURL)
	database.ConnectToDB(postgresURL)
	defer database.DB.Close()
	//grpc connection
	grpclient.GrpcConnection()

	// router setting

	router := routes.Router()

	router.Run("localhost:8080")
}
