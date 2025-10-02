package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"

	"github.com/medidew/ApplicationTracker/types"
)

func main() {

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var id int
	err = conn.QueryRow(context.Background(), "select id from test where id=$1", 1).Scan(&id)

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	test, err := types.NewJobApplication("company", types.SoftwareEngineer, types.Active)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
		fmt.Printf("id: %v\n", id)
		fmt.Printf("test: %v\n", test)
	})

	http.ListenAndServe(":3000", r)
}