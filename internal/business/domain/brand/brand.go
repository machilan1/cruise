package brand

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/machilan1/cruise/internal/business/sdk/tran"
)

var (
	ErrNotFound           = errors.New("brand not found")
	ErrConflict           = errors.New("request data conflicts with existing data")
	ErrDuplicatedDishName = errors.New("duplicated brand name")
)

type Storer interface {
	NewWithTx(txM tran.TxManager) (Storer, error)
	Query(ctx context.Context, filter QueryFilter) ([]Brand, error)
	QueryByID(ctx context.Context, id int) (Brand, error)
	Create(ctx context.Context, brd Brand) (Brand, error)
	Update(ctx context.Context, brd Brand) (Brand, error)
	Delete(ctx context.Context, id int) error
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

func (c *Core) Query(ctx context.Context, filter QueryFilter) ([]Brand, error) {
	ds, err := c.storer.Query(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return ds, nil
}

func (c *Core) QueryByID(ctx context.Context, id int) (Brand, error) {
	d, err := c.storer.QueryByID(ctx, id)
	if err != nil {
		return Brand{}, fmt.Errorf("querybyid: %w", err)
	}
	return d, nil
}

func (c *Core) Create(ctx context.Context, nb NewBrand) (Brand, error) {
	now := time.Now()
	d := Brand{
		Name:      nb.Name,
		Logo:      nb.Logo,
		CreatedAt: now,
		UpdatedAt: now,
	}

	d, err := c.storer.Create(ctx, d)
	if err != nil {
		return Brand{}, fmt.Errorf("create: %w", err)
	}

	return d, nil
}

func (c *Core) Update(ctx context.Context, d Brand, ub UpdateBrand) (Brand, error) {
	now := time.Now()

	d.UpdatedAt = now
	d.Logo = &ub.Logo

	d, err := c.storer.Update(ctx, d)
	if err != nil {
		return Brand{}, fmt.Errorf("update: %w", err)
	}

	d, err = c.storer.QueryByID(ctx, d.ID)
	if err != nil {
		return Brand{}, fmt.Errorf("querybyid: %w", err)
	}

	return d, nil
}

func (c *Core) Delete(ctx context.Context, id int) error {

	if err := c.storer.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
