package healthapi

import (
	"net/http"

	"github.com/machilan1/cruise/internal/business/sdk/sqldb"
	"github.com/machilan1/cruise/internal/framework/logger"
	"github.com/machilan1/cruise/internal/framework/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log *logger.Logger
	DB  *sqldb.DB
}

func Routes(app *web.App, cfg Config) {
	const version = "v1"

	api := newHandlers(cfg.Log, cfg.DB)

	app.HandleFuncNoMid(http.MethodGet, version, "/healthz", api.readiness)
}
