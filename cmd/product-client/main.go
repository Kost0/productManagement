package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb "products/gen/go/product"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewProductServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createResp, err := client.CreateProduct(ctx, &pb.CreateProductRequest{
		Name:        "Test1",
		Description: "Test1 description",
		PriceBuy:    int32(5),
		PriceSell:   int32(10),
		SupplierId:  int32(1),
		Weight:      int32(200),
	})

	if err != nil {
		log.Fatalf("could not create: %v", err)
	}

	log.Printf("createResp: %v", createResp)

	getResp, err := client.GetProduct(ctx, &pb.GetProductRequest{
		ProductId: "10",
	})

	if err != nil {
		log.Fatalf("could not get: %v", err)
	}

	log.Printf("getResp: %v", getResp)

	createResp2, err := client.CreateProduct(ctx, &pb.CreateProductRequest{
		Name:        "Test2",
		Description: "Test2 description",
		PriceBuy:    int32(10),
		PriceSell:   int32(11),
		SupplierId:  int32(1),
		Weight:      int32(100),
	})

	if err != nil {
		log.Fatalf("could not create: %v", err)
	}

	log.Printf("createResp: %v", createResp2)

	list, err := client.GetList(ctx, &pb.GetListRequest{
		Limit:  int32(10),
		Offset: int32(0),
	})

	if err != nil {
		log.Fatalf("could not get: %v", err)
	}

	log.Printf("list:")

	for _, product := range list.Products {
		log.Printf("Product: %v", product)
	}

	updateResp, err := client.UpdateProduct(ctx, &pb.UpdateProductRequest{
		ProductId:   "0",
		Name:        "Test1",
		Description: "Test1 description",
		PriceBuy:    int32(5),
		PriceSell:   int32(20),
		SupplierId:  int32(1),
		Weight:      int32(200),
	})

	if err != nil {
		log.Fatalf("could not update: %v", err)
	}

	log.Printf("updateResp: %v", updateResp)

	deleted, err := client.DeleteProduct(ctx, &pb.GetProductRequest{
		ProductId: "11",
	})

	if err != nil {
		log.Fatalf("could not delete: %v", err)
	}

	if !deleted.Success {
		log.Fatalf("could not delete")
	}

	log.Printf("deleted: %v", deleted.Message)

	log.Printf("All checks of product service completed successfully")
}
