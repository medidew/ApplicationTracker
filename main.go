package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"

	"github.com/medidew/ApplicationTracker/types"
)

type Config struct {
	Database struct {
		User string `yaml:"user`
		Host string `yaml:"host`
		Port string `yaml:"port`
		Name string `yaml:"name`
	} `yaml:"database`
}

func main() {
	database_url := os.Getenv("DATABASE_URL")
	
	if database_url == "" {
		cfg, err := loadConfig("config.yaml")
		if err != nil { errAndExit(err, 1, "failed to open config.yaml") }

		fmt.Println("Enter the database password, then press enter.")
		fmt.Println("The password will not be visible as you type it:")
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil { errAndExit(err, 1, "failed to scan password from input") }

		database_url = "postgresql://" + cfg.Database.User + ":" + string(password) + "@" + cfg.Database.Host + ":" + cfg.Database.Port + "/" + cfg.Database.Name
	}

	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil { errAndExit(err, 1, "failed to connect to db") }
	defer conn.Close(context.Background())

	var id int
	err = conn.QueryRow(context.Background(), "select id from test where id=$1", 1).Scan(&id)
	if err != nil { errAndExit(err, 1, "query failed") }

	test, err := types.NewJobApplication("company", types.SoftwareEngineer, types.Active)
	if err != nil { errAndExit(err, 1, "failed to create job application") }

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
		fmt.Printf("id: %v\n", id)
		fmt.Printf("test: %v\n", test)
	})

	http.ListenAndServe(":3000", r)
}

func loadConfig(path string) (Config, error) {
	cfg := Config{}

	config_file, err := os.ReadFile(path)
	if err != nil { return cfg, err }

	err = yaml.Unmarshal(config_file, &cfg)
	if err != nil { return cfg, err }

	return cfg, nil
}

// Prints the given error to terminal and exits with `code`.
func errAndExit(err error, code int, msg string) {
	if msg != "" {
		fmt.Printf("%v: %v\n", msg, err)
	} else {
		fmt.Printf("%v\n", err)
	}
	os.Exit(code)
}