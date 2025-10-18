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
func loadConfig(path string) (types.DBConfig, error) {
	cfg := types.DBConfig{}

	config_file, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(config_file, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func main() {

	/*
		Setup log
	*/

	log_file_write_syncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    10,
		MaxBackups: 2,
		MaxAge:     14,
		Compress:   false,
	})

	console_write_syncer := zapcore.Lock(os.Stdout)

	encoder_cfg := zap.NewDevelopmentEncoderConfig()

	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoder_cfg),
			log_file_write_syncer,
			zap.InfoLevel,
		),

		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoder_cfg),
			console_write_syncer,
			zap.WarnLevel,
		),
	)

	logger := zap.New(core)
	defer logger.Sync()

	sugared_logger := logger.Sugar()

	sugared_logger.Info("Logger constructed")
	sugared_logger.Info("App starting")

	/*
		DB Connect
	*/

	database_url := os.Getenv("DATABASE_URL")

	if database_url != "" {
		sugared_logger.Info("DATABASE_URL env var found, skipping dbconfig...")
	} else {
		cfg, err := loadConfig("configs/dbconfig.yaml")
		if err != nil {
			sugared_logger.Panic(err)
		}
		sugared_logger.Info("Config file loaded")

		fmt.Println("Enter the database password:")
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			sugared_logger.Panic(err)
		}
		sugared_logger.Info("Password accepted")

		database_url = "postgresql://" + cfg.Database.User + ":" + string(password) + "@" + cfg.Database.Host + ":" + cfg.Database.Port + "/" + cfg.Database.Name
	}

	pool, err := pgxpool.New(context.Background(), database_url)
	if err != nil {
		sugared_logger.Panic(err)
	}
	defer pool.Close()

	/*
		App
	*/

	app := &handlers.App{
		DB:     pool,
		Logger: sugared_logger,
	}

	router := chi.NewRouter()

	router.Use(ZapLoggerMiddleware(sugared_logger))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {

	})

	router.Route("/applications", func(router chi.Router) {
		router.Get("/", app.ListApplications)
		router.Get("/{companyID}", app.GetApplication)

		//r.Post("/", app.CreateApplication)
		//r.Put("/{companyID}", app.UpdateApplication)
		//r.Delete("/{companyID}", app.DeleteApplication)
	})

	http.ListenAndServe(":3000", router)
}
