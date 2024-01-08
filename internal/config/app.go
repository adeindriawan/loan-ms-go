package config

import (
	"database/sql"
)

type BootstrapConfig struct {
	DB *sql.DB
}

func Bootstrap(config *BootstrapConfig) {
	
}