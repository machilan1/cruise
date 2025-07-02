package vehiclemodel

type VehicleModel struct {
	ID                 int
	SeriesName         string
	CommercialName     string
	ModelYear          int
	BrandID            int
	Nickname           *string
	EngineDisplacement int
	DriveType          DriveType
	FuelType           FuelType
	BodyStyle          BodyStyle
	TransmissionType   TransmissionType
}

type NewVehicleModel struct {
	SeriesName         string
	CommercialName     string
	ModelYear          int
	BrandID            int
	Nickname           *string
	EngineDisplacement int
	DriveType          DriveType
	FuelType           FuelType
	BodyStyle          BodyStyle
	TransmissionType   TransmissionType
}

type UpdateVehicleModel struct {
	SeriesName         *string
	CommercialName     *string
	ModelYear          *int
	Nickname           *string
	EngineDisplacement *int
	DriveType          *DriveType
	FuelType           *FuelType
	BodyStyle          *BodyStyle
	TransmissionType   *TransmissionType
}
