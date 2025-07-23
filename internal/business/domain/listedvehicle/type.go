package listedvehicle

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
		return "", ErrInvalidFuelType
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
		return "", ErrInvalidTransmissionType
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
		return "", ErrInvalidEngineType
	}
}

type HeadlightType string

const (
	HeadlightTypeHID         HeadlightType = "hid"
	HeadlightTypeLED         HeadlightType = "led"
	HeadlightTypeTungsten    HeadlightType = "tungsten"
	HeadlightTypeUnspecified HeadlightType = "unspecified"
)

func ParseHeadlightType(et string) (HeadlightType, error) {
	switch et {
	case string(HeadlightTypeHID):
		return HeadlightTypeHID, nil
	case string(HeadlightTypeLED):
		return HeadlightTypeLED, nil
	case string(HeadlightTypeTungsten):
		return HeadlightTypeTungsten, nil
	case string(HeadlightTypeUnspecified):
		return HeadlightTypeUnspecified, nil
	default:
		return "", ErrInvalidHeadlightType
	}
}

type WheelSideType string

const (
	WheelSideTypeLeft        WheelSideType = "left"
	WheelSideTypeRight       WheelSideType = "right"
	WheelSideTypeUnspecified WheelSideType = "unspecified"
)

func ParseWheelSideType(et string) (WheelSideType, error) {
	switch et {
	case string(WheelSideTypeLeft):
		return WheelSideTypeLeft, nil
	case string(WheelSideTypeRight):
		return WheelSideTypeRight, nil
	case string(WheelSideTypeUnspecified):
		return WheelSideTypeUnspecified, nil
	default:
		return "", ErrInvalidWheelSideType
	}
}

type VehicleSourceType string

const (
	VehicleSourceTypeJudicial    VehicleSourceType = "judicial"
	VehicleSourceTypeCommission  VehicleSourceType = "commission"
	VehicleSourceTypeOverseas    VehicleSourceType = "overseas"
	VehicleSourceTypeUnspecified VehicleSourceType = "unspecified"
)

func ParseVehicleSourceType(et string) (VehicleSourceType, error) {
	switch et {
	case string(VehicleSourceTypeJudicial):
		return VehicleSourceTypeJudicial, nil
	case string(VehicleSourceTypeCommission):
		return VehicleSourceTypeCommission, nil
	case string(VehicleSourceTypeOverseas):
		return VehicleSourceTypeOverseas, nil
	case string(VehicleSourceTypeUnspecified):
		return VehicleSourceTypeUnspecified, nil
	default:
		return "", ErrInvalidVehicleSourceType
	}
}

type SpecialIncidentType string

const (
	SpecialIncidentTypeCasualty    SpecialIncidentType = "casualty"
	SpecialIncidentTypeSuicide     SpecialIncidentType = "suicide"
	SpecialIncidentTypeWatered     SpecialIncidentType = "watered"
	SpecialIncidentTypeUnspecified SpecialIncidentType = "unspecified"
)

func ParseSpecialIncidentType(et string) (SpecialIncidentType, error) {
	switch et {
	case string(SpecialIncidentTypeCasualty):
		return SpecialIncidentTypeCasualty, nil
	case string(SpecialIncidentTypeSuicide):
		return SpecialIncidentTypeSuicide, nil
	case string(SpecialIncidentTypeWatered):
		return SpecialIncidentTypeWatered, nil
	case string(SpecialIncidentTypeUnspecified):
		return SpecialIncidentTypeWatered, nil
	default:
		return "", ErrInvalidSpecialIncidentType
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
	case string(BodyStyleUnspecified):
		return BodyStyleUnspecified, nil
	default:
		return "", ErrInvalidBodyStyle
	}
}

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
		return "", ErrInvalidDriveType
	}
}
