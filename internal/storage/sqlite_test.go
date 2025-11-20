package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"alerta_climatica/internal/processing"
)

func TestSQLiteSaveAndList(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test_alerts.db")

	s, err := NewSQLite(dbPath)
	if err != nil {
		t.Fatalf("NewSQLite failed: %v", err)
	}
	defer func() {
		_ = s.Close()
		_ = os.Remove(dbPath)
	}()

	now := time.Now().UTC().Truncate(time.Second)
	a := processing.Alert{ID: "a1", Zone: "Zona Test", Type: "lluvia", Severity: "alta", Message: "Prueba", Extract: "lluvia", Timestamp: now}

	if err := s.SaveAlert(a); err != nil {
		t.Fatalf("SaveAlert failed: %v", err)
	}

	list, err := s.ListAlerts()
	if err != nil {
		t.Fatalf("ListAlerts failed: %v", err)
	}

	if len(list) == 0 {
		t.Fatalf("expected at least one alert, got 0")
	}

	found := false
	for _, it := range list {
		if it.ID == a.ID {
			found = true
			if !it.Timestamp.Equal(a.Timestamp) {
				t.Logf("timestamps differ: got %v want %v", it.Timestamp, a.Timestamp)
			}
		}
	}
	if !found {
		t.Fatalf("saved alert not found in list")
	}
}
