// core/genesis.go
package core

func CreateGenesisBlock() *Block {
	// Podríamos configurar un "alloc" de cuentas, balances iniciales, etc.
	genesisTxs := []*Transaction{
		// Transacciones especiales para asignar tokens a ciertas direcciones
	}

	// Normalmente su "parentHash" es 0x00... y su blockNumber es 0
	genesis := NewBlock(
		"0x0000000000000000",
		0,
		genesisTxs,
		"0xGENESIS_STATE_ROOT", // algo predefinido o calculado
	)

	return genesis
}
