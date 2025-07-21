package town

var (
	SampleCity = City{
		ID:   1,
		Name: "臺北市",
	}
	SampleTown = Town{
		ID:       1,
		Name:     "中正區",
		City:     SampleCity,
		PostCode: 100,
	}
	SampleTown2 = Town{
		ID:       2,
		Name:     "大同區",
		City:     SampleCity,
		PostCode: 103,
	}
)
