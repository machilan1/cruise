// Package auctiondb contains auction related CRUD functionality.
package auctiondb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/machilan1/cruise/internal/business/domain/auction"
	"github.com/machilan1/cruise/internal/business/sdk/order"
	"github.com/machilan1/cruise/internal/business/sdk/paging"
	"github.com/machilan1/cruise/internal/business/sdk/sqldb"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
)

// Store manages the set of APIs for auction database access.
type Store struct {
	db *sqldb.DB
}

// NewStore constructs the api for data access.
func NewStore(db *sqldb.DB) *Store {
	return &Store{
		db: db,
	}
}

// NewWithTx constructs a new Store which replaces the underlying database connection with the provided transaction.
func (s *Store) NewWithTx(txM tran.TxManager) (auction.Storer, error) {
	ec, err := tran.GetExtContext(txM)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: ec,
	}, nil
}

// Count returns the total number of auctions in the DB.
func (s *Store) Count(ctx context.Context, filter auction.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
        SELECT COUNT(*)
        FROM auctions
    `

	var sb strings.Builder
	sb.WriteString(q)
	s.applyFilter(filter, data, &sb)

	var dest struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.db, sb.String(), data, &dest); err != nil {
		return 0, fmt.Errorf("namedquerystruct: %w", err)
	}

	return dest.Count, nil
}

// Query retrieves a list of existing auctions from the database.
func (s *Store) Query(ctx context.Context, filter auction.QueryFilter, orderBy order.By, page paging.Page) ([]auction.Auction, error) {
	data := map[string]any{
		"offset":        page.Offset(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
        SELECT auction_id,
               created_at,
               updated_at,
               deleted_at
        FROM auctions
    `

	var sb strings.Builder
	sb.WriteString(q)
	s.applyFilter(filter, data, &sb)

	if err := s.orderByClause(orderBy, &sb); err != nil {
		return nil, err
	}

	sb.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbaucs []dbAuction
	if err := sqldb.NamedQuerySlice(ctx, s.db, sb.String(), data, &dbaucs); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreAuctions(dbaucs)
}

// QueryByID finds the auction identified by a given ID.
func (s *Store) QueryByID(ctx context.Context, auctionID int) (auction.Auction, error) {
	data := struct {
		ID int `db:"auction_id"`
	}{
		ID: auctionID,
	}

	const q = `
        SELECT auction_id,
               created_at,
               updated_at,
               deleted_at
        FROM auctions
        WHERE auction_id = :auction_id
    `

	var dbauc dbAuction
	if err := sqldb.NamedQueryStruct(ctx, s.db, q, data, &dbauc); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return auction.Auction{}, auction.ErrNotFound
		}
		return auction.Auction{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreAuction(dbauc)
}

// Create adds a Auction to the database. It returns an error if something went wrong
func (s *Store) Create(ctx context.Context, auc auction.Auction) (auction.Auction, error) {
	dbauc := toDBAuction(auc)

	const q = `
        INSERT INTO auctions
            (created_at, updated_at)
        VALUES
            (:created_at, :updated_at)
        RETURNING auction_id
    `

	if err := sqldb.NamedQueryStruct(ctx, s.db, q, dbauc, &dbauc); err != nil {
		if errors.Is(err, sqldb.ErrDBIntegrity) || errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return auction.Auction{}, auction.ErrDataConflict
		}
		return auction.Auction{}, fmt.Errorf("namedexeccontext: %w", err)
	}

	return toCoreAuction(dbauc)
}

// Update modifies data about a Auction. It will error if the specified ID is
// invalid or does not reference an existing Auction.
func (s *Store) Update(ctx context.Context, auc auction.Auction) (auction.Auction, error) {
	dbauc := toDBAuction(auc)

	const q = `
        UPDATE auctions
        SET updated_at = :updated_at
            deleted_at = :deleted_at,
        WHERE auction_id = :auction_id
        RETURNING auction_id
    `

	if err := sqldb.NamedQueryStruct(ctx, s.db, q, dbauc, &dbauc); err != nil {
		if errors.Is(err, sqldb.ErrDBIntegrity) || errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return auction.Auction{}, auction.ErrDataConflict
		}
		return auction.Auction{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreAuction(dbauc)
}

// Delete removes the Auction identified by a given ID.
func (s *Store) Delete(ctx context.Context, auc auction.Auction) error {
	data := struct {
		ID int `db:"auction_id"`
	}{
		ID: auc.ID,
	}

	const q = `
        DELETE
        FROM auctions
        WHERE auction_id = :auction_id
    `

	if err := sqldb.NamedExecContext(ctx, s.db, q, data); err != nil {
		if errors.Is(err, sqldb.ErrDBIntegrity) || errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return auction.ErrDataConflict
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
