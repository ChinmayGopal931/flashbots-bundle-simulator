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

	bundleFilePath := os.Args[1]
	log.Printf("Reading bundle file: %s", bundleFilePath)

	bundleJSON, err := ioutil.ReadFile(bundleFilePath)
	if err != nil {
		log.Fatalf("Failed to read bundle file: %v", err)
	}
	log.Printf("Bundle JSON content: %s", string(bundleJSON))

	flashbotsBundle, err := bundle.ParseAndValidateBundle(bundleJSON)
	if err != nil {
		log.Fatalf("Failed to parse and validate bundle: %v", err)
	}
	log.Printf("Parsed bundle: %+v", flashbotsBundle)

	// Connect to Sepolia testnet
	client, err := ethereum.NewClient("https://sepolia.infura.io/v3/c978b74938064a98b67a150e4ade294d")
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
