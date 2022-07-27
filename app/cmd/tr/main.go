package main

import (
	"path"

	"github.com/owicca/tr/internal/infra"
	"go.uber.org/zap"
)

//_ "github.com/owicca/tr/internal/routes/backend"
//_ "github.com/owicca/tr/internal/routes/frontend"
//_ "github.com/owicca/tr/internal/routes/middleware"

func main() {
	cfg, store, conn, logger := infra.Setup(path.Join(path.Dir(filename), "./config/config.json"))
	infra.LoggerSync = logger.Sync
	infra.Undo = zap.ReplaceGlobals(logger)

	infra.S = infra.NewServer(
		cfg,
		store,
		conn,
		infra.NewTemplate(),
	)

	infra.S.Run()
}
