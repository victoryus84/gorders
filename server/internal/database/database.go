package database

import (
	"log"
	"github.com/victoryus84/gorders/internal/config" // Importă pachetul tău de config
	"github.com/victoryus84/gorders/migrations"
	"github.com/victoryus84/gorders/internal/seeds"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect inițializează conexiunea la baza de date
func Connect(cfg *config.Config) *gorm.DB {

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		// În main e ok să folosim log.Fatal dacă baza e critică pentru aplicație
		log.Fatal("❌ failed to connect database:", err)
	}
	log.Println("✅ Database connected")

	// Migrate schema
	if err := db.AutoMigrate(migrations.GetAllModels()...); err != nil {
		log.Fatal("migration failed:", err)
	}
	log.Println("✅ Migration completed successfully")

	// Analyze schema differences
	migrations.AnalyzeSchemaSync(db)
	migrations.PrintSyncCommands(db)

	// Clean up orphaned columns
	if err := migrations.DropUnusedColumns(db); err != nil {
		log.Fatal("cleanup failed:", err)
	}
	log.Println("✅ Cleanup completed successfully")

	// Seed initial data
	seeds.RunAllSeeds(db)
    log.Println("✅ Seeding completed")

	return db
}