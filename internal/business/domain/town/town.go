// Package town provides support for the town domain.
package town

import (
	"context"
	"fmt"

	"github.com/machilan1/cruise/internal/business/sdk/tran"
)

type Storer interface {
	NewWithTx(txM tran.TxManager) (Storer, error)
	QueryByCityID(ctx context.Context, cityID int) ([]Town, error)
	QueryCities(ctx context.Context) ([]City, error)
}

// ====================================================================================

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

func (c *Core) QueryByCityID(ctx context.Context, cityID int) ([]Town, error) {
	twns, err := c.storer.QueryByCityID(ctx, cityID)
	if err != nil {
		return nil, fmt.Errorf("query: cityID[%d]: %w", cityID, err)
	}

	return twns, nil
}

func (c *Core) QueryCities(ctx context.Context) ([]City, error) {
	cities, err := c.storer.QueryCities(ctx)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return cities, nil
}
