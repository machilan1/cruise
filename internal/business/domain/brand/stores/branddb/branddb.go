package branddb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/machilan1/cruise/internal/business/domain/brand"
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

func (s *Store) NewWithTx(txM tran.TxManager) (brand.Storer, error) {
	ec, err := tran.GetExtContext(txM)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: ec,
	}, nil
}

func (s *Store) Query(ctx context.Context, filter brand.QueryFilter) ([]brand.Brand, error) {
	data := map[string]any{}

	const q = `
		SELECT 
			brand_id,
			brand_name,
			logo,
			created_at,
			updated_at
		FROM brands
	`

	var sb strings.Builder
	applyFilter(filter, data, &sb)

	var brds []dbBrand
	if err := sqldb.NamedQuerySlice(ctx, s.db, q, data, &brds); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreBrands(brds), nil
}

func (s *Store) QueryByID(ctx context.Context, id int) (brand.Brand, error) {
	data := map[string]any{
		"brand_id": id,
	}

	const q = `
		SELECT 
			brand_id,
			brand_name,
			logo,
			created_at,
			updated_at
		FROM brands
		WHERE brand_id = :brand_id
	`

	var dbbr dbBrand
	if err := sqldb.NamedQueryStruct(ctx, s.db, q, data, &dbbr); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return brand.Brand{}, brand.ErrNotFound
		}
		return brand.Brand{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreBrand(dbbr), nil
}

func (s *Store) Create(ctx context.Context, brd brand.Brand) (brand.Brand, error) {
	dbbrd := toDBBrand(brd)

	const q = `
		INSERT INTO brands
			(
				brand_name,
				logo,
				created_at,
				updated_at
			)VALUES(	
				:brand_name,
				:logo,
				:created_at,
				:updated_at	
			)
		RETURNING brand_id
	`

	if err := sqldb.NamedQueryStruct(ctx, s.db, q, dbbrd, &dbbrd); err != nil {
		if errors.Is(err, sqldb.ErrDBIntegrity) {
			return brand.Brand{}, brand.ErrConflict
		}
		return brand.Brand{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreBrand(dbbrd), nil
}

func (s *Store) Update(ctx context.Context, brd brand.Brand) (brand.Brand, error) {
	dbbrd := toDBBrand(brd)
	const q = `
		UPDATE brands
		SET logo = :logo
		WHERE brand_id = :brand_id
		RETURNING brand_id
	`

	if err := sqldb.NamedQueryStruct(ctx, s.db, q, dbbrd, &dbbrd); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return brand.Brand{}, brand.ErrNotFound
		}

		if errors.Is(err, sqldb.ErrDBIntegrity) {
			return brand.Brand{}, brand.ErrConflict
		}

		return brand.Brand{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreBrand(dbbrd), nil
}

func (s *Store) Delete(ctx context.Context, id int) error {
	data := map[string]any{
		"brand_id": id,
	}

	const q = `
		DELETE FROM brands WHERE brand_id = :brand_id
	`

	if err := sqldb.NamedExecContext(ctx, s.db, q, data); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return brand.ErrNotFound
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
