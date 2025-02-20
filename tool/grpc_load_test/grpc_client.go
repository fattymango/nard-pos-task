package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"multitenant/pkg/config"
	pb "multitenant/proto/multitenant"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	// Establish gRPC connection
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", config.GRPC.Host, config.GRPC.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	defer conn.Close()

	c := pb.NewMultiTenantClient(conn)

	// Prepare transaction request
	req := &pb.CrtTransaction{
		TenantId:     1,
		BranchId:     1,
		ProductId:    1,
		QuantitySold: 2,
		PricePerUnit: 19.99,
	}

	// Use a WaitGroup to wait for all goroutines to complete
	for {
		for i := 0; i < 100; i++ {
			go func(i int) {

				// Create a NEW timeout context for each request
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				// Call CreateTransaction RPC
				res, err := c.CreateTransaction(ctx, req)
				if err != nil {
					log.Printf("[Request %d] Error calling CreateTransaction: %v", i, err)
					return
				}

				// Print response
				fmt.Printf("[Request %d] Response: %s (Success: %v)\n", i, res.Message, res.Success)
			}(i)

		}

		time.Sleep(2 * time.Second)
	}

}
