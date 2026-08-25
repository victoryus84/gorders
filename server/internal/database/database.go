package database

import (
	"github.com/victoryus84/gorders/internal/config" // Importă pachetul tău de config
	"github.com/victoryus84/gorders/internal/logger"
	"github.com/victoryus84/gorders/internal/migrations"
	"github.com/victoryus84/gorders/internal/seeds"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger" // Îl importăm doar ca să-i dăm comanda de "Mute"
)

// Connect inițializează conexiunea la baza de date
func Connect(cfg *config.Config) *gorm.DB {

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		// 3. Folosim porecla "gormlogger" în loc de "logger"
		Logger: gormlogger.Default.LogMode(gormlogger.Error),
	})

	if err != nil {
		logger.LogFatal("❌ failed to connect database", err)
	}
	logger.LogInfo("✅ Database connected")

	// Migrate schema
	if err := db.AutoMigrate(migrations.GetAllModels()...); err != nil {
		logger.LogFatal("❌ migration failed", err)
	}
	logger.LogInfo("✅ Migration completed successfully")

	// Analyze schema differences
	migrations.AnalyzeSchemaSync(db)
	migrations.PrintSyncCommands(db)

	// Clean up orphaned columns
	if err := migrations.DropUnusedColumns(db); err != nil {
		logger.LogFatal("❌ cleanup failed", err)
	}
	logger.LogInfo("✅ Cleanup completed successfully")

	// Seed initial data
	seeds.RunAllSeeds(db)
	logger.LogInfo("✅ Seeding completed")

	return db
}
