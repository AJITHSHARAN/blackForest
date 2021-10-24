package models

type Vehicle struct {
	Year  string `json:"year"`
	Makes []Make `json:"makes"`
}

// Vehicle Make represents data about a record vehicle.
type Make struct {
	Make   string  `json:"make"`
	Models []Model `json:"models"`
}

// Vehicle Model represents data about a record vehicle.
type Model struct {
	Model  string `json:"model"`
	Engine string `json:"engine"`
}

type vehchileTable struct {
	Year   string `json:"year"`
	Make   string `json:"make"`
	Model  string `json:"model"`
	Engine string `json:"engine"`
}
