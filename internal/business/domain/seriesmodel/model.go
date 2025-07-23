package seriesmodel

import "time"

type SeriesModel struct {
	ID                 int
	Name               string
	Version            string
	ModelYear          int
	BodyStyle          BodyStyle
	DriveType          DriveType
	FuelType           FuelType
	EngineType         EngineType
	EngineDisplacement int
	ValveCount         int
	HasTurbo           bool
	TransmissionType   TransmissionType
	HorsePower         int
	Series             SeriesModelSeries
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type SeriesModelSeries struct {
	ID             int
	Name           string
	BrandID        int
	BrandName      string
	BrandLogoImage string
}

type NewSeriesModel struct {
	Name               string
	Version            *string
	ModelYear          int
	BodyStyle          BodyStyle
	DriveType          DriveType
	FuelType           FuelType
	EngineType         EngineType
	EngineDisplacement int
	ValveCount         int
	HasTurbo           bool
	TransmissionType   TransmissionType
	HorsePower         int
	SeriesID           int
}

type UpdateSeriesModel struct {
	Name               *string
	Version            *string
	ModelYear          *int
	BodyStyle          *BodyStyle
	DriveType          *DriveType
	FuelType           *FuelType
	EngineType         *EngineType
	EngineDisplacement *int
	ValveCount         *int
	HasTurbo           *bool
	TransmissionType   *TransmissionType
	HorsePower         *int
	SeriesID           *int
}
