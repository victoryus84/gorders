package migrations

import (
	"fmt"

	"github.com/victoryus84/gorders/internal/logger"
	"gorm.io/gorm"
)

// ColumnInfo holds information about a column
type ColumnInfo struct {
	Name    string
	Type    string
	InModel bool
	InDB    bool
}

// GetDBColumns retrieves column names from database for a specific table
func GetDBColumns(db *gorm.DB, tableName string) ([]string, error) {
	var columns []string
	result := db.Raw(`
		SELECT column_name 
		FROM information_schema.columns 
		WHERE table_name = ? 
		ORDER BY ordinal_position
	`, tableName).Scan(&columns)
	return columns, result.Error
}

// AnalyzeSchemaSync compares DB schema with models and reports discrepancies
func AnalyzeSchemaSync(db *gorm.DB) {
	// Aici ne folosim doar de mesaje statice
	logger.LogInfo("\n======= 🔍 DATABASE SCHEMA ANALYSIS =======")

	tables := GetAllModels()

	for _, table := range tables {
		analyzeTable(db, table)
	}

	logger.LogInfo("\n========================================")
}

func analyzeTable(db *gorm.DB, table interface{}) {
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(table)

	tableName := stmt.Table

	// Get columns from model - FILTRĂM DOAR COLOANELE REALE
	modelColumns := make(map[string]bool)
	for _, field := range stmt.Schema.Fields {
		// Dacă DBName e gol, înseamnă că e o relație virtuală (ex: Client Client)
		// sau un câmp ignorat, deci nu îl numărăm ca și coloană în DB
		if field.DBName != "" {
			modelColumns[field.DBName] = true
		}
	}

	// Get columns from database
	dbColumns, err := GetDBColumns(db, tableName)
	if err != nil {
		logger.LogFatal(fmt.Sprintf("❌ Error reading table %s", tableName), err)
		return
	}

	dbColumnsMap := make(map[string]bool)
	for _, col := range dbColumns {
		dbColumnsMap[col] = true
	}

	// Find discrepancies
	orphanedColumns := []string{}
	missingColumns := []string{}

	for col := range dbColumnsMap {
		if !modelColumns[col] {
			orphanedColumns = append(orphanedColumns, col)
		}
	}

	for col := range modelColumns {
		if !dbColumnsMap[col] {
			missingColumns = append(missingColumns, col)
		}
	}

	// Print report
	if len(orphanedColumns) == 0 && len(missingColumns) == 0 {
		msg := fmt.Sprintf("✅ %-20s - SYNCED (%d physical columns)", tableName, len(dbColumnsMap))
		logger.LogInfo(msg)
	} else {
		msg := fmt.Sprintf("⚠️  TABLE: %s (DB: %d | Model: %d)", tableName, len(dbColumnsMap), len(modelColumns))
		logger.LogWarn(msg)

		if len(orphanedColumns) > 0 {
			logger.LogWarn("  🗑️  ORPHANED (Există în DB, dar NU în Model):")
			for _, col := range orphanedColumns {
				logger.LogWarn(fmt.Sprintf("      - %s", col))
			}
		}

		if len(missingColumns) > 0 {
			logger.LogWarn("  ❌ MISSING (Lipsesc din DB, trebuie create):")
			for _, col := range missingColumns {
				logger.LogWarn(fmt.Sprintf("      - %s", col))
			}
		}
	}
}

// PrintSyncCommands generates DROP COLUMN commands for cleanup
func PrintSyncCommands(db *gorm.DB) {
	logger.LogInfo("\n======= 📋 AUTO-FIX COMMANDS =======")

	tables := GetAllModels()

	for _, table := range tables {
		stmt := &gorm.Statement{DB: db}
		stmt.Parse(table)

		tableName := stmt.Table

		// Get columns from model
		modelColumns := make(map[string]bool)
		for _, field := range stmt.Schema.Fields {
			if field.DBName != "" {
				modelColumns[field.DBName] = true
			}
		}

		// Get columns from database
		dbColumns, err := GetDBColumns(db, tableName)
		if err != nil {
			logger.LogError(fmt.Sprintf("❌ Error reading table %s", tableName), err)
			continue
		}

		// Find orphaned columns
		for _, col := range dbColumns {
			if !modelColumns[col] {
				cmd := fmt.Sprintf("db.Migrator().DropColumn(&models.%s{}, \"%s\")", TableNameToModel(tableName), col)
				logger.LogInfo(cmd)
			}
		}
	}

	logger.LogInfo("\n====================================")
}
