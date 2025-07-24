package auctionapi

import (
	"net/http"

	"github.com/machilan1/cruise/internal/app/sdk/mid"
	"github.com/machilan1/cruise/internal/business/domain/auction"
	"github.com/machilan1/cruise/internal/business/domain/auth"
	"github.com/machilan1/cruise/internal/business/sdk/sess"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
	"github.com/machilan1/cruise/internal/framework/logger"
	"github.com/machilan1/cruise/internal/framework/web"
)

// Config contains all the mandatory dependencies for this group of handlers.
type Config struct {
	Log     *logger.Logger
	TxM     tran.TxManager
	Sess    *sess.Manager
	Auth    *auth.Core
	Auction *auction.Core
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	aucCtx := auctionCtx(cfg.Auction)

	hdl := newHandlers(cfg.Log, cfg.TxM, cfg.Auction)

	authn := mid.AuthN(cfg.Sess, cfg.Auth)
	authAdmin := mid.AuthNAndIsOneOfAdmin(cfg.Sess, cfg.Auth)

	authApp := *app
	authApp.PostMid = []web.MidFunc{
		authn,
	}

	authIsOneOfAdmin := *app
	authIsOneOfAdmin.PostMid = []web.MidFunc{
		authAdmin,
	}

	app.HandleFunc(http.MethodGet, version, "/auctions", hdl.query)
	app.HandleFunc(http.MethodGet, version, "/auctions/{auctionID}", hdl.queryByID, aucCtx)
	app.HandleFunc(http.MethodPost, version, "/auctions", hdl.create)
	app.HandleFunc(http.MethodPatch, version, "/auctions/{auctionID}", hdl.update, aucCtx)
	// TODO: If restore is ever needed, please isolate a new queryByID without using view to filter deletedAt, otherwise the precondition is not available.
	app.HandleFunc(http.MethodDelete, version, "/auctions/{auctionID}", hdl.softDelete, aucCtx)

	// Ready for testing:
}
