package migrations

import (
	"fmt"

	"github.com/victoryus84/gorders/internal/logger"
	"gorm.io/gorm"
)

// DropUnusedColumns automatically detects and removes columns from DB that no longer exist in models
func DropUnusedColumns(db *gorm.DB) error {
	logger.LogInfo("\n🧹 Cleaning up orphaned columns from database...")

	tables := GetAllModels()

	dropCount := 0

	for _, table := range tables {
		stmt := &gorm.Statement{DB: db}
		stmt.Parse(table)

		tableName := stmt.Table
		migrator := db.Migrator()

		// Get columns from model
		modelColumns := make(map[string]bool)
		for _, field := range stmt.Schema.Fields {
			modelColumns[field.DBName] = true
		}

		// Get columns from database
		dbColumns, err := GetDBColumns(db, tableName)
		if err != nil {
			logger.LogError(fmt.Sprintf("⚠️  Could not read table %s", tableName), err)
			continue
		}

		// Find and drop orphaned columns
		for _, dbCol := range dbColumns {
			if !modelColumns[dbCol] {
				// This column exists in DB but not in model - drop it
				if err := migrator.DropColumn(table, dbCol); err != nil {
					logger.LogError(fmt.Sprintf("❌ Failed to drop column %s from %s:", dbCol, tableName), err)
				} else {
					logger.LogInfo(fmt.Sprintf("✅ Dropped orphaned column: %s.%s\n", tableName, dbCol))
					dropCount++
				}
			}
		}
	}

	if dropCount == 0 {
		logger.LogInfo("✅ No orphaned columns found - database is clean!")
	} else {
		logger.LogInfo(fmt.Sprintf("✅ Successfully removed %d orphaned column(s)\n", dropCount))
	}
	logger.LogInfo("")
	return nil
}
