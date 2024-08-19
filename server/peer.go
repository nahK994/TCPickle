package server

import (
	"log/slog"
	"net"

	"github.com/nahK994/TCPickle/handlers"
	"github.com/nahK994/TCPickle/models"
	"github.com/nahK994/TCPickle/utils"
)

type Peer struct {
	conn net.Conn
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		conn: conn,
	}
}

func (p *Peer) readConn() {
	buf := make([]byte, 1024)
	n, err := p.conn.Read(buf)
	if err != nil {
		slog.Error("peer read error", "err", err, "remoteAddr", p.conn.RemoteAddr())
		p.conn.Close()
		return
	}

	request := handlers.HandleRequest(buf[:n])
	// fmt.Println("Request from", p.conn.RemoteAddr(), " ==>", request)
	response := new(models.Response)
	requestHandler := utils.RouteMapper[models.HttpUrlPath(request.UrlPath)]
	if requestHandler.Method != request.Method {
		response.StatusCode = 405
		response.Body = ""
	} else {
		requestHandler.Func(*request, response)
	}
	handlers.HandleResponse(response, p.conn)
	p.conn.Close()
}
