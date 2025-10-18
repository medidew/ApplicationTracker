package handlers

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type App struct {
	DB     *pgxpool.Pool
	Logger *zap.SugaredLogger
}
