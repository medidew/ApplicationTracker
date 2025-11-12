package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"golang.org/x/term"

	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"

	"github.com/medidew/ApplicationTracker/internal/http/handlers"
	"github.com/medidew/ApplicationTracker/internal/store"
)

const LOG_TO_CLI bool = true

// Loads the .yaml config file from `path`.
func loadConfig(path string) (store.DBConfig, error) {
	cfg := store.DBConfig{}

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

	logger.Info("Logger constructed")
	logger.Info("App starting")

	/*
		DB Connect
	*/

	database_url := os.Getenv("DATABASE_URL")

	if database_url != "" {
		logger.Info("DATABASE_URL env var found, skipping dbconfig...")
	} else {
		cfg, err := loadConfig("configs/dbconfig.yaml") // TODO: Fix for deploy
		if err != nil {
			logger.Panic(err.Error())
		}
		logger.Info("Config file loaded")

		fmt.Println("Enter the database password:")
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			logger.Panic(err.Error())
		}
		logger.Info("Password accepted")

		database_url = "postgresql://" + cfg.Database.User + ":" + string(password) + "@" + cfg.Database.Host + ":" + cfg.Database.Port + "/" + cfg.Database.Name
	}

	pool, err := pgxpool.New(context.Background(), database_url)
	if err != nil {
		logger.Panic(err.Error())
	}
	defer pool.Close()

	/*
		App
	*/
	session_manager := scs.New()
	session_manager.Lifetime = 12 * time.Hour
	session_manager.IdleTimeout = 1 * time.Hour
	session_manager.Cookie.Secure = true
	session_manager.Cookie.HttpOnly = true
	session_manager.Cookie.SameSite = http.SameSiteLaxMode

	app := &handlers.App{
		DB:     &store.DB{Pool: pool},
		Logger: logger,
		SessionManager: session_manager,
	}

	router := handlers.SetupRouter(app)

	logger.Info("Starting HTTP server on :4000")

	http.ListenAndServe(":4000", router)
}