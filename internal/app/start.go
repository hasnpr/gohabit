package app

import (
	"log/slog"
)

func Start() {
	slog.Info("app started")
	// start simplevisor
	// start http server with business-app.WithDB.WithRedis.WithAnalyticsCMQ
	// handle simplevisor shutdown
}

// ctx := context.Background()
// http.NewHTTPServer(app.Sattar.
// 	WithDB().
// 	WithRedis().
// 	WithAnalyticsCMQ(),
// ).Start(ctx)
