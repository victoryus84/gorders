package repository

import (
	"go.uber.org/fx"
)

// Module grupează toate repozitoriile pentru Uber Fx
var Module = fx.Provide(
	// ------------------------------------------------------------------------
	// 1. REPO-URILE SPECIFICE (Echipamentul greu din rep_client.go, rep_contract.go)
	// Le trecem pur și simplu pe nume, Fx se prinde singur ce au nevoie.
	// ------------------------------------------------------------------------
	NewUserRepository,
	NewClientRepository,
	NewContractRepository,
	NewProductRepository,
	NewUnitRepository,
	NewVatTaxRepository,	
)
