package handlers

import (
	"context"
	"io"
	"strings"

	"net/http"
	"time"

	"golang.org/x/text/encoding"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/mgutz/ansi"

	"github.com/flowqio/flowqlet/service"
	"github.com/flowqio/flowqlet/version"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second

	welcomeMsg = "Welcome to FlowQ, Interactive learning platform!\n\r"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSWriter struct {
	Conn *websocket.Conn
}

func (s *WSWriter) Write(p []byte) (n int, err error) {
	s.Conn.WriteMessage(websocket.TextMessage, p)
	return len(p), nil
}

func ServeWS(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	oid, hasOID := vars["oid"]
	cid, hasCID := vars["cid"]

	var sid = ""

	if !hasOID || !hasCID {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if _sid, ok := service.ContainerExits(oid, cid); !ok {
		log.Errorf("OID: %s , Container %s not exits", oid, cid)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	} else {
		sid = _sid
	}
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Errorf("upgrade: %s", err.Error())
		return
	}

	defer ws.Close()

	writer := &WSWriter{Conn: ws}
	ctx := context.Background()

	conn, err := service.CreateExecAttachConnection(ctx, cid, "")

	if err != nil {
		log.Error(err)
		ws.WriteMessage(websocket.TextMessage, []byte("Internal Server Error , Please contact service@flowq.io"))
		return
	}

	//print welcome message
	phosphorize := ansi.ColorFunc("green+h")

	// ws.WriteMessage(websocket.TextMessage, []byte(phosphorize(` ________  __                           ______  `+"\n\r")))
	// ws.WriteMessage(websocket.TextMessage, []byte(phosphorize(`/        |/  |                         /      \ `+"\n\r")))
	// ws.WriteMessage(websocket.TextMessage, []byte(phosphorize(`########/ ## |  ______   __   __   __ /######  |`+"\n\r")))
	// ws.WriteMessage(websocket.TextMessage, []byte(phosphorize(`## |__    ## | /      \ /  | /  | /  |## |  ## |`+"\n\r")))
	// ws.WriteMessage(websocket.TextMessage, []byte(phosphorize(`##    |   ## |/######  |## | ## | ## |## |  ## |`+"\n\r")))
	// ws.WriteMessage(websocket.TextMessage, []byte(phosphorize(`#####/    ## |## |  ## |## | ## | ## |## |_ ## |`+"\n\r")))
	// ws.WriteMessage(websocket.TextMessage, []byte(phosphorize(`## |      ## |## \__## |## \_## \_## |## / \## |`+"\n\r")))
	// ws.WriteMessage(websocket.TextMessage, []byte(phosphorize(`## |      ## |##    ##/ ##   ##   ##/ ## ## ##< `+"\n\r")))
	// ws.WriteMessage(websocket.TextMessage, []byte(phosphorize(`##/       ##/  ######/   #####/####/   ######  |`+"\n\r")))
	// ws.WriteMessage(websocket.TextMessage, []byte(phosphorize(`                                           ###/ `+"\n\r")))
	// ws.WriteMessage(websocket.TextMessage, []byte(phosphorize(`                                      ver 0.1.1 `+"\n\r")))
	// ws.WriteMessage(websocket.TextMessage, []byte(phosphorize(welcomeMsg)))

	ws.WriteMessage(websocket.TextMessage, []byte(phosphorize("Current Env: "+service.GetFlowqLetConfig().NodeID+"\n\r")))
	ws.WriteMessage(websocket.TextMessage, []byte(phosphorize("Scenario: "+sid+"\n\r")))
	ws.WriteMessage(websocket.TextMessage, []byte(phosphorize("flowqlet v"+version.Version+"\n\r\n\r\n\r")))

	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}

	encoder := encoding.Replacement.NewEncoder()
	go func() {
		_, err := io.Copy(encoder.Writer(writer), conn.Reader)
		if err != nil {
			log.Error(err)
		}
	}()
	//conn.Conn.Write([]byte("ls\n"))
	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			break
		}

		log.Debugf("Data: %s", string(data))
		if strings.Index(string(data), "ping") == -1 {
			_, err = conn.Conn.Write([]byte(string(data) + ""))
			if err != nil {
				ws.WriteMessage(websocket.TextMessage, []byte("Container not exits"))
			}
		}

	}

}
