// Package auction provides support for the auction domain.
package auction

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/machilan1/cruise/internal/business/sdk/order"
	"github.com/machilan1/cruise/internal/business/sdk/paging"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
)

// Set of errors that are known to the business.
var (
	ErrNotFound     = errors.New("auction not found")
	ErrDataConflict = errors.New("request data conflict with current data")
)

type Storer interface {
	NewWithTx(txM tran.TxManager) (Storer, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page paging.Page) ([]Auction, error)
	QueryByID(ctx context.Context, auctionID int) (Auction, error)
	Create(ctx context.Context, auc Auction) (Auction, error)
	Update(ctx context.Context, auc Auction) (Auction, error)
	Delete(ctx context.Context, auc Auction) error
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

// Count returns the total number of auctions.
func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	cnt, err := c.storer.Count(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("count: %w", err)
	}

	return cnt, nil
}

// Query retrieves a list of existing auctions.
func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page paging.Page) ([]Auction, error) {
	aucs, err := c.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return aucs, nil
}

// QueryByID finds the auction by the specified ID.
func (c *Core) QueryByID(ctx context.Context, auctionID int) (Auction, error) {
	auc, err := c.storer.QueryByID(ctx, auctionID)
	if err != nil {
		return Auction{}, fmt.Errorf("query: auctionID[%d]: %w", auctionID, err)
	}

	return auc, nil
}

// Create adds a new auction to the system.
func (c *Core) Create(ctx context.Context, nAuc NewAuction) (Auction, error) {
	now := time.Now()
	auc := Auction{
		Name:             nAuc.Name,
		ExecutedCount:    nAuc.ExecutedCount,
		NotExecutedCount: nAuc.NotExecutedCount,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	auc, err := c.storer.Create(ctx, auc)
	if err != nil {
		return Auction{}, fmt.Errorf("create: %w", err)
	}

	result, err := c.storer.QueryByID(ctx, auc.ID)
	if err != nil {
		return Auction{}, fmt.Errorf("query after create: %w", err)
	}

	return result, nil
}

// Update modifies information about a auction.
func (c *Core) Update(ctx context.Context, auc Auction, uAuc UpdateAuction) (Auction, error) {
	auc.UpdatedAt = time.Now()
	if uAuc.Name != nil {
		auc.Name = *uAuc.Name
	}
	if uAuc.ExecutedCount != nil {
		auc.ExecutedCount = *uAuc.ExecutedCount
	}
	if uAuc.NotExecutedCount != nil {
		auc.NotExecutedCount = *uAuc.NotExecutedCount
	}
	auc, err := c.storer.Update(ctx, auc)
	if err != nil {
		return Auction{}, fmt.Errorf("update: %w", err)
	}

	result, err := c.storer.QueryByID(ctx, auc.ID)
	if err != nil {
		return Auction{}, fmt.Errorf("query after update: %w", err)
	}

	return result, nil
}

func (c *Core) Archive(ctx context.Context, auc Auction) (Auction, error) {
	now := time.Now()
	auc.UpdatedAt = now
	auc.DeletedAt = &now

	auc, err := c.storer.Update(ctx, auc)
	if err != nil {
		return Auction{}, fmt.Errorf("update: %w", err)
	}

	return auc, nil
}

func (c *Core) Restore(ctx context.Context, auc Auction) (Auction, error) {
	auc.UpdatedAt = time.Now()
	auc.DeletedAt = nil

	auc, err := c.storer.Update(ctx, auc)
	if err != nil {
		return Auction{}, fmt.Errorf("update: %w", err)
	}

	return auc, nil
}

// Delete removes the specified auction.
func (c *Core) Delete(ctx context.Context, auc Auction) error {
	if err := c.storer.Delete(ctx, auc); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
