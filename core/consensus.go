// core/consensus.go
package core

import (
	"math/rand"
	"time"
)

func SelectProposer(stakeMap map[string]uint64) string {
	// Suma total del stake
	var totalStake uint64
	for _, s := range stakeMap {
		totalStake += s
	}
	if totalStake == 0 {
		// Nadie hace stake, elegimos un validador predeterminado
		return "0xValidatorDefault"
	}

	// Escoger aleatoriamente en proporción al stake
	rand.Seed(time.Now().UnixNano())
	r := rand.Uint64() % totalStake
	var cumulative uint64
	for addr, s := range stakeMap {
		cumulative += s
		if r < cumulative {
			return addr
		}
	}
	return "0xValidatorDefault"
}
