package main

import (
	"log"
	"net/http"
	"os"
	"sleeve/graph"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

const (
	defaultPort             = "8080"
	queryCacheSize          = 1000
	persistedQueryCacheSize = 100
	serverReadTimeout       = 15 * time.Second
	serverWriteTimeout      = 15 * time.Second
	serverIdleTimeout       = 60 * time.Second
	serverReadHeaderTimeout = 10 * time.Second
)

func main() {
	var port string

	port = os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](queryCacheSize))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](persistedQueryCacheSize),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	server := &http.Server{
		Addr:              ":" + port,
		ReadTimeout:       serverReadTimeout,
		WriteTimeout:      serverWriteTimeout,
		IdleTimeout:       serverIdleTimeout,
		ReadHeaderTimeout: serverReadHeaderTimeout,
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(server.ListenAndServe())
}
