// core/state.go
package core

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/cbergoon/merkletree" // ejemplo de librería Merkle Tree (3rd party)
)

// State representaría la estructura de balances y su Merkle Trie
type State struct {
	Balances   map[string]uint64
	Nonces     map[string]uint64
	merkleTree *merkletree.MerkleTree
}

// Leaf implementa la interfaz merkletree.Content
type Leaf struct {
	Key   string
	Value uint64
}

func (l Leaf) CalculateHash() ([]byte, error) {
	h := sha256.Sum256([]byte(l.Key + string(rune(l.Value))))
	return h[:], nil
}
func (l Leaf) Equals(other merkletree.Content) (bool, error) {
	return l.Key == other.(Leaf).Key && l.Value == other.(Leaf).Value, nil
}

// NewState crea un state inicial
func NewState() *State {
	return &State{
		Balances:   make(map[string]uint64),
		Nonces:     make(map[string]uint64),
		merkleTree: nil,
	}
}

// Métodos para manipular Nonces
func (s *State) GetNonce(address string) uint64 {
	return s.Nonces[address]
}

func (s *State) IncrementNonce(address string) {
	s.Nonces[address] = s.Nonces[address] + 1
}

// UpdateMerkle actualiza la Merkle Trie del State
func (s *State) UpdateMerkle() error {
	var list []merkletree.Content
	for k, v := range s.Balances {
		list = append(list, Leaf{Key: k, Value: v})
	}
	tree, err := merkletree.NewTree(list)
	if err != nil {
		return err
	}
	s.merkleTree = tree
	return nil
}

// Root devuelve la raíz Merkle del estado
func (s *State) Root() (string, error) {
	if s.merkleTree == nil {
		return "", errors.New("Merkle Tree not initialized")
	}
	return hex.EncodeToString(s.merkleTree.MerkleRoot()), nil
}
