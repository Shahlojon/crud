package main

import (
	"github.com/gorilla/mux"
	"time"
	"net"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/Shahlojon/crud/cmd/app"
	"github.com/Shahlojon/crud/pkg/customers"
	"net/http"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"os"
	"go.uber.org/dig"
)

func main() {
	host :="0.0.0.0"
	port := "9999"
	dbConnectionString :="postgres://app:pass@localhost:5432/db"
	if err := execute(host, port, dbConnectionString); err != nil{
		log.Print(err)
		os.Exit(1)
	}
}

func execute(host, port, dbConnectionString string) (err error){
	deps := []interface{}{
		app.NewServer,
		mux.NewRouter,
		func() (*pgxpool.Pool, error) {
			ctx,_:=context.WithTimeout(context.Background(), time.Second*5)
			return pgxpool.Connect(ctx, dbConnectionString)
		},
		customers.NewService,
		func (server *app.Server) *http.Server{
			return &http.Server{
				Addr: net.JoinHostPort(host, port),
				Handler: server,
			}
		},
	}

	container:=dig.New()
	for _,dep:=range deps{
		err = container.Provide(dep)
		if err!=nil {
			return err
		}
	}
	err = container.Invoke(func(server *app.Server) {
		server.Init()
	})

	if err!=nil {
		return err
	}

	return container.Invoke(func(server *http.Server) error {
		return server.ListenAndServe()
	})
}