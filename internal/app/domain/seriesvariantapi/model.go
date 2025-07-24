package seriesvariantapi

import (
	"time"

	"github.com/machilan1/cruise/internal/business/domain/seriesvariant"
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

func toAppSeriesModel(sm seriesvariant.SeriesVariant) AppSeriesModel {
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

func toAppSeriesModels(sms []seriesvariant.SeriesVariant) []AppSeriesModel {
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

func toCoreNewSeriesModel(ansm AppNewSeriesModel) (seriesvariant.NewSeriesVariant, error) {
	bodyStyle, err := seriesvariant.ParseBodyStyle(ansm.BodyStyle)
	if err != nil {
		return seriesvariant.NewSeriesVariant{}, err
	}
	driveType, err := seriesvariant.ParseDriveType(ansm.DriveType)
	if err != nil {
		return seriesvariant.NewSeriesVariant{}, err
	}
	fuelType, err := seriesvariant.ParseFuelType(ansm.FuelType)
	if err != nil {
		return seriesvariant.NewSeriesVariant{}, err
	}
	engineType, err := seriesvariant.ParseEngineType(ansm.EngineType)
	if err != nil {
		return seriesvariant.NewSeriesVariant{}, err
	}
	transmissionType, err := seriesvariant.ParseTransmissionType(ansm.TransmissionType)
	if err != nil {
		return seriesvariant.NewSeriesVariant{}, err
	}

	return seriesvariant.NewSeriesVariant{
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

func toCoreUpdateSeriesModel(ausm AppUpdateSeriesModel) (seriesvariant.UpdateSeriesVariant, error) {
	var bodyStyle *seriesvariant.BodyStyle
	var driveType *seriesvariant.DriveType
	var fuelType *seriesvariant.FuelType
	var engineType *seriesvariant.EngineType
	var transmissionType *seriesvariant.TransmissionType

	if ausm.BodyStyle != nil {
		bs, err := seriesvariant.ParseBodyStyle(*ausm.BodyStyle)
		if err != nil {
			return seriesvariant.UpdateSeriesVariant{}, err
		}
		bodyStyle = &bs
	}

	if ausm.DriveType != nil {
		dt, err := seriesvariant.ParseDriveType(*ausm.DriveType)
		if err != nil {
			return seriesvariant.UpdateSeriesVariant{}, err
		}
		driveType = &dt
	}

	if ausm.FuelType != nil {
		ft, err := seriesvariant.ParseFuelType(*ausm.FuelType)
		if err != nil {
			return seriesvariant.UpdateSeriesVariant{}, err
		}
		fuelType = &ft
	}

	if ausm.EngineType != nil {
		et, err := seriesvariant.ParseEngineType(*ausm.EngineType)
		if err != nil {
			return seriesvariant.UpdateSeriesVariant{}, err
		}
		engineType = &et
	}

	if ausm.TransmissionType != nil {
		tt, err := seriesvariant.ParseTransmissionType(*ausm.TransmissionType)
		if err != nil {
			return seriesvariant.UpdateSeriesVariant{}, err
		}
		transmissionType = &tt
	}

	return seriesvariant.UpdateSeriesVariant{
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
