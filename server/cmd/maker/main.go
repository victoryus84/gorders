package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Field struct {
	Name     string
	Type     string
	JSONName string
	Binding  string
}

type Entity struct {
	Name      string
	LowerName string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("❌ Error: Please provide the module name. Example: go run cmd/maker/main.go product")
		os.Exit(1)
	}

	rawName := os.Args[1]
	entity := Entity{
		Name:      toTitle(strings.ToLower(rawName)),
		LowerName: strings.ToLower(rawName),
	}

	fmt.Printf("🤖 GOrders Maker: Creating module '%s' with Logger included!\n", entity.Name)
	fmt.Println("👉 Input (format: name:type). Type 'exit' when you're done.")

	scanner := bufio.NewScanner(os.Stdin)
	var fields []Field

	for {
		fmt.Print("Column: ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		if strings.ToLower(input) == "exit" {
			break
		}

		parts := strings.Split(input, ":")

		// Permitem fie 2 parametri (nume:tip), fie 3 (nume:tip:req)
		if len(parts) < 2 || len(parts) > 3 {
			fmt.Println("⚠️ Format invalid. Use name:type or name:type:req (e.g., price:float64:req)")
			continue
		}

		// Setăm binding-ul pe gol by default
		bindingTag := ""

		// Dacă avem 3 bucăți și a 3-a e "req"
		if len(parts) == 3 && strings.ToLower(strings.TrimSpace(parts[2])) == "req" {
			bindingTag = ` binding:"required"` // Lăsăm un spațiu la început intenționat!
		}

		fields = append(fields, Field{
			Name:     toTitle(parts[0]),
			Type:     strings.TrimSpace(parts[1]),
			JSONName: strings.ToLower(strings.TrimSpace(parts[0])),
			Binding:  bindingTag, // Îl pasăm șablonului
		})
	}

	data := struct {
		Entity Entity
		Fields []Field
	}{
		Entity: entity,
		Fields: fields,
	}

	generateFile("internal/models/mod_"+entity.LowerName+".go", modelTpl, data)
	generateFile("internal/dto/dto_"+entity.LowerName+".go", dtoTpl, data)
	generateFile("internal/service/svc_"+entity.LowerName+".go", serviceTpl, data)
	generateFile("internal/handler/hdl_"+entity.LowerName+".go", handlerTpl, data)

	fmt.Println("\n🎉 Success, files were generated.")
}

func toTitle(s string) string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func generateFile(path string, tplContent string, data interface{}) {
	os.MkdirAll(filepath.Dir(path), 0755)
	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("❌ Error creating file %s: %v\n", path, err)
		return
	}
	defer f.Close()

	t := template.Must(template.New("").Parse(tplContent))
	if err := t.Execute(f, data); err != nil {
		fmt.Printf("❌ Error writing template to %s: %v\n", path, err)
		return
	}
	fmt.Printf("✅ Generated: %s\n", path)
}

// =====================================================================
// ȘABLOANELE (TEMPLATES) DE COD - ACUM CU LOGGER INCLUS!
// =====================================================================

var modelTpl = `package models

import "time"

type {{.Entity.Name}} struct {
	ID        uint      ` + "`gorm:\"primarykey\"`" + `
{{range .Fields}}	{{.Name}} {{.Type}} ` + "`gorm:\"column:{{.JSONName}}\"`" + `
{{end}}	CreatedAt time.Time
	UpdatedAt time.Time
}
`

var dtoTpl = `package dto

type {{.Entity.Name}}DTO struct {
{{range .Fields}}	{{.Name}} {{.Type}} ` + "`json:\"{{.JSONName}}\" xml:\"{{.JSONName}}\"{{.Binding}}`" + `
{{end}}}
`

var serviceTpl = `package service

import (
	"github.com/victoryus84/gorders/internal/logger"
	"github.com/victoryus84/gorders/internal/models"
	"github.com/victoryus84/gorders/internal/repository"
)

type {{.Entity.Name}}Service struct {
	repo *repository.GenericRepo[models.{{.Entity.Name}}]
}

func New{{.Entity.Name}}Service(repo *repository.GenericRepo[models.{{.Entity.Name}}]) *{{.Entity.Name}}Service {
	return &{{.Entity.Name}}Service{repo: repo}
}

func (s *{{.Entity.Name}}Service) Create(item *models.{{.Entity.Name}}) error {
	err := s.repo.Create(item)
	if err != nil {
		logger.LogError("❌ Database error during {{.Entity.Name}} creation", err)
		return err
	}
	
	logger.LogInfo("✅ {{.Entity.Name}} successfully created", logger.Uint("id", item.ID))
	return nil
}

func (s *{{.Entity.Name}}Service) FindAll() ([]models.{{.Entity.Name}}, error) {
	items, err := s.repo.FindManyWhere("1 = 1")
	if err != nil {
		logger.LogError("❌ Failed to fetch {{.Entity.Name}} records", err)
		return nil, err
	}
	
	logger.LogInfo("🔍 Fetched {{.Entity.Name}} records", logger.Int("count", len(items)))
	return items, nil
}
`

var handlerTpl = `package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victoryus84/gorders/internal/dto"
	"github.com/victoryus84/gorders/internal/logger"
	"github.com/victoryus84/gorders/internal/models"
	"github.com/victoryus84/gorders/internal/service"
)

type {{.Entity.Name}}Handler struct {
	svc *service.{{.Entity.Name}}Service
}

func New{{.Entity.Name}}Handler(svc *service.{{.Entity.Name}}Service) *{{.Entity.Name}}Handler {
	return &{{.Entity.Name}}Handler{svc: svc}
}

func (h *{{.Entity.Name}}Handler) Create(c *gin.Context) {
	var input dto.{{.Entity.Name}}DTO
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.LogError("⚠️ Invalid {{.Entity.Name}} payload received", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}

	// Mapează DTO -> Model
	item := models.{{.Entity.Name}}{
{{range .Fields}}		{{.Name}}: input.{{.Name}},
{{end}}	}

	if err := h.svc.Create(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create {{.Entity.LowerName}}"})
		return
	}

	c.JSON(http.StatusCreated, item)
}
`
