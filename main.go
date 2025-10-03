package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"

	"golang.org/x/term"
	"gopkg.in/yaml.v3"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/medidew/ApplicationTracker/types"
)

const LOG_TO_CLI bool = true

// Loads the .yaml config file from `path`.
func loadConfig(path string) (types.Config, error) {
	cfg := types.Config{}

	config_file, err := os.ReadFile(path)
	if err != nil { return cfg, err }

	err = yaml.Unmarshal(config_file, &cfg)
	if err != nil { return cfg, err }

	return cfg, nil
}

func main() {

	/*
		Setup log
	*/

	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename: "./logs/app.log",
		MaxSize: 10,
		MaxBackups: 2,
		MaxAge: 14,
		Compress: false,
	})

	encoder_cfg := zap.NewDevelopmentEncoderConfig()

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoder_cfg),
		ws,
		zap.InfoLevel,
	)

	logger := zap.New(core)
	defer logger.Sync()

	sugar := logger.Sugar()

	sugar.Info("Logger Constructed")
	sugar.Info("App Starting")
	
	/*
		DB Connect
	*/

	database_url := os.Getenv("DATABASE_URL")

	if database_url == "" {
		cfg, err := loadConfig("config.yaml")
		if err != nil { sugar.Panic(err) }
		sugar.Info("Config file loaded")

		fmt.Println("Enter the database password, then press enter.")
		fmt.Println("The password will not be visible as you type it:")
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil { sugar.Panic(err) }
		sugar.Info("Password accepted")

		database_url = "postgresql://" + cfg.Database.User + ":" + string(password) + "@" + cfg.Database.Host + ":" + cfg.Database.Port + "/" + cfg.Database.Name
	} else {
		sugar.Info("DATABASE_URL env var found")
	}

	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil { sugar.Panic(err) }
	defer conn.Close(context.Background())

	row, err := conn.Query(context.Background(), "elect company from applications where status=$1", 1)
	if err != nil { sugar.Panic(err) }
	defer row.Close()

	test, err := types.NewJobApplication("company", types.SoftwareEngineer, types.Active)
	if err != nil { sugar.Panic(err) }
	test.AddNote("badabing")
	test.AddNote("badaboom")
	test.AddNote("b b")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
		log.Printf("test: %v\n", test)
	})

	http.ListenAndServe(":3000", r)
}