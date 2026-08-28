package handler

import (
	"github.com/victoryus84/gorders/internal/config"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Module grupează absolut tot ce ține de Handlers.
var Module = fx.Provide(
	// 1. Constructorii standard
	NewUserHandler,
	NewClientHandler,
	NewContractHandler,

	// 2. Cazul special: CoreHandler
	// În main.go aveai nevoie de variabilele globale Version și Commit.
	// Acum le extragem elegant direct din config.Config!
	func(db *gorm.DB, cfg *config.Config) *CoreHandler {
		return NewCoreHandler(db, cfg.Version, cfg.Commit)
	},

	// 3. Piesa de rezistență: Agregatorul
	// Când Routerul va cere "Dă-mi toți handlerii", Fx va apela funcția asta.
	func(core *CoreHandler, user *UserHandler, client *ClientHandler, contract *ContractHandler) *Handlers {
		return &Handlers{
			Core:     core,
			User:     user,
			Client:   client,
			Contract: contract,
		}
	},
)
