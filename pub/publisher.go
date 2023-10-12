package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"projectZero/storage"
	"time"

	"github.com/nats-io/stan.go"
)

func generateRandomText(characters string, length int) string {
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		result[i] = characters[rand.Intn(len(characters))]
	}

	return string(result)
}

func generateOrder() (storage.Order, error) {
	var order storage.Order

	file, err := os.Open("../model.json")
	if err != nil {
		return order, errors.New("MODEL NOT FOUND")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&order); err != nil {
		return order, errors.New("DECODE ERROR")
	}

	lowerSymbols := "abcdefghijklmnopqrstuvwxyz0123456789"
	upperLetters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	order.OrderUID = generateRandomText(lowerSymbols, 19)
	order.Entry = generateRandomText(upperLetters, 4)
	order.TrackNumber = generateRandomText(upperLetters, 16)
	order.Payment.Transaction = generateRandomText(lowerSymbols, 19)
	for i := 0; i < len(order.Items); i++ {
		order.Items[i].ChrtID = rand.Intn(9000000) + 1000000
	}
	return order, nil
}

func main() {
	opts := []stan.Option{
		stan.NatsURL("nats://localhost:4222"),
		stan.ConnectWait(20 * time.Second),
	}
	sc, err := stan.Connect("test-cluster", "publisher-client", opts...)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer sc.Close()

	for i := 0; i < 5; i++ {
		order, err := generateOrder()
		if err != nil {
			log.Fatalf("Failed to generate order: %v", err)
		}

		data, _ := json.Marshal(order)
		if err := sc.Publish("my-stream", data); err != nil {
			log.Fatalf("Error sending the message: %v", err)
		}

		fmt.Printf("Message published\n")
		time.Sleep(2 * time.Second)
	}
}
