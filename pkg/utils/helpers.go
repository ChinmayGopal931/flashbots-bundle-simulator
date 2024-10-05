package utils

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

// WeiToEther converts wei to ether
func WeiToEther(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetString(wei.String())
	return new(big.Float).Quo(f, big.NewFloat(1e18))
}

// FormatEther formats ether value to a string with 18 decimal places
func FormatEther(ether *big.Float) string {
	return fmt.Sprintf("%.18f", ether)
}

// HandleError is a generic error handler that logs the error and exits the program
func HandleError(err error, message string) {
	if err != nil {
		fmt.Printf("%s: %v\n", message, err)
		panic(err)
	}
}

func TestWeiToEther(t *testing.T) {
	testCases := []struct {
		name     string
		wei      *big.Int
		expected *big.Float
	}{
		{
			name:     "1 Ether",
			wei:      big.NewInt(1e18),
			expected: big.NewFloat(1),
		},
		{
			name:     "0.5 Ether",
			wei:      big.NewInt(5e17),
			expected: big.NewFloat(0.5),
		},
		{
			name:     "1 Gwei",
			wei:      big.NewInt(1e9),
			expected: big.NewFloat(1e-9),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := WeiToEther(tc.wei)
			assert.Equal(t, tc.expected.String(), result.String())
		})
	}
}

func TestFormatEther(t *testing.T) {
	testCases := []struct {
		name     string
		ether    *big.Float
		expected string
	}{
		{
			name:     "1 Ether",
			ether:    big.NewFloat(1),
			expected: "1.000000000000000000",
		},
		{
			name:     "0.5 Ether",
			ether:    big.NewFloat(0.5),
			expected: "0.500000000000000000",
		},
		{
			name:     "Small amount",
			ether:    big.NewFloat(1e-18),
			expected: "0.000000000000000001",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatEther(tc.ether)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestHandleError(t *testing.T) {
	t.Run("No Error", func(t *testing.T) {
		assert.NotPanics(t, func() {
			HandleError(nil, "This should not panic")
		})
	})

	t.Run("With Error", func(t *testing.T) {
		assert.Panics(t, func() {
			HandleError(errors.New("test error"), "This should panic")
		})
	})
}
