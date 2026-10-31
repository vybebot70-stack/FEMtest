package state

import (
	"math/big"
	"testing"

	"github.com/0xPolygon/polygon-edge/chain"
	"github.com/0xPolygon/polygon-edge/state/runtime"
	"github.com/0xPolygon/polygon-edge/types"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"
)

func TestSplitGasFee(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		fee             int64
		wantTreasuryFee int64
		wantProposerFee int64
	}{
		{name: "zero", fee: 0, wantTreasuryFee: 0, wantProposerFee: 0},
		{name: "one wei rounds to proposer", fee: 1, wantTreasuryFee: 0, wantProposerFee: 1},
		{name: "even fee", fee: 100, wantTreasuryFee: 50, wantProposerFee: 50},
		{name: "odd fee remainder goes to proposer", fee: 101, wantTreasuryFee: 50, wantProposerFee: 51},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			proposer := types.StringToAddress("0x1234")
			transition := newTestTransition(map[types.Address]*PreState{
				proposer:           {Balance: 7},
				femTreasuryAddress: {Balance: 11},
			})
			transition.ctx.Coinbase = proposer

			transition.creditGasFee(big.NewInt(tt.fee))

			require.Equal(t, big.NewInt(11+tt.wantTreasuryFee), transition.state.GetBalance(femTreasuryAddress))
			require.Equal(t, big.NewInt(7+tt.wantProposerFee), transition.state.GetBalance(proposer))
			require.Equal(
				t,
				big.NewInt(tt.fee),
				new(big.Int).Add(big.NewInt(tt.wantTreasuryFee), big.NewInt(tt.wantProposerFee)),
			)
		})
	}
}

func TestApplySplitsGasFeesBetweenTreasuryAndProposer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		chainID int64
		to      *types.Address
		data    []byte
	}{
		{
			name:    "mainnet ordinary transfer",
			chainID: chain.FEMChainID,
			to:      addressPtr(types.StringToAddress("0x200")),
		},
		{
			name:    "testnet contract deployment",
			chainID: chain.FEMTestnetChainID,
			data:    []byte{0x00}, // STOP init code returns empty runtime code.
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sender := types.StringToAddress("0x100")
			proposer := types.StringToAddress("0x300")
			preState := map[types.Address]*PreState{
				sender: {Balance: 1_000_000_000_000_000_000},
			}
			snapshot := newStateWithPreState(preState)
			transition := NewTransition(chain.ForksInTime{Homestead: true, EIP158: true}, snapshot, newTestTxn(preState))
			transition.logger = hclog.NewNullLogger()
			transition.ctx = runtime.TxContext{
				Coinbase: proposer,
				BaseFee:  big.NewInt(0),
				ChainID:  tt.chainID,
				GasLimit: 200_000,
			}
			transition.gasPool = 200_000

			tx := &types.Transaction{
				From:     sender,
				To:       tt.to,
				Input:    tt.data,
				Gas:      100_000,
				GasPrice: new(big.Int).SetUint64(chain.FEMMinGasPrice),
				Value:    big.NewInt(0),
			}

			result, err := transition.Apply(tx)
			require.NoError(t, err)
			require.NotZero(t, result.GasUsed)

			fee := new(big.Int).Mul(
				new(big.Int).SetUint64(result.GasUsed),
				new(big.Int).SetUint64(chain.FEMMinGasPrice),
			)
			wantTreasury := new(big.Int).Div(new(big.Int).Set(fee), big.NewInt(2))
			wantProposer := new(big.Int).Sub(new(big.Int).Set(fee), wantTreasury)

			require.Equal(t, wantTreasury, transition.state.GetBalance(femTreasuryAddress))
			require.Equal(t, wantProposer, transition.state.GetBalance(proposer))
			require.Equal(t, fee, new(big.Int).Add(wantTreasury, wantProposer))
		})
	}
}

func TestCheckAndProcessTxEnforcesFEMGasPriceFloor(t *testing.T) {
	t.Parallel()

	sender := types.StringToAddress("0x100")

	tests := []struct {
		name    string
		chainID int64
		tx      *types.Transaction
	}{
		{
			name:    "mainnet legacy transaction below floor",
			chainID: chain.FEMChainID,
			tx: &types.Transaction{
				From:     sender,
				Gas:      21_000,
				GasPrice: new(big.Int).SetUint64(chain.FEMMinGasPrice - 1),
			},
		},
		{
			name:    "testnet dynamic fee transaction cannot bypass floor while London is off",
			chainID: chain.FEMTestnetChainID,
			tx: &types.Transaction{
				Type:      types.DynamicFeeTx,
				From:      sender,
				Gas:       21_000,
				GasFeeCap: new(big.Int).SetUint64(chain.FEMMinGasPrice),
				GasTipCap: new(big.Int).SetUint64(chain.FEMMinGasPrice),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			transition := newTestTransition(map[types.Address]*PreState{
				sender: {Balance: 1_000_000_000_000_000_000},
			})
			transition.ctx.ChainID = tt.chainID
			err := checkAndProcessTx(tt.tx, transition)
			require.ErrorContains(t, err, "gas price below FEM minimum")
		})
	}
}

func TestApplySplitsGasFeesForRevertedTransaction(t *testing.T) {
	t.Parallel()

	const initialBalance = uint64(1_000_000_000_000_000_000)
	sender := types.StringToAddress("0x100")
	proposer := types.StringToAddress("0x300")
	contract := types.StringToAddress("0x400")
	initialTreasuryBalance := big.NewInt(11)
	initialProposerBalance := big.NewInt(7)
	preState := map[types.Address]*PreState{
		sender: {
			Balance: initialBalance,
		},
		femTreasuryAddress: {
			Balance: uint64(initialTreasuryBalance.Int64()),
		},
		proposer: {
			Balance: uint64(initialProposerBalance.Int64()),
		},
	}

	transition := NewTransition(
		chain.ForksInTime{Homestead: true, EIP158: true, Byzantium: true, Istanbul: true},
		newStateWithPreState(preState),
		newTestTxn(preState),
	)
	transition.logger = hclog.NewNullLogger()
	transition.ctx = runtime.TxContext{
		Coinbase: proposer,
		BaseFee:  big.NewInt(0),
		ChainID:  chain.FEMTestnetChainID,
		GasLimit: 100_000,
	}
	transition.gasPool = 100_000
	transition.state.SetCode(contract, []byte{0x60, 0x00, 0x60, 0x00, 0xfd})

	tx := &types.Transaction{
		From:     sender,
		To:       addressPtr(contract),
		Gas:      50_000,
		GasPrice: new(big.Int).SetUint64(chain.FEMMinGasPrice),
		Value:    big.NewInt(0),
	}

	result, err := transition.Apply(tx)
	require.NoError(t, err)
	require.True(t, result.Reverted())

	fee := new(big.Int).Mul(
		new(big.Int).SetUint64(result.GasUsed),
		new(big.Int).SetUint64(chain.FEMMinGasPrice),
	)
	treasuryShare := new(big.Int).Div(new(big.Int).Set(fee), big.NewInt(2))
	proposerShare := new(big.Int).Sub(new(big.Int).Set(fee), treasuryShare)

	expectedTreasuryBalance := new(big.Int).Add(
		new(big.Int).Set(initialTreasuryBalance),
		treasuryShare,
	)
	expectedProposerBalance := new(big.Int).Add(
		new(big.Int).Set(initialProposerBalance),
		proposerShare,
	)

	require.Equal(t, expectedTreasuryBalance, transition.state.GetBalance(femTreasuryAddress))
	require.Equal(t, expectedProposerBalance, transition.state.GetBalance(proposer))
	require.Equal(t, fee, new(big.Int).Add(treasuryShare, proposerShare))
	require.Equal(
		t,
		new(big.Int).Sub(new(big.Int).SetUint64(initialBalance), fee),
		transition.state.GetBalance(sender),
	)
}

func TestApplySplitsRefundedGasFeesExactly(t *testing.T) {
	t.Parallel()

	const initialBalance = uint64(1_000_000_000_000_000_000)
	sender := types.StringToAddress("0x100")
	proposer := types.StringToAddress("0x300")
	contract := types.StringToAddress("0x400")
	storageKey := types.ZeroHash
	preState := map[types.Address]*PreState{
		sender: {
			Balance: initialBalance,
		},
		contract: {
			State: map[types.Hash]types.Hash{
				storageKey: types.StringToHash("1"),
			},
		},
	}

	transition := NewTransition(
		chain.ForksInTime{Homestead: true, EIP158: true, Byzantium: true, Istanbul: true},
		newStateWithPreState(preState),
		newTestTxn(preState),
	)
	transition.logger = hclog.NewNullLogger()
	transition.ctx = runtime.TxContext{
		Coinbase: proposer,
		BaseFee:  big.NewInt(0),
		ChainID:  chain.FEMTestnetChainID,
		GasLimit: 100_000,
	}
	transition.gasPool = 100_000
	transition.state.SetCode(contract, []byte{0x60, 0x00, 0x60, 0x00, 0x55, 0x00})

	tx := &types.Transaction{
		From:     sender,
		To:       addressPtr(contract),
		Gas:      50_000,
		GasPrice: new(big.Int).SetUint64(chain.FEMMinGasPrice),
		Value:    big.NewInt(0),
	}

	result, err := transition.Apply(tx)
	require.NoError(t, err)
	require.True(t, result.Succeeded())
	require.NotZero(t, transition.state.GetRefund(), "SSTORE clear should earn a gas refund")
	require.Zero(t, transition.state.GetState(contract, storageKey))

	fee := new(big.Int).Mul(
		new(big.Int).SetUint64(result.GasUsed),
		new(big.Int).SetUint64(chain.FEMMinGasPrice),
	)
	wantTreasury := new(big.Int).Div(new(big.Int).Set(fee), big.NewInt(2))
	wantProposer := new(big.Int).Sub(new(big.Int).Set(fee), wantTreasury)

	require.Equal(t, wantTreasury, transition.state.GetBalance(femTreasuryAddress))
	require.Equal(t, wantProposer, transition.state.GetBalance(proposer))
	require.Equal(t, fee, new(big.Int).Add(wantTreasury, wantProposer))
	require.Equal(
		t,
		new(big.Int).Sub(new(big.Int).SetUint64(initialBalance), fee),
		transition.state.GetBalance(sender),
	)
}

func addressPtr(address types.Address) *types.Address {
	return &address
}
