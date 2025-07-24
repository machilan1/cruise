package seriesvariant

import "fmt"

type DriveType string

const (
	DriveTypeFWD         DriveType = "fwd"
	DriveTypeRWD         DriveType = "rwd"
	DriveTypeAWD         DriveType = "awd"
	DriveType4WD         DriveType = "4wd"
	DriveTypeUnspecified DriveType = "unspecified"
)

func ParseDriveType(d string) (DriveType, error) {
	switch d {
	case string(DriveTypeFWD):
		return DriveTypeFWD, nil
	case string(DriveTypeRWD):
		return DriveTypeRWD, nil
	case string(DriveTypeAWD):
		return DriveTypeAWD, nil
	case string(DriveType4WD):
		return DriveType4WD, nil
	case string(DriveTypeUnspecified):
		return DriveTypeUnspecified, nil
	default:
		return "", fmt.Errorf("invalid drive type: %s", d)
	}
}

type FuelType string

const (
	FuelTypeGasoline    FuelType = "gasoline"
	FuelTypeDiesel      FuelType = "diesel"
	FuelTypeElectric    FuelType = "electric"
	FuelTypeGas         FuelType = "gas"
	FuelTypeUnspecified FuelType = "unspecified"
)

func ParseFuelType(ft string) (FuelType, error) {
	switch ft {
	case string(FuelTypeGasoline):
		return FuelTypeGasoline, nil
	case string(FuelTypeDiesel):
		return FuelTypeDiesel, nil
	case string(FuelTypeElectric):
		return FuelTypeElectric, nil
	case string(FuelTypeGas):
		return FuelTypeGas, nil
	case string(FuelTypeUnspecified):
		return FuelTypeUnspecified, nil
	default:
		return "", fmt.Errorf("invalid fuel type: %s", ft)
	}
}

type BodyStyle string

const (
	BodyStyleSedan       BodyStyle = "sedan"
	BodyStyleWagon       BodyStyle = "wagon"
	BodyStyleHatchback   BodyStyle = "hatchback"
	BodyStyleGt          BodyStyle = "gt"
	BodyStyleSports      BodyStyle = "sports"
	BodyStyleVan         BodyStyle = "van"
	BodyStyleTruck       BodyStyle = "truck"
	BodyStyleSuv         BodyStyle = "suv"
	BodyStyleConvertible BodyStyle = "convertible"
	BodyStyleUnspecified BodyStyle = "unspecified"
)

func ParseBodyStyle(bs string) (BodyStyle, error) {
	switch bs {
	case string(BodyStyleSedan):
		return BodyStyleSedan, nil
	case string(BodyStyleWagon):
		return BodyStyleWagon, nil
	case string(BodyStyleHatchback):
		return BodyStyleHatchback, nil
	case string(BodyStyleGt):
		return BodyStyleGt, nil
	case string(BodyStyleSports):
		return BodyStyleSports, nil
	case string(BodyStyleVan):
		return BodyStyleVan, nil
	case string(BodyStyleTruck):
		return BodyStyleTruck, nil
	case string(BodyStyleSuv):
		return BodyStyleSuv, nil
	case string(BodyStyleConvertible):
		return BodyStyleConvertible, nil
	case string(BodyStyleUnspecified):
		return BodyStyleUnspecified, nil
	default:
		return "", fmt.Errorf("invalid body style: %s", bs)
	}
}

type TransmissionType string

const (
	TransmissionTypeAuto        TransmissionType = "automatic"
	TransmissionTypeManual      TransmissionType = "manual"
	TransmissionTypeUnspecified TransmissionType = "unspecified"
)

func ParseTransmissionType(tm string) (TransmissionType, error) {
	switch tm {
	case string(TransmissionTypeAuto):
		return TransmissionTypeAuto, nil
	case string(TransmissionTypeManual):
		return TransmissionTypeManual, nil
	case string(TransmissionTypeUnspecified):
		return TransmissionTypeUnspecified, nil
	default:
		return "", fmt.Errorf("invalid transmission type: %s", tm)
	}
}

type EngineType string

const (
	EngineTypeV           EngineType = "v"
	EngineTypeInline      EngineType = "inline"
	EngineTypeBoxer       EngineType = "boxer"
	EngineTypeRotary      EngineType = "rotary"
	EngineTypeUnspecified EngineType = "unspecified"
)

func ParseEngineType(et string) (EngineType, error) {
	switch et {
	case string(EngineTypeV):
		return EngineTypeV, nil
	case string(EngineTypeInline):
		return EngineTypeInline, nil
	case string(EngineTypeBoxer):
		return EngineTypeBoxer, nil
	case string(EngineTypeRotary):
		return EngineTypeRotary, nil
	case string(EngineTypeUnspecified):
		return EngineTypeUnspecified, nil
	default:
		return "", fmt.Errorf("invalid engine type: %s", et)
	}
}
