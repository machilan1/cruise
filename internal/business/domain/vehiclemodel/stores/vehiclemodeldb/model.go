package vehiclemodeldb

import "github.com/machilan1/cruise/internal/business/domain/vehiclemodel"

type dbVehicleModel struct {
	ID                 int     `db:"vehicle_model_id"`
	SeriesName         string  `db:"model_series_name"`
	CommercialName     string  `db:"model_commercial_name"`
	ModelYear          int     `db:"model_year"`
	BrandID            int     `db:"brand_id"`
	Nickname           *string `db:"nickname"`
	EngineDisplacement int     `db:"engine_displacement"`
	DriveType          string  `db:"drive_type"`
	FuelType           string  `db:"fuel_type"`
	BodyStyle          string  `db:"body_style"`
	TransmissionType   string  `db:"transmission_type"`
}

func toCoreVehicleModel(dbv dbVehicleModel) (vehiclemodel.VehicleModel, error) {
	driveType, err := vehiclemodel.ParseDriveType(dbv.DriveType)
	if err != nil {
		return vehiclemodel.VehicleModel{}, err
	}

	fuelType, err := vehiclemodel.ParseFuelType(dbv.FuelType)
	if err != nil {
		return vehiclemodel.VehicleModel{}, err
	}

	bodyStyle, err := vehiclemodel.ParseBodyStyle(dbv.BodyStyle)
	if err != nil {
		return vehiclemodel.VehicleModel{}, err
	}

	transmissionType, err := vehiclemodel.ParseTransmissionType(dbv.TransmissionType)
	if err != nil {
		return vehiclemodel.VehicleModel{}, err
	}

	return vehiclemodel.VehicleModel{
		ID:                 dbv.ID,
		SeriesName:         dbv.SeriesName,
		CommercialName:     dbv.CommercialName,
		ModelYear:          dbv.ModelYear,
		BrandID:            dbv.BrandID,
		Nickname:           dbv.Nickname,
		EngineDisplacement: dbv.EngineDisplacement,
		DriveType:          driveType,
		FuelType:           fuelType,
		BodyStyle:          bodyStyle,
		TransmissionType:   transmissionType,
	}, nil
}

func toCoreVehicleModels(dbvs []dbVehicleModel) ([]vehiclemodel.VehicleModel, error) {
	vms := make([]vehiclemodel.VehicleModel, len(dbvs))
	for i, v := range dbvs {
		vm, err := toCoreVehicleModel(v)
		if err != nil {
			return nil, err
		}
		vms[i] = vm
	}
	return vms, nil
}

func toDBVehicleModel(vm vehiclemodel.VehicleModel) dbVehicleModel {
	return dbVehicleModel{
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
