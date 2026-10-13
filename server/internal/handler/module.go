package handler

import "go.uber.org/fx"

// Module conține toți handlerii care au un constructor simplu.
// ATENȚIE: NewCoreHandler NU este aici, pentru că îl furnizăm manual în main.go!
var Module = fx.Options(
    // 1. Aici doar oferim constructorii
    fx.Provide(
        NewCoreHandler,
        NewClientHandler,
        NewContractHandler,
        NewProductHandler,
        NewUserHandler,
    ),
    // 2. Aici îi spunem lui Fx să execute funcțiile de mapare a rutelor!
    fx.Invoke(
        RegisterCoreRoutes,
        RegisterClientRoutes,
        RegisterContractRoutes,
        RegisterProductRoutes,
        RegisterUserRoutes,
    ),
)