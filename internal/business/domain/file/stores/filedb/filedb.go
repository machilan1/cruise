package filedb

import (
	"context"
	"errors"
	"fmt"

	"github.com/machilan1/cruise/internal/business/domain/file"
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

// NewWithTx constructs a new Store which replaces the underlying database connection with the provided transaction.
func (s *Store) NewWithTx(txM tran.TxManager) (file.Storer, error) {
	ec, err := tran.GetExtContext(txM)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: ec,
	}, nil
}

func (s *Store) QueryByID(ctx context.Context, id int) (file.File, error) {
	data := struct {
		ID int `db:"file_id"`
	}{
		ID: id,
	}

	const q = `
		SELECT
			file_id, path, original_filename, mime_type, SIZE, created_at, updated_at
		FROM files
		WHERE file_id = :file_id
	`
	var dbF dbFile
	if err := sqldb.NamedQueryStruct(ctx, s.db, q, data, &dbF); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return file.File{}, file.ErrNotFound
		}
		return file.File{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreFile(dbF)
}

func (s *Store) QueryByPath(ctx context.Context, path file.Path) (file.File, error) {
	data := struct {
		Path string `db:"path"`
	}{
		Path: path.String(),
	}

	const q = `
		SELECT
			file_id, path, original_filename, mime_type, SIZE, created_at, updated_at
		FROM files
		WHERE PATH = :PATH
	`
	var dbF dbFile
	if err := sqldb.NamedQueryStruct(ctx, s.db, q, data, &dbF); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return file.File{}, file.ErrNotFound
		}
		return file.File{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreFile(dbF)
}

func (s *Store) Create(ctx context.Context, fl file.File) (file.File, error) {
	dbF := toDBFile(fl)

	const q = `
		INSERT INTO
			files (path, original_filename, mime_type, size, created_at, updated_at)
		VALUES (:path, :original_filename, :mime_type, :size, :created_at, :updated_at)
		RETURNING file_id
	`
	var dest struct {
		ID int `db:"file_id"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.db, q, dbF, &dest); err != nil {
		return file.File{}, fmt.Errorf("namedquerystruct: %w", err)
	}
	dbF.ID = dest.ID

	return toCoreFile(dbF)
}
