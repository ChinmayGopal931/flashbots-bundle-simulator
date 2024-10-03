package bundle

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAndValidateBundle(t *testing.T) {
	validBundle := FlashbotsBundle{
		Txs:         []string{"0xf86d8202b28477359400825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d880de0b6b3a7640000802ba0f3514e458983a5b4644f99d88ae2e24489fe2b1aade6d86724bdb3fae5bbd5c9a03a309017057d36a3e9ea92626cd816dbf63bec6f29a43171a05b1a92eb790bea"},
		BlockNumber: big.NewInt(12345678),
	}

	invalidBundle := FlashbotsBundle{
		Txs:         []string{"invalid_tx_hex"},
		BlockNumber: big.NewInt(12345678),
	}

	t.Run("Valid Bundle", func(t *testing.T) {
		input, _ := json.Marshal(validBundle)
		parsedBundle, err := ParseAndValidateBundle(input)
		assert.NoError(t, err)
		assert.NotNil(t, parsedBundle)
		assert.Equal(t, validBundle.BlockNumber, parsedBundle.BlockNumber)
	})

	t.Run("Invalid Bundle", func(t *testing.T) {
		input, _ := json.Marshal(invalidBundle)
		_, err := ParseAndValidateBundle(input)
		assert.Error(t, err)
	})
}

func TestValidateBundle(t *testing.T) {
	t.Run("Empty Bundle", func(t *testing.T) {
		bundle := &FlashbotsBundle{}
		err := validateBundle(bundle)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must contain at least one transaction")
	})

	t.Run("Missing Block Number", func(t *testing.T) {
		bundle := &FlashbotsBundle{
			Txs: []string{"0xf86d8202b28477359400825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d880de0b6b3a7640000802ba0f3514e458983a5b4644f99d88ae2e24489fe2b1aade6d86724bdb3fae5bbd5c9a03a309017057d36a3e9ea92626cd816dbf63bec6f29a43171a05b1a92eb790bea"},
		}
		err := validateBundle(bundle)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "block number must be specified")
	})
}
