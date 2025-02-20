package main

import (
	"fmt"
	"multitenant/handler"
	"multitenant/pkg/config"
	"multitenant/pkg/rabbitmq"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %s", err))
	}

	a, err := rabbitmq.NewRabbitMQ(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to create RabbitMQ connection: %s", err))
	}

	fmt.Println("RabbitMQ connection successful")

	ch, err := a.Connection.Channel()
	if err != nil {
		panic(fmt.Errorf("failed to open a channel: %s", err))
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		handler.AMQP_TRANSACTION_QUEUE, // name
		false,                          // durable
		false,                          // delete when unused
		false,                          // exclusive
		false,                          // no-wait
		nil,                            // arguments
	)
	if err != nil {
		panic(fmt.Errorf("failed to declare a queue: %s", err))
	}

	payload := []byte(`{
    "tenant_id": 1,
    "branch_id": 2,
    "product_id": 2,
    "quantity_sold": 4,
    "price_per_unit":5.5
}`)
	for {
		for i := 0; i < 100; i++ {
			err = ch.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        payload,
				})
			if err != nil {
				panic(fmt.Errorf("failed to publish a message: %s", err))
			}
			fmt.Printf("Message %d sent\n", i+1)
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		panic(fmt.Errorf("failed to publish a message: %s", err))
	}

	fmt.Println("Message sent")

}
