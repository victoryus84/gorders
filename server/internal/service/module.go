package service

import "go.uber.org/fx"

// Module grupează toți constructorii din stratul de Service.
// De fiecare dată când adaugi un domeniu nou (ex: ProductService),
// îl treci doar pe lista asta, și Fx se ocupă de restul.
var Module = fx.Provide(
	NewUserService,
	NewClientService,
	NewContractService,
	// NewProductService, // <-- Așa de simplu va fi în viitor
)
