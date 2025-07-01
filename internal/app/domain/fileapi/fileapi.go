package fileapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/machilan1/cruise/internal/business/domain/auth"
	"github.com/machilan1/cruise/internal/business/sdk/sess"

	"github.com/machilan1/cruise/internal/app/sdk/errs"
	"github.com/machilan1/cruise/internal/business/domain/file"
	"github.com/machilan1/cruise/internal/business/sdk/blobstore"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
	"github.com/machilan1/cruise/internal/framework/logger"
	"github.com/machilan1/cruise/internal/framework/web"
)

type handlers struct {
	log     *logger.Logger
	txM     tran.TxManager
	storage blobstore.BlobStore
	file    *file.Core
	Sess    *sess.Manager
	Auth    *auth.Core
}

func newHandlers(log *logger.Logger, txM tran.TxManager, storage blobstore.BlobStore, file *file.Core, sess *sess.Manager, auth *auth.Core) *handlers {
	return &handlers{
		log:     log,
		txM:     txM,
		storage: storage,
		file:    file,
		Sess:    sess,
		Auth:    auth,
	}
}

func (h *handlers) newWithTx(txM tran.TxManager) (*handlers, error) {
	fl, err := h.file.NewWithTx(txM)
	if err != nil {
		return nil, err
	}

	return &handlers{
		log:     h.log,
		txM:     txM,
		storage: h.storage,
		file:    fl,
		Sess:    h.Sess,
		Auth:    h.Auth,
	}, nil
}

func (h *handlers) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var input AppFileInput
	if err := web.Decode(r, &input); err != nil {
		return err
	}

	dir, err := file.NewDir("uploads")
	if err != nil {
		return err
	}

	path, err := file.NewRandomPath(dir, filepath.Ext(input.Filename))
	if err != nil {
		return err
	}

	nf := file.NewFile{
		Original: input.Filename,
		Path:     path,
		Size:     input.Size,
	}
	f, err := h.file.Create(ctx, nf)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}

	// get the path from the file
	signedURL, err := h.storage.SignedURL(ctx, f.Path.String(), blobstore.Write)
	if err != nil {
		return fmt.Errorf("signedURL: %w", err)
	}

	signedFl := AppSignedFile{
		Path:      f.Path.String(),
		SignedURL: signedURL,
	}
	return web.Respond(ctx, w, signedFl, http.StatusCreated)
}

// redirect handles the redirection to the signed URL of the file.
func (h *handlers) redirect(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	path := web.Param(r, "path")
	if path == "" {
		return errs.NewTrustedError(errors.New("missing path"), http.StatusBadRequest)
	}

	signedURL, err := h.storage.SignedURL(ctx, path, blobstore.Read)
	if err != nil {
		return fmt.Errorf("signed url: %w", err)
	}

	// Redirect to signedURL.
	// Note that the blob store must be configured with the correct CORS settings,
	// which should be `null` if the request is coming from a different origin
	// since we're using `credentials: 'include'` in the frontend.
	if err := web.Redirect(ctx, w, r, signedURL, http.StatusFound); err != nil {
		return fmt.Errorf("redirect: %w", err)
	}

	return nil
}
