package chain

import (
	"math/big"
)

// Treasury Optional treasury for launching PulseChain testnets
type Treasury struct {
	Addr    string `json:"addr"`
	Balance string `json:"balance"`
}

// PulseChainTTDOffset A trivially small amount of work to add to the Ethereum Mainnet TTD
// to allow for un-merging and merging with the PulseChain beacon chain
var PulseChainTTDOffset = big.NewInt(131_072)
