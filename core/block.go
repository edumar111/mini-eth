// core/block.go
package core

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Block representa un bloque muy básico.
type Block struct {
	Header       *BlockHeader
	Transactions []*Transaction
	// El Merkle Root podría estar en el header o aquí según se prefiera
}

// BlockHeader información esencial de cabecera.
type BlockHeader struct {
	ParentHash  string
	Timestamp   int64
	StateRoot   string // Raíz de la Merkle Trie de estado
	BlockNumber uint64
	// Podríamos tener más campos como ExtraData, Difficulty, etc.
}

// NewBlock crea un nuevo bloque y calculamos su StateRoot simplificado.
func NewBlock(parentHash string, blockNumber uint64, txs []*Transaction, stateRoot string) *Block {
	header := &BlockHeader{
		ParentHash:  parentHash,
		Timestamp:   time.Now().Unix(),
		BlockNumber: blockNumber,
		StateRoot:   stateRoot,
	}
	return &Block{
		Header:       header,
		Transactions: txs,
	}
}

// Hash del bloque (ejemplo simplificado usando sha256)
func (b *Block) Hash() string {
	data := []byte(
		b.Header.ParentHash +
			string(b.Header.Timestamp) +
			string(b.Header.BlockNumber) +
			b.Header.StateRoot,
	)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
