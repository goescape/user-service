package usecases

import (
	"context"
	"log"
	"user-svc/proto/product"
)

type ProductUsecases interface {
	InsertProduct(ctx context.Context, req *product.ProductInsertRequest) (*product.ProductInsertResponse, error)
	ListProduct(ctx context.Context, req *product.ListProductRequest) (*product.ListProductResponse, error)
}

var _ ProductUsecases = &productUsecase{}

type productUsecase struct {
	grpcServer product.ProductServiceClient
}

func NewProductUsecase(grpcServer product.ProductServiceClient) *productUsecase {
	return &productUsecase{
		grpcServer: grpcServer,
	}
}

func (s *productUsecase) InsertProduct(ctx context.Context, req *product.ProductInsertRequest) (*product.ProductInsertResponse, error) {
	res, err := s.grpcServer.InsertProduct(ctx, req)
	if err != nil {
		log.Default().Println("error insert product", err)
		return nil, err
	}

	return res, nil
}

func (s *productUsecase) ListProduct(ctx context.Context, req *product.ListProductRequest) (*product.ListProductResponse, error) {
	res, err := s.grpcServer.ListProduct(ctx, req)
	if err != nil {
		log.Default().Println("error list product", err)
		return nil, err
	}

	return res, nil
}
