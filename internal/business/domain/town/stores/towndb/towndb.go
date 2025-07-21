// Package towndb contains town related CRUD functionality.
package towndb

import (
	"context"
	"fmt"

	"github.com/machilan1/cruise/internal/business/domain/town"
	"github.com/machilan1/cruise/internal/business/sdk/sqldb"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
)

// Store manages the set of APIs for town database access.
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
func (s *Store) NewWithTx(txM tran.TxManager) (town.Storer, error) {
	ec, err := tran.GetExtContext(txM)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: ec,
	}, nil
}

func (s *Store) QueryByCityID(ctx context.Context, cityID int) ([]town.Town, error) {
	data := struct {
		ID int `db:"city_id"`
	}{
		ID: cityID,
	}

	const q = `
		SELECT 	town_id,
				city_id,
		    	city_name,
				town_name,
			   	post_code
		FROM m_v_towns
		WHERE city_id = :city_id
		ORDER BY town_id
	`

	var dbtwns []dbTown
	if err := sqldb.NamedQuerySlice(ctx, s.db, q, data, &dbtwns); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	twns, err := toCoreTowns(dbtwns)
	if err != nil {
		return nil, err
	}

	return twns, nil
}

func (s *Store) QueryCities(ctx context.Context) ([]town.City, error) {
	const q = `
		SELECT city_id, city_name
		FROM cities
		ORDER BY city_id
	`

	var dbCities []dbCity
	if err := sqldb.NamedQuerySlice(ctx, s.db, q, struct{}{}, &dbCities); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}
	return toCoreCities(dbCities), nil
}
