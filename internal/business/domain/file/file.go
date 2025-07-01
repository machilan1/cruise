package file

import (
	"context"
	"fmt"
	"time"

	"github.com/machilan1/cruise/internal/business/sdk/tran"
)

var ErrNotFound = fmt.Errorf("file not found")

type Storer interface {
	NewWithTx(txM tran.TxManager) (Storer, error)
	QueryByID(ctx context.Context, id int) (File, error)
	QueryByPath(ctx context.Context, path Path) (File, error)
	Create(ctx context.Context, fl File) (File, error)
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

func (c *Core) QueryByID(ctx context.Context, id int) (File, error) {
	fl, err := c.storer.QueryByID(ctx, id)
	if err != nil {
		return File{}, fmt.Errorf("query: id[%d]: %w", id, err)
	}

	return fl, nil
}

func (c *Core) QueryByPath(ctx context.Context, path Path) (File, error) {
	fl, err := c.storer.QueryByPath(ctx, path)
	if err != nil {
		return File{}, fmt.Errorf("query: path[%s]: %w", path, err)
	}

	return fl, nil
}

func (c *Core) Create(ctx context.Context, nFl NewFile) (File, error) {
	now := time.Now()
	fl := File{
		Path:      nFl.Path,
		Original:  nFl.Original,
		MimeType:  nFl.Path.Mime(),
		Size:      nFl.Size,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if fl.Path.IsZero() {
		return File{}, fmt.Errorf("path is required")
	}
	if fl.Original == "" {
		fl.Original = fl.Path.Base() // Default to the base of the path, if not provided.
	}
	if fl.Size < 0 {
		return File{}, fmt.Errorf("size must be zero or greater")
	}

	fl, err := c.storer.Create(ctx, fl)
	if err != nil {
		return File{}, fmt.Errorf("create: %w", err)
	}

	return fl, nil
}
