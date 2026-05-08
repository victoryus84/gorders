package migrations

import (
	"log"

	"gorm.io/gorm"
)

// DropUnusedColumns automatically detects and removes columns from DB that no longer exist in models
func DropUnusedColumns(db *gorm.DB) error {
	log.Println("\n🧹 Cleaning up orphaned columns from database...")

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
			log.Printf("⚠️  Could not read table %s: %v\n", tableName, err)
			continue
		}

		// Find and drop orphaned columns
		for _, dbCol := range dbColumns {
			if !modelColumns[dbCol] {
				// This column exists in DB but not in model - drop it
				if err := migrator.DropColumn(table, dbCol); err != nil {
					log.Printf("❌ Failed to drop column %s from %s: %v\n", dbCol, tableName, err)
				} else {
					log.Printf("✅ Dropped orphaned column: %s.%s\n", tableName, dbCol)
					dropCount++
				}
			}
		}
	}

	if dropCount == 0 {
		log.Println("✅ No orphaned columns found - database is clean!")
	} else {
		log.Printf("✅ Successfully removed %d orphaned column(s)\n", dropCount)
	}
	log.Println()

	return nil
}
