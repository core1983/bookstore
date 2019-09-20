package main

import (
	"log"
	"order"
)

func main() {

	var r order.Repository
	r, _ = order.NewOrderRepository(
		"host = localhost port = 5432 user = paul password = kredka12 dbname = order_db sslmode = disable")

	defer r.Close()

	log.Println("Listening on port 7003...")
	s := order.NewService(r)
	log.Fatal(order.ListenGRPC(s, ":7001", ":7002", 7003))
}
