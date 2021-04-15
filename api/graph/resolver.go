package graph

import (
	"github.com/Shelex/split-test/domain"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	SplitService domain.SplitService
}

func NewResolver(svc domain.SplitService) *Resolver {
	return &Resolver{
		SplitService: svc,
	}

}
