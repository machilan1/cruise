package seriesvariantdb

import (
	"time"

	"github.com/machilan1/cruise/internal/business/domain/seriesvariant"
	"github.com/machilan1/cruise/internal/business/sdk/sqldb/dbjson"
)

type dbSeriesVariant struct {
	ID                 int                                      `db:"series_variant_id"`
	Name               string                                   `db:"series_variant_name"`
	Version            string                                   `db:"version"`
	ModelYear          int                                      `db:"model_year"`
	BodyStyle          string                                   `db:"body_style"`
	DriveType          string                                   `db:"drive_type"`
	FuelType           string                                   `db:"fuel_type"`
	EngineType         string                                   `db:"engine_type"`
	EngineDisplacement int                                      `db:"engine_displacement"`
	ValveCount         int                                      `db:"valve_count"`
	HasTurbo           bool                                     `db:"has_turbo"`
	TransmissionType   string                                   `db:"transmission_type"`
	HorsePower         int                                      `db:"horse_power"`
	Series             dbjson.JSONColumn[dbSeriesVariantSeries] `db:"series"`
	SeriesID           int                                      `db:"series_id"`
	CreatedAt          time.Time                                `db:"created_at"`
	UpdatedAt          time.Time                                `db:"updated_at"`
}

type dbSeriesVariantSeries struct {
	ID             int    `json:"brand_series_id"`
	Name           string `json:"brand_series_name"`
	BrandID        int    `json:"brand_id"`
	BrandName      string `json:"brand_name"`
	BrandLogoImage string `json:"brand_logo_image"`
}

func toCoreSeriesVariant(dsv dbSeriesVariant) (seriesvariant.SeriesVariant, error) {
	bodyStyle, err := seriesvariant.ParseBodyStyle(dsv.BodyStyle)
	if err != nil {
		return seriesvariant.SeriesVariant{}, err
	}
	driveType, err := seriesvariant.ParseDriveType(dsv.DriveType)
	if err != nil {
		return seriesvariant.SeriesVariant{}, err
	}
	fuelType, err := seriesvariant.ParseFuelType(dsv.FuelType)
	if err != nil {
		return seriesvariant.SeriesVariant{}, err
	}
	engineType, err := seriesvariant.ParseEngineType(dsv.EngineType)
	if err != nil {
		return seriesvariant.SeriesVariant{}, err
	}
	tranmissionType, err := seriesvariant.ParseTransmissionType(dsv.TransmissionType)
	if err != nil {
		return seriesvariant.SeriesVariant{}, err
	}

	return seriesvariant.SeriesVariant{
		ID:                 dsv.ID,
		Name:               dsv.Name,
		Version:            dsv.Version,
		ModelYear:          dsv.ModelYear,
		BodyStyle:          bodyStyle,
		DriveType:          driveType,
		FuelType:           fuelType,
		EngineType:         engineType,
		EngineDisplacement: dsv.EngineDisplacement,
		ValveCount:         dsv.ValveCount,
		HasTurbo:           dsv.HasTurbo,
		TransmissionType:   tranmissionType,
		HorsePower:         dsv.HorsePower,
		Series: seriesvariant.SeriesVariantSeries{
			ID:             dsv.Series.Get().ID,
			Name:           dsv.Series.Get().Name,
			BrandID:        dsv.Series.Get().BrandID,
			BrandName:      dsv.Series.Get().BrandName,
			BrandLogoImage: dsv.Series.Get().BrandLogoImage,
		},
		CreatedAt: dsv.CreatedAt,
		UpdatedAt: dsv.UpdatedAt,
	}, nil
}

func toCoreSeriesVariants(dsvs []dbSeriesVariant) ([]seriesvariant.SeriesVariant, error) {
	svs := make([]seriesvariant.SeriesVariant, len(dsvs))
	for i, v := range dsvs {
		sv, err := toCoreSeriesVariant(v)
		if err != nil {
			return nil, err
		}

		svs[i] = sv
	}
	return svs, nil
}

func toDBSeriesVariant(sv seriesvariant.SeriesVariant) dbSeriesVariant {
	return dbSeriesVariant{
		ID:                 sv.ID,
		Name:               sv.Name,
		Version:            sv.Version,
		ModelYear:          sv.ModelYear,
		BodyStyle:          string(sv.BodyStyle),
		DriveType:          string(sv.DriveType),
		FuelType:           string(sv.FuelType),
		EngineType:         string(sv.EngineType),
		EngineDisplacement: sv.EngineDisplacement,
		ValveCount:         sv.ValveCount,
		HasTurbo:           sv.HasTurbo,
		TransmissionType:   string(sv.TransmissionType),
		HorsePower:         sv.HorsePower,
		SeriesID:           sv.Series.ID,
		CreatedAt:          sv.CreatedAt,
		UpdatedAt:          sv.UpdatedAt,
	}
}
