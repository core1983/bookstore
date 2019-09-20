package main

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"log"
	"user"
)

func main() {
	config.AddDriver(yaml.Driver)
	err := config.LoadFiles("config/config.yml")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	port := config.Int("server.port")
	connectionString := config.String("database.connectionstring")

	repo, err := user.NewUserRepository(connectionString)
	defer repo.Close()

	service := user.NewUserService(repo)
	log.Printf("Server GRPC Running on port %d", port)
	log.Fatal(user.NewGrpcServer(service, port))
}
