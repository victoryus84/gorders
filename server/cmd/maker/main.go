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
	GormTag  string
	Comment  string
}

type Entity struct {
	Name      string
	LowerName string
	HasGroup  bool
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
	generateFile("internal/repository/rep_"+entity.LowerName+".go", repositoryTpl, data)

	fmt.Println("\n✅ Fisiere generate cu succes!")
	fmt.Println("👉 NU UITA SA ADAUGI CONSTRUCTORII IN module.go:")
	fmt.Printf("\n--- internal/repository/module.go ---\nfx.Provide(\n    New%sRepository,\n)\n", entity.Name)
	fmt.Printf("\n--- internal/service/module.go ---\nfx.Provide(\n    New%sService,\n)\n", entity.Name)
	fmt.Printf("\n--- internal/handler/module.go ---\nfx.Provide(\n    New%sHandler,\n)\nfx.Invoke(\n    Register%sRoutes,\n)\n", entity.Name, entity.Name)
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

import (
	"gorm.io/gorm"
)

// ********** {{.Entity.Name}} **********
type {{.Entity.Name}} struct {
	gorm.Model
	UUIDModel 	` + "`gorm:\"embedded\"`" + `
{{range .Fields}}	{{.Name}} {{.Type}} ` + "{{if .GormTag}}`{{.GormTag}}`{{end}}{{if .Comment}} // {{.Comment}}{{end}}" + `
{{end}}
{{- if .Entity.HasGroup}}
	{{.Entity.Name}}Group   {{.Entity.Name}}Group ` + "`gorm:\"foreignKey:{{.Entity.Name}}GroupID;references:ID\"`" + ` // Grupa din care face parte
	{{.Entity.Name}}GroupID *uint                 ` + "`gorm:\"column:{{.Entity.LowerName}}_group_id\"`" + `         // ID-ul grupei
{{- end}}
}

{{if .Entity.HasGroup}}
// ********** {{.Entity.Name}}Group - Grupa de {{.Entity.Name}} **********
type {{.Entity.Name}}Group struct {
	gorm.Model
	UUIDModel   ` + "`gorm:\"embedded\"`" + `
	Code        string ` + "`gorm:\"type:varchar(15);unique;not null\"`" + ` // Codul grupei (ex: "001")
	Name        string ` + "`gorm:\"type:varchar(100);not null\"`" + `       // Numele grupei
	Description string ` + "`gorm:\"type:text\"`" + `                        // Descrierea grupei
	
	ParentID    *uint 
	Parent      *{{.Entity.Name}}Group   ` + "`gorm:\"foreignKey:ParentID\"`" + `             // Legătură către grupa părinte (dacă există)
	Children    []{{.Entity.Name}}Group  ` + "`gorm:\"foreignKey:ParentID\"`" + `             // Legătură către grupele copil
	{{.Entity.Name}}s    []{{.Entity.Name}}       ` + "`gorm:\"foreignKey:{{.Entity.Name}}GroupID\"`" + ` // Elementele din această grupă
}
{{end}}
`

var dtoTpl = `package dto

// DTO pentru {{.Entity.Name}}
type {{.Entity.Name}}DTO struct {
	SyncID string ` + "`json:\"sync_id\" xml:\"sync_id\"`" + ` // ID opțional pentru sincronizare
{{range .Fields}}	{{.Name}} {{.Type}} ` + "`json:\"{{.JSONName}}\" xml:\"{{.JSONName}}\"{{.Binding}}`" + `
{{end}}}
{{if .Entity.HasGroup}}
	// DTO pentru Grupa de {{.Entity.Name}}
	type {{.Entity.Name}}GroupDTO struct {
	SyncID      string ` + "`json:\"sync_id\" xml:\"sync_id\"`" + `
	Code        string ` + "`json:\"code\" xml:\"code\" binding:\"required\"`" + `
	Name        string ` + "`json:\"name\" xml:\"name\" binding:\"required\"`" + `
	Description string ` + "`json:\"description,omitempty\" xml:\"description,omitempty\"`" + `
	ParentCode  string ` + "`json:\"parent_code,omitempty\" xml:\"parent_code,omitempty\"`" + `
}
{{end}}
`

var serviceTpl = `package service

import (
	"strings"

	"github.com/victoryus84/gorders/internal/dto"
	"github.com/victoryus84/gorders/internal/models"
	"github.com/victoryus84/gorders/internal/repository"
)

// 1. INTERFAȚA PUBLICĂ
type {{.Entity.Name}}Service interface {
	ProcessImport(requests []dto.{{.Entity.Name}}DTO) dto.ImportResult
	FindAll() ([]models.{{.Entity.Name}}, error)
{{- if .Entity.HasGroup}}
	ProcessGroupImport(requests []dto.{{.Entity.Name}}GroupDTO) dto.ImportResult
{{- end}}
}

// 2. STRUCTURA PRIVATĂ
type {{.Entity.LowerName}}Service struct {
	repo repository.{{.Entity.Name}}Repository
}

func New{{.Entity.Name}}Service(repo repository.{{.Entity.Name}}Repository) {{.Entity.Name}}Service {
	return &{{.Entity.LowerName}}Service{repo: repo}
}

// 3. IMPLEMENTAREA LOGICII PRINCIPALE
func (svc *{{.Entity.LowerName}}Service) ProcessImport(requests []dto.{{.Entity.Name}}DTO) dto.ImportResult {
	if len(requests) == 0 {
		return dto.ImportResult{Status: "success", Message: "Niciun element primit."}
	}

	itemsToSave := make([]models.{{.Entity.Name}}, 0, len(requests))
	skipped := make([]map[string]string, 0)

{{if .Entity.HasGroup}}	// Încărcăm dicționarul pentru maparea grupei
	groupMap, _ := svc.repo.GetGroupIDMap()
{{end}}
	for _, req := range requests {
		// Validare de bază
		if strings.TrimSpace(req.Code) == "" {
			skipped = append(skipped, map[string]string{"code": req.Code, "reason": "missing_code"})
			continue
		}

		item := models.{{.Entity.Name}}{
{{range .Fields}}			{{.Name}}: req.{{.Name}},
{{end}}		}

{{if .Entity.HasGroup}}		// Mapare dinamică a grupei
		if groupID, ok := groupMap[req.GroupCode]; ok && groupID != 0 {
			item.{{.Entity.Name}}GroupID = &groupID
		}
{{end}}
		itemsToSave = append(itemsToSave, item)
	}

	if err := svc.repo.UpsertBatch(itemsToSave, 500); err != nil {
		return dto.ImportResult{Status: "error", Message: "Eroare la salvare: " + err.Error()}
	}

	return dto.ImportResult{
		Status:         "success",
		TotalProcessed: len(itemsToSave),
		TotalSkipped:   len(skipped),
		Message:        "Sincronizare {{.Entity.Name}} finalizată!",
	}
}

func (svc *{{.Entity.LowerName}}Service) FindAll() ([]models.{{.Entity.Name}}, error) {
	return svc.repo.FindAll()
}

{{if .Entity.HasGroup}}
// 4. IMPLEMENTAREA LOGICII PENTRU GRUPE
func (svc *{{.Entity.LowerName}}Service) ProcessGroupImport(requests []dto.{{.Entity.Name}}GroupDTO) dto.ImportResult {
	if len(requests) == 0 {
		return dto.ImportResult{Status: "success", Message: "Nicio grupă primită."}
	}

	groupsToSave := make([]models.{{.Entity.Name}}Group, 0, len(requests))
	skipped := make([]map[string]string, 0)

	// Aducem dicționarul ca să legăm sub-categoriile de părinții lor (ParentCode -> ParentID)
	groupMap, _ := svc.repo.GetGroupIDMap()

	for _, req := range requests {
		if strings.TrimSpace(req.Code) == "" || strings.TrimSpace(req.Name) == "" {
			skipped = append(skipped, map[string]string{"code": req.Code, "reason": "missing_required"})
			continue
		}

		group := models.{{.Entity.Name}}Group{
			Code:        req.Code,
			Name:        req.Name,
			Description: req.Description,
		}

		// Verificăm dacă are un părinte
		if req.ParentCode != "" {
			if parentID, ok := groupMap[strings.TrimSpace(req.ParentCode)]; ok && parentID != 0 {
				group.ParentID = &parentID
			}
		}

		groupsToSave = append(groupsToSave, group)
	}

	if err := svc.repo.UpsertGroupsBatch(groupsToSave, 500); err != nil {
		return dto.ImportResult{Status: "error", Message: "Eroare la salvarea grupelor: " + err.Error()}
	}

	return dto.ImportResult{
		Status:         "success",
		TotalProcessed: len(groupsToSave),
		TotalSkipped:   len(skipped),
		Message:        "Sincronizare grupe {{.Entity.Name}} finalizată!",
	}
}
{{end}}
`

var handlerTpl = `package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victoryus84/gorders/internal/dto"
	"github.com/victoryus84/gorders/internal/logger"
	"github.com/victoryus84/gorders/internal/service"
	"github.com/victoryus84/gorders/internal/utils"
)

type {{.Entity.Name}}Handler struct {
	svc service.{{.Entity.Name}}Service
}

func New{{.Entity.Name}}Handler(svc service.{{.Entity.Name}}Service) *{{.Entity.Name}}Handler {
	return &{{.Entity.Name}}Handler{svc: svc}
}

// Înregistrarea rutelor pentru Uber Fx (module.go)
func Register{{.Entity.Name}}Routes(r *gin.Engine, h *{{.Entity.Name}}Handler) {
	api := r.Group("/api/{{.Entity.LowerName}}s")
	{
		api.POST("/import", h.Create)
{{- if .Entity.HasGroup}}
		api.POST("/groups/import", h.CreateGroup)
{{- end}}
	}
}

// Sincronizare Entități Principale
func (h *{{.Entity.Name}}Handler) Create(c *gin.Context) {
	requests, err := utils.ParseBody[dto.{{.Entity.Name}}DTO](c)
	if err != nil {
		logger.LogError("⚠️ Format invalid primit de la 1C pentru {{.Entity.Name}}", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format invalid: " + err.Error()})
		return
	}

	result := h.svc.ProcessImport(requests)
	
	status := http.StatusCreated
	if result.Status == "error" {
		status = http.StatusInternalServerError
	}
	c.JSON(status, result)
}

{{if .Entity.HasGroup}}
// Sincronizare Grupe
func (h *{{.Entity.Name}}Handler) CreateGroup(c *gin.Context) {
	requests, err := utils.ParseBody[dto.{{.Entity.Name}}GroupDTO](c)
	if err != nil {
		logger.LogError("⚠️ Format invalid primit de la 1C pentru Grupe de {{.Entity.Name}}", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format invalid: " + err.Error()})
		return
	}

	result := h.svc.ProcessGroupImport(requests)
	
	status := http.StatusCreated
	if result.Status == "error" {
		status = http.StatusInternalServerError
	}
	c.JSON(status, result)
}
{{end}}
`

var repositoryTpl = `package repository

import (
	"github.com/victoryus84/gorders/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 1. INTERFAȚA
type {{.Entity.Name}}Repository interface {
	UpsertBatch(items []models.{{.Entity.Name}}, batchSize int) error
	FindAll() ([]models.{{.Entity.Name}}, error)
{{- if .Entity.HasGroup}}
	UpsertGroupsBatch(groups []models.{{.Entity.Name}}Group, batchSize int) error
	GetGroupIDMap() (map[string]uint, error)
{{- end}}
}

// 2. STRUCTURA PRIVATĂ
type {{.Entity.LowerName}}Repository struct {
	db *gorm.DB
}

// 3. CONSTRUCTORUL UBER FX
func New{{.Entity.Name}}Repository(db *gorm.DB) {{.Entity.Name}}Repository {
	return &{{.Entity.LowerName}}Repository{db: db}
}

// 4. METODELE PRINCIPALE
func (r *{{.Entity.LowerName}}Repository) UpsertBatch(items []models.{{.Entity.Name}}, batchSize int) error {
	return r.db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(items, batchSize).Error
}

func (r *{{.Entity.LowerName}}Repository) FindAll() ([]models.{{.Entity.Name}}, error) {
	var items []models.{{.Entity.Name}}
	err := r.db.Find(&items).Error
	return items, err
}

{{if .Entity.HasGroup}}
// 5. METODELE PENTRU GRUPE
func (r *{{.Entity.LowerName}}Repository) UpsertGroupsBatch(groups []models.{{.Entity.Name}}Group, batchSize int) error {
	return r.db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(groups, batchSize).Error
}

func (r *{{.Entity.LowerName}}Repository) GetGroupIDMap() (map[string]uint, error) {
	var groups []models.{{.Entity.Name}}Group
	if err := r.db.Select("id", "code").Find(&groups).Error; err != nil {
		return nil, err
	}
	
	groupMap := make(map[string]uint)
	for _, g := range groups {
		groupMap[g.Code] = g.ID
	}
	return groupMap, nil
}
{{end}}
`