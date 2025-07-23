package gateway

import (
	"encoding/json"
	"fmt" // Import fmt for Sprintf
	"net/http"
	"websocket-service/utils"

	log "github.com/sirupsen/logrus"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

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

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})
	server.OnConnect("/order", func(s socketio.Conn) error {
		s.SetContext("order")
		remoteAddr := s.RemoteAddr().String()
		utils.InfoLog(nil, fmt.Sprintf("[Websocket] Order Client:%s Connected: %s", remoteAddr, s.ID()))
		s.Join("global_join_order")
		return nil
	})
	server.OnEvent("/order", "global_join_order", func(s socketio.Conn, data map[string]interface{}) {
		log.Infof("[Websocket] : Received event: %v", data)
		server.BroadcastToRoom("/", "global_updates", "ws:new_data", map[string]interface{}{
			"type": "global_join_order",
		})
	})
	server.OnEvent("/order", "ws:join_order_channel", func(s socketio.Conn, data map[string]interface{}) {
		log.Infof("[Websocket] Received client ws:join_order_channel order event: %v", data)
		jsonData, err := json.Marshal(data)
		if err != nil {
			handleErrorAndRespond(s, "join_order_channel", data, err, "Internal server error: Failed to marshal data for validation.", "error")
			return
		}

		var payload WSJoinOrderChannelPayload
		err = json.Unmarshal(jsonData, &payload)
		if err != nil {
			handleErrorAndRespond(s, "join_order_channel", data, err, "Invalid data format for ws:join_order_channel.", "error")
			return
		}

		err = Validate.Struct(payload)
		if err != nil {
			handleErrorAndRespond(s, "join_order_channel", data, err, "Validation failed for ws:join_order_channel data.", "error")
			return
		}

		s.Join(payload.Channel)
		utils.InfoLog(data, fmt.Sprintf("Client %s joined channel : %s", s.ID(), payload.Channel))
		s.Emit(`ws:join_order_channel_response`, map[string]interface{}{"channel": payload.Channel, "status": "success"})
	})
	server.OnEvent("/order", "tableSelected", func(s socketio.Conn, data map[string]interface{}) {
		log.Infof("[Websocket] tableSelected event: %v", data)
		jsonData, err := json.Marshal(data)
		if err != nil {
			handleErrorAndRespond(s, "tableSelected", data, err, "Internal server error: Failed to marshal data for validation.", "error")
			return
		}
		var payload TableSelectedPayload
		err = json.Unmarshal(jsonData, &payload)
		if err != nil {
			handleErrorAndRespond(s, "tableSelected", data, err, "Invalid data format for tableSelected.", "error")
			return
		}

		err = Validate.Struct(payload)
		if err != nil {
			handleErrorAndRespond(s, "tableSelected", data, err, "Validation failed for tableSelected data.", "error")
			return
		}

		server.BroadcastToRoom("/order", "table_updates", "ws:tableStatusUpdated", map[string]interface{}{
			"tableId":   int(payload.TableID),
			"newStatus": payload.NewStatus,
			"reason":    "table selected by another client/pos",
		})
		utils.InfoLog(data, fmt.Sprintf("Broadcasted ws:tableStatusUpdated for table %v to 'table_updates' room.", payload.TableID))
	})
	server.OnEvent("/order", "billSelected", func(s socketio.Conn, data map[string]interface{}) {

		utils.InfoLog(data, "[Websocket] Event : billSelected")
		jsonData, err := json.Marshal(data)
		if err != nil {
			handleErrorAndRespond(s, "billSelected", data, err, "Internal server error: Failed to marshal data for validation.", "error")
			return
		}

		var payload BillSelectedPayload
		err = json.Unmarshal(jsonData, &payload)
		if err != nil {
			handleErrorAndRespond(s, "billSelected", data, err, "Invalid data format for billSelected.", "error")
			return
		}
		err = Validate.Struct(payload)
		if err != nil {
			handleErrorAndRespond(s, "billSelected", data, err, "Validation failed for billSelected data.", "error")
			return
		}
		server.BroadcastToRoom("/order", "table_updates", "ws:tableStatusUpdated", map[string]interface{}{
			"tableId":   int(payload.TableID),
			"invoiceId": payload.InvoiceID,
			"invoiceNo": payload.InvoiceNo, // Cast to int if it's supposed to be an int, otherwise keep as float
		})
		utils.InfoLog(data, fmt.Sprintf("[Websocket] Broadcasted ws:tableStatusUpdated to table_updates:%v : ", payload.TableID))

		//log.Infof("Broadcasted ws:tableStatusUpdated for table %v to 'table_updates' room.", tableID)
	})
	server.OnEvent("/order", "client:joinBill", func(s socketio.Conn, data map[string]interface{}) {
		utils.InfoLog(data, "[Websocket] Event : client:joinBill")
		jsonData, err := json.Marshal(data)
		if err != nil {
			handleErrorAndRespond(s, "client:joinBill", data, err, "Internal server error: Failed to marshal data for validation.", "error")
			return
		}
		var payload ClientJoinBillPayload
		err = json.Unmarshal(jsonData, &payload)
		if err != nil {
			handleErrorAndRespond(s, "client:joinBill", data, err, "Invalid data format for client:joinBill.", "error")
			return
		}
		err = Validate.Struct(payload)
		if err != nil {
			handleErrorAndRespond(s, "client:joinBill", data, err, "Validation failed for client:joinBill data.", "error")
			return
		}
		s.Join(payload.InvoiceNo)
		s.Emit("ws:joinBillResponse", map[string]interface{}{"invoiceNo": payload.InvoiceNo, "status": "success"})
		utils.InfoLog(data, fmt.Sprintf("[Websocket] Emit ws:joinBillResponse for invoice: %s", payload.InvoiceNo))

	})
	server.OnEvent("/order", "orderSaved", func(s socketio.Conn, data map[string]interface{}) {
		utils.InfoLog(nil, "[Websocket] Received orderSaved")
		jsonData, err := json.Marshal(data)
		if err != nil {
			handleErrorAndRespond(s, "orderSaved", data, err, "Internal server error: Failed to marshal data for validation.", "error")
			return
		}
		var payload OrderSavedPayload
		err = json.Unmarshal(jsonData, &payload)
		if err != nil {
			handleErrorAndRespond(s, "orderSaved", data, err, "Invalid data format for orderSaved.", "error")
			return
		}
		err = Validate.Struct(payload)
		if err != nil {
			handleErrorAndRespond(s, "orderSaved", nil, err, "Validation failed for orderSaved data.", "error")
			return
		}

		server.BroadcastToRoom("/order", payload.InvoiceNo, "ws:orderItemsUpdated", map[string]interface{}{
			"tableId":      int(payload.TableID),
			"invoiceId":    int(payload.InvoiceID),
			"invoiceNo":    payload.InvoiceNo,
			"updatedItems": payload.UpdatedItems,
			"clientId":     payload.ClientID,
		})
		utils.InfoLog(nil, fmt.Sprintf("[Websocket] Broadcasted ws:orderItemsUpdated to 'invoice_%v : ", payload.InvoiceNo))

	})
	server.OnEvent("/order", "orderProcessed", func(s socketio.Conn, data map[string]interface{}) {
		log.Infof("Received orderProcessed event: %v", data)
		jsonData, err := json.Marshal(data)
		if err != nil {
			handleErrorAndRespond(s, "orderProcessed", data, err, "Internal server error: Failed to marshal data for validation.", "error")
			return
		}

		var payload OrderProcessedPayload
		err = json.Unmarshal(jsonData, &payload)
		if err != nil {
			handleErrorAndRespond(s, "orderProcessed", data, err, "Invalid data format for orderProcessed.", "error")
			return
		}
		err = Validate.Struct(payload)
		if err != nil {
			handleErrorAndRespond(s, "orderProcessed", data, err, "Validation failed for orderProcessed data.", "error")
			return
		}
		server.BroadcastToRoom("/order", payload.InvoiceNo, "ws:billFinalized", map[string]interface{}{
			"invoiceId": int(payload.InvoiceID),
			"invoiceNo": payload.InvoiceNo,
			"tableId":   int(payload.TableID),
			"status":    payload.Status,
			"clientId":  payload.ClientID,
		})
		utils.InfoLog(data, fmt.Sprintf("Broadcasted ws:billFinalized for invoice %v, table %v to %s room.", payload.InvoiceID, payload.TableID, payload.InvoiceNo))

		server.BroadcastToRoom("/", "table_updates", "ws:tableStatusUpdated", map[string]interface{}{
			"tableId":   int(payload.TableID),
			"newStatus": payload.Status, // Should be 'available'
			"reason":    "order processed",
		})
		utils.InfoLog(data, fmt.Sprintf("Broadcasted ws:tableStatusUpdated for table %v to 'table_updates' room (order processed).", payload.TableID))
	})
	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		remoteAddr := s.RemoteAddr().String()
		log.Printf("Client:%s Disconnected: %s, Reason: %s", remoteAddr, s.ID(), reason)
	})

	return server
}
