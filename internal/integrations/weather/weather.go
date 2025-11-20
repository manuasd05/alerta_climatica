package weather

// Paquete weather: stub para una futura integración con API meteorológica.

// Client define los métodos esperados de un proveedor de datos del tiempo.
type Client interface {
	// CurrentConditions obtiene condiciones actuales para una ubicación.
	CurrentConditions(location string) (Conditions, error)
	// Forecast obtiene pronóstico breve para una ubicación.
	Forecast(location string) ([]Conditions, error)
}

// Conditions es un modelo simplificado de datos climáticos.
type Conditions struct {
	TemperatureC float64
	RainMM       float64
	WindKPH      float64
	HumidityPct  float64
	Summary      string
}

// Se podrían añadir adaptadores a APIs públicas (e.g., con claves de API).
