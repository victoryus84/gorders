package utils

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
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