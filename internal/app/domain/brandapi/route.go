package brandapi

import (
	"net/http"

	"github.com/machilan1/cruise/internal/app/sdk/mid"
	"github.com/machilan1/cruise/internal/business/domain/auth"
	"github.com/machilan1/cruise/internal/business/domain/brand"
	"github.com/machilan1/cruise/internal/business/sdk/sess"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
	"github.com/machilan1/cruise/internal/framework/logger"
	"github.com/machilan1/cruise/internal/framework/web"
)

type Config struct {
	Log   *logger.Logger
	TxM   tran.TxManager
	Sess  *sess.Manager
	Auth  *auth.Core
	Brand *brand.Core
}

func Routes(app *web.App, cfg Config) {
	const version = "v1"

	brdCtx := brandCtx(cfg.Brand)

	hdl := newHandlers(cfg.Log, cfg.TxM, cfg.Brand)

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

	authApp.HandleFunc(http.MethodGet, version, "/brands", hdl.query)
	authApp.HandleFunc(http.MethodGet, version, "/brands/{brandID}", hdl.queryByID, brdCtx)

	authIsOneofAdmin.HandleFunc(http.MethodPost, version, "/brands", hdl.create)
	authIsOneofAdmin.HandleFunc(http.MethodPut, version, "/brands/{brandID}", hdl.update, brdCtx)
	authIsOneofAdmin.HandleFunc(http.MethodDelete, version, "/brands/{brandID}", hdl.delete, brdCtx)
}
