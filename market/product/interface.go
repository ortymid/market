package product

import "context"

//go:generate mockgen -destination=../../mock/product_service.go -package mock -mock_names=Interface=ProductService . Interface

type Interface interface {
	Find(ctx context.Context, r FindRequest) ([]*Product, error)
	FindOne(ctx context.Context, id string) (*Product, error)
	Create(ctx context.Context, r CreateRequest) (*Product, error)
	Update(ctx context.Context, r UpdateRequest) (*Product, error)
	Delete(ctx context.Context, id string) (*Product, error)
}
