package app

import (
	"database/sql"

	"github.com/ariiiiph/ecommerce/internal/config"
	"github.com/ariiiiph/ecommerce/internal/redis"
)

type Dependencies struct {
	Config *config.Config
	DB     *sql.DB
	Redis  *redis.Client
}
