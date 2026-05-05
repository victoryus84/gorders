package utils

import (
	"crypto/rand"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/gin-gonic/gin"
)

func ParseBody[T any](c *gin.Context) ([]T, error) {
	var result []T
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, fmt.Errorf("eroare la citirea body-ului: %v", err)
	}

	contentType := strings.ToLower(c.GetHeader("Content-Type"))

	// --- LOGICA XML ---
	if strings.Contains(contentType, "xml") {
		// 1. Încercăm varianta de listă cu wrapper <items><item>...</item></items>
		var wrapper struct {
			Items []T `xml:"item"`
		}
		if err := xml.Unmarshal(data, &wrapper); err == nil && len(wrapper.Items) > 0 {
			return wrapper.Items, nil
		}

		// 2. Încercăm obiect singur sau listă fără wrapper (direct <item>...</item>)
		var single T
		if err := xml.Unmarshal(data, &single); err == nil {
			return []T{single}, nil
		}
		
		return nil, errors.New("format XML invalid sau incompatibil")
	}

	// --- LOGICA JSON (Sau default dacă nu e XML) ---
	// 1. Încercăm listă [{}, {}]
	if err := json.Unmarshal(data, &result); err == nil {
		return result, nil
	}

	// 2. Încercăm obiect singur {}
	var single T
	if err := json.Unmarshal(data, &single); err == nil {
		return []T{single}, nil
	}

	// Dacă am ajuns aici, înseamnă că niciun format nu a mers
	return nil, errors.New("nu s-au putut decoda datele (nici JSON, nici XML)")
}
func GeneratePassword(length int) (string, error) {

	// Setul de caractere pe care îl vom folosi
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"
	// 1. Verificăm lungimea (Protecție la erori)
	if length < 8 {
		return "", fmt.Errorf("lungimea %d este prea mică pentru o parolă sigură", length)
	}

	// Creăm o felie (slice) de caractere de lungimea dorită
	password := make([]byte, length)

	for i := 0; i < length; i++ {
		// Alegem un index aleatoriu din charset
		// Folosim crypto/rand pentru securitate maximă
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err // Returnăm eroarea dacă "mașinăria" de noroc s-a stricat
		}

		// Punem caracterul ales în parola noastră
		password[i] = charset[index.Int64()]
	}

	// Transformăm felia de bytes într-un string frumos
	return string(password), nil
}