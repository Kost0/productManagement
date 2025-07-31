package service

import (
	pb "products/gen/go/product"
	"products/internal/repository"
)

type ProductServer struct {
	pb.ProductServiceServer
	repo repository.ProductRepository
}
