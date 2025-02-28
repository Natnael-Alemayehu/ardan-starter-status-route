package mux

import (
	"context"
	"net/http"

	"github.com/Natnael-Alemayehu/chat_clone/chat/app/domain/chatapp"
	"github.com/Natnael-Alemayehu/chat_clone/chat/app/sdk/mid"
	"github.com/Natnael-Alemayehu/chat_clone/chat/foundation/logger"
	"github.com/Natnael-Alemayehu/chat_clone/chat/foundation/web"
)

type Config struct {
	Log *logger.Logger
}

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(cfg Config) http.Handler {

	logger := func(ctx context.Context, msg string, args ...any) {
		cfg.Log.Info(ctx, msg, args...)
	}

	app := web.NewApp(
		logger,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Panics(),
	)

	chatapp.Routes(app, chatapp.Config{
		Log: cfg.Log,
	})

	return app
}
