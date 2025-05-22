package state

import (
	"github.com/marekmchl/aggreGATOR/internal/config"
	"github.com/marekmchl/aggreGATOR/internal/database"
)

type State struct {
	DB     *database.Queries
	Config *config.Config
}
