package vehiclemodel

type DriveType string

const (
	DriveTypeFWD    DriveType = "fwd"
	DriveTypeRWD    DriveType = "rwd"
	DriveTypeAWD    DriveType = "awd"
	DriveType4WD    DriveType = "4wd"
	DriveTypeOthers DriveType = "others"
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
	case string(DriveTypeOthers):
		return DriveTypeOthers, nil
	default:
		return "", ErrInvalidDriveType
	}
}

type FuelType string

const (
	FuelTypeGasoline FuelType = "gasoline"
	FuelTypeDiesel   FuelType = "diesel"
	FuelTypeElectric FuelType = "electric"
	FuelTypeGas      FuelType = "gas"
	FuelTypeOthers   FuelType = "others"
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
	case string(FuelTypeOthers):
		return FuelTypeOthers, nil
	default:
		return "", ErrInvalidFuelType
	}
}

type BodyStyle string

const (
	BodyStyleSedan     BodyStyle = "sedan"
	BodyStyleWagon     BodyStyle = "wagon"
	BodyStyleHatchback BodyStyle = "hatchback"
	BodyStyleGt        BodyStyle = "gt"
	BodyStyleSports    BodyStyle = "sports"
	BodyStyleVan       BodyStyle = "van"
	BodyStyleTruck     BodyStyle = "truck"
	BodyStyleSuv       BodyStyle = "suv"
	BodyStyleOthers    BodyStyle = "others"
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
	case string(BodyStyleOthers):
		return BodyStyleOthers, nil
	default:
		return "", ErrInvalidBodyStyle
	}
}

type TransmissionType string

const (
	TransmissionTypeAuto   TransmissionType = "automatic"
	TransmissionTypeManual TransmissionType = "manual"
)

func ParseTransmissionType(tm string) (TransmissionType, error) {
	switch tm {
	case string(TransmissionTypeAuto):
		return TransmissionTypeAuto, nil
	case string(TransmissionTypeManual):
		return TransmissionTypeManual, nil
	default:
		return "", ErrInvalidTransmissionType
	}
}

type EngineType string

const (
	EngineTypeV      EngineType = "v"
	EngineTypeInline EngineType = "inline"
	EngineTypeBoxer  EngineType = "boxer"
	EngineTypeRotary EngineType = "rotary"
	EngineTypeOthers EngineType = "others"
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
	case string(EngineTypeOthers):
		return EngineTypeOthers, nil
	default:
		return "", ErrInvalidEngineType
	}
}
