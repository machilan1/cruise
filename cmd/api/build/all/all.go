package all

import (
	"github.com/machilan1/cruise/internal/app/domain/auctionapi"
	"github.com/machilan1/cruise/internal/app/domain/auctionhouseapi"
	"github.com/machilan1/cruise/internal/app/domain/authapi"
	"github.com/machilan1/cruise/internal/app/domain/brandapi"
	"github.com/machilan1/cruise/internal/app/domain/brandseriesapi"
	"github.com/machilan1/cruise/internal/app/domain/fileapi"
	"github.com/machilan1/cruise/internal/app/domain/healthapi"
	"github.com/machilan1/cruise/internal/app/domain/seriesvariantapi"
	"github.com/machilan1/cruise/internal/app/sdk/mux"
	"github.com/machilan1/cruise/internal/business/domain/auction"
	"github.com/machilan1/cruise/internal/business/domain/auction/stores/auctiondb"
	"github.com/machilan1/cruise/internal/business/domain/auctionhouse"
	"github.com/machilan1/cruise/internal/business/domain/auctionhouse/stores/auctionhousedb"
	"github.com/machilan1/cruise/internal/business/domain/auth"
	"github.com/machilan1/cruise/internal/business/domain/auth/stores/authdb"
	"github.com/machilan1/cruise/internal/business/domain/brand"
	"github.com/machilan1/cruise/internal/business/domain/brand/stores/branddb"
	"github.com/machilan1/cruise/internal/business/domain/brandseries"
	"github.com/machilan1/cruise/internal/business/domain/brandseries/stores/brandseriesdb"
	"github.com/machilan1/cruise/internal/business/domain/file"
	"github.com/machilan1/cruise/internal/business/domain/file/stores/filedb"
	"github.com/machilan1/cruise/internal/business/domain/notification"
	"github.com/machilan1/cruise/internal/business/domain/seriesvariant"
	"github.com/machilan1/cruise/internal/business/domain/seriesvariant/stores/seriesvariantdb"
	"github.com/machilan1/cruise/internal/framework/web"
)

func Routes() add { // nolint: revive
	return add{}
}

type add struct{}

func (add) Add(app *web.App, cfg mux.Config) {
	fileCore := file.NewCore(filedb.NewStore(cfg.DB))
	authCore := auth.NewCore(authdb.NewStore(cfg.DB), cfg.JWTKey)
	notifyCore := notification.NewCore(cfg.Mailer, cfg.FrontendOrigin)
	brandCore := brand.NewCore(branddb.NewStore(cfg.DB))
	brandSeriesCore := brandseries.NewCore(brandseriesdb.NewStore(cfg.DB))
	seriesVariantCore := seriesvariant.NewCore(seriesvariantdb.NewStore(cfg.DB))
	// vehicleModelCore := vehiclemodel.NewCore(vehiclemodeldb.NewStore(cfg.DB))
	auctionHouseCore := auctionhouse.NewCore(auctionhousedb.NewStore(cfg.DB))
	auctionCore := auction.NewCore(auctiondb.NewStore(cfg.DB))

	healthapi.Routes(app, healthapi.Config{
		Log: cfg.Log,
		DB:  cfg.DB,
	})
	fileapi.Routes(app, fileapi.Config{
		Log:     cfg.Log,
		TxM:     cfg.TxM,
		Storage: cfg.Storage,
		File:    fileCore,
		Sess:    cfg.Sess,
		Auth:    authCore,
	})

	authapi.Routes(app, authapi.Config{
		Log:          cfg.Log,
		TxM:          cfg.TxM,
		Sess:         cfg.Sess,
		Auth:         authCore,
		Notification: notifyCore,
	})

	brandapi.Routes(app, brandapi.Config{
		Log:   cfg.Log,
		TxM:   cfg.TxM,
		Sess:  cfg.Sess,
		Auth:  authCore,
		Brand: brandCore,
	})

	brandseriesapi.Routes(app, brandseriesapi.Config{
		Log:         cfg.Log,
		TxM:         cfg.TxM,
		Sess:        cfg.Sess,
		Auth:        authCore,
		BrandSeries: brandSeriesCore,
	})

	seriesvariantapi.Routes(app, seriesvariantapi.Config{
		Log:           cfg.Log,
		TxM:           cfg.TxM,
		Sess:          cfg.Sess,
		Auth:          authCore,
		SeriesVariant: seriesVariantCore,
	})

	// vehiclemodelapi.Routes(app, vehiclemodelapi.Config{
	// 	Log:          cfg.Log,
	// 	TxM:          cfg.TxM,
	// 	Sess:         cfg.Sess,
	// 	Auth:         authCore,
	// 	VehicleModel: vehicleModelCore,
	// })

	auctionhouseapi.Routes(app, auctionhouseapi.Config{
		Log:          cfg.Log,
		TxM:          cfg.TxM,
		Sess:         cfg.Sess,
		Auth:         authCore,
		AuctionHouse: auctionHouseCore,
	})

	auctionapi.Routes(app, auctionapi.Config{
		Log:     cfg.Log,
		TxM:     cfg.TxM,
		Sess:    cfg.Sess,
		Auth:    authCore,
		Auction: auctionCore,
	})
}
