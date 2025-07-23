package seriesmodelapi

import (
	"net/http"

	"github.com/machilan1/cruise/internal/app/sdk/mid"
	"github.com/machilan1/cruise/internal/business/domain/auth"
	"github.com/machilan1/cruise/internal/business/domain/seriesmodel"
	"github.com/machilan1/cruise/internal/business/sdk/sess"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
	"github.com/machilan1/cruise/internal/framework/logger"
	"github.com/machilan1/cruise/internal/framework/web"
)

type Config struct {
	Log         *logger.Logger
	TxM         tran.TxManager
	Sess        *sess.Manager
	Auth        *auth.Core
	SeriesModel *seriesmodel.Core
}

func Routes(app *web.App, cfg Config) {
	const version = "v1"

	smCtx := seriesModelCtx(cfg.SeriesModel)

	hdl := newHandlers(cfg.Log, cfg.TxM, cfg.SeriesModel)

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

	authApp.HandleFunc(http.MethodGet, version, "/series-models", hdl.query)
	authApp.HandleFunc(http.MethodGet, version, "/series-models/{seriesModelID}", hdl.queryByID, smCtx)

	authIsOneofAdmin.HandleFunc(http.MethodPost, version, "/series-models", hdl.create)
	authIsOneofAdmin.HandleFunc(http.MethodPut, version, "/series-models/{seriesModelID}", hdl.update, smCtx)
	authIsOneofAdmin.HandleFunc(http.MethodDelete, version, "/series-models/{seriesModelID}", hdl.delete, smCtx)

}
