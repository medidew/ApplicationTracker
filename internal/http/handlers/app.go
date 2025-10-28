package handlers

import (
	"github.com/alexedwards/scs/v2"
	"github.com/medidew/ApplicationTracker/internal/store"
	"go.uber.org/zap"
)

type App struct {
	DB     store.Store
	Logger *zap.Logger
	SessionManager *scs.SessionManager
}