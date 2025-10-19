package handlers

import (
	"go.uber.org/zap"
	"github.com/medidew/ApplicationTracker/internal/types"
)

type App struct {
	DB     types.Store
	Logger *zap.Logger
}