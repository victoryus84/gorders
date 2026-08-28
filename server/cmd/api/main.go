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
	"github.com/victoryus84/gorders/internal/middleware"
	"github.com/victoryus84/gorders/internal/router"
	"github.com/victoryus84/gorders/internal/service"
	"go.uber.org/fx"
	"gorm.io/gorm"
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

			// 4. Repositories & Services
			service.NewUserService,
			service.NewClientService,
			service.NewContractService,

			// 5. Handlers (Micile ajustări unde avem nevoie de mai mulți parametri)
			func(db *gorm.DB) *handler.CoreHandler {
				// CoreHandler cere Version și Commit, care sunt variabile globale aici
				return handler.NewCoreHandler(db, Version, Commit)
			},
			handler.NewUserHandler,
			handler.NewClientHandler,
			handler.NewContractHandler,

			// 6. Grupăm toți handlerii în structura ta "Handlers" (ca să o putem da la Router)
			func(core *handler.CoreHandler, user *handler.UserHandler, client *handler.ClientHandler, contract *handler.ContractHandler) *handler.Handlers {
				return &handler.Handlers{
					Core:     core,
					User:     user,
					Client:   client,
					Contract: contract,
				}
			},

			// 7. Gin Engine (Router-ul principal cu Middlewares)
			func(c *config.Config) *gin.Engine {
				if c.AppEnv == "production" {
					gin.SetMode(gin.ReleaseMode)
				}
				r := gin.New()
				r.Use(middleware.RequestLogging())
				r.Use(middleware.PanicRecovery())
				r.Use(middleware.CORS())
				r.Use(middleware.RateLimit())
				return r
			},
		),

		// --- B. INVOKE (Pornirea efectivă) ---
		fx.Invoke(startHTTPServer),
	).Run()
}

// Această funcție trage automat routerul, handlerii și config-ul din cutia Fx
func startHTTPServer(lc fx.Lifecycle, r *gin.Engine, allHandlers *handler.Handlers, cfg *config.Config) {
	// Setup API routes
	router.SetupRoutes(r, allHandlers)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Îi spunem lui Fx cum să pornească serverul fără să blocheze restul proceselor
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.LogInfo("🎯 Server starting",
				logger.String("port", "8080"),
				logger.String("env", cfg.AppEnv),
				logger.String("version", Version),
				logger.String("commit", Commit),
				logger.String("buildTime", BuildTime),
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
