package usecases

import (
	"context"
	"fmt"
	"net/http"
	"user-svc/helpers/fault"
	"user-svc/proto/product"
)

type productUsecase struct {
	serverRPC product.ProductServiceClient
}

func NewProductUsecase(serverRPC product.ProductServiceClient) *productUsecase {
	return &productUsecase{
		serverRPC: serverRPC,
	}
}

type ProductUsecases interface {
	InsertProduct(ctx context.Context, req *product.ProductInsertRequest) (*product.ProductInsertResponse, error)
	ListProduct(ctx context.Context, req *product.ListProductRequest) (*product.ListProductResponse, error)
}

func (s *productUsecase) InsertProduct(ctx context.Context, req *product.ProductInsertRequest) (*product.ProductInsertResponse, error) {
	insertOK, err := s.serverRPC.InsertProduct(ctx, req)
	if err != nil {
		return nil, fault.Custom(
			http.StatusUnprocessableEntity,
			fault.ErrUnprocessable,
			fmt.Sprintf("failed inserted product: %v", err.Error()),
		)
	}

	return insertOK, nil
}

func (s *productUsecase) ListProduct(ctx context.Context, req *product.ListProductRequest) (*product.ListProductResponse, error) {
	product, err := s.serverRPC.ListProduct(ctx, req)
	if err != nil {
		return nil, fault.Custom(
			http.StatusUnprocessableEntity,
			fault.ErrUnprocessable,
			fmt.Sprintf("failed retrieve list product: %v", err.Error()),
		)
	}

	return product, nil
}
