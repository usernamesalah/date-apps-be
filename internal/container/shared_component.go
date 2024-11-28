package container

import (
	"date-apps-be/infrastructure/config"
	"date-apps-be/infrastructure/database"

	"go.uber.org/zap"
)

type SharedComponent struct {
	Conf *config.Config
	Log  *zap.Logger
	DB   *database.DB
}
