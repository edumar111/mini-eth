// p2p/server.go
package p2p

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"log"
	"time"
)

type P2PServer struct {
	Host host.Host
	Ctx  context.Context
}

func NewP2PServer(listenPort int) (*P2PServer, error) {
	ctx := context.Background()

	cm, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater
		connmgr.WithGracePeriod(2*time.Minute),
	)
	if err != nil {
		return nil, err
	}

	h, err := libp2p.New(
		libp2p.ConnectionManager(cm),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
	)
	if err != nil {
		return nil, err
	}

	server := &P2PServer{
		Host: h,
		Ctx:  ctx,
	}

	// Podríamos registrar "handlers" para mensajes, etc.
	// h.SetStreamHandler("/mini-eth/1.0.0", server.handleStream)

	log.Printf("P2P node started. Listening on: %v\n", h.Addrs())
	return server, nil
}

// Cerrar conexiones
func (s *P2PServer) Shutdown() {
	s.Host.Close()
}
