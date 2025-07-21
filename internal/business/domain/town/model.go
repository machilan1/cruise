package town

type Town struct {
	ID       int
	Name     string
	City     City
	PostCode int
}

type City struct {
	ID   int
	Name string
}
