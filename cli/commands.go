package cli

import (
	"github.com/edumar111/my-geth-edu/core"
	"github.com/edumar111/my-geth-edu/p2p"
	"github.com/edumar111/my-geth-edu/rpc"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

func InitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Inicializa el bloque génesis",
		Run: func(cmd *cobra.Command, args []string) {
			// Lógica para crear y almacenar el bloque génesis en disco
			log.Println("Genesis block created!")
		},
	}
}

func RunCmd() *cobra.Command {
	var p2pPort int
	var rpcHTTPPort int
	var rpcWSPort int

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Inicia el nodo",
		Run: func(cmd *cobra.Command, args []string) {
			// 1. Cargar génesis (o cargar de disco)
			genesis := core.CreateGenesisBlock()
			log.Printf("Genesis block hash: %s\n", genesis.Hash())

			// 2. Creamos un State (o cargamos balances, nonces, etc.)
			state := core.NewState()
			// (Opcional) Asignar balances o nonces iniciales

			// 3. Iniciar P2P
			server, err := p2p.NewP2PServer(p2pPort)
			if err != nil {
				log.Fatal("Error al iniciar P2P:", err)
			}
			defer server.Shutdown()

			// 4. Creamos el RPCServer con referencia a nuestro State
			rpcServer := &rpc.RPCServer{
				State: state,
			}
			rpcServer.StartRPC(strconv.Itoa(rpcHTTPPort))
			//go rpcServer.StartRPC(strconv.Itoa(rpcHTTPPort))

			// 5. Servidor WebSocket (en "/")
			wsServer := &rpc.RPCWSServer{
				State: state,
			}
			wsServer.StartWS(strconv.Itoa(rpcWSPort))
			//go rpcServer.StartWS(strconv.Itoa(rpcWSPort))

			log.Printf("Node running on P2P port %d, RPC HTTP %d, WS %d", p2pPort, rpcHTTPPort, rpcWSPort)
			select {}
		},
	}

	// Definimos los flags
	cmd.Flags().IntVar(&p2pPort, "p2p-port", 30303, "Puerto para P2P")
	cmd.Flags().IntVar(&rpcHTTPPort, "rpc-http-port", 4045, "Puerto para RPC HTTP")
	cmd.Flags().IntVar(&rpcWSPort, "rpc-ws-port", 4046, "Puerto para RPC WebSocket")

	return cmd
}
