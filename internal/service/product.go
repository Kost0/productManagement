package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	pb "products/gen/go/product"
	"products/internal/domain"
	"products/internal/repository"
	"strconv"
)

type ProductServer struct {
	pb.ProductServiceServer
	repo repository.ProductRepository
}

func NewProductServer(repo repository.ProductRepository) *ProductServer {
	return &ProductServer{repo: repo}
}

func (s *ProductServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Description == "" {
		return nil, status.Error(codes.InvalidArgument, "description is required")
	}
	if req.PriceBuy < 0 || req.PriceSell < 0 {
		return nil, status.Error(codes.InvalidArgument, "price_buy or price_sell less then zero")
	}
	if req.SupplierId < 0 {
		return nil, status.Error(codes.InvalidArgument, "supplier_id is required")
	}
	if req.Weight < 0 {
		return nil, status.Error(codes.InvalidArgument, "weight is required")
	}

	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		PriceBuy:    int(req.PriceBuy),
		PriceSell:   int(req.PriceSell),
		SupplierID:  int(req.SupplierId),
		Weight:      int(req.Weight),
	}

	createdProduct, err := s.repo.Create(ctx, product)
	if err != nil {
		log.Printf("CreateProduct failed: %v", err)
		return nil, status.Error(codes.Internal, "CreateProduct failed")
	}

	return &pb.ProductResponse{
		ProductId:   strconv.Itoa(createdProduct.ID),
		Name:        createdProduct.Name,
		Description: createdProduct.Description,
		PriceBuy:    int32(createdProduct.PriceBuy),
		PriceSell:   int32(createdProduct.PriceSell),
		SupplierId:  int32(createdProduct.SupplierID),
		Weight:      int32(createdProduct.Weight),
	}, nil
}

func (s *ProductServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	receivedProduct, err := s.repo.GetByID(ctx, req.ProductId)
	if err != nil {
		log.Printf("GetProduct failed: %v", err)
		return nil, status.Error(codes.Internal, "Error")
	}
	if receivedProduct == nil {
		return nil, status.Error(codes.NotFound, "Product not found")
	}

	return &pb.ProductResponse{
		ProductId:   strconv.Itoa(receivedProduct.ID),
		Name:        receivedProduct.Name,
		Description: receivedProduct.Description,
		PriceBuy:    int32(receivedProduct.PriceBuy),
		PriceSell:   int32(receivedProduct.PriceSell),
		SupplierId:  int32(receivedProduct.SupplierID),
		Weight:      int32(receivedProduct.Weight),
	}, nil
}
