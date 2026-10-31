package server

import (
	"testing"

	"github.com/0xPolygon/polygon-edge/chain"
	"github.com/0xPolygon/polygon-edge/command/server/config"
	"github.com/stretchr/testify/require"
)

func TestApplyFEMGasPriceFloor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		chainID        int64
		priceLimit     uint64
		wantPriceLimit uint64
		wantErr        bool
	}{
		{
			name:           "FEM default",
			chainID:        chain.FEMChainID,
			wantPriceLimit: chain.FEMMinGasPrice,
		},
		{
			name:           "FEM exact floor",
			chainID:        chain.FEMChainID,
			priceLimit:     chain.FEMMinGasPrice,
			wantPriceLimit: chain.FEMMinGasPrice,
		},
		{
			name:       "FEM below floor is rejected",
			chainID:    chain.FEMChainID,
			priceLimit: chain.FEMMinGasPrice - 1,
			wantErr:    true,
		},
		{
			name:           "other chains keep their configured value",
			chainID:        1,
			priceLimit:     0,
			wantPriceLimit: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			txPool := &config.TxPool{PriceLimit: tt.priceLimit}
			err := applyFEMGasPriceFloor(tt.chainID, txPool)

			if tt.wantErr {
				require.ErrorContains(t, err, "tx pool price limit must be at least 1000000000 wei")
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantPriceLimit, txPool.PriceLimit)
		})
	}
}

func TestApplyFEMGasPriceFloorForTestnet(t *testing.T) {
	t.Parallel()

	txPool := &config.TxPool{}
	require.NoError(t, applyFEMGasPriceFloor(chain.FEMTestnetChainID, txPool))
	require.Equal(t, chain.FEMMinGasPrice, txPool.PriceLimit)

	txPool.PriceLimit = chain.FEMMinGasPrice - 1
	require.Error(t, applyFEMGasPriceFloor(chain.FEMTestnetChainID, txPool))
}

func TestApplyBootnodesOverride(t *testing.T) {
	t.Parallel()

	genesis := &chain.Chain{
		Bootnodes: []string{"/ip4/127.0.0.1/tcp/30301/p2p/peer-a"},
	}
	override := []string{
		"/dns4/validator-1.internal/tcp/30303/p2p/peer-b",
		"/dns4/validator-2.internal/tcp/30303/p2p/peer-c",
	}

	require.NoError(t, applyBootnodesOverride(genesis, override))
	require.Equal(t, override, genesis.Bootnodes)
}

func TestApplyBootnodesOverrideRejectsInvalidInput(t *testing.T) {
	t.Parallel()

	require.Error(t, applyBootnodesOverride(&chain.Chain{}, []string{""}))
	require.Error(t, applyBootnodesOverride(nil, []string{"bootnode"}))
}
