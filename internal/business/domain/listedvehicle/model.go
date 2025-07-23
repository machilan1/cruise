package listedvehicle

import "time"

type ListedVehicle struct {
	ID               int
	ManufacturedAt   time.Time
	LicensedAt       *time.Time
	Model            ListedVehicleModel
	DoorCount        int
	KeyCount         int
	Color            string
	BodySerial       string
	TransmissionType TransmissionType
	HeadlightType    HeadlightType
	WheelSide        WheelSideType
	EngineSerial     string
	BodyModified     bool
	VehicleSource    VehicleSourceType
	SpecialIncident  SpecialIncidentType
	Note             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

type ListedVehicleModel struct {
	ID               int
	BrandID          int
	BrandName        string
	BrandLogoImage   string
	SeriesName       string
	CommercialName   string
	ModelYear        int
	Nickname         string
	DriveType        DriveType
	FuelType         FuelType
	BodyStyle        BodyStyle
	TransmissionType TransmissionType
	EngineType       EngineType
}

type NewListedVehicle struct {
	ManufacturedAt   time.Time
	LicensedAt       *time.Time
	ModelID          int
	DoorCount        int
	KeyCount         int
	Color            string
	BodySerial       string
	TransmissionType TransmissionType
	HeadlightType    HeadlightType
	WheelSide        WheelSideType
	EngineSerial     string
	BodyModified     bool
	VehicleSource    VehicleSourceType
	SpecialIncident  SpecialIncidentType
	Note             string
}

type UpdateListedVehicle struct {
	ManufacturedAt   *time.Time
	LicensedAt       *time.Time
	ModelID          *int
	DoorCount        *int
	KeyCount         *int
	Color            *string
	BodySerial       *string
	TransmissionType *TransmissionType
	HeadlightType    *HeadlightType
	WheelSide        *WheelSideType
	EngineSerial     *string
	BodyModified     *bool
	VehicleSource    *VehicleSourceType
	SpecialIncident  *SpecialIncidentType
	Note             *string
}
