package graph

import (
	"sleeve/ent"
	"sleeve/usecase/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Client              *ent.Client
	RegisterUserUseCase *user.RegisterUserUseCase
}
