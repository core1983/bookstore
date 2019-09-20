package main

import (
	"catalog"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"log"
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

	repo, err := catalog.NewBookRepository(connectionString)
	defer repo.Close()

	service := catalog.NewBookService(repo)
	log.Printf("Catalog GRPC Server running on port %d", port)
	log.Fatal(catalog.NewGrpcServer(service, port))
}
