package messagehandler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"ctraderapi/messages/github.com/Carlosokumu/messages"

	"ctraderapi/middlewares"
	"ctraderapi/models"
	"ctraderapi/service"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 60 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

var (
	CanAcessAccount *bool
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// writePump pumps messages from the hub to the  app's websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// // executing all writes from this goroutine.
func writePump(c *middlewares.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.Hub.Unregister <- c
		ticker.Stop()
		//c.Appconn.Close()
	}()

	for {
		select {

		case <-ticker.C:
			c.Appconn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Appconn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case accountModel, ok := <-c.Hub.AccountModelChannel:
			accountModelUsecase := models.AccountModelUseCase{}
			c.Appconn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Appconn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			accountModelUsecase.Balance = accountModel.Balance
			accountModelUsecase.Equity = accountModel.Equity
			accountModelUsecase.Positions = accountModel.Positions
			accountModelUsecase.Symbols = accountModel.Symbols

			err := c.Appconn.WriteJSON(accountModelUsecase)
			if err != nil {
				log.Fatal(err)
			}

			c.Hub.SubChannel <- accountModel

		}
	}

}

func ConnectToOpen(host string, port int, hub *middlewares.Hub, w http.ResponseWriter, r *http.Request) {

	Appconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Set up a dialer with the desired options
	dialer := websocket.DefaultDialer
	dialer.EnableCompression = true

	// Connect to the  Ctrader WebSocket endpoint
	url := fmt.Sprintf("wss://%s:%d", host, port)
	fmt.Println(url)
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &middlewares.Client{Hub: hub, Conn: conn, Protomessages: make(chan messages.ProtoMessage), Appconn: Appconn}

	hub.Register <- client
	appAuth := &service.AppAuth{}
	// accountAuth := &service.AccountAuth{}
	// appAuth.SetNext(accountAuth)

	appAuth.Execute(conn, hub)

	go service.ReadCtraderMessages(conn, *client)
	go writePump(client)
	service.CollectAllMessages(hub, conn, Appconn)
	go service.ListenSpots(hub, client)

}
