package chatapp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Natnael-Alemayehu/chat_clone/chat/app/sdk/errs"
	"github.com/Natnael-Alemayehu/chat_clone/chat/foundation/logger"
	"github.com/Natnael-Alemayehu/chat_clone/chat/foundation/web"
	"github.com/gorilla/websocket"
)

type app struct {
	log *logger.Logger
}

func newApp(log *logger.Logger) *app {
	return &app{
		log: log,
	}
}

var upgrader = websocket.Upgrader{}

func (a *app) connect(ctx context.Context, r *http.Request) web.Encoder {

	c, err := upgrader.Upgrade(web.GetWriter(ctx), r, nil)
	if err != nil {
		return errs.Newf(errs.FailedPrecondition, "Failed to connect to websocker: %v", err)
	}

	_, err = a.handshake(ctx, c)
	if err != nil {
		return errs.Newf(errs.FailedPrecondition, "Failed to handshake: %v", err)
	}

	return web.NewNoResponse()
}

func (a *app) handshake(ctx context.Context, c *websocket.Conn) (user, error) {

	a.log.Info(ctx, "Handshake Started")
	defer a.log.Info(ctx, "Handshake Ended")

	err := c.WriteMessage(websocket.TextMessage, []byte("Hello"))
	if err != nil {
		return user{}, errs.Newf(errs.FailedPrecondition, "Failed to send hello message: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var usr user

	msg, err := a.readmessage(ctx, c)
	if err != nil {
		return user{}, errs.Newf(errs.FailedPrecondition, "Error Reading message: %v", err)
	}
	err = json.Unmarshal(msg, &usr)
	if err != nil {
		return user{}, errs.Newf(errs.InvalidArgument, "Error unmarshaling data: %v", err)
	}

	welcome_msg := fmt.Sprintf("Welcome %v", usr.Name)
	err = c.WriteMessage(websocket.TextMessage, []byte(welcome_msg))
	if err != nil {
		return user{}, errs.Newf(errs.FailedPrecondition, "Error Writing message: %v", err)
	}

	return usr, nil

}

func (a *app) readmessage(ctx context.Context, c *websocket.Conn) ([]byte, error) {

	type response struct {
		msg []byte
		err error
	}

	ch := make(chan response, 1)

	func() {
		a.log.Info(ctx, "Statred Reading message")
		defer a.log.Info(ctx, "Message Reading Ended")
		_, mes, err := c.ReadMessage()
		if err != nil {
			ch <- response{nil, err}
		}

		ch <- response{mes, nil}

	}()

	var resp response

	select {
	case <-ctx.Done():
		c.Close()
		return nil, errs.Newf(errs.DeadlineExceeded, "Context Deadline Reached")
	case resp = <-ch:
		return resp.msg, nil
	}

}
