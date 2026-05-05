package utils

import (
	"testing"
)

func TestGeneratePassword(t *testing.T) {
    parola, err := GeneratePassword(12)
    if err != nil {
        t.Errorf("Eroare neașteptată: %v", err)
    }

    if len(parola) != 12 {
        t.Errorf("Lungime incorectă! Am primit %d, voiam 12", len(parola))
    }
    
    // Printăm parola ca să o vedem în consolă
    t.Logf("Parola generată la test: %s", parola)
}