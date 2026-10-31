package state

import (
	"math/big"

	"github.com/0xPolygon/polygon-edge/types"
)

var femTreasuryAddress = types.StringToAddress("0x38B9bbECA3bB1437edeBA2c6298024e1c1C8E997")

// splitGasFee rounds the Treasury share down; the proposer receives the remainder.
func splitGasFee(fee *big.Int) (treasuryFee, proposerFee *big.Int) {
	treasuryFee = new(big.Int).Div(new(big.Int).Set(fee), big.NewInt(2))
	proposerFee = new(big.Int).Sub(new(big.Int).Set(fee), treasuryFee)

	return treasuryFee, proposerFee
}

func (t *Transition) creditGasFee(fee *big.Int) {
	treasuryFee, proposerFee := splitGasFee(fee)

	if treasuryFee.Sign() > 0 {
		t.state.AddBalance(femTreasuryAddress, treasuryFee)
	}

	t.state.AddBalance(t.ctx.Coinbase, proposerFee)
}
