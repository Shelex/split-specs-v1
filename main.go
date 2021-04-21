package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Shelex/split-specs/api/graph"
	"github.com/Shelex/split-specs/api/graph/generated"
	"github.com/Shelex/split-specs/domain"
	"github.com/Shelex/split-specs/internal/auth"
	"github.com/Shelex/split-specs/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const defaultPort = "8080"

func lineSeparator() {
	fmt.Println("========================")
}

func startMessage(port string) {
	lineSeparator()
	log.Printf("connect to http://localhost:%s/playground for GraphQL playground\n", port)
	lineSeparator()
}

func main() {
	if err := Start(); err != nil {
		os.Exit(1)
	}
}

func Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
		log.Printf("Defaulting to port %s", port)
	}

	db, err := InitDb()
	if err != nil {
		return fmt.Errorf("failed to initialize db: %s", err)
	}

	svc := domain.NewSplitService(db)

	router := chi.NewRouter()
	router.Use(auth.Middleware(), middleware.Logger)

	gql := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: graph.NewResolver(svc),
	}))

	gql.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		return fmt.Errorf("internal server error: %s", err)
	})

	router.Handle("/query", (gql))
	router.Handle("/playground", playground.Handler("GraphQL playground", "/query"))

	startMessage(port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Printf("error %s:\n", err)
		return err
	}
	return nil
}

func InitDb() (storage.Storage, error) {
	env := os.Getenv("ENV")

	var repo storage.Storage
	var err error

	switch env {
	case "dev":
		repo, err = storage.NewInMemStorage()
	default:
		repo, err = storage.NewDataStore()
	}
	if err != nil {
		return nil, err
	}

	return repo, nil
}
