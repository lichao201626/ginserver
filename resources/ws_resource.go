package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

// WsResource ...
type WsResource struct {
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var cons = []*websocket.Conn{}

// NewWsResource ...
func NewWsResource(e *gin.Engine) {
	u := WsResource{}
	// Setup Routes
	e.GET("/ws", u.Ping)
}

func (r *WsResource) Ping(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	cons = append(cons, ws)
	if err != nil {
		return
	}
	defer ws.Close()
	for {
		//读取ws中的数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if string(message) == "ping" {
			message = []byte("pong")
		}
		//写入ws数据
		for _, con := range cons {
			err = con.WriteMessage(mt, message)
			if err != nil {
				break
			}
		}
	}
}
