package server

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const graphqlUrl = "/graphql"

// Options required to start the server.
type Options struct {
	Port          string
	Schema        graphql.ExecutableSchema
	AntiCSRF      bool
	Playground    bool
	Introspection bool
}

type Server struct {
	Echo *echo.Echo
	Port string
}

// Initializes a new graphql server.
func Init(options Options) Server {
	server := echo.New()
	gqlHandler := customNewServer(options.Schema, options.Introspection)
	playgroundHandler := getPlaygroundHandler(options.Playground)

	server.Use(middleware.Recover())

	if options.AntiCSRF {
		server.Use(middleware.CSRF())
	}

	server.GET(graphqlUrl, func(ctx echo.Context) error {
		if ctx.Request().Header.Get("Upgrade") == "websocket" {
			gqlHandler.ServeHTTP(ctx.Response(), ctx.Request())

			return nil
		}

		if playgroundHandler != nil {
			playgroundHandler.ServeHTTP(ctx.Response(), ctx.Request())

			return nil
		}

		return echo.ErrNotFound
	})

	server.POST(graphqlUrl, func(ctx echo.Context) error {
		gqlHandler.ServeHTTP(ctx.Response(), ctx.Request())

		return nil
	})

	return Server{
		Echo: server,
		Port: options.Port,
	}
}

// Starts the server.
// Supports graceful shutdown.
func (server Server) Start() error {
	// Gratefully borrowed from https://echo.labstack.com/cookbook/graceful-shutdown
	// Start server
	go func() {
		if err := server.Echo.Start(":" + server.Port); err != nil {
			log.Fatalf("error starting server: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return server.Echo.Shutdown(ctx)
}

// Create a new server, allowing to disable introspection.
// Borrowed from https://github.com/99designs/gqlgen/blob/v0.13.0/handler/handler.go
func customNewServer(es graphql.ExecutableSchema, introspection bool) *handler.Server {
	server := handler.New(es)

	server.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	server.AddTransport(transport.Options{})
	server.AddTransport(transport.POST{})

	server.SetQueryCache(lru.New(1000))

	server.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	if introspection {
		server.Use(extension.Introspection{})
	}

	return server
}

func getPlaygroundHandler(enable bool) http.HandlerFunc {
	if !enable {
		return nil
	}

	return playground.Handler("GraphQL playground", graphqlUrl)
}
