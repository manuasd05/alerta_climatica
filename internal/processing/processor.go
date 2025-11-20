package processing

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// Processor coordina el procesamiento concurrente de mensajes entrantes
// usando goroutines y canales. Asigna trabajos a workers simulando zonas.
type Processor struct {
	zones   []string
	inCh    chan IncomingMessage
	wg      sync.WaitGroup
	onAlert func(Alert) // callback para notificar alertas detectadas
}

// NewProcessor crea un nuevo procesador con las zonas provistas.
func NewProcessor(zones []string, onAlert func(Alert)) *Processor {
	return &Processor{
		zones:   zones,
		inCh:    make(chan IncomingMessage, 64),
		onAlert: onAlert,
	}
}

// StartWorkers inicia n workers que consumen del canal y procesan mensajes.
func (p *Processor) StartWorkers(n int) {
	if n <= 0 {
		n = 1
	}
	for i := 0; i < n; i++ {
		p.wg.Add(1)
		go func(workerID int) {
			defer p.wg.Done()
			for msg := range p.inCh {
				typ, sev, extract := detect(msg.Text)
				alert := Alert{
					ID:        newID(),
					Zone:      msg.Zone,
					Type:      typ,
					Severity:  sev,
					Message:   msg.Text,
					Extract:   extract,
					Timestamp: time.Now(),
				}
				// Entregar al callback para que el servidor actualice estado.
				if p.onAlert != nil {
					p.onAlert(alert)
				}
				// Simular latencia variable entre sensores/zonas.
				time.Sleep(50 * time.Millisecond)
			}
		}(i)
	}
}

// Submit envía un mensaje entrante al pool de workers.
func (p *Processor) Submit(msg IncomingMessage) {
	// Asignación simple: se envía al canal común y los workers compiten.
	p.inCh <- msg
}

// Close cierra el canal y espera a que terminen los workers.
func (p *Processor) Close() {
	close(p.inCh)
	p.wg.Wait()
}

func newID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
