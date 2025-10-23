package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/oasisprotocol/oasis-sdk/client-sdk/go/client"
	"github.com/oasisprotocol/oasis-sdk/client-sdk/go/config"
	"github.com/oasisprotocol/oasis-sdk/client-sdk/go/connection"
	"github.com/oasisprotocol/oasis-sdk/client-sdk/go/modules/rofl"

	consensus "github.com/oasisprotocol/oasis-core/go/consensus/api"
)

func main() {
	// Parse command-line arguments.
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <testnet|mainnet> [app_id]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  - Network: testnet or mainnet\n")
		fmt.Fprintf(os.Stderr, "  - App ID (optional): specific ROFL app ID to query\n")
		os.Exit(1)
	}

	networkName := os.Args[1]
	var specificAppID string
	if len(os.Args) >= 3 {
		specificAppID = os.Args[2]
	}

	// Setup client.
	if networkName != "testnet" && networkName != "mainnet" {
		log.Fatalf("Invalid network: %s (must be 'testnet' or 'mainnet')", networkName)
	}
	ctx := context.Background()
	network := config.DefaultNetworks.All[networkName]
	if network == nil {
		log.Fatalf("Network '%s' not found in configuration", networkName)
	}
	fmt.Printf("Connecting to %s...\n", networkName)
	conn, err := connection.Connect(ctx, network)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	consensusBackend := conn.Consensus().Core()

	// Fetch and display latest block.
	block, err := consensusBackend.GetBlock(ctx, consensus.HeightLatest)
	if err != nil {
		log.Fatalf("Failed to fetch latest block: %v", err)
	}
	fmt.Printf("\nLatest Block Information:\n")
	fmt.Printf("  Height: %d\n", block.Height)
	fmt.Printf("  Hash: %s\n", block.Hash)
	fmt.Printf("  Time: %s\n", block.Time)

	fmt.Printf("\n--- Querying Sapphire ROFL Apps on %s ---\n", networkName)

	// Setup ROFL client.
	sapphirePT := network.ParaTimes.All["sapphire"]
	if sapphirePT == nil {
		log.Fatalf("Sapphire paratime not found in network configuration")
	}
	runtimeClient := conn.Runtime(sapphirePT)
	roflClient := rofl.NewV1(runtimeClient)

	// Query ROFL apps.
	var apps []*rofl.AppConfig
	if specificAppID != "" {
		appID := rofl.NewAppIDFromBech32(specificAppID)
		app, err := roflClient.App(ctx, client.RoundLatest, appID)
		if err != nil {
			log.Fatalf("Failed to fetch ROFL app %s: %v", specificAppID, err)
		}
		apps = []*rofl.AppConfig{app}
		fmt.Printf("\nQuerying specific app: %s\n\n", specificAppID)
	} else {
		allApps, err := roflClient.Apps(ctx, client.RoundLatest)
		if err != nil {
			log.Fatalf("Failed to fetch ROFL apps: %v", err)
		}
		// Filter for active apps only.
		for _, app := range allApps {
			instances, err := roflClient.AppInstances(ctx, client.RoundLatest, app.ID)
			if err == nil && len(instances) > 0 {
				apps = append(apps, app)
			}
		}
		fmt.Printf("\nTotal active ROFL apps found: %d\n\n", len(apps))
	}

	// Print app details.
	for i, app := range apps {
		if i > 0 {
			time.Sleep(500 * time.Millisecond) // Rate limiting.
		}
		fmt.Printf("ROFL App #%d:\n", i+1)
		fmt.Printf("  ID: %s\n", app.ID)
		if app.Admin != nil {
			fmt.Printf("  Admin: %s\n", *app.Admin)
		} else {
			fmt.Printf("  Admin: <none>\n")
		}
		fmt.Printf("  Stake: %s\n", app.Stake)
		fmt.Printf("  SEK: %s\n", app.SEK)

		// Display policy.
		fmt.Printf("  Policy:\n")
		fmt.Printf("    Fee Policy: %d\n", app.Policy.Fees)
		fmt.Printf("    Max Expiration: %d\n", app.Policy.MaxExpiration)

		// Display enclaves.
		fmt.Printf("    Enclaves: %d\n", len(app.Policy.Enclaves))
		for j, enclave := range app.Policy.Enclaves {
			fmt.Printf("      Enclave #%d:\n", j+1)
			fmt.Printf("        MrEnclave: %s\n", enclave.MrEnclave)
			fmt.Printf("        MrSigner: %s\n", enclave.MrSigner)
		}

		// Display endorsements.
		fmt.Printf("    Endorsements: %d\n", len(app.Policy.Endorsements))
		for j, endorsement := range app.Policy.Endorsements {
			fmt.Printf("      Endorsement #%d:\n", j+1)
			if endorsement.Any != nil {
				fmt.Printf("        Type: Any\n")
			}
			if endorsement.ComputeRole != nil {
				fmt.Printf("        Type: ComputeRole\n")
			}
			if endorsement.ObserverRole != nil {
				fmt.Printf("        Type: ObserverRole\n")
			}
			if endorsement.Entity != nil {
				fmt.Printf("        Type: Entity\n")
				fmt.Printf("        Entity: %s\n", *endorsement.Entity)
			}
			if endorsement.Node != nil {
				fmt.Printf("        Type: Node\n")
				fmt.Printf("        Node: %s\n", *endorsement.Node)
			}
			if endorsement.Provider != nil {
				fmt.Printf("        Type: Provider\n")
				fmt.Printf("        Provider: %s\n", *endorsement.Provider)
			}
			if endorsement.ProviderInstanceAdmin != nil {
				fmt.Printf("        Type: ProviderInstanceAdmin\n")
				fmt.Printf("        Admin: %s\n", *endorsement.ProviderInstanceAdmin)
			}
			if len(endorsement.And) > 0 {
				fmt.Printf("        Type: And (%d policies)\n", len(endorsement.And))
			}
			if len(endorsement.Or) > 0 {
				fmt.Printf("        Type: Or (%d policies)\n", len(endorsement.Or))
			}
		}

		// Display app metadata.
		if len(app.Metadata) > 0 {
			fmt.Printf("  App Metadata:\n")
			for key, value := range app.Metadata {
				fmt.Printf("    %s: %s\n", key, value)
			}
		}

		// Display secrets (values are encrypted).
		if len(app.Secrets) > 0 {
			fmt.Printf("  Secrets: %d\n", len(app.Secrets))
			for key := range app.Secrets {
				fmt.Printf("    %s: <encrypted>\n", key)
			}
		}

		// Query instances for this app.
		instances, err := roflClient.AppInstances(ctx, client.RoundLatest, app.ID)
		if err != nil {
			log.Printf("  Warning: Failed to fetch instances: %v\n", err)
		} else {
			fmt.Printf("  Instances: %d\n", len(instances))
			for j, inst := range instances {
				fmt.Printf("    Instance #%d:\n", j+1)
				fmt.Printf("      App: %s\n", inst.App)
				fmt.Printf("      Node ID: %s\n", inst.NodeID)
				if inst.EntityID != nil {
					fmt.Printf("      Entity ID: %s\n", *inst.EntityID)
				} else {
					fmt.Printf("      Entity ID: <none>\n")
				}
				fmt.Printf("      RAK: %s\n", inst.RAK)
				fmt.Printf("      REK: %s\n", inst.REK)
				fmt.Printf("      Expiration: %d\n", inst.Expiration)
				fmt.Printf("      Extra Keys: %d\n", len(inst.ExtraKeys))
				for k, key := range inst.ExtraKeys {
					fmt.Printf("        Key #%d: %s\n", k+1, key)
				}
				if len(inst.Metadata) > 0 {
					fmt.Printf("      Instance Metadata:\n")
					for key, value := range inst.Metadata {
						fmt.Printf("        %s: %s\n", key, value)
					}
				} else {
					fmt.Printf("      Instance Metadata: <none>\n")
				}
			}
		}
		fmt.Println()
	}

	if len(apps) == 0 {
		fmt.Println("No active ROFL apps found.")
	}
}
