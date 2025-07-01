package mux

import (
	"net/http"
	"os"

	"github.com/machilan1/cruise/internal/app/sdk/mid"
	"github.com/machilan1/cruise/internal/business/sdk/blobstore"
	"github.com/machilan1/cruise/internal/business/sdk/mailer"
	"github.com/machilan1/cruise/internal/business/sdk/sess"
	"github.com/machilan1/cruise/internal/business/sdk/sqldb"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
	"github.com/machilan1/cruise/internal/framework/logger"
	"github.com/machilan1/cruise/internal/framework/web"
	"go.opentelemetry.io/otel/trace"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build          string
	Shutdown       chan os.Signal
	Log            *logger.Logger
	DB             *sqldb.DB
	Tracer         trace.Tracer
	TxM            tran.TxManager
	Sess           *sess.Manager
	Storage        blobstore.BlobStore
	Mailer         mailer.Mailer
	JWTKey         []byte
	FrontendOrigin string
}

// RouteAdder defines behavior that sets the routes to bind for an instance
// of the service.
type RouteAdder interface {
	Add(app *web.App, cfg Config)
}

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(cfg Config, routeAdder RouteAdder) http.Handler {
	app := web.NewApp(
		cfg.Shutdown,
		cfg.Tracer,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Metrics(),
		mid.Panics(),
	)

	routeAdder.Add(app, cfg)

	return app
}
