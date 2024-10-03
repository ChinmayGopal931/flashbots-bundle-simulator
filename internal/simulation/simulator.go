package simulation

import (
	"context"
	"errors"
	"math/big"

	"github.com/ChinmayGopal931/flashbots-bundle-simulator/internal/bundle"
	"github.com/ChinmayGopal931/flashbots-bundle-simulator/internal/ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

type Simulator struct {
	client *ethereum.Client
}

func NewSimulator(client *ethereum.Client) *Simulator {
	return &Simulator{client: client}
}

func (s *Simulator) SimulateBundle(ctx context.Context, bundle *bundle.FlashbotsBundle) (*SimulationResult, error) {
	block, err := s.client.GetBlockByNumber(ctx, bundle.BlockNumber)
	if err != nil {
		return nil, err
	}

	statedb, err := s.createStateDB(ctx, block)
	if err != nil {
		return nil, err
	}

	result := &SimulationResult{
		Success: true,
		Profit:  big.NewInt(0),
		GasUsed: 0,
	}

	for _, txHex := range bundle.Txs {
		tx, err := bundle.HexToTx(txHex)
		if err != nil {
			return nil, err
		}

		receipt, err := s.applyTransaction(statedb, tx)
		if err != nil {
			result.Success = false
			result.Error = err.Error()
			break
		}

		result.GasUsed += receipt.GasUsed
		// Calculate profit (this is a simplified version)
		profit := new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(receipt.GasUsed))
		result.Profit.Add(result.Profit, profit)
	}

	return result, nil
}

func (s *Simulator) createStateDB(ctx context.Context, block *types.Block) (*StateDB, error) {
	// TODO
	return nil, errors.New("createStateDB TODO")
}

func (s *Simulator) applyTransaction(statedb *StateDB, tx *types.Transaction) (*types.Receipt, error) {
	// TODO
	return nil, errors.New("applyTransaction TODO")
}

type SimulationResult struct {
	Success bool
	Profit  *big.Int
	GasUsed uint64
	Error   string
}
