package listedvehicledb

import (
	"time"

	"github.com/machilan1/cruise/internal/business/sdk/sqldb/dbjson"
)

type dbListedVehicle struct {
	ID               int                                     `db:"listed_vehicle_id"`
	ManufacturedAt   time.Time                               `db:"manufactured_at"`
	LicensedAt       *time.Time                              `db:"licensed_at"`
	Model            dbjson.JSONColumn[dbListedVehicleModel] `db:"model"`
	ModelID          int                                     `db:"model_id"`
	DoorCount        int                                     `db:"door_count"`
	KeyCount         int                                     `db:"key_count"`
	Color            string                                  `db:"color"`
	FuelType         string                                  `db:"fuel_type"`
	BodySerial       string                                  `db:"body_serial"`
	TransmissionType string                                  `db:"transmission_type"`
	HeadlightType    string                                  `db:"headlight_type"`
	WheelSide        string                                  `db:"wheel_side"`
	EngineSerial     *string                                 `db:"engine_serial"`
	BodyModified     bool                                    `db:"body_modified"`
	VehicleSource    string                                  `db:"vehicle_source"`
	SpecialIncident  string                                  `db:"special_incident"`
	Note             string                                  `db:"note"`
	CreatedAt        time.Time                               `db:"created_at"`
	UpdatedAt        time.Time                               `db:"updated_at"`
	DeletedAt        *time.Time                              `db:"deleted_at"`
	// EngineDisplacement *int                                    `db:"engine_displacement"`
	// ValvesCount        *int                                    `db:"valves_count"`
	// EngineType         *string                                 `db:"engine_type"`
}

type dbListedVehicleBrand struct {
	ID        int     `json:"brand_id"`
	Name      string  `json:"brand_name"`
	LogoImage *string `json:"logo_image"`
}

type dbListedVehicleModel struct {
	ID               int    `json:"vehicle_model_id"`
	BrandID          int    `json:"brand_id"`
	BrandName        string `json:"brand_name"`
	BrandLogoImage   string `json:"brand_logo"`
	SeriesName       string `json:"model_series_name"`
	CommercialName   string `json:"model_commercial_name"`
	ModelYear        int    `json:"model_year"`
	Nickname         string `json:"nickname"`
	DriveType        string `json:"drive_type"`
	FuelType         string `json:"fuel_type"`
	BodyStyle        string `json:"body_style"`
	TransmissionType string `json:"transmission_type"`
	EngineType       string `json:"engine_type"`
}
