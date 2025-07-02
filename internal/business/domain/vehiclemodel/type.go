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
	BodyStyleSedan     = "sedan"
	BodyStyleWagon     = "wagon"
	BodyStyleHatchback = "hatchback"
	BodyStyleGt        = "gt"
	BodyStyleSports    = "sports"
	BodyStyleVan       = "van"
	BodyStyleTruck     = "truck"
	BodyStyleSuv       = "suv"
	BodyStyleOthers    = "others"
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
	TransmissionTypeAuto   = "automatic"
	TransmissionTypeManual = "manual"
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
