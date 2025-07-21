package vehiclemodelapi

import (
	"net/http"

	"github.com/machilan1/cruise/internal/app/sdk/errs"
	"github.com/machilan1/cruise/internal/business/domain/vehiclemodel"
	"github.com/machilan1/cruise/internal/framework/validate"
)

type AppVehicleModel struct {
	ID                 int     `json:"id"`
	SeriesName         string  `json:"seriesName"`
	CommercialName     string  `json:"commercialName"`
	ModelYear          int     `json:"modelYear"`
	BrandID            int     `json:"brandId"`
	Nickname           *string `json:"nickname"`
	EngineDisplacement int     `json:"engineDisplacement"`
	DriveType          string  `json:"driveType"`
	FuelType           string  `json:"fuelType"`
	BodyStyle          string  `json:"bodyStyle"`
	TransmissionType   string  `json:"transmissionType"`
}

func toAppVehicleModel(vm vehiclemodel.VehicleModel) AppVehicleModel {
	return AppVehicleModel{
		ID:                 vm.ID,
		SeriesName:         vm.SeriesName,
		CommercialName:     vm.CommercialName,
		ModelYear:          vm.ModelYear,
		BrandID:            vm.BrandID,
		Nickname:           vm.Nickname,
		EngineDisplacement: vm.EngineDisplacement,
		DriveType:          string(vm.DriveType),
		FuelType:           string(vm.FuelType),
		BodyStyle:          string(vm.BodyStyle),
		TransmissionType:   string(vm.TransmissionType),
	}
}

func toAppVehicleModels(vms []vehiclemodel.VehicleModel) []AppVehicleModel {
	avms := make([]AppVehicleModel, len(vms))

	for i, v := range vms {
		avms[i] = toAppVehicleModel(v)
	}
	return avms
}

type AppNewVehicleModel struct {
	SeriesName         string  `json:"seriesName" validate:"required"`
	CommercialName     string  `json:"commercialName" validate:"required"`
	ModelYear          int     `json:"modelYear" validate:"required"`
	BrandID            int     `json:"brandId" validate:"required"`
	Nickname           *string `json:"nickname" validate:"required"`
	EngineDisplacement int     `json:"engineDisplacement" validate:"required"`
	DriveType          string  `json:"driveType" validate:"oneof=fwd rwd 4wd awd others"`
	FuelType           string  `json:"fuelType" validate:"oneof=diesel gasoline electric gas others"`
	BodyStyle          string  `json:"bodyStyle" validate:"oneof=sedan wagon hatchback gt sports van truck suv others"`
	TransmissionType   string  `json:"transmissionType" validate:"oneof=automatic manual"`
}

func (app AppNewVehicleModel) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}
	return nil
}

func toCoreNewVehicleModel(anvm AppNewVehicleModel) (vehiclemodel.NewVehicleModel, error) {
	driveType, err := vehiclemodel.ParseDriveType(anvm.DriveType)
	if err != nil {
		return vehiclemodel.NewVehicleModel{}, errs.NewTrustedError(err, http.StatusBadRequest)
	}

	fuelType, err := vehiclemodel.ParseFuelType(anvm.FuelType)
	if err != nil {
		return vehiclemodel.NewVehicleModel{}, errs.NewTrustedError(err, http.StatusBadRequest)
	}

	bodyStyle, err := vehiclemodel.ParseBodyStyle(anvm.BodyStyle)
	if err != nil {
		return vehiclemodel.NewVehicleModel{}, errs.NewTrustedError(err, http.StatusBadRequest)
	}

	transmissionType, err := vehiclemodel.ParseTransmissionType(anvm.TransmissionType)
	if err != nil {
		return vehiclemodel.NewVehicleModel{}, errs.NewTrustedError(err, http.StatusBadRequest)
	}

	return vehiclemodel.NewVehicleModel{
		SeriesName:         anvm.SeriesName,
		CommercialName:     anvm.CommercialName,
		ModelYear:          anvm.ModelYear,
		BrandID:            anvm.BrandID,
		Nickname:           anvm.Nickname,
		EngineDisplacement: anvm.EngineDisplacement,
		DriveType:          driveType,
		FuelType:           fuelType,
		BodyStyle:          bodyStyle,
		TransmissionType:   transmissionType,
	}, nil
}

type AppUpdateVehicleModel struct {
	SeriesName         *string `json:"seriesName"`
	CommercialName     *string `json:"commercialName"`
	ModelYear          *int    `json:"modelYear"`
	Nickname           *string `json:"nickname"`
	EngineDisplacement *int    `json:"engineDisplacement"`
	DriveType          *string `json:"driveType"`
	FuelType           *string `json:"fuelType"`
	BodyStyle          *string `json:"bodyStyle"`
	TransmissionType   *string `json:"transmissionType"`
}

func toCoreUpdateVehicleModel(auvm AppUpdateVehicleModel) (vehiclemodel.UpdateVehicleModel, error) {
	var driveType *vehiclemodel.DriveType
	var fuelType *vehiclemodel.FuelType
	var bodyStyle *vehiclemodel.BodyStyle
	var transmissionType *vehiclemodel.TransmissionType

	if auvm.DriveType != nil {
		dt, err := vehiclemodel.ParseDriveType(*auvm.DriveType)
		if err != nil {
			return vehiclemodel.UpdateVehicleModel{}, errs.NewTrustedError(err, http.StatusBadRequest)
		}
		driveType = &dt
	}

	if auvm.FuelType != nil {
		ft, err := vehiclemodel.ParseFuelType(*auvm.FuelType)
		if err != nil {
			return vehiclemodel.UpdateVehicleModel{}, errs.NewTrustedError(err, http.StatusBadRequest)
		}
		fuelType = &ft
	}

	if auvm.BodyStyle != nil {
		bs, err := vehiclemodel.ParseBodyStyle(*auvm.BodyStyle)
		if err != nil {
			return vehiclemodel.UpdateVehicleModel{}, errs.NewTrustedError(err, http.StatusBadRequest)
		}
		bodyStyle = &bs
	}

	if auvm.TransmissionType != nil {

		tt, err := vehiclemodel.ParseTransmissionType(*auvm.TransmissionType)
		if err != nil {
			return vehiclemodel.UpdateVehicleModel{}, errs.NewTrustedError(err, http.StatusBadRequest)
		}
		transmissionType = &tt
	}

	return vehiclemodel.UpdateVehicleModel{
		SeriesName:         auvm.SeriesName,
		CommercialName:     auvm.CommercialName,
		ModelYear:          auvm.ModelYear,
		Nickname:           auvm.Nickname,
		EngineDisplacement: auvm.EngineDisplacement,
		DriveType:          driveType,
		FuelType:           fuelType,
		BodyStyle:          bodyStyle,
		TransmissionType:   transmissionType,
	}, nil
}
