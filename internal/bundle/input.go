package bundle

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// HexBigInt is a custom type for handling hexadecimal strings that represent big integers
type HexBigInt big.Int

// UnmarshalJSON implements the json.Unmarshaler interface for HexBigInt
func (hbi *HexBigInt) UnmarshalJSON(data []byte) error {
	var hexString string
	if err := json.Unmarshal(data, &hexString); err != nil {
		return err
	}

	hexString = strings.TrimPrefix(hexString, "0x")
	bi, success := new(big.Int).SetString(hexString, 16)
	if !success {
		return errors.New("failed to parse hexadecimal string")
	}

	*hbi = HexBigInt(*bi)
	return nil
}

// BigInt returns the underlying big.Int value
func (hbi *HexBigInt) BigInt() *big.Int {
	return (*big.Int)(hbi)
}

// rawFlashbotsBundle represents the JSON structure of the input
type rawFlashbotsBundle struct {
	Txs               []string `json:"txs"`
	BlockNumber       string   `json:"blockNumber"`
	MinTimestamp      uint64   `json:"minTimestamp"`
	MaxTimestamp      uint64   `json:"maxTimestamp"`
	RevertingTxHashes []string `json:"revertingTxHashes,omitempty"`
}

// FlashbotsBundle represents a Flashbots bundle with parsed values
type FlashbotsBundle struct {
	Txs               []string
	BlockNumber       *big.Int
	MinTimestamp      uint64
	MaxTimestamp      uint64
	RevertingTxHashes []common.Hash
}

// ParseAndValidateBundle parses and validates a Flashbots bundle from JSON input
func ParseAndValidateBundle(input []byte) (*FlashbotsBundle, error) {
	var raw rawFlashbotsBundle
	if err := json.Unmarshal(input, &raw); err != nil {
		return nil, errors.New("failed to parse bundle JSON: " + err.Error())
	}

	blockNumber, ok := new(big.Int).SetString(raw.BlockNumber, 10)
	if !ok {
		return nil, errors.New("invalid block number")
	}

	bundle := &FlashbotsBundle{
		Txs:          raw.Txs,
		BlockNumber:  blockNumber,
		MinTimestamp: raw.MinTimestamp,
		MaxTimestamp: raw.MaxTimestamp,
	}

	for _, hash := range raw.RevertingTxHashes {
		bundle.RevertingTxHashes = append(bundle.RevertingTxHashes, common.HexToHash(hash))
	}

	if err := validateBundle(bundle); err != nil {
		return nil, err
	}

	return bundle, nil
}

func validateBundle(bundle *FlashbotsBundle) error {
	if len(bundle.Txs) == 0 {
		return errors.New("bundle must contain at least one transaction")
	}

	if bundle.BlockNumber == nil || bundle.BlockNumber.Sign() <= 0 {
		return errors.New("block number must be a positive integer")
	}

	if bundle.MinTimestamp > bundle.MaxTimestamp && bundle.MaxTimestamp != 0 {
		return errors.New("minTimestamp cannot be greater than maxTimestamp")
	}

	return nil
}

// // HexToTx converts a hex string to a transaction
// func HexToTx(hexStr string) (*types.Transaction, error) {
// 	return hexToTx(hexStr)
// }

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

func HexToTx(hexStr string) (*types.Transaction, error) {
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
