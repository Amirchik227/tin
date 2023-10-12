package main

import (
	"encoding/json"
	"fmt"
	"log"
	"projectZero/database"
	"projectZero/handle"
	"projectZero/storage"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/stan.go"
)

func runRouter(handler *handle.Handler) {
	router := gin.Default()
	router.POST("/order", handler.CreateOrderFromJSON)
	router.GET("/order/:id", handler.GetOrder)
	router.Run()
}

func main() {
	database.ConnectDatabase()
	defer database.Close_db()

	memoryStorage := storage.NewMemoryStorage()
	database.GetCache(memoryStorage)
	handler := handle.NewHandler(memoryStorage)

	go runRouter(handler)

	opts := []stan.Option{
		stan.NatsURL("nats://localhost:4222"),
		stan.ConnectWait(20 * time.Second),
	}
	sc, err := stan.Connect("test-cluster", "subscriber-client", opts...)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer sc.Close()

	sub, err := sc.Subscribe("my-stream", func(m *stan.Msg) {
		fmt.Printf("Received a NATS message\n")
		var order storage.Order
		json.Unmarshal(m.Data, &order)
		handler.CreateOrder(order)
	})
	if err != nil {
		log.Fatalf("Error subscribing: %v", err)
	}
	defer sub.Close()

	var input string
	fmt.Scanln(&input)
}
