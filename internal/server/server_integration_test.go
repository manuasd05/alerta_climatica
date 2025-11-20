package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"alerta_climatica/internal/processing"
	srvpkg "alerta_climatica/internal/server"
	"alerta_climatica/internal/storage"
)

// TestEndToEnd levanta el servidor en memoria, envía un SMS simulado y verifica
// que la alerta aparece en GET /api/alerts (persistida en SQLite).
func TestEndToEnd(t *testing.T) {
	dir := t.TempDir()
	dbPath := dir + "/itest_alerts.db"

	store, err := storage.NewSQLite(dbPath)
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	defer store.Close()

	st := srvpkg.NewState(store)

	zones := []string{"Zona Test"}
	proc := processing.NewProcessor(zones, st.AddAlert)
	proc.StartWorkers(1)
	defer proc.Close()

	srv := srvpkg.NewServer(st, proc)
	ts := httptest.NewServer(srv.Router())
	defer ts.Close()

	// Enviar POST /api/sms
	payload := map[string]string{"zona": "Zona Test", "texto": "Lluvia intensa en el area"}
	b, _ := json.Marshal(payload)
	resp, err := http.Post(ts.URL+"/api/sms", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("post request failed: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status from POST /api/sms: %v", resp.Status)
	}

	// Poll GET /api/alerts hasta encontrar la alerta o timeout
	deadline := time.Now().Add(3 * time.Second)
	found := false
	for time.Now().Before(deadline) {
		time.Sleep(100 * time.Millisecond)
		resp, err := http.Get(ts.URL + "/api/alerts")
		if err != nil {
			t.Fatalf("get alerts failed: %v", err)
		}
		var list []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
			resp.Body.Close()
			t.Fatalf("decode alerts failed: %v", err)
		}
		resp.Body.Close()
		for _, a := range list {
			if a["zona"] == "Zona Test" {
				if a["tipo"] == "lluvia" || a["severidad"] == "alta" || a["mensaje"] != nil {
					found = true
					break
				}
			}
		}
		if found {
			break
		}
	}

	if !found {
		t.Fatal("alert not found in GET /api/alerts within timeout")
	}
}
