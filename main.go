package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"golang.org/x/term"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"

	"github.com/medidew/ApplicationTracker/handlers"
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

	sugar.Info("Logger constructed")
	sugar.Info("App starting")
	
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

	conn, err := pgxpool.New(context.Background(), database_url)
	if err != nil { sugar.Panic(err) }
	defer conn.Close()

	/*
		App
	*/

	app := &handlers.App{DB: conn, Logger: sugar}

	r := chi.NewRouter()

	r.Use(ZapLoggerMiddleware(sugar))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

	})

	r.Route("/applications", func(r chi.Router) {
		r.Get("/", app.ListApplications)
		//r.Get("/{company}", handlers.GetApplication)

		//r.Post("/", handlers.CreateApplication)
		//r.Put("/{company}", handlers.UpdateApplication)
		//r.Delete("/{company}", handlers.DeleteApplication)
	})

	http.ListenAndServe(":3000", r)
}