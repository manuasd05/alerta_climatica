package sms

// Paquete sms: stub para una futura integración con proveedor de SMS.
// Define interfaces y estructuras básicas para desacoplar implementación.

// Sender representa un cliente capaz de enviar SMS reales.
type Sender interface {
	Send(to string, message string) error
}

// Receiver representa un webhook/cliente para recibir SMS entrantes.
type Receiver interface {
	// ParseAndAck procesa la carga entrante del proveedor y confirma recepción.
	ParseAndAck(payload []byte) error
}

// En el futuro, puede añadirse una implementación concreta (Twilio, etc.).
