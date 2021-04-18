package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Shelex/split-specs/api/graph"
	"github.com/Shelex/split-specs/api/graph/generated"
	"github.com/Shelex/split-specs/domain"
	"github.com/Shelex/split-specs/storage"
)

const defaultPort = "8080"

func lineSeparator() {
	fmt.Println("========================")
}

func startMessage(port string) {
	lineSeparator()
	log.Printf("connect to http://localhost:%s/ for GraphQL playground\n", port)
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

	svc, err := InitDomainService()

	if err != nil {
		log.Printf("failed to initiate service %s:\n", err)
		return err
	}

	gql := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: graph.NewResolver(svc),
	}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", gql)

	startMessage(port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("error %s:\n", err)
		return err
	}
	return nil
}

func InitDomainService() (domain.SplitService, error) {
	env := os.Getenv("ENV")

	var repo storage.Storage
	var err error

	switch env {
	case "dev":
		repo, err = storage.NewInMemStorage()
	default:
		//TODO change to datastore
		repo, err = storage.NewInMemStorage()
	}
	if err != nil {
		return domain.SplitService{}, err
	}
	return domain.NewSplitService(repo), nil
}
