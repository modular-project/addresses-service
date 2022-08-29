package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/modular-project/address-service/model"
	pf "github.com/modular-project/protobuffers/address/address"
)

type AddressServicer interface {
	CreateDelivery(context.Context, *model.Delivery) (string, error)
	User(c context.Context, uID uint64) ([]model.Address, error)
	GetByID(c context.Context, uID uint64, aID string) (model.Address, error)
	GetAddByID(c context.Context, aID string) (model.Address, error)
	DeleteByUser(c context.Context, uID uint64, aID string) (int64, error)
	Create(context.Context, *model.Address) (string, error)
	DeleteByID(c context.Context, aID string) (int64, error)
	Search(context.Context, *model.Search) ([]model.Address, error)
	Nearest(c context.Context, uID uint64, aID string) (string, error)
}

//
type AddressUC struct {
	pf.UnimplementedAddressServiceServer
	as AddressServicer
}

func NewAddressUC(as AddressServicer) AddressUC {
	return AddressUC{as: as}
}

func protoAddress(m *model.Address) pf.Address {
	return pf.Address{
		Id:      m.ID.String(),
		Line1:   m.Street,
		Line2:   m.Suburb,
		City:    m.City,
		Pc:      m.PostalCode,
		State:   m.State,
		Country: m.Country,
	}
}

func modelAddress(p *pf.Address) model.Address {
	return model.Address{
		Street:     p.Line1,
		Suburb:     p.Line2,
		City:       p.City,
		PostalCode: p.Pc,
		State:      p.State,
		Country:    p.Country,
	}
}

func (uc AddressUC) CreateDelivery(c context.Context, d *pf.Delivery) (*pf.ID, error) {
	if d.Address == nil {
		return nil, errors.New("empty address")
	}
	m := model.Delivery{
		UserID:  d.UserId,
		Address: modelAddress(d.Address),
	}
	id, err := uc.as.CreateDelivery(c, &m)
	if err != nil {
		return nil, fmt.Errorf("create delivery: %w", err)
	}
	return &pf.ID{Id: id}, nil
}

func (uc AddressUC) GetAllByUser(c context.Context, u *pf.User) (*pf.ResponseAll, error) {
	ads, err := uc.as.User(c, u.Id)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	if ads == nil {
		return nil, nil
	}
	resA := make([]*pf.Address, len(ads))
	for i := range ads {
		pa := protoAddress(&ads[i])
		resA[i] = &pa
	}
	return &pf.ResponseAll{Address: resA}, nil
}

func (uc AddressUC) DeleteByID(c context.Context, u *pf.User) (*pf.ResponseDelete, error) {
	_, err := uc.as.DeleteByUser(c, u.Id, u.AddressId)
	if err != nil {
		return nil, fmt.Errorf("delete by user: %w", err)
	}
	return &pf.ResponseDelete{}, nil
}

func (uc AddressUC) GetByID(c context.Context, u *pf.User) (*pf.Address, error) {
	ma, err := uc.as.GetByID(c, u.Id, u.AddressId)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	pa := protoAddress(&ma)
	return &pa, nil
}

func (uc AddressUC) GetAddByID(c context.Context, ID *pf.ID) (*pf.Address, error) {
	ma, err := uc.as.GetAddByID(c, ID.Id)
	if err != nil {
		return nil, fmt.Errorf("get add by id: %w", err)
	}
	pa := protoAddress(&ma)
	return &pa, nil
}

func (uc AddressUC) GetByUser(c context.Context, u *pf.User) (*pf.Address, error) {
	ma, err := uc.as.GetByID(c, u.Id, u.AddressId)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	pa := protoAddress(&ma)
	return &pa, nil
}

func (uc AddressUC) CreateEstablishment(c context.Context, pa *pf.Address) (*pf.ID, error) {
	ma := modelAddress(pa)
	id, err := uc.as.Create(c, &ma)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &pf.ID{Id: id}, nil
}

func (uc AddressUC) DeleteEstablishment(c context.Context, id *pf.ID) (*pf.ResponseDelete, error) {
	_, err := uc.as.DeleteByID(c, id.Id)
	if err != nil {
		return nil, fmt.Errorf("delete by id: %w", err)
	}
	return &pf.ResponseDelete{}, nil
}

func (uc AddressUC) Search(c context.Context, ps *pf.SearchAddress) (*pf.ResponseAll, error) {
	ms := model.Search{
		Limit:  int64(ps.Default.Limit),
		Offset: int64(ps.Default.Offset),
	}
	mas, err := uc.as.Search(c, &ms)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	if mas == nil {
		return nil, nil
	}
	pas := make([]*pf.Address, len(mas))
	for i := range mas {
		pa := protoAddress(&mas[i])
		pas[i] = &pa
	}
	return &pf.ResponseAll{Address: pas}, nil
}

func (uc AddressUC) Nearest(c context.Context, u *pf.User) (*pf.ID, error) {
	id, err := uc.as.Nearest(c, u.Id, u.AddressId)
	if err != nil {
		return nil, fmt.Errorf("nearest: %w", err)
	}
	return &pf.ID{Id: id}, nil
}
