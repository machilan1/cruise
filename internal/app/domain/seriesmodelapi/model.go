package seriesmodelapi

import (
	"time"

	"github.com/machilan1/cruise/internal/business/domain/seriesmodel"
	"github.com/machilan1/cruise/internal/framework/validate"
)

type AppSeriesModel struct {
	ID                 int                  `json:"id"`
	Name               string               `json:"name"`
	Version            string               `json:"version"`
	ModelYear          int                  `json:"modelYear"`
	BodyStyle          string               `json:"bodyStyle"`
	DriveType          string               `json:"driveType"`
	FuelType           string               `json:"fuelType"`
	EngineType         string               `json:"engineType"`
	EngineDisplacement int                  `json:"engineDisplacement"`
	ValveCount         int                  `json:"valveCount"`
	HasTurbo           bool                 `json:"hasTurbo"`
	TransmissionType   string               `json:"transmissionType"`
	HorsePower         int                  `json:"horsePower"`
	Series             AppSeriesModelSeries `json:"series"`
	CreatedAt          time.Time            `json:"createdAt"`
	UpdatedAt          time.Time            `json:"updatedAt"`
}

type AppSeriesModelSeries struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BrandID        int    `json:"brandId"`
	BrandName      string `json:"brandName"`
	BrandLogoImage string `json:"brandLogoImage"`
}

func toAppSeriesModel(sm seriesmodel.SeriesModel) AppSeriesModel {
	return AppSeriesModel{
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
		Series: AppSeriesModelSeries{
			ID:             sm.Series.ID,
			Name:           sm.Series.Name,
			BrandID:        sm.Series.BrandID,
			BrandName:      sm.Series.BrandName,
			BrandLogoImage: sm.Series.BrandLogoImage,
		},
		CreatedAt: sm.CreatedAt,
		UpdatedAt: sm.UpdatedAt,
	}
}

func toAppSeriesModels(sms []seriesmodel.SeriesModel) []AppSeriesModel {
	asms := make([]AppSeriesModel, len(sms))
	for i, v := range sms {
		asms[i] = toAppSeriesModel(v)
	}
	return asms
}

type AppNewSeriesModel struct {
	Name               string  `json:"name" validate:"required"`
	Version            *string `json:"version" validate:"required"`
	ModelYear          int     `json:"modelYear" validate:"required"`
	BodyStyle          string  `json:"bodyStyle" validate:"oneof=sedan wagon hatchback gt sports van truck suv convertible unspecified"`
	DriveType          string  `json:"driveType" validate:"oneof=fwd rwd 4wd awd unspecified"`
	FuelType           string  `json:"fuelType" validate:"oneof=diesel gasoline electric gas unspecified"`
	EngineType         string  `json:"engineType" validate:"oneof=v inline boxer rotary unspecified"`
	EngineDisplacement int     `json:"engineDisplacement" validate:"required"`
	ValveCount         int     `json:"valveCount" validate:"required"`
	HasTurbo           bool    `json:"hasTurbo" validate:"required"`
	TransmissionType   string  `json:"transmissionType" validate:"oneof=automatic manual unspecified"`
	HorsePower         int     `json:"horsePower" validate:"required"`
	SeriesID           int     `json:"seriesId" validate:"required"`
}

func (app AppNewSeriesModel) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}
	return nil
}

func toCoreNewSeriesModel(ansm AppNewSeriesModel) (seriesmodel.NewSeriesModel, error) {
	bodyStyle, err := seriesmodel.ParseBodyStyle(ansm.BodyStyle)
	if err != nil {
		return seriesmodel.NewSeriesModel{}, ErrInvalidBodyStyle
	}
	driveType, err := seriesmodel.ParseDriveType(ansm.DriveType)
	if err != nil {
		return seriesmodel.NewSeriesModel{}, ErrInvalidDriveType
	}
	fuelType, err := seriesmodel.ParseFuelType(ansm.FuelType)
	if err != nil {
		return seriesmodel.NewSeriesModel{}, ErrInvalidFuelType
	}
	engineType, err := seriesmodel.ParseEngineType(ansm.EngineType)
	if err != nil {
		return seriesmodel.NewSeriesModel{}, ErrInvalidEngineType
	}
	transmissionType, err := seriesmodel.ParseTransmissionType(ansm.TransmissionType)
	if err != nil {
		return seriesmodel.NewSeriesModel{}, ErrInvalidTransmissionType
	}

	return seriesmodel.NewSeriesModel{
		Name:               ansm.Name,
		Version:            ansm.Version,
		ModelYear:          ansm.ModelYear,
		BodyStyle:          bodyStyle,
		DriveType:          driveType,
		FuelType:           fuelType,
		EngineType:         engineType,
		EngineDisplacement: ansm.EngineDisplacement,
		ValveCount:         ansm.ValveCount,
		HasTurbo:           ansm.HasTurbo,
		TransmissionType:   transmissionType,
		HorsePower:         ansm.HorsePower,
		SeriesID:           ansm.SeriesID,
	}, nil
}

type AppUpdateSeriesModel struct {
	Name               *string `json:"name"`
	Version            *string `json:"version"`
	ModelYear          *int    `json:"modelYear"`
	BodyStyle          *string `json:"bodyStyle" validate:"omitempty,oneof=sedan wagon hatchback gt sports van truck suv convertible unspecified"`
	DriveType          *string `json:"driveType" validate:"omitempty,oneof=fwd rwd 4wd awd unspecified"`
	FuelType           *string `json:"fuelType" validate:"omitempty,oneof=diesel gasoline electric gas unspecified"`
	EngineType         *string `json:"engineType" validate:"omitempty,oneof=v inline boxer rotary unspecified"`
	EngineDisplacement *int    `json:"engineDisplacement"`
	ValveCount         *int    `json:"valveCount"`
	HasTurbo           *bool   `json:"hasTurbo"`
	TransmissionType   *string `json:"transmissionType" validate:"omitempty,oneof=automatic manual unspecified"`
	HorsePower         *int    `json:"horsePower"`
	SeriesID           *int    `json:"seriesId"`
}

func (app AppUpdateSeriesModel) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}
	return nil
}

func toCoreUpdateSeriesModel(ausm AppUpdateSeriesModel) (seriesmodel.UpdateSeriesModel, error) {
	var bodyStyle *seriesmodel.BodyStyle
	var driveType *seriesmodel.DriveType
	var fuelType *seriesmodel.FuelType
	var engineType *seriesmodel.EngineType
	var transmissionType *seriesmodel.TransmissionType

	if ausm.BodyStyle != nil {
		bs, err := seriesmodel.ParseBodyStyle(*ausm.BodyStyle)
		if err != nil {
			return seriesmodel.UpdateSeriesModel{}, ErrInvalidBodyStyle
		}
		bodyStyle = &bs
	}

	if ausm.DriveType != nil {
		dt, err := seriesmodel.ParseDriveType(*ausm.DriveType)
		if err != nil {
			return seriesmodel.UpdateSeriesModel{}, ErrInvalidDriveType
		}
		driveType = &dt
	}

	if ausm.FuelType != nil {
		ft, err := seriesmodel.ParseFuelType(*ausm.FuelType)
		if err != nil {
			return seriesmodel.UpdateSeriesModel{}, ErrInvalidFuelType
		}
		fuelType = &ft
	}

	if ausm.EngineType != nil {
		et, err := seriesmodel.ParseEngineType(*ausm.EngineType)
		if err != nil {
			return seriesmodel.UpdateSeriesModel{}, ErrInvalidEngineType
		}
		engineType = &et
	}

	if ausm.TransmissionType != nil {
		tt, err := seriesmodel.ParseTransmissionType(*ausm.TransmissionType)
		if err != nil {
			return seriesmodel.UpdateSeriesModel{}, ErrInvalidTransmissionType
		}
		transmissionType = &tt
	}

	return seriesmodel.UpdateSeriesModel{
		Name:               ausm.Name,
		Version:            ausm.Version,
		ModelYear:          ausm.ModelYear,
		BodyStyle:          bodyStyle,
		DriveType:          driveType,
		FuelType:           fuelType,
		EngineType:         engineType,
		EngineDisplacement: ausm.EngineDisplacement,
		ValveCount:         ausm.ValveCount,
		HasTurbo:           ausm.HasTurbo,
		TransmissionType:   transmissionType,
		HorsePower:         ausm.HorsePower,
		SeriesID:           ausm.SeriesID,
	}, nil
}
