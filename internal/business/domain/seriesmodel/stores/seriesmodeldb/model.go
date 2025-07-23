package seriesmodeldb

import (
	"time"

	"github.com/machilan1/cruise/internal/business/domain/seriesmodel"
	"github.com/machilan1/cruise/internal/business/sdk/sqldb/dbjson"
)

type dbSeriesModel struct {
	ID                 int                                    `db:"series_model_id"`
	Name               string                                 `db:"series_model_name"`
	Version            string                                 `db:"version"`
	ModelYear          int                                    `db:"model_year"`
	BodyStyle          string                                 `db:"body_style"`
	DriveType          string                                 `db:"drive_type"`
	FuelType           string                                 `db:"fuel_type"`
	EngineType         string                                 `db:"engine_type"`
	EngineDisplacement int                                    `db:"engine_displacement"`
	ValveCount         int                                    `db:"valve_count"`
	HasTurbo           bool                                   `db:"has_turbo"`
	TransmissionType   string                                 `db:"transmission_type"`
	HorsePower         int                                    `db:"horse_power"`
	Series             dbjson.JSONColumn[dbSeriesModelSeries] `db:"series"`
	SeriesID           int                                    `db:"series_id"`
	CreatedAt          time.Time                              `db:"created_at"`
	UpdatedAt          time.Time                              `db:"updated_at"`
}

type dbSeriesModelSeries struct {
	ID             int    `json:"brand_series_id"`
	Name           string `json:"brand_series_name"`
	BrandID        int    `json:"brand_id"`
	BrandName      string `json:"brand_name"`
	BrandLogoImage string `json:"brand_logo_image"`
}

func toCoreSeriesModel(dsm dbSeriesModel) (seriesmodel.SeriesModel, error) {
	bodyStyle, err := seriesmodel.ParseBodyStyle(dsm.BodyStyle)
	if err != nil {
		return seriesmodel.SeriesModel{}, err
	}
	driveType, err := seriesmodel.ParseDriveType(dsm.DriveType)
	if err != nil {
		return seriesmodel.SeriesModel{}, err
	}
	fuelType, err := seriesmodel.ParseFuelType(dsm.FuelType)
	if err != nil {
		return seriesmodel.SeriesModel{}, err
	}
	engineType, err := seriesmodel.ParseEngineType(dsm.EngineType)
	if err != nil {
		return seriesmodel.SeriesModel{}, err
	}
	tranmissionType, err := seriesmodel.ParseTransmissionType(dsm.TransmissionType)
	if err != nil {
		return seriesmodel.SeriesModel{}, err
	}

	return seriesmodel.SeriesModel{
		ID:                 dsm.ID,
		Name:               dsm.Name,
		Version:            dsm.Version,
		ModelYear:          dsm.ModelYear,
		BodyStyle:          bodyStyle,
		DriveType:          driveType,
		FuelType:           fuelType,
		EngineType:         engineType,
		EngineDisplacement: dsm.EngineDisplacement,
		ValveCount:         dsm.ValveCount,
		HasTurbo:           dsm.HasTurbo,
		TransmissionType:   tranmissionType,
		HorsePower:         dsm.HorsePower,
		Series: seriesmodel.SeriesModelSeries{
			ID:             dsm.Series.Get().ID,
			Name:           dsm.Series.Get().Name,
			BrandID:        dsm.Series.Get().BrandID,
			BrandName:      dsm.Series.Get().BrandName,
			BrandLogoImage: dsm.Series.Get().BrandLogoImage,
		},
		CreatedAt: dsm.CreatedAt,
		UpdatedAt: dsm.UpdatedAt,
	}, nil
}

func toCoreSeriesModels(dsms []dbSeriesModel) ([]seriesmodel.SeriesModel, error) {
	sms := make([]seriesmodel.SeriesModel, len(dsms))
	for i, v := range dsms {
		sm, err := toCoreSeriesModel(v)
		if err != nil {
			return nil, err
		}

		sms[i] = sm
	}
	return sms, nil
}

func toDBSeriesModels(sm seriesmodel.SeriesModel) dbSeriesModel {
	return dbSeriesModel{
		ID:                 sm.ID,
		Name:               sm.Name,
		Version:            sm.Version,
		ModelYear:          sm.ModelYear,
		BodyStyle:          string(sm.BodyStyle),
		DriveType:          string(sm.DriveType),
		FuelType:           string(sm.FuelType),
		EngineType:         string(sm.EngineType),
		EngineDisplacement: sm.EngineDisplacement,
		ValveCount:         sm.ValveCount,
		HasTurbo:           sm.HasTurbo,
		TransmissionType:   string(sm.TransmissionType),
		HorsePower:         sm.HorsePower,
		SeriesID:           sm.Series.ID,
		CreatedAt:          sm.CreatedAt,
		UpdatedAt:          sm.UpdatedAt,
	}
}
