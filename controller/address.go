package controller

import (
	"context"
	"fmt"

	"github.com/modular-project/address-service/model"
)

type GeoCoder interface {
	GeoCode(context.Context, string) (model.Location, error)
}

type AddressStorager interface {
	Create(context.Context, *model.Address) (string, error)
	DeleteByID(context.Context, string) (int64, error)
	GetByID(context.Context, string) (model.Address, error)
	Search(context.Context, *model.Search) ([]model.Address, error)
	Nearest(context.Context, []float64) (string, error)
}

type DeliveryStorager interface {
	Create(context.Context, *model.Delivery) (string, error)
	GetAll(context.Context, uint64) ([]model.Address, error)
	GetByID(context.Context, uint64, string) (model.Address, error)
	DeleteByID(context.Context, uint64, string) (int64, error)
}

type AddressService struct {
	ast AddressStorager
	dst DeliveryStorager
	gc  GeoCoder
}

func NewAddressService(as AddressStorager, ds DeliveryStorager, gc GeoCoder) AddressService {
	return AddressService{ast: as, dst: ds, gc: gc}
}

func (as AddressService) CreateDelivery(ctx context.Context, d *model.Delivery) (string, error) {
	loc, err := as.gc.GeoCode(ctx, d.Address.String())
	if err != nil {
		return "", fmt.Errorf("gc.GeoCode: %w", err)
	}
	d.Location = loc
	id, err := as.dst.Create(ctx, d)
	if err != nil {
		return "", fmt.Errorf("dst.Create: %w", err)
	}
	return id, nil
}

func (as AddressService) User(ctx context.Context, uID uint64) ([]model.Address, error) {
	a, err := as.dst.GetAll(ctx, uID)
	if err != nil {
		return nil, fmt.Errorf("dst.getAll: %w", err)
	}
	return a, nil
}

func (as AddressService) GetAddByID(ctx context.Context, aID string) (model.Address, error) {
	a, err := as.ast.GetByID(ctx, aID)
	if err != nil {
		return model.Address{}, fmt.Errorf("ast.GetByID: %w", err)
	}
	return a, nil
}

func (as AddressService) DeleteByUser(ctx context.Context, uID uint64, aID string) (int64, error) {
	d, err := as.dst.DeleteByID(ctx, uID, aID)
	if err != nil {
		return 0, fmt.Errorf("dst.DeleteByID: %w", err)
	}
	return d, nil
}

func (as AddressService) Nearest(ctx context.Context, uID uint64, aID string) (string, error) {
	add, err := as.dst.GetByID(ctx, uID, aID)
	if err != nil {
		return "", fmt.Errorf("dst.GetByID: %w", err)
	}
	id, err := as.ast.Nearest(ctx, add.Location.Coordinates)
	if err != nil {
		return "", fmt.Errorf("ast.Nearest: %w", err)
	}
	return id, nil
}

func (as AddressService) Create(ctx context.Context, a *model.Address) (string, error) {
	loc, err := as.gc.GeoCode(ctx, a.String())
	if err != nil {
		return "", fmt.Errorf("gc.GeoCode: %w", err)
	}
	a.Location = loc
	id, err := as.ast.Create(ctx, a)
	if err != nil {
		return "", fmt.Errorf("ast.Create: %w", err)
	}
	return id, nil
}

func (as AddressService) DeleteByID(ctx context.Context, aID string) (int64, error) {
	d, err := as.ast.DeleteByID(ctx, aID)
	if err != nil {
		return 0, fmt.Errorf("ast.DeleteByID: %w", err)
	}
	return d, nil
}

func (as AddressService) GetByID(ctx context.Context, uID uint64, aID string) (model.Address, error) {
	a, err := as.dst.GetByID(ctx, uID, aID)
	if err != nil {
		return model.Address{}, fmt.Errorf("dst.GetByID: %w", err)
	}
	return a, nil
}

func (as AddressService) Search(ctx context.Context, s *model.Search) ([]model.Address, error) {
	a, err := as.ast.Search(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("ast.Search: %w", err)
	}
	return a, nil
}
