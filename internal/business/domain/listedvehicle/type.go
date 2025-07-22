package listedvehicle

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

type HeadlightType string

const (
	HeadlightTypeHID      HeadlightType = "hid"
	HeadlightTypeLED      HeadlightType = "led"
	HeadlightTypeTungsten HeadlightType = "tungsten"
	HeadlightTypeOthers   HeadlightType = "others"
)

func ParseHeadlightType(et string) (HeadlightType, error) {
	switch et {
	case string(HeadlightTypeHID):
		return HeadlightTypeHID, nil
	case string(HeadlightTypeLED):
		return HeadlightTypeLED, nil
	case string(HeadlightTypeTungsten):
		return HeadlightTypeTungsten, nil
	case string(HeadlightTypeOthers):
		return HeadlightTypeOthers, nil
	default:
		return "", ErrInvalidHeadlightType
	}
}

type WheelSideType string

const (
	WheelSideTypeLeft   WheelSideType = "left"
	WheelSideTypeRight  WheelSideType = "right"
	WheelSideTypeOthers WheelSideType = "others"
)

func ParseWheelSideType(et string) (WheelSideType, error) {
	switch et {
	case string(WheelSideTypeLeft):
		return WheelSideTypeLeft, nil
	case string(WheelSideTypeRight):
		return WheelSideTypeRight, nil
	case string(WheelSideTypeOthers):
		return WheelSideTypeOthers, nil
	default:
		return "", ErrInvalidWheelSideType
	}
}

type VehicleSourceType string

const (
	VehicleSourceTypeJudicial   VehicleSourceType = "judicial"
	VehicleSourceTypeCommission VehicleSourceType = "commission"
	VehicleSourceTypeOverseas   VehicleSourceType = "overseas"
	VehicleSourceTypeUnknown    VehicleSourceType = "unknown"
)

func ParseVehicleSourceType(et string) (VehicleSourceType, error) {
	switch et {
	case string(VehicleSourceTypeJudicial):
		return VehicleSourceTypeJudicial, nil
	case string(VehicleSourceTypeCommission):
		return VehicleSourceTypeCommission, nil
	case string(VehicleSourceTypeOverseas):
		return VehicleSourceTypeOverseas, nil
	case string(VehicleSourceTypeUnknown):
		return VehicleSourceTypeUnknown, nil
	default:
		return "", ErrInvalidVehicleSourceType
	}
}

type SpecialIncidentType string

const (
	SpecialIncidentTypeCasualty SpecialIncidentType = "casualty"
	SpecialIncidentTypeSuicide  SpecialIncidentType = "suicide"
	SpecialIncidentTypeWatered  SpecialIncidentType = "watered"
)

func ParseSpecialIncidentType(et string) (SpecialIncidentType, error) {
	switch et {
	case string(SpecialIncidentTypeCasualty):
		return SpecialIncidentTypeCasualty, nil
	case string(SpecialIncidentTypeSuicide):
		return SpecialIncidentTypeSuicide, nil
	case string(SpecialIncidentTypeWatered):
		return SpecialIncidentTypeWatered, nil
	default:
		return "", ErrInvalidSpecialIncidentType
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
