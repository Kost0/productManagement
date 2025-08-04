package server

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

func (s *ProductServer) GetList(ctx context.Context, req *pb.GetListRequest) (*pb.ListResponse, error) {
	products, err := s.repo.GetAll(ctx, req.Limit, req.Offset)

	if err != nil {
		log.Printf("GetProducts failed: %v", err)
	}

	pbProducts := make([]*pb.ProductResponse, 0, len(products))

	for _, product := range products {
		pbProducts = append(pbProducts, &pb.ProductResponse{
			ProductId:   strconv.Itoa(product.ID),
			Name:        product.Name,
			Description: product.Description,
			PriceBuy:    int32(product.PriceBuy),
			PriceSell:   int32(product.PriceSell),
			SupplierId:  int32(product.SupplierID),
			Weight:      int32(product.Weight),
		})
	}

	return &pb.ListResponse{
		Products:   pbProducts,
		TotalCount: int32(len(products)),
	}, nil
}

func (s *ProductServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}
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

	id, err := strconv.Atoi(req.ProductId)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "product_id is invalid")
	}

	product := &domain.Product{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		PriceBuy:    int(req.PriceBuy),
		PriceSell:   int(req.PriceSell),
		SupplierID:  int(req.SupplierId),
		Weight:      int(req.Weight),
	}

	updatedProduct, err := s.repo.Update(ctx, product)
	if err != nil {
		log.Printf("CreateProduct failed: %v", err)
		return nil, status.Error(codes.Internal, "CreateProduct failed")
	}

	return &pb.ProductResponse{
		ProductId:   strconv.Itoa(updatedProduct.ID),
		Name:        updatedProduct.Name,
		Description: updatedProduct.Description,
		PriceBuy:    int32(updatedProduct.PriceBuy),
		PriceSell:   int32(updatedProduct.PriceSell),
		SupplierId:  int32(updatedProduct.SupplierID),
		Weight:      int32(updatedProduct.Weight),
	}, nil
}

func (s *ProductServer) DeleteProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.DeleteProductResponse, error) {
	err := s.repo.Delete(ctx, req.ProductId)
	if err != nil {
		log.Printf("DeleteProduct failed: %v", err)
		return &pb.DeleteProductResponse{
			Success: false,
			Message: "DeleteProduct failed",
		}, err
	}
	return &pb.DeleteProductResponse{
		Success: true,
		Message: "DeleteProduct success",
	}, nil
}
