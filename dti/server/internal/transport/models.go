package transport

type SuscribeDTO struct {
	Suscribe bool `json:"suscribe"`
}

type SuscribeResponseDTO struct {
	Suscribed bool `json:"suscribed"`
}

type HealthCheckDTO struct {
	HealthCheck bool `json:"health-check"`
}
