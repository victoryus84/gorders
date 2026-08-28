package repository

import (
	"github.com/victoryus84/gorders/internal/models" // Aici ai structurile (Product, Client, etc)
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Module grupează toate repozitoriile pentru Uber Fx
var Module = fx.Provide(
	// ------------------------------------------------------------------------
	// 1. REPO-URILE SPECIFICE (Echipamentul greu din rep_client.go, rep_contract.go)
	// Le trecem pur și simplu pe nume, Fx se prinde singur ce au nevoie.
	// ------------------------------------------------------------------------
	NewClientRepository,
	NewContractRepository,

	// ------------------------------------------------------------------------
	// 2. REPO-URILE GENERICE (Echipamentul standard pentru tabele simple)
	// Aici punem DOAR obiectele pentru care NU am creat un fișier rep_... dedicat.
	// ------------------------------------------------------------------------

	// Ex: Grupuri de clienți (pe care l-am cerut în ClientService mai devreme)
	func(db *gorm.DB) *GenericRepo[models.ClientGroup] { return NewGenericRepo[models.ClientGroup](db) },

	// Ex: Adrese
	func(db *gorm.DB) *GenericRepo[models.ClientAddress] { return NewGenericRepo[models.ClientAddress](db) },

	// Ex: Când vei face Produse pe viitor, doar decomentezi linia asta:
	// func(db *gorm.DB) *GenericRepo[models.Product] { return NewGenericRepo[models.Product](db) },
)
