package graph

//go:generate go run github.com/99designs/gqlgen

import (
	"github.com/llamadeus/iot-logger/graph/generated"
	"github.com/llamadeus/iot-logger/graph/resolvers"
)

func New() generated.Config {
	return generated.Config{Resolvers: &resolvers.Resolver{}}
}
