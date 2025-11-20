package processing

import "time"

// IncomingMessage representa un "SMS" simulado recibido por el sistema.
type IncomingMessage struct {
	Zone       string    `json:"zona"`
	Text       string    `json:"texto"`
	ReceivedAt time.Time `json:"recibido_en"`
}

// Alert representa una alerta resultante del análisis del mensaje.
type Alert struct {
	ID        string    `json:"id"`
	Zone      string    `json:"zona"`
	Type      string    `json:"tipo"`
	Severity  string    `json:"severidad"`
	Message   string    `json:"mensaje"`
	Extract   string    `json:"extracto"`
	Timestamp time.Time `json:"timestamp"`
}
