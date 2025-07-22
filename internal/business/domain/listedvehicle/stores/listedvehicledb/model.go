package listedvehicledb

import (
	"time"

	"github.com/machilan1/cruise/internal/business/domain/listedvehicle"
	"github.com/machilan1/cruise/internal/business/sdk/sqldb/dbjson"
)

type dbListedVehicle struct {
	ID                 int                                     `db:"listed_vehicle_id"`
	ManufacturedAt     time.Time                               `db:"manufactured_at"`
	LicensedAt         *time.Time                              `db:"licensed_at"`
	Brand              dbjson.JSONColumn[dbListedVehicleBrand] `db:"brand"`
	Model              dbjson.JSONColumn[dbListedVehicleModel] `db:"model"`
	ModelID            int                                     `db:"model_id"`
	DoorCount          int                                     `db:"door_count"`
	KeyCount           int                                     `db:"key_count"`
	Color              string                                  `db:"color"`
	FuelType           string                                  `db:"fuel_type"`
	BodySerial         *string                                 `db:"body_serial"`
	TransmissionType   *string                                 `db:"transmission_type"`
	HeadlightType      *string                                 `db:"headlight_type"`
	WheelSide          string                                  `db:"wheel_side"`
	EngineDisplacement *int                                    `db:"engine_displacement"`
	ValvesCount        *int                                    `db:"valves_count"`
	EngineSerial       *string                                 `db:"engine_serial"`
	EngineType         *string                                 `db:"engine_type"`
	HasTurbo           *bool                                   `db:"has_turbo"`
	BodyModified       bool                                    `db:"body_modified"`
	VehicleSource      *string                                 `db:"vehicle_source"`
	SpecialIncident    *string                                 `db:"special_incident"`
	Note               string                                 `db:"note"`
	CreatedAt          time.Time                               `db:"created_at"`
	UpdatedAt          time.Time                               `db:"updated_at"`
	DeletedAt          *time.Time                              `db:"deleted_at"`
}

func toDBListedVehicle(lvh listedvehicle.ListedVehicle) dbListedVehicle {
	return dbListedVehicle{
		ID:                 lvh.ID,
		ManufacturedAt:     lvh.ManufacturedAt,
		LicensedAt:         lvh.LicensedAt,
		ModelID:            lvh.Model.ID,
		DoorCount:          lvh.DoorCount,
		KeyCount:           lvh.KeyCount,
		Color:              lvh.Color,
		FuelType:           string(lvh.FuelType),
		BodySerial:         lvh.BodySerial,
		TransmissionType:   (*string)(lvh.TransmissionType),
		HeadlightType:      (*string)(lvh.HeadlightType),
		WheelSide:          string(lvh.WheelSide),
		EngineDisplacement: lvh.EngineDisplacement,
		ValvesCount:        lvh.ValvesCount,
		EngineSerial:       lvh.EngineSerial,
		EngineType:         (*string)(lvh.EngineType),
		HasTurbo:           lvh.HasTurbo,
		BodyModified:       lvh.BodyModified,
		VehicleSource:      (*string)(lvh.VehicleSource),
		SpecialIncident:    (*string)(lvh.SpecialIncident),
		Note:               lvh.Note,
		CreatedAt:          lvh.CreatedAt,
		UpdatedAt:          lvh.UpdatedAt,
		DeletedAt:          lvh.DeletedAt,
	}
}

type dbListedVehicleBrand struct {
	ID        int     `json:"brand_id"`
	Name      string  `json:"brand_name"`
	LogoImage *string `json:"logo_image"`
}

type dbListedVehicleModel struct {
	ID               int    `json:"vehicle_model_id"`
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
