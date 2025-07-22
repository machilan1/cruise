package listedvehicle

import "errors"

var (
	ErrNotfound                   = errors.New("listed vehicle model not found")
	ErrConflict                   = errors.New("provided input conflicts with existing data")
	ErrInvalidDriveType           = errors.New("invalid drive type")
	ErrInvalidFuelType            = errors.New("invalid fuel type")
	ErrInvalidBodyStyle           = errors.New("invalid body style")
	ErrInvalidTransmissionType    = errors.New("invalid transmission style")
	ErrInvalidEngineType          = errors.New("invalid engine type")
	ErrInvalidHeadlightType       = errors.New("invalid headlight type")
	ErrInvalidWheelSideType       = errors.New("invalid wheel side type")
	ErrInvalidVehicleSourceType   = errors.New("invalid vehicle source type")
	ErrInvalidSpecialIncidentType = errors.New("invalid special incident type")
)
