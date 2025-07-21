package townapi

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/machilan1/cruise/internal/app/sdk/errs"
	"github.com/machilan1/cruise/internal/business/domain/town"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
	"github.com/machilan1/cruise/internal/framework/logger"
	"github.com/machilan1/cruise/internal/framework/web"
)

type handlers struct {
	log  *logger.Logger
	txM  tran.TxManager
	town *town.Core
}

func newHandlers(log *logger.Logger, txM tran.TxManager, town *town.Core) *handlers {
	return &handlers{
		log:  log,
		txM:  txM,
		town: town,
	}
}

func (h *handlers) newWithTx(txM tran.TxManager) (*handlers, error) {
	townCore, err := h.town.NewWithTx(txM)
	if err != nil {
		return nil, err
	}

	return &handlers{
		log:  h.log,
		txM:  txM,
		town: townCore,
	}, nil
}

func (h *handlers) queryByCityID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "cityID")
	if id == "" {
		return errs.NewTrustedError(fmt.Errorf("cityID is required"), http.StatusBadRequest)
	}

	cityID, err := strconv.Atoi(id)
	if err != nil {
		return errs.NewTrustedError(fmt.Errorf("cityID must be a number"), http.StatusBadRequest)
	}

	twns, err := h.town.QueryByCityID(ctx, cityID)
	if err != nil {
		return fmt.Errorf("query: cityID[%d]: %w", cityID, err)
	}

	return web.Respond(ctx, w, ToAppTowns(twns), http.StatusOK)
}

func (h *handlers) queryCities(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	cities, err := h.town.QueryCities(ctx)
	if err != nil {
		return fmt.Errorf("query cities: %w", err)
	}

	return web.Respond(ctx, w, toAppCities(cities), http.StatusOK)
}
