package product

import (
	"context"
	"fmt"
	"github.com/ortymid/market/market/auth"
)

type Service struct {
	Storage Storage
}

func (s *Service) List(ctx context.Context, r ListRequest) ([]*Product, error) {
	ps, err := s.Storage.List(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	return ps, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Product, error) {
	p, err := s.Storage.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get product: %w", err)
	}

	return p, nil
}

func (s *Service) Create(ctx context.Context, r CreateRequest) (p *Product, err error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}
	if user == nil {
		err := auth.ErrPermission{Reason: "user not provided"}
		return nil, fmt.Errorf("create product: %w", err)
	}

	r.Seller = user.ID
	p, err = s.Storage.Create(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}

	return p, nil
}

func (s *Service) Update(ctx context.Context, r UpdateRequest) (*Product, error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}
	if user == nil {
		err := auth.ErrPermission{Reason: "user not provided"}
		return nil, fmt.Errorf("update product: %w", err)
	}

	p, err := s.Storage.Get(ctx, r.ID)
	if err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}

	if user.ID != p.Seller {
		err := auth.ErrPermission{Reason: "only own products allowed to update"}
		return nil, fmt.Errorf("update product: %w", err)
	}

	p, err = s.Storage.Update(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}

	return p, nil
}

func (s *Service) Delete(ctx context.Context, id string) (*Product, error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("delete product: %w", err)
	}
	if user == nil {
		err := auth.ErrPermission{Reason: "user not provided"}
		return nil, fmt.Errorf("delete product: %w", err)
	}

	p, err := s.Storage.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("delete product: %w", err)
	}

	if user.ID != p.Seller {
		err := auth.ErrPermission{Reason: "only own products allowed to delete"}
		return nil, fmt.Errorf("delete product: %w", err)
	}

	p, err = s.Storage.Delete(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("delete product: %w", err)
	}

	return p, nil
}