package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/dulchik/learn-pub-sub-starter/internal/routing"
	"github.com/dulchik/learn-pub-sub-starter/cmd/internal/pubsub"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	rabbitChan, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not get the channel from RabbitMQ: %v", err)
	}

	err = pubsub.PublishJSON(rabbitChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
	if err != nil {
		log.Fatalf("could not publish json: %v", err)
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("RabbitMQ connection closed.")

}
