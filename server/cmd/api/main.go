package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/victoryus84/gorders/internal/config"
	"github.com/victoryus84/gorders/internal/database"
	"github.com/victoryus84/gorders/internal/handler"
	"github.com/victoryus84/gorders/internal/kafka"
	"github.com/victoryus84/gorders/internal/logger"
	"github.com/victoryus84/gorders/internal/service"
	"github.com/victoryus84/gorders/internal/repository"
	"go.uber.org/fx"
)

// Build flags
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func main() {
	// 1. INIȚIALIZĂRI PRE-FX (Config & Logger)
	// Vrem ca Zap să pornească ÎNAINTE de Fx, ca să putem loga eventualele erori de start.
	_ = godotenv.Load()
	cfg := config.Load()

	Version = cfg.Version
	Commit = cfg.Commit

	if cfg.AppEnv == "production" {
		logger.Init("info")
		logger.LogInfo("⚠️ RUNNING IN PRODUCTION MODE", logger.String("version", Version), logger.String("commit", Commit))
	} else {
		logger.Init("debug")
		logger.LogInfo("🛠️ Running in Development mode", logger.String("version", Version))
	}
	defer logger.Logger.Sync()

	printBanner()

	// 2. MAGIA UBER FX
	fx.New(
		// --- A. PROVIDERS (Toate componentele tale) ---
		fx.Provide(
			// 1. Configurația (îi dăm direct variabila deja încărcată mai sus)
			func() *config.Config { return cfg },

			// 2. Baza de date
			database.Connect,

			// 3. Kafka (Cu tot cu funcția de închidere la oprirea serverului!)
			func(lc fx.Lifecycle, c *config.Config) *kafka.Producer {
				kp := kafka.NewProducer(c.KafkaAddr)
				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						logger.LogInfo("🛑 Închidem conexiunea Kafka...")
						kp.Close()
						return nil
					},
				})
				return kp
			},

        	repository.Module,
        	service.Module,
        	handler.Module,
		),

		// --- B. INVOKE (Pornirea efectivă) ---
		fx.Invoke(startHTTPServer),
	).Run()
}

func startHTTPServer(lc fx.Lifecycle, r *gin.Engine, cfg *config.Config) {
	// ❌ Nu mai avem nevoie de router.SetupRoutes(r, allHandlers)
	// Rutele sunt deja înregistrate de funcțiile Register*Routes din module.go!

	srv := &http.Server{
		Addr:    ":8080", // Opțional: poți trage și portul din cfg.AppPort dacă îl ai
		Handler: r,
	}

	// Îi spunem lui Fx cum să pornească serverul fără să blocheze restul proceselor
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.LogInfo("🎯 Server starting",
				logger.String("port", "8080"),
				logger.String("env", cfg.AppEnv),
				logger.String("version", Version), // Dacă Version nu e global, îl poți pune în config
				logger.String("commit", Commit),
			)

			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.LogFatal("Server failed to start", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.LogInfo("🛑 Oprim serverul HTTP grațios...")
			return srv.Shutdown(ctx)
		},
	})
}

// printBanner prints startup banner
func printBanner() {
	banner := fmt.Sprintf(`
╔════════════════════════════════════════╗
║     🚀 GOrders Backend Server 🚀      ║
╠════════════════════════════════════════╣
║  Version:  %-26s  					 ║
║  Commit:   %-26s  					 ║
║  Built:    %-26s  					 ║
╚════════════════════════════════════════╝
`, Version, Commit, BuildTime)

	logger.LogInfo(banner)
}
