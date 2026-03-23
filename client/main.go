package main

import (
	"context"
	"log"
	pb "master/pkg/api/test"
	"time"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	respAdd, err := client.CreateOrder(ctx, &pb.CreateOrderRequest{
		Item:     "Get a new Ram",
		Quantity: 5,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Result:", respAdd)

	respAdd1, err1 := client.CreateOrder(ctx, &pb.CreateOrderRequest{
		Item:     "Get a new Porsche",
		Quantity: 5,
	})
	if err1 != nil {
		log.Fatal(err)
	}
	log.Println("Result:", respAdd1)

	respList, err2 := client.ListOrders(ctx, &pb.ListOrdersRequest{})
	if err2 != nil {
		log.Fatal(err)
	}
	log.Println(respList)
}
