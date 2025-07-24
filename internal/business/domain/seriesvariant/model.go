package seriesvariant

import "time"

type SeriesVariant struct {
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
	Series             SeriesVariantSeries
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type SeriesVariantSeries struct {
	ID             int
	Name           string
	BrandID        int
	BrandName      string
	BrandLogoImage string
}

type NewSeriesVariant struct {
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

type UpdateSeriesVariant struct {
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
