package gateway

import (
	"fmt"
	"net/http"
	"websocket-service/utils"

	log "github.com/sirupsen/logrus"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

func CreateServer() *socketio.Server {
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
			&websocket.Transport{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		},
	})

	// General namespace for all events
	const namespace = "/"

	// Event handlers for the main namespace
	server.OnEvent(namespace, "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})
	server.OnConnect(namespace, func(s socketio.Conn) error {
		s.SetContext("main")
		remoteAddr := s.RemoteAddr().String()
		utils.InfoLog(nil, fmt.Sprintf("[Websocket] Main Client:%s Connected: %s", remoteAddr, s.ID()))
		s.Join("global_updates")
		return nil
	})
	server.OnError(namespace, func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})
	server.OnDisconnect(namespace, func(s socketio.Conn, reason string) {
		remoteAddr := s.RemoteAddr().String()
		log.Printf("Client:%s Disconnected: %s, Reason: %s", remoteAddr, s.ID(), reason)
	})

	// Specific namespace for /order events
	server.OnConnect("/order", func(s socketio.Conn) error {
		s.SetContext("order")
		remoteAddr := s.RemoteAddr().String()
		utils.InfoLog(nil, fmt.Sprintf("[Websocket] Order Client:%s Connected: %s", remoteAddr, s.ID()))
		s.Join("global_join_order")
		return nil
	})
	server.OnEvent("/order", "global_join_order", func(s socketio.Conn, data map[string]interface{}) {
		log.Infof("[Websocket] : Received event: %v", data)
		server.BroadcastToRoom("/order", "global_join_order", "ws:new_data", map[string]interface{}{
			"type": "global_join_order",
		})
	})

	return server
}
