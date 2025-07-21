package auctionhouse

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/machilan1/cruise/internal/business/sdk/tran"
)

var (
	ErrNotFound = errors.New("auction house not found")
	ErrConflict = errors.New("input data conflicts with existing data")
)

type Storer interface {
	NewWithTx(txM tran.TxManager) (Storer, error)
	Query(ctx context.Context, filter QueryFilter) ([]AuctionHouse, error)
	QueryByID(ctx context.Context, auID int) (AuctionHouse, error)
	Create(ctx context.Context, au AuctionHouse) (AuctionHouse, error)
	Update(ctx context.Context, au AuctionHouse) (AuctionHouse, error)
	Archive(ctx context.Context, auID int) error
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

func (c *Core) Query(ctx context.Context, filter QueryFilter) ([]AuctionHouse, error) {
	ahs, err := c.storer.Query(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return ahs, nil
}

func (c *Core) QueryByID(ctx context.Context, ahID int) (AuctionHouse, error) {
	ah, err := c.storer.QueryByID(ctx, ahID)
	if err != nil {
		return AuctionHouse{}, fmt.Errorf("querybyid: %w", err)
	}

	return ah, nil
}

func (c *Core) Create(ctx context.Context, nah NewAuctionHouse) (AuctionHouse, error) {
	now := time.Now()
	au := AuctionHouse{
		Name: nah.Name,
		Location: AuctionHouseLocation{
			Address: nah.Address,
			TownID:  nah.TownID,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	ah, err := c.storer.Create(ctx, au)
	if err != nil {
		return AuctionHouse{}, fmt.Errorf("create: %w", err)
	}

	return ah, nil
}

func (c *Core) Update(ctx context.Context, uah UpdateAuctionHouse, ah AuctionHouse) (AuctionHouse, error) {

	if uah.Name != nil {
		ah.Name = *uah.Name
	}

	if uah.Address != nil {
		ah.Location.Address = *uah.Address
	}

	if uah.TownID != nil {
		ah.Location.TownID = *uah.TownID
	}

	ah, err := c.storer.Update(ctx, ah)
	if err != nil {
		return AuctionHouse{}, fmt.Errorf("update: %w", err)
	}

	return ah, nil
}

func (c *Core) Archive(ctx context.Context, ahID int) error {
	if err := c.storer.Archive(ctx, ahID); err != nil {
		return fmt.Errorf("archive: %w", err)
	}

	return nil
}
