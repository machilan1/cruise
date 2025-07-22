package auctionhouseapi

import (
	"net/http"

	"github.com/machilan1/cruise/internal/app/sdk/mid"
	"github.com/machilan1/cruise/internal/business/domain/auctionhouse"
	"github.com/machilan1/cruise/internal/business/domain/auth"
	"github.com/machilan1/cruise/internal/business/sdk/sess"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
	"github.com/machilan1/cruise/internal/framework/logger"
	"github.com/machilan1/cruise/internal/framework/web"
)

type Config struct {
	Log          *logger.Logger
	TxM          tran.TxManager
	Sess         *sess.Manager
	Auth         *auth.Core
	AuctionHouse *auctionhouse.Core
}

func Routes(app *web.App, cfg Config) {
	const version = "v1"

	ahCtx := auctionHouseCtx(cfg.AuctionHouse)

	hdl := newHandlers(cfg.Log, cfg.TxM, cfg.AuctionHouse)

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

	authApp.HandleFunc(http.MethodGet, version, "/auction-houses", hdl.query)
	authApp.HandleFunc(http.MethodGet, version, "/auction-houses/{auctionHouseID}", hdl.queryByID, ahCtx)

	authIsOneOfAdmin.HandleFunc(http.MethodPost, version, "/auction-houses", hdl.create)
	authIsOneOfAdmin.HandleFunc(http.MethodPut, version, "/auction-houses/{auctionHouseID}", hdl.update, ahCtx)
	authIsOneOfAdmin.HandleFunc(http.MethodDelete, version, "/auction-houses/{auctionHouseID}", hdl.Archive, ahCtx)
}
