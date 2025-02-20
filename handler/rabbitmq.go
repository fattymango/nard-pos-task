package handler

import (
	"encoding/json"
	"fmt"
	"multitenant/model"

	"github.com/streadway/amqp"
)

const (
	AMQP_TRANSACTION_QUEUE = "transaction_queue"
)

func (s *MultiTanentServer) ConsumeTransactions() error {
	s.logger.Debug("ConsumeTransaction")

	ch, err := s.amqp.Connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %s", err)
	}

	q, err := ch.QueueDeclare(
		AMQP_TRANSACTION_QUEUE, // name
		false,                  // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %s", err)
	}

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Debug("closing transaction consumer")
			return nil
		case d := <-msgs:
			// s.logger.Debugf("received transaction: %s", d.Body)
			if d.Body == nil {
				continue
			}
			go s.processTransactionMessage(d)
		}
	}

}

// Process each transaction message independently to avoid blocking.
func (s *MultiTanentServer) processTransactionMessage(d amqp.Delivery) {
	// s.logger.Debugf("received transaction: %s", d.Body)
	t := &model.Transaction{}
	err := json.Unmarshal(d.Body, t)
	if err != nil {
		s.logger.Errorf("failed to unmarshal transaction: %s", err)
		return
	}

	if err := s.validator.Struct(t); err != nil {
		s.logger.Errorf("error validating transaction: %s", err)
		return
	}

	s.engine.CreateTransaction(t)
}
