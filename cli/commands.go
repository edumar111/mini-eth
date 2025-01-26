package cli

import (
	"github.com/edumar111/my-geth-edu/core"
	"github.com/edumar111/my-geth-edu/p2p"
	"github.com/edumar111/my-geth-edu/rpc"
	"github.com/spf13/cobra"
	"log"
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
	var rpcPort string
	var wsPort string
	var p2pPort int

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Inicia el nodo",
		Run: func(cmd *cobra.Command, args []string) {
			// 1. Cargar génesis y cadena
			genesis := core.CreateGenesisBlock()
			log.Printf("Genesis block hash: %s\n", genesis.Hash())

			// 2. Iniciar P2P
			server, err := p2p.NewP2PServer(p2pPort)
			if err != nil {
				log.Fatal("Error al iniciar P2P:", err)
			}
			defer server.Shutdown()

			// 3. Iniciar RPC
			rpcServer := &rpc.RPCServer{}
			go rpcServer.StartRPC(rpcPort)

			// 4. Iniciar WS
			go rpcServer.StartWS(wsPort)

			// Bloquear para que la aplicación no termine
			select {}
		},
	}

	cmd.Flags().StringVar(&rpcPort, "rpc", "4045", "Puerto para RPC HTTP")
	cmd.Flags().StringVar(&wsPort, "ws", "4046", "Puerto para WS")
	cmd.Flags().IntVar(&p2pPort, "p2p", 30303, "Puerto para P2P")

	return cmd
}
