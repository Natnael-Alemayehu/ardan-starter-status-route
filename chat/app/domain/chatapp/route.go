package chatapp

import (
	"net/http"

	"github.com/Natnael-Alemayehu/chat_clone/chat/foundation/logger"
	"github.com/Natnael-Alemayehu/chat_clone/chat/foundation/web"
)

type Config struct {
	Log *logger.Logger
}

func Routes(app *web.App, cfg Config) {
	const version = "v1"

	api := newApp()

	app.HandlerFunc(http.MethodGet, version, "/test", api.test)
}
