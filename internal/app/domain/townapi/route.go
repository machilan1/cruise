package townapi

import (
	"net/http"

	"github.com/machilan1/cruise/internal/app/sdk/mid"
	"github.com/machilan1/cruise/internal/business/domain/auth"
	"github.com/machilan1/cruise/internal/business/domain/town"
	"github.com/machilan1/cruise/internal/business/sdk/sess"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
	"github.com/machilan1/cruise/internal/framework/logger"
	"github.com/machilan1/cruise/internal/framework/web"
)

// Config contains all the mandatory dependencies for this group of handlers.
type Config struct {
	Log  *logger.Logger
	TxM  tran.TxManager
	Sess *sess.Manager
	Auth *auth.Core
	Town *town.Core
}

func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authn := mid.AuthN(cfg.Sess, cfg.Auth)

	hdl := newHandlers(cfg.Log, cfg.TxM, cfg.Town)

	authApp := *app
	authApp.PostMid = []web.MidFunc{
		authn,
	}

	authApp.HandleFunc(http.MethodGet, version, "/cities", hdl.queryCities)
	authApp.HandleFunc(http.MethodGet, version, "/cities/{cityID}/towns", hdl.queryByCityID)
}
