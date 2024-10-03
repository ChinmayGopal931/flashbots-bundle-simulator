package bundle

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// FlashbotsBundle represents a Flashbots bundle
type FlashbotsBundle struct {
	Txs               []string      `json:"txs"`
	BlockNumber       *big.Int      `json:"blockNumber"`
	MinTimestamp      uint64        `json:"minTimestamp,omitempty"`
	MaxTimestamp      uint64        `json:"maxTimestamp,omitempty"`
	RevertingTxHashes []common.Hash `json:"revertingTxHashes,omitempty"`
}

// ParseAndValidateBundle parses and validates a Flashbots bundle from JSON input
func ParseAndValidateBundle(input []byte) (*FlashbotsBundle, error) {
	var bundle FlashbotsBundle
	if err := json.Unmarshal(input, &bundle); err != nil {
		return nil, errors.New("failed to parse bundle JSON: " + err.Error())
	}

	if err := validateBundle(&bundle); err != nil {
		return nil, err
	}

	return &bundle, nil
}

func validateBundle(bundle *FlashbotsBundle) error {
	if len(bundle.Txs) == 0 {
		return errors.New("bundle must contain at least one transaction")
	}

	if bundle.BlockNumber == nil {
		return errors.New("block number must be specified")
	}

	if bundle.MinTimestamp > bundle.MaxTimestamp && bundle.MaxTimestamp != 0 {
		return errors.New("minTimestamp cannot be greater than maxTimestamp")
	}

	for _, txHex := range bundle.Txs {
		if _, err := hexToTx(txHex); err != nil {
			return errors.New("invalid transaction in bundle: " + err.Error())
		}
	}

	return nil
}

func hexToTx(hexStr string) (*types.Transaction, error) {
	// Remove '0x' prefix if present
	if len(hexStr) > 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}

	txBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}

	tx := new(types.Transaction)
	err = rlp.DecodeBytes(txBytes, tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
