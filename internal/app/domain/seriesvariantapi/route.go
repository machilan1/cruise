package seriesvariantapi

import (
	"net/http"

	"github.com/machilan1/cruise/internal/app/sdk/mid"
	"github.com/machilan1/cruise/internal/business/domain/auth"
	"github.com/machilan1/cruise/internal/business/domain/seriesvariant"
	"github.com/machilan1/cruise/internal/business/sdk/sess"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
	"github.com/machilan1/cruise/internal/framework/logger"
	"github.com/machilan1/cruise/internal/framework/web"
)

type Config struct {
	Log           *logger.Logger
	TxM           tran.TxManager
	Sess          *sess.Manager
	Auth          *auth.Core
	SeriesVariant *seriesvariant.Core
}

func Routes(app *web.App, cfg Config) {
	const version = "v1"

	svCtx := seriesVariantCtx(cfg.SeriesVariant)

	hdl := newHandlers(cfg.Log, cfg.TxM, cfg.SeriesVariant)

	authn := mid.AuthN(cfg.Sess, cfg.Auth)
	authAdmin := mid.AuthNAndIsOneOfAdmin(cfg.Sess, cfg.Auth)

	authApp := *app
	authApp.PostMid = []web.MidFunc{
		authn,
	}
	authIsOneofAdmin := *app
	authIsOneofAdmin.PostMid = []web.MidFunc{
		authAdmin,
	}

	authApp.HandleFunc(http.MethodGet, version, "/series-variants", hdl.query)
	authApp.HandleFunc(http.MethodGet, version, "/series-variants/{seriesVariantID}", hdl.queryByID, svCtx)

	authIsOneofAdmin.HandleFunc(http.MethodPost, version, "/series-variants", hdl.create)
	authIsOneofAdmin.HandleFunc(http.MethodPut, version, "/series-variants/{seriesVariantID}", hdl.update, svCtx)
	authIsOneofAdmin.HandleFunc(http.MethodDelete, version, "/series-variants/{seriesVariantID}", hdl.delete, svCtx)
}
