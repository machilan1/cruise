package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/gorilla/sessions"
	"github.com/machilan1/cruise/cmd/api/build/all"
	"github.com/machilan1/cruise/internal/app/sdk/debug"
	"github.com/machilan1/cruise/internal/app/sdk/mux"
	"github.com/machilan1/cruise/internal/business/sdk/blobstore/gcs"
	"github.com/machilan1/cruise/internal/business/sdk/mailer"
	"github.com/machilan1/cruise/internal/business/sdk/mailer/sendgrid"
	"github.com/machilan1/cruise/internal/business/sdk/mailer/smtp"
	"github.com/machilan1/cruise/internal/business/sdk/sess"
	"github.com/machilan1/cruise/internal/business/sdk/sess/stores/sessdb"
	"github.com/machilan1/cruise/internal/business/sdk/sqldb"
	"github.com/machilan1/cruise/internal/business/sdk/tran"
	"github.com/machilan1/cruise/internal/framework/cors"
	"github.com/machilan1/cruise/internal/framework/logger"
	"github.com/machilan1/cruise/internal/framework/web"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var build = "develop"

func main() {
	var log *logger.Logger

	traceIDFn := func(ctx context.Context) string {
		return web.GetTraceID(ctx)
	}

	log = logger.New(os.Stdout, logger.LevelDebug, "API", traceIDFn)

	// -------------------------------------------------------------------------

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:20s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s"`
			APIHost            string        `conf:"default:0.0.0.0:3000"`
			DebugHost          string        `conf:"default:0.0.0.0:3010"`
			CORSAllowedOrigins []string      `conf:"default:*"`
		}
		DB struct {
			User         string `conf:"default:postgres"`
			Password     string `conf:"default:postgres,mask"`
			Host         string `conf:"default:localhost:5432"`
			Name         string `conf:"default:postgres"`
			MaxIdleConns int    `conf:"default:0"`
			MaxOpenConns int    `conf:"default:3"`
			DisableTLS   bool   `conf:"default:false"`
			TimeZone     string `conf:"default:Asia/Taipei"`
			// Only set if using CloudSQL. When set, the CloudSQL connector will be used, and `Password`, `Host`, `DisableTLS`
			// will be ignored.
			// Example: "project:region:instance"
			CloudSQLConnectionName string
		}
		Session struct {
			SecretKey  string `conf:"default:secret-key,mask"`
			CleanupMin int    `conf:"default:5"` // not used for now
			MaxAge     int    `conf:"default:0"`
			Secure     bool   `conf:"default:true"`
			SameSite   string `conf:"default:Lax"`
		}
		Storage struct {
			Bucket string `conf:"default:demo-bucket"`
		}
		JWT struct {
			Key string `conf:"default:secret-key,mask"`
		}
		SMTP struct {
			Host     string `conf:"default:localhost"`
			Port     int    `conf:"default:1025"`
			User     string `conf:"default:user"`
			Password string `conf:"default:password,mask"`
		}
		Sendgrid struct {
			// If this is set, the SMTP configuration is ignored.
			APIKey string `conf:"mask"`
		}
		Mail struct {
			FromName    string `conf:"default:go-tool"`
			FromAddress string `conf:"default:noreply@gotool.com.tw"`
		}
		Frontend struct {
			// Origin is the origin of the frontend service, should contain the protocol, host and port.
			Origin string `conf:"default:http://localhost:4200"`
		}
	}{
		Version: conf.Version{
			Build: build,
		},
	}

	const prefix = ""
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// -------------------------------------------------------------------------
	// App Starting

	log.Info(ctx, "starting service", "version", cfg.Build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "startup", "config", out)

	expvar.NewString("build").Set(cfg.Build)

	// -------------------------------------------------------------------------
	// Main database Support

	log.Info(ctx, "startup", "status", "initializing main database support", "hostport", cfg.DB.Host)

	db, err := sqldb.Open(log, sqldb.Config{
		User:                     cfg.DB.User,
		Password:                 cfg.DB.Password,
		Host:                     cfg.DB.Host,
		Name:                     cfg.DB.Name,
		MaxIdleConns:             cfg.DB.MaxIdleConns,
		MaxOpenConns:             cfg.DB.MaxOpenConns,
		DisableTLS:               cfg.DB.DisableTLS,
		TimeZone:                 cfg.DB.TimeZone,
		CloudSQLInstanceConnName: cfg.DB.CloudSQLConnectionName,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer db.Close(ctx)

	// -------------------------------------------------------------------------
	// Main Session Core Support

	sessOpts := sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30, // 30 days
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}
	if cfg.Session.MaxAge > 0 {
		sessOpts.MaxAge = cfg.Session.MaxAge
	}
	if cfg.Session.Secure {
		sessOpts.Secure = true
	}
	switch cfg.Session.SameSite {
	case "Lax":
		sessOpts.SameSite = http.SameSiteLaxMode
	case "Strict":
		sessOpts.SameSite = http.SameSiteStrictMode
	case "None":
		sessOpts.SameSite = http.SameSiteNoneMode
	default:
		sessOpts.SameSite = http.SameSiteDefaultMode
	}

	sessCore, err := sess.NewManager(sessdb.NewStore(db), sessOpts, []byte(cfg.Session.SecretKey))
	if err != nil {
		return fmt.Errorf("new session core: %w", err)
	}

	// -------------------------------------------------------------------------
	// Cloud Storage Support

	log.Info(ctx, "startup", "status", "initializing cloud storage support", "bucket", cfg.Storage.Bucket)

	storage, err := gcs.NewGoogleCloudStorage(ctx, gcs.Config{
		Bucket: cfg.Storage.Bucket,
	})
	if err != nil {
		return fmt.Errorf("new storage: %w", err)
	}
	defer storage.Close()

	// -------------------------------------------------------------------------
	// Mailer Support

	var ml mailer.Mailer
	if cfg.Sendgrid.APIKey != "" {
		log.Info(ctx, "startup", "status", "initializing sendgrid support")

		ml = sendgrid.New(log, sendgrid.Config{
			APIKey:      cfg.Sendgrid.APIKey,
			FromName:    cfg.Mail.FromName,
			FromAddress: cfg.Mail.FromAddress,
		})
	} else {
		log.Info(ctx, "startup", "status", "initializing mailer support", "host", cfg.SMTP.Host, "port", cfg.SMTP.Port)

		ml = smtp.New(smtp.Config{
			From: cfg.Mail.FromAddress,
			Host: cfg.SMTP.Host,
			Port: cfg.SMTP.Port,
			User: cfg.SMTP.User,
			Pass: cfg.SMTP.Password,
		})
	}

	// -------------------------------------------------------------------------
	// Start Tracing Support

	log.Info(ctx, "startup", "status", "initializing tracing support")

	// TODO: this is just a mock tracer for now, uncomment the above block when ready to use tempo
	tracer := sdktrace.NewTracerProvider().Tracer(cfg.Build)

	// -------------------------------------------------------------------------
	// Start Debug Service

	go func() {
		log.Info(ctx, "startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)

		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux()); err != nil { // nolint: gosec
			log.Error(ctx, "shutdown", "status", "debug v1 router closed", "host", cfg.Web.DebugHost, "msg", err)
		}
	}()

	// -------------------------------------------------------------------------
	// Start API Service

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	cfgMux := mux.Config{
		Build:          build,
		Shutdown:       shutdown,
		Log:            log,
		DB:             db,
		Tracer:         tracer,
		TxM:            tran.NewTxManager(db),
		Sess:           sessCore,
		Storage:        storage,
		Mailer:         ml,
		JWTKey:         []byte(cfg.JWT.Key),
		FrontendOrigin: cfg.Frontend.Origin,
	}

	corsOpts := cors.Options{
		AllowedOrigins:   cfg.Web.CORSAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
	}

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      cors.Handler(corsOpts)(mux.WebAPI(cfgMux, buildRoutes())),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)

		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "server shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "server shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

func buildRoutes() mux.RouteAdder {
	// The idea here is that we can build different versions of the binary
	// with different sets of exposed web APIs. By default, we build a single
	// an instance with all the web APIs.
	//
	// Here is the scenario. It would be nice to build two binaries, one for the
	// transactional APIs (CRUD) and one for the reporting APIs. This would allow
	// the system to run two instances of the database. One instance tuned for the
	// transactional database calls and the other tuned for the reporting calls.
	// Tuning meaning indexing and memory requirements. The two databases can be
	// kept in sync with replication.

	return all.Routes()
}
