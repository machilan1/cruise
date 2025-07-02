package vehiclemodelapi

import (
	"net/http"

	"github.com/machilan1/cruise/internal/app/sdk/mid"
	"github.com/machilan1/cruise/internal/business/domain/auth"
	"github.com/machilan1/cruise/internal/business/domain/vehiclemodel"
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
	VehicleModel *vehiclemodel.Core
}

func Routes(app *web.App, cfg Config) {
	const version = "v1"

	vmCtx := vehicleModelCtx(cfg.VehicleModel)

	hdl := newHandlers(cfg.Log, cfg.TxM, cfg.VehicleModel)

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

	authApp.HandleFunc(http.MethodGet, version, "/vehicle-models", hdl.query)
	authApp.HandleFunc(http.MethodGet, version, "/vehicle-models/{vehicleModelID}", hdl.queryByID, vmCtx)

	authIsOneOfAdmin.HandleFunc(http.MethodPost, version, "/vehicle-models", hdl.create)
	authIsOneOfAdmin.HandleFunc(http.MethodPut, version, "/vehicle-models/{vehicleModelID}", hdl.update, vmCtx)
	authIsOneOfAdmin.HandleFunc(http.MethodDelete, version, "/vehicle-models/{vehicleModelID}", hdl.delete, vmCtx)
}
