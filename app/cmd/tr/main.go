package main

import (
	"log"
	"time"
)

//"github.com/owicca/tr/internal/infra"
//_ "github.com/owicca/tr/internal/routes/backend"
//_ "github.com/owicca/tr/internal/routes/frontend"
//_ "github.com/owicca/tr/internal/routes/middleware"

func main() {
	for {
		log.Println(time.Now().Unix())
		time.Sleep(time.Second)
	}
	//cfg, store, conn, logger := infra.Setup(path.Join(path.Dir(filename), "./config/config.json"))
	//infra.LoggerSync = logger.Sync
	//infra.Undo = zap.ReplaceGlobals(logger)

	//	infra.S = infra.NewServer(
	//		cfg,
	//		store,
	//		conn,
	//		NewTemplate(),
	//	)

	//infra.S.Run()
}
