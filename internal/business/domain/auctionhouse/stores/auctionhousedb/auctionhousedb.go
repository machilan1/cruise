package auctionhousedb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/machilan1/cruise/internal/business/domain/auctionhouse"
	"github.com/machilan1/cruise/internal/business/sdk/sqldb"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
)

type Store struct {
	db *sqldb.DB
}

func NewStore(db *sqldb.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) NewWithTx(txM tran.TxManager) (auctionhouse.Storer, error) {
	ec, err := tran.GetExtContext(txM)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: ec,
	}, nil
}

func (s *Store) Query(ctx context.Context, filter auctionhouse.QueryFilter) ([]auctionhouse.AuctionHouse, error) {
	data := map[string]any{}

	const q = `
		SELECT
			ah.auction_house_id,
			ah.auction_house_name,
			json_build_object(
					'address_detail',	ah.address_detail,
					'city_id',			c.city_id,
					'city_name', 		c.city_name,
					'town_id', 			ct.town_id,
					'town_name', 		ct.town_name,
				) as location
			ah.town_id,
			ah.address_detail,
			ah.created_at,
			ah.updated_at,
			ah.deleted_at
		FROM auction_houses ah
		LEFT JOIN city_towns ct ON ah.town_id = ct.town_id
		LEFT JOIN cities c ON ct.city_id = c.city_id
	`

	var sb strings.Builder
	applyFilter(filter, data, &sb)

	var dbahs []dbAuctionHouse
	if err := sqldb.NamedQuerySlice(ctx, s.db, q, data, &dbahs); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreAuctionHouses(dbahs)
}

func (s *Store) QueryByID(ctx context.Context, ahID int) (auctionhouse.AuctionHouse, error) {
	data := map[string]any{
		"auction_house_id": ahID,
	}

	const q = `
				SELECT
			ah.auction_house_id,
			ah.auction_house_name,
			json_build_object(
					'address_detail',	ah.address_detail,
					'city_id',			c.city_id,
					'city_name', 		c.city_name,
					'town_id', 			ct.town_id,
					'town_name', 		ct.town_name,
				) as location
			ah.town_id,
			ah.address_detail,
			ah.created_at,
			ah.updated_at,
			ah.deleted_at
		FROM auction_houses ah
		LEFT JOIN city_towns ct ON ah.town_id = ct.town_id
		LEFT JOIN cities c ON ct.city_id = c.city_id
		WHERE ah.auction_house_id  = :auction_house_id
	`

	var dbah dbAuctionHouse
	if err := sqldb.NamedQueryStruct(ctx, s.db, q, data, &dbah); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return auctionhouse.AuctionHouse{}, auctionhouse.ErrNotFound
		}
		return auctionhouse.AuctionHouse{}, fmt.Errorf("namedquerystruct: %w", err)
	}
	return toCoreAuctionHouse(dbah)
}

func (s *Store) Create(ctx context.Context, ah auctionhouse.AuctionHouse) (auctionhouse.AuctionHouse, error) {
	dba := toDBAuctionHouse(ah)

	const q = `
		INSERT INTO auction_houses(
			auction_house_name,
			address_detail,
			town_id,
			created_at,
			updated_at
		)VALUES(
			:auction_house_name,
			:address_detail,
			:town_id,
			:created_at,
			:updated_at
		)
		RETURNING action_house_id
	`

	if err := sqldb.NamedQueryStruct(ctx, s.db, q, dba, &dba); err != nil {
		if errors.Is(err, auctionhouse.ErrConflict) {
			return auctionhouse.AuctionHouse{}, auctionhouse.ErrConflict
		}
		return auctionhouse.AuctionHouse{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	ah, err := s.QueryByID(ctx, dba.ID)
	if err != nil {
		return auctionhouse.AuctionHouse{}, fmt.Errorf("querybyid: %w", err)
	}

	return ah, nil
}

func (s *Store) Update(ctx context.Context, ah auctionhouse.AuctionHouse) (auctionhouse.AuctionHouse, error) {
	dba := toDBAuctionHouse(ah)

	const q = `
		UPDATE auction_houses
		SET auction_house_name = :auction_house_name,
			address_detail = :address_detail,
			town_id = :town_id,
			updated_at = :updated_at
		WHERE auction_house_id = :auction_house_id
	`

	if err := sqldb.NamedQueryStruct(ctx, s.db, q, dba, &dba); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return auctionhouse.AuctionHouse{}, auctionhouse.ErrNotFound
		}

		if errors.Is(err, sqldb.ErrDBIntegrity) {
			return auctionhouse.AuctionHouse{}, auctionhouse.ErrConflict
		}

		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return auctionhouse.AuctionHouse{}, auctionhouse.ErrConflict
		}
		return auctionhouse.AuctionHouse{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	ah, err := s.QueryByID(ctx, dba.ID)
	if err != nil {
		return auctionhouse.AuctionHouse{}, fmt.Errorf("querybyid: %w", err)
	}

	return ah, nil
}

func (s *Store) Archive(ctx context.Context, ahID int) error {
	data := map[string]any{
		"auction_house_id": ahID,
	}

	const q = `
		UPDATE auction_house
		SET updated_at = CURRENT_TIMESTAMP,
			deleted_at = CURRENT_TIMESTAMP
		WHERE auction_house_id = :auction_house_id AND deleted_at IS NULL
	`

	if err := sqldb.NamedExecContext(ctx, s.db, q, data); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return auctionhouse.ErrNotFound
		}

		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
