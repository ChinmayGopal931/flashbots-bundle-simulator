package simulation

import (
	"context"
	"math/big"
	"testing"

	"github.com/ChinmayGopal931/flashbots-bundle-simulator/internal/bundle"
	"github.com/ChinmayGopal931/flashbots-bundle-simulator/internal/ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEthereumClient is a mock of the ethereum.Client interface
type MockEthereumClient struct {
	mock.Mock
}

func (m *MockEthereumClient) GetBlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	args := m.Called(ctx, number)
	return args.Get(0).(*types.Block), args.Error(1)
}

func (m *MockEthereumClient) GetBalance(ctx context.Context, address common.Address, blockNumber *big.Int) (*big.Int, error) {
	args := m.Called(ctx, address, blockNumber)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *MockEthereumClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockEthereumClient) Close() {
	m.Called()
}

func TestSimulateBundle(t *testing.T) {
	mockClient := new(MockEthereumClient)
	simulator := NewSimulator(mockClient)

	validTxHex := "0xf86d8202b28477359400825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d880de0b6b3a7640000802ba0f3514e458983a5b4644f99d88ae2e24489fe2b1aade6d86724bdb3fae5bbd5c9a03a309017057d36a3e9ea92626cd816dbf63bec6f29a43171a05b1a92eb790bea"

	testCases := []struct {
		name          string
		bundle        *bundle.FlashbotsBundle
		expectedError bool
	}{
		{
			name: "Valid Bundle",
			bundle: &bundle.FlashbotsBundle{
				Txs:         []string{validTxHex},
				BlockNumber: big.NewInt(12345678),
			},
			expectedError: false,
		},
		{
			name: "Invalid Transaction",
			bundle: &bundle.FlashbotsBundle{
				Txs:         []string{"invalid_tx_hex"},
				BlockNumber: big.NewInt(12345678),
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockBlock := types.NewBlockWithHeader(&types.Header{Number: tc.bundle.BlockNumber})
			mockClient.On("GetBlockByNumber", mock.Anything, tc.bundle.BlockNumber).Return(mockBlock, nil)

			result, err := simulator.SimulateBundle(context.Background(), tc.bundle)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.True(t, result.Success)
				assert.Equal(t, uint64(21000), result.GasUsed) // This is the placeholder value we set in estimateGas
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestNewSimulator(t *testing.T) {
	mockClient := new(MockEthereumClient)
	simulator := NewSimulator(mockClient)
	assert.NotNil(t, simulator)
	assert.Equal(t, ethereum.Client(mockClient), simulator.client)
}

func TestEstimateGas(t *testing.T) {
	mockClient := new(MockEthereumClient)
	simulator := NewSimulator(mockClient)

	tx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 21000, big.NewInt(1), nil)

	gasUsed, err := simulator.estimateGas(context.Background(), tx)
	assert.NoError(t, err)
	assert.Equal(t, uint64(21000), gasUsed)
}
