package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ChinmayGopal931/flashbots-bundle-simulator/internal/bundle"
	"github.com/ChinmayGopal931/flashbots-bundle-simulator/internal/ethereum"
	"github.com/ChinmayGopal931/flashbots-bundle-simulator/internal/simulation"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a path to a bundle JSON file")
	}

	bundleJSON, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to read bundle file: %v", err)
	}

	flashbotsBundle, err := bundle.ParseAndValidateBundle(bundleJSON)
	if err != nil {
		log.Fatalf("Failed to parse and validate bundle: %v", err)
	}

	// Connect to Goerli testnet
	client, err := ethereum.NewClient("https://goerli.infura.io/v3/YOUR-PROJECT-ID")
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	defer client.Close()

	simulator := simulation.NewSimulator(client)

	result, err := simulator.SimulateBundle(context.Background(), flashbotsBundle)
	if err != nil {
		log.Fatalf("Failed to simulate bundle: %v", err)
	}

	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal result to JSON: %v", err)
	}

	fmt.Println(string(jsonResult))
}
