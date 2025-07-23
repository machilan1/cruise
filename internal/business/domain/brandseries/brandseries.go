package brandseries

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/machilan1/cruise/internal/business/sdk/tran"
)

var (
	ErrNotFound = errors.New("brand series not found")
	ErrConflict = errors.New("input data conflicts with existing data")
)

type Storer interface {
	NewWithTx(txM tran.TxManager) (Storer, error)
	Query(ctx context.Context, filter QueryFilter) ([]BrandSeries, error)
	QueryByID(ctx context.Context, bsID int) (BrandSeries, error)
	Create(ctx context.Context, bs BrandSeries) (BrandSeries, error)
	Update(ctx context.Context, bs BrandSeries) (BrandSeries, error)
	Delete(ctx context.Context, bsID int) error
}

type Core struct {
	storer Storer
}

func NewCore(storer Storer) *Core {
	return &Core{
		storer: storer,
	}
}

func (c *Core) NewWithTx(txM tran.TxManager) (*Core, error) {
	storer, err := c.storer.NewWithTx(txM)
	if err != nil {
		return nil, err
	}

	return &Core{
		storer: storer,
	}, nil
}

func (c *Core) Query(ctx context.Context, filter QueryFilter) ([]BrandSeries, error) {
	bss, err := c.storer.Query(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return bss, nil
}

func (c *Core) QueryByID(ctx context.Context, bsID int) (BrandSeries, error) {
	bs, err := c.storer.QueryByID(ctx, bsID)
	if err != nil {
		return BrandSeries{}, fmt.Errorf("querybyid: %w", err)
	}
	return bs, nil
}

func (c *Core) Create(ctx context.Context, nbs NewBrandSeries) (BrandSeries, error) {
	now := time.Now()
	bs := BrandSeries{
		Name: nbs.Name,
		Brand: BrandSeriesBrand{
			ID: nbs.BrandID,
		},
		CreatedAt: now,
		UpdateAt:  now,
	}

	bs, err := c.storer.Create(ctx, bs)
	if err != nil {
		return BrandSeries{}, fmt.Errorf("create: %w", err)
	}

	return bs, nil
}

func (c *Core) Update(ctx context.Context, ubs UpdateBrandSeries, bs BrandSeries) (BrandSeries, error) {
	if ubs.BrandID != nil {
		bs.Brand.ID = *ubs.BrandID
	}

	if ubs.Name != nil {
		bs.Name = *ubs.Name
	}

	bs, err := c.storer.Update(ctx, bs)
	if err != nil {
		return BrandSeries{}, fmt.Errorf("update: %w", err)
	}

	return bs, nil
}

func (c *Core) Delete(ctx context.Context, bsID int) error {
	if err := c.storer.Delete(ctx, bsID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
