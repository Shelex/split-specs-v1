package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Shelex/split-test/api/graph"
	"github.com/Shelex/split-test/api/graph/generated"
	"github.com/Shelex/split-test/domain"
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

	svc, err := domain.NewSplitService()

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
