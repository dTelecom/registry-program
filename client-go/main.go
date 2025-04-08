package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/joho/godotenv"

	"solana-registry-client/registry"
)

const LAMPORTS_PER_SOL = 1000000000

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	rpcURL := os.Getenv("SOLANA_RPC_URL")
	if rpcURL == "" {
		log.Fatal("SOLANA_RPC_URL is required")
	}

	wsURL := os.Getenv("SOLANA_WS_URL")
	if wsURL == "" {
		log.Fatal("SOLANA_WS_URL is required")
	}

	programID := os.Getenv("PROGRAM_ID")
	if programID == "" {
		log.Fatal("PROGRAM_ID is required")
	}

	privateKey := os.Getenv("WALLET_PRIVATE_KEY")
	if privateKey == "" {
		log.Fatal("WALLET_PRIVATE_KEY is required")
	}

	client, err := registry.NewRegistryClient(rpcURL, wsURL, programID, privateKey)
	if err != nil {
		log.Fatalf("Failed to create registry client: %v", err)
	}
	defer client.Close()

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	ctx := context.Background()

	switch os.Args[1] {
	case "create":
		if len(os.Args) != 3 {
			log.Fatal("Usage: create <registry_name>")
		}
		name := os.Args[2]
		sig, err := client.CreateRegistry(ctx, name)
		if err != nil {
			log.Fatalf("Failed to create registry: %v", err)
		}
		fmt.Printf("Registry created. Transaction signature: %s\n", sig)

	case "add-client":
		if len(os.Args) != 6 {
			log.Fatal("Usage: add-client <registry_name> <account_to_add> <valid_days> <limit>")
		}
		registryName := os.Args[2]
		account, err := solana.PublicKeyFromBase58(os.Args[3])
		if err != nil {
			log.Fatalf("Invalid account address: %v", err)
		}
		validDays := 0
		if _, err := fmt.Sscanf(os.Args[4], "%d", &validDays); err != nil {
			log.Fatalf("Invalid valid days: %v", err)
		}
		limit := uint32(0)
		if _, err := fmt.Sscanf(os.Args[5], "%d", &limit); err != nil {
			log.Fatalf("Invalid limit: %v", err)
		}
		validUntil := time.Now().AddDate(0, 0, validDays)

		sig, err := client.AddClientToRegistry(ctx, registryName, account, validUntil, limit)
		if err != nil {
			log.Fatalf("Failed to add client to registry: %v", err)
		}
		fmt.Printf("Client account added to registry. Transaction signature: %s\n", sig)

	case "add-node":
		if len(os.Args) != 5 {
			log.Fatal("Usage: add-node <registry_name> <account_to_add> <domain>")
		}
		registryName := os.Args[2]
		account, err := solana.PublicKeyFromBase58(os.Args[3])
		if err != nil {
			log.Fatalf("Invalid account address: %v", err)
		}
		domain := os.Args[4]
		if len(domain) > 64 {
			log.Fatalf("Domain name must be 64 characters or less")
		}

		sig, err := client.AddNodeToRegistry(ctx, registryName, account, domain)
		if err != nil {
			log.Fatalf("Failed to add node to registry: %v", err)
		}
		fmt.Printf("Node account added to registry. Transaction signature: %s\n", sig)

	case "get-client":
		if len(os.Args) != 4 {
			log.Fatal("Usage: get-client <registry_name> <account_to_check>")
		}
		registryName := os.Args[2]
		account, err := solana.PublicKeyFromBase58(os.Args[3])
		if err != nil {
			log.Fatalf("Invalid account address: %v", err)
		}
		entry, err := client.GetClientFromRegistry(ctx, registryName, account)
		if err != nil {
			log.Fatalf("Failed to get client from registry: %v", err)
		}
		if entry == nil {
			fmt.Println("Client account not found in registry")
		} else {
			fmt.Printf("Client registry entry:\n")
			fmt.Printf("  Parent: %s\n", entry.Parent)
			fmt.Printf("  Registered: %s\n", entry.Registred)
			fmt.Printf("  Valid until: %s\n", time.Unix(entry.Until, 0))
			fmt.Printf("  Limit: %d\n", entry.Limit)
		}

	case "get-node":
		if len(os.Args) != 4 {
			log.Fatal("Usage: get-node <registry_name> <account_to_check>")
		}
		registryName := os.Args[2]
		account, err := solana.PublicKeyFromBase58(os.Args[3])
		if err != nil {
			log.Fatalf("Invalid account address: %v", err)
		}
		entry, err := client.GetNodeFromRegistry(ctx, registryName, account)
		if err != nil {
			log.Fatalf("Failed to get node from registry: %v", err)
		}
		if entry == nil {
			fmt.Println("Node account not found in registry")
		} else {
			fmt.Printf("Node registry entry:\n")
			fmt.Printf("  Parent: %s\n", entry.Parent)
			fmt.Printf("  Registered: %s\n", entry.Registred)
			fmt.Printf("  Domain: %s\n", entry.Domain)
			fmt.Printf("  Online: %d\n", entry.Online)
			fmt.Printf("  Active: %t\n", entry.Active)
		}

	case "delete-client":
		if len(os.Args) != 4 {
			log.Fatal("Usage: delete-client <registry_name> <account_to_delete>")
		}
		registryName := os.Args[2]
		account, err := solana.PublicKeyFromBase58(os.Args[3])
		if err != nil {
			log.Fatalf("Invalid account address: %v", err)
		}
		sig, err := client.DeleteClientFromRegistry(ctx, registryName, account)
		if err != nil {
			log.Fatalf("Failed to delete client from registry: %v", err)
		}
		fmt.Printf("Client account deleted from registry. Transaction signature: %s\n", sig)

	case "delete-node":
		if len(os.Args) != 4 {
			log.Fatal("Usage: delete-node <registry_name> <account_to_delete>")
		}
		registryName := os.Args[2]
		account, err := solana.PublicKeyFromBase58(os.Args[3])
		if err != nil {
			log.Fatalf("Invalid account address: %v", err)
		}
		sig, err := client.DeleteNodeFromRegistry(ctx, registryName, account)
		if err != nil {
			log.Fatalf("Failed to delete node from registry: %v", err)
		}
		fmt.Printf("Node account deleted from registry. Transaction signature: %s\n", sig)

	case "balance":
		balance, err := client.GetBalance(ctx)
		if err != nil {
			log.Fatalf("Failed to get balance: %v", err)
		}
		fmt.Printf("Wallet balance: %.9f SOL (%d lamports)\n", float64(balance)/LAMPORTS_PER_SOL, balance)

	case "airdrop":
		amount := uint64(LAMPORTS_PER_SOL) // Default 1 SOL
		if len(os.Args) > 2 {
			solAmount, err := strconv.ParseFloat(os.Args[2], 64)
			if err != nil {
				log.Fatalf("Invalid amount: %v", err)
			}
			amount = uint64(solAmount * LAMPORTS_PER_SOL)
		}

		sig, err := client.RequestAirdrop(ctx, amount)
		if err != nil {
			log.Fatalf("Failed to request airdrop: %v", err)
		}
		fmt.Printf("Airdrop requested. Transaction signature: %s\n", sig)

		// Get and display new balance
		balance, err := client.GetBalance(ctx)
		if err != nil {
			log.Fatalf("Failed to get new balance: %v", err)
		}
		fmt.Printf("New wallet balance: %.9f SOL (%d lamports)\n", float64(balance)/LAMPORTS_PER_SOL, balance)

	case "list-clients":
		if len(os.Args) != 3 {
			log.Fatal("Usage: list-clients <registry_name>")
		}
		registryName := os.Args[2]
		entries, err := client.ListClientsInRegistry(ctx, registryName)
		if err != nil {
			log.Fatalf("Failed to list clients: %v", err)
		}

		if len(entries) == 0 {
			fmt.Println("No clients found in registry")
			return
		}

		fmt.Printf("Found %d clients in registry:\n", len(entries))
		for i, entry := range entries {
			fmt.Printf("\nClient #%d:\n", i+1)
			fmt.Printf("  Registered: %s\n", entry.Registred)
			fmt.Printf("  Valid until: %s\n", time.Unix(entry.Until, 0))
			fmt.Printf("  Limit: %d\n", entry.Limit)
		}

	case "list-nodes":
		if len(os.Args) != 3 {
			log.Fatal("Usage: list-nodes <registry_name>")
		}
		registryName := os.Args[2]
		entries, err := client.ListNodesInRegistry(ctx, registryName)
		if err != nil {
			log.Fatalf("Failed to list nodes: %v", err)
		}

		if len(entries) == 0 {
			fmt.Println("No nodes found in registry")
			return
		}

		fmt.Printf("Found %d nodes in registry:\n", len(entries))
		for i, entry := range entries {
			fmt.Printf("\nNode #%d:\n", i+1)
			fmt.Printf("  Registered: %s\n", entry.Registred)
			fmt.Printf("  Domain: %s\n", entry.Domain)
			fmt.Printf("  Online: %d\n", entry.Online)
			fmt.Printf("  Active: %t\n", entry.Active)
		}

	case "update-node-online":
		if len(os.Args) != 6 {
			log.Fatal("Usage: update-node-online <registry_name> <authority> <account_to_update> <value>")
		}
		registryName := os.Args[2]
		authority, err := solana.PublicKeyFromBase58(os.Args[3])
		if err != nil {
			log.Fatalf("Invalid authority address: %v", err)
		}
		account, err := solana.PublicKeyFromBase58(os.Args[4])
		if err != nil {
			log.Fatalf("Invalid account address: %v", err)
		}
		value := int32(0)
		if _, err := fmt.Sscanf(os.Args[5], "%d", &value); err != nil {
			log.Fatalf("Invalid online value: %v", err)
		}

		sig, err := client.UpdateNodeOnline(ctx, registryName, authority, account, value)
		if err != nil {
			log.Fatalf("Failed to update node online status: %v", err)
		}
		fmt.Printf("Node online status updated. Transaction signature: %s\n", sig)

	case "update-node-active":
		if len(os.Args) != 6 {
			log.Fatal("Usage: update-node-active <registry_name> <authority> <account_to_update> <active>")
		}
		registryName := os.Args[2]
		authority, err := solana.PublicKeyFromBase58(os.Args[3])
		if err != nil {
			log.Fatalf("Invalid authority address: %v", err)
		}
		account, err := solana.PublicKeyFromBase58(os.Args[4])
		if err != nil {
			log.Fatalf("Invalid account address: %v", err)
		}
		active := false
		if _, err := fmt.Sscanf(os.Args[5], "%t", &active); err != nil {
			log.Fatalf("Invalid active value (must be true/false): %v", err)
		}

		sig, err := client.UpdateNodeActive(ctx, registryName, authority, account, active)
		if err != nil {
			log.Fatalf("Failed to update node active status: %v", err)
		}
		fmt.Printf("Node active status updated. Transaction signature: %s\n", sig)

	case "transfer":
		if len(os.Args) != 4 {
			log.Fatal("Usage: transfer <to_address> <amount_in_sol>")
		}
		toAddress, err := solana.PublicKeyFromBase58(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid destination address: %v", err)
		}
		solAmount, err := strconv.ParseFloat(os.Args[3], 64)
		if err != nil {
			log.Fatalf("Invalid amount: %v", err)
		}
		lamports := uint64(solAmount * LAMPORTS_PER_SOL)

		// Get current balance before transfer
		balance, err := client.GetBalance(ctx)
		if err != nil {
			log.Fatalf("Failed to get balance: %v", err)
		}
		if balance < lamports {
			log.Fatalf("Insufficient balance: have %.9f SOL, need %.9f SOL", float64(balance)/LAMPORTS_PER_SOL, solAmount)
		}

		sig, err := client.TransferSol(ctx, toAddress, lamports)
		if err != nil {
			log.Fatalf("Failed to transfer SOL: %v", err)
		}
		fmt.Printf("Transferred %.9f SOL to %s\n", solAmount, toAddress)
		fmt.Printf("Transaction signature: %s\n", sig)

		// Get and display new balance
		newBalance, err := client.GetBalance(ctx)
		if err != nil {
			log.Fatalf("Failed to get new balance: %v", err)
		}
		fmt.Printf("New wallet balance: %.9f SOL\n", float64(newBalance)/LAMPORTS_PER_SOL)
	case "delegate-node":
		if len(os.Args) != 4 {
			log.Fatal("Usage: delegate-node <registry_name> <account>")
		}
		registryName := os.Args[2]
		account, err := solana.PublicKeyFromBase58(os.Args[3])
		if err != nil {
			log.Fatalf("Invalid account address: %v", err)
		}

		sig, err := client.DelegateNode(ctx, registryName, account)
		if err != nil {
			log.Fatalf("Failed to delegate node: %v", err)
		}
		fmt.Printf("Node account delegated. Transaction signature: %s\n", sig)
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  create <registry_name>")
	fmt.Println("  add-client <registry_name> <account_to_add> <valid_days> <limit>")
	fmt.Println("  add-node <registry_name> <account_to_add> <domain>")
	fmt.Println("  delegate-node <registry_name> <account_to_add>")
	fmt.Println("  get-client <registry_name> <account_to_check>")
	fmt.Println("  get-node <registry_name> <account_to_check>")
	fmt.Println("  delete-client <registry_name> <account_to_delete>")
	fmt.Println("  delete-node <registry_name> <account_to_delete>")
	fmt.Println("  list-clients <registry_name>")
	fmt.Println("  list-nodes <registry_name>")
	fmt.Println("  update-node-online <registry_name> <authority> <account_to_update> <value>")
	fmt.Println("  update-node-active <registry_name> <authority> <account_to_update> <active>")
	fmt.Println("  transfer <to_address> <amount_in_sol>")
	fmt.Println("  balance")
	fmt.Println("  airdrop [amount_in_sol]")
}
