package registry

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

// RegistryClient represents a client for interacting with the registry program
type RegistryClient struct {
	programID solana.PublicKey
	client    *rpc.Client
	wsClient  *ws.Client
	signer    solana.PrivateKey
}

// ClientEntry represents a client entry in the registry
type ClientEntry struct {
	Parent    solana.PublicKey
	Registred solana.PublicKey
	Until     int64
	Limit     uint32
}

// NodeEntry represents a node entry in the registry
type NodeEntry struct {
	Parent    solana.PublicKey
	Registred solana.PublicKey
	Domain    string
	Online    int32
	Active    bool
}

// NewRegistryClient creates a new instance of the registry client
func NewRegistryClient(rpcEndpoint string, wsEndpoint string, programID string, privateKey string) (*RegistryClient, error) {
	client := rpc.New(rpcEndpoint)

	wsClient, err := ws.Connect(context.Background(), wsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to websocket: %v", err)
	}

	programPubkey, err := solana.PublicKeyFromBase58(programID)
	if err != nil {
		return nil, fmt.Errorf("invalid program ID: %v", err)
	}

	privateKeyBytes, err := solana.PrivateKeyFromBase58(privateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}

	return &RegistryClient{
		programID: programPubkey,
		client:    client,
		wsClient:  wsClient,
		signer:    privateKeyBytes,
	}, nil
}

// CreateRegistry creates a new registry with the given name
func (c *RegistryClient) CreateRegistry(ctx context.Context, name string) (solana.Signature, error) {
	// Build the instruction
	instruction, _, err := buildInitRegistryInstruction(
		c.programID,
		c.signer.PublicKey(),
		name,
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to build instruction: %v", err)
	}

	// Create the transaction
	recent, err := c.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to get recent blockhash: %v", err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		recent.Value.Blockhash,
		solana.TransactionPayer(c.signer.PublicKey()),
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to create transaction: %v", err)
	}

	// Sign and send the transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(c.signer.PublicKey()) {
			return &c.signer
		}
		return nil
	})
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to sign transaction: %v", err)
	}

	sig, err := confirm.SendAndConfirmTransaction(
		ctx,
		c.client,
		c.wsClient,
		tx,
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to send and confirm transaction: %v", err)
	}

	return sig, nil
}

// AddClientToRegistry adds a client account to the registry
func (c *RegistryClient) AddClientToRegistry(ctx context.Context, registryName string, accountToAdd solana.PublicKey, validUntil time.Time, limit uint32) (solana.Signature, error) {
	// Find the registry PDA
	registryPDA, _, err := findRegistryPDA(c.programID, c.signer.PublicKey(), registryName)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to find registry PDA: %v", err)
	}

	// Build the instruction
	instruction, err := buildAddClientToRegistryInstruction(
		c.programID,
		c.signer.PublicKey(),
		registryPDA,
		accountToAdd,
		validUntil,
		limit,
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to build instruction: %v", err)
	}

	// Create the transaction
	recent, err := c.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to get recent blockhash: %v", err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		recent.Value.Blockhash,
		solana.TransactionPayer(c.signer.PublicKey()),
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to create transaction: %v", err)
	}

	// Sign and send the transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(c.signer.PublicKey()) {
			return &c.signer
		}
		return nil
	})
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to sign transaction: %v", err)
	}

	sig, err := confirm.SendAndConfirmTransaction(
		ctx,
		c.client,
		c.wsClient,
		tx,
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to send and confirm transaction: %v", err)
	}

	return sig, nil
}

// AddNodeToRegistry adds a node account to the registry
func (c *RegistryClient) AddNodeToRegistry(ctx context.Context, registryName string, accountToAdd solana.PublicKey, domain string) (solana.Signature, error) {
	if len(domain) > 253 {
		return solana.Signature{}, fmt.Errorf("domain name must be 253 characters or less")
	}

	// Find the registry PDA
	registryPDA, _, err := findRegistryPDA(c.programID, c.signer.PublicKey(), registryName)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to find registry PDA: %v", err)
	}

	// Build the instruction
	instruction, err := buildAddNodeToRegistryInstruction(
		c.programID,
		c.signer.PublicKey(),
		registryPDA,
		accountToAdd,
		domain,
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to build instruction: %v", err)
	}

	// Create the transaction
	recent, err := c.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to get recent blockhash: %v", err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		recent.Value.Blockhash,
		solana.TransactionPayer(c.signer.PublicKey()),
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to create transaction: %v", err)
	}

	// Sign and send the transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(c.signer.PublicKey()) {
			return &c.signer
		}
		return nil
	})
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to sign transaction: %v", err)
	}

	sig, err := confirm.SendAndConfirmTransaction(
		ctx,
		c.client,
		c.wsClient,
		tx,
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to send and confirm transaction: %v", err)
	}

	return sig, nil
}

func (c *RegistryClient) DelegateNode(ctx context.Context, registryName string, account solana.PublicKey) (solana.Signature, error) {
	// Find the registry PDA
	registryPDA, _, err := findRegistryPDA(c.programID, c.signer.PublicKey(), registryName)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to find registry PDA: %v", err)
	}

	instruction, err := buildDelegateNodeAccountInstruction(
		c.programID,
		c.signer.PublicKey(),
		registryPDA,
		account,
	)

	// Create the transaction
	recent, err := c.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to get recent blockhash: %v", err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		recent.Value.Blockhash,
		solana.TransactionPayer(c.signer.PublicKey()),
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to create transaction: %v", err)
	}

	// Sign and send the transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(c.signer.PublicKey()) {
			return &c.signer
		}
		return nil
	})
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to sign transaction: %v", err)
	}

	sig, err := confirm.SendAndConfirmTransaction(
		ctx,
		c.client,
		c.wsClient,
		tx,
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to send and confirm transaction: %v", err)
	}

	return sig, nil
}

// GetClientFromRegistry retrieves a client entry from the registry
func (c *RegistryClient) GetClientFromRegistry(ctx context.Context, registryName string, accountToCheck solana.PublicKey) (*ClientEntry, error) {
	// Find the registry PDA
	registryPDA, _, err := findRegistryPDA(c.programID, c.signer.PublicKey(), registryName)
	if err != nil {
		return nil, fmt.Errorf("failed to find registry PDA: %v", err)
	}

	return getClientEntry(ctx, c.client, c.programID, registryPDA, accountToCheck)
}

// GetNodeFromRegistry retrieves a node entry from the registry
func (c *RegistryClient) GetNodeFromRegistry(ctx context.Context, registryName string, accountToCheck solana.PublicKey) (*NodeEntry, error) {
	// Find the registry PDA
	registryPDA, _, err := findRegistryPDA(c.programID, c.signer.PublicKey(), registryName)
	if err != nil {
		return nil, fmt.Errorf("failed to find registry PDA: %v", err)
	}

	return getNodeEntry(ctx, c.client, c.programID, registryPDA, accountToCheck)
}

// DeleteClientFromRegistry removes a client account from the registry
func (c *RegistryClient) DeleteClientFromRegistry(ctx context.Context, registryName string, accountToDelete solana.PublicKey) (solana.Signature, error) {
	// Find the registry PDA
	registryPDA, _, err := findRegistryPDA(c.programID, c.signer.PublicKey(), registryName)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to find registry PDA: %v", err)
	}

	// Build the instruction
	instruction, err := buildRemoveClientFromRegistryInstruction(
		c.programID,
		c.signer.PublicKey(),
		registryPDA,
		accountToDelete,
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to build instruction: %v", err)
	}

	// Create the transaction
	recent, err := c.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to get recent blockhash: %v", err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		recent.Value.Blockhash,
		solana.TransactionPayer(c.signer.PublicKey()),
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to create transaction: %v", err)
	}

	// Sign and send the transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(c.signer.PublicKey()) {
			return &c.signer
		}
		return nil
	})
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to sign transaction: %v", err)
	}

	sig, err := confirm.SendAndConfirmTransaction(
		ctx,
		c.client,
		c.wsClient,
		tx,
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to send and confirm transaction: %v", err)
	}

	return sig, nil
}

// DeleteNodeFromRegistry removes a node account from the registry
func (c *RegistryClient) DeleteNodeFromRegistry(ctx context.Context, registryName string, accountToDelete solana.PublicKey) (solana.Signature, error) {
	// Find the registry PDA
	registryPDA, _, err := findRegistryPDA(c.programID, c.signer.PublicKey(), registryName)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to find registry PDA: %v", err)
	}

	// Build the instruction
	instruction, err := buildRemoveNodeFromRegistryInstruction(
		c.programID,
		c.signer.PublicKey(),
		registryPDA,
		accountToDelete,
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to build instruction: %v", err)
	}

	// Create the transaction
	recent, err := c.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to get recent blockhash: %v", err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		recent.Value.Blockhash,
		solana.TransactionPayer(c.signer.PublicKey()),
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to create transaction: %v", err)
	}

	// Sign and send the transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(c.signer.PublicKey()) {
			return &c.signer
		}
		return nil
	})
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to sign transaction: %v", err)
	}

	sig, err := confirm.SendAndConfirmTransaction(
		ctx,
		c.client,
		c.wsClient,
		tx,
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to send and confirm transaction: %v", err)
	}

	return sig, nil
}

// RequestAirdrop requests an airdrop of SOL to the signer's wallet
func (c *RegistryClient) RequestAirdrop(ctx context.Context, amount uint64) (solana.Signature, error) {
	sig, err := c.client.RequestAirdrop(
		ctx,
		c.signer.PublicKey(),
		amount,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to request airdrop: %v", err)
	}

	// Wait for confirmation
	confirmed, err := confirm.WaitForConfirmation(
		ctx,
		c.wsClient,
		sig,
		nil,
	)
	if err != nil {
		return sig, fmt.Errorf("failed to confirm airdrop: %v", err)
	}

	if !confirmed {
		return sig, fmt.Errorf("airdrop transaction was not confirmed")
	}

	return sig, nil
}

// GetBalance returns the current balance of the signer's wallet in lamports
func (c *RegistryClient) GetBalance(ctx context.Context) (uint64, error) {
	balance, err := c.client.GetBalance(
		ctx,
		c.signer.PublicKey(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %v", err)
	}

	return balance.Value, nil
}

// Close closes the websocket connection
func (c *RegistryClient) Close() {
	if c.wsClient != nil {
		c.wsClient.Close()
	}
}

// ListClientsInRegistry retrieves all client entries in the given registry
func (c *RegistryClient) ListClientsInRegistry(ctx context.Context, registryName string) ([]*ClientEntry, error) {
	// Find the registry PDA
	registryPDA, _, err := findRegistryPDA(c.programID, c.signer.PublicKey(), registryName)
	if err != nil {
		return nil, fmt.Errorf("failed to find registry PDA: %v", err)
	}

	// Get all program accounts of type ClientEntry
	filters := []rpc.RPCFilter{
		{
			Memcmp: &rpc.RPCFilterMemcmp{
				Offset: 8, // Skip discriminator
				Bytes:  registryPDA.Bytes(),
			},
		},
		{
			DataSize: ClientEntrySize,
		},
	}

	accounts, err := c.client.GetProgramAccountsWithOpts(
		ctx,
		c.programID,
		&rpc.GetProgramAccountsOpts{
			Filters:    filters,
			Commitment: rpc.CommitmentFinalized,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get program accounts: %v", err)
	}

	entries := make([]*ClientEntry, 0, len(accounts))
	for _, acc := range accounts {
		data := acc.Account.Data.GetBinary()
		if len(data) != ClientEntrySize {
			continue
		}

		// Skip the 8-byte discriminator
		data = data[8:]

		entry := &ClientEntry{
			Parent:    solana.PublicKeyFromBytes(data[:32]),
			Registred: solana.PublicKeyFromBytes(data[32:64]),
			Until:     int64(binary.LittleEndian.Uint64(data[64:72])),
			Limit:     binary.LittleEndian.Uint32(data[72:76]),
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// ListNodesInRegistry retrieves all node entries in the given registry
func (c *RegistryClient) ListNodesInRegistry(ctx context.Context, registryName string) ([]*NodeEntry, error) {
	// Find the registry PDA
	registryPDA, _, err := findRegistryPDA(c.programID, c.signer.PublicKey(), registryName)
	if err != nil {
		return nil, fmt.Errorf("failed to find registry PDA: %v", err)
	}

	// Get all program accounts of type NodeEntry
	filters := []rpc.RPCFilter{
		{
			Memcmp: &rpc.RPCFilterMemcmp{
				Offset: 8, // Skip discriminator
				Bytes:  registryPDA.Bytes(),
			},
		},
		{
			DataSize: NodeEntrySize,
		},
	}

	accounts, err := c.client.GetProgramAccountsWithOpts(
		ctx,
		c.programID,
		&rpc.GetProgramAccountsOpts{
			Filters:    filters,
			Commitment: rpc.CommitmentFinalized,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get program accounts: %v", err)
	}

	entries := make([]*NodeEntry, 0, len(accounts))
	for _, acc := range accounts {
		data := acc.Account.Data.GetBinary()
		if len(data) != NodeEntrySize {
			continue
		}

		// Skip the 8-byte discriminator
		data = data[8:]

		// Read domain string length (4 bytes)
		domainLen := binary.LittleEndian.Uint32(data[64:68])

		entry := &NodeEntry{
			Parent:    solana.PublicKeyFromBytes(data[:32]),
			Registred: solana.PublicKeyFromBytes(data[32:64]),
			Domain:    string(data[68 : 68+domainLen]),
			Online:    int32(binary.LittleEndian.Uint32(data[68+domainLen : 72+domainLen])),
			Active:    data[72+domainLen] == 1,
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// UpdateNodeOnline updates the online status of a node in the registry
func (c *RegistryClient) UpdateNodeOnline(ctx context.Context, registryName string, authority solana.PublicKey, accountToUpdate solana.PublicKey, value int32) (solana.Signature, error) {
	if value < 0 {
		return solana.Signature{}, fmt.Errorf("online value must be non-negative")
	}

	// Find the registry PDA using the provided authority
	registryPDA, _, err := findRegistryPDA(c.programID, authority, registryName)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to find registry PDA: %v", err)
	}

	// Build the instruction
	instruction, err := buildUpdateNodeOnlineInstruction(
		c.programID,
		accountToUpdate,
		registryPDA,
		accountToUpdate,
		value,
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to build instruction: %v", err)
	}

	// Create the transaction
	recent, err := c.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to get recent blockhash: %v", err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		recent.Value.Blockhash,
		solana.TransactionPayer(c.signer.PublicKey()),
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to create transaction: %v", err)
	}

	// Sign and send the transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(c.signer.PublicKey()) {
			return &c.signer
		}
		return nil
	})
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to sign transaction: %v", err)
	}

	sig, err := confirm.SendAndConfirmTransaction(
		ctx,
		c.client,
		c.wsClient,
		tx,
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to send and confirm transaction: %v", err)
	}

	return sig, nil
}

// UpdateNodeActive updates the active status of a node in the registry
func (c *RegistryClient) UpdateNodeActive(ctx context.Context, registryName string, authority solana.PublicKey, accountToUpdate solana.PublicKey, active bool) (solana.Signature, error) {
	// Find the registry PDA using the provided authority
	registryPDA, _, err := findRegistryPDA(c.programID, authority, registryName)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to find registry PDA: %v", err)
	}

	// Build the instruction
	instruction, err := buildUpdateNodeActiveInstruction(
		c.programID,
		accountToUpdate,
		registryPDA,
		c.signer.PublicKey(),
		active,
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to build instruction: %v", err)
	}

	// Create the transaction
	recent, err := c.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to get recent blockhash: %v", err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		recent.Value.Blockhash,
		solana.TransactionPayer(c.signer.PublicKey()),
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to create transaction: %v", err)
	}

	// Sign and send the transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(c.signer.PublicKey()) {
			return &c.signer
		}
		return nil
	})
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to sign transaction: %v", err)
	}

	sig, err := confirm.SendAndConfirmTransaction(
		ctx,
		c.client,
		c.wsClient,
		tx,
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to send and confirm transaction: %v", err)
	}

	return sig, nil
}

// TransferSol transfers SOL from the signer's wallet to the target address
func (c *RegistryClient) TransferSol(ctx context.Context, to solana.PublicKey, amount uint64) (solana.Signature, error) {
	// Create the transfer instruction
	instruction := system.NewTransferInstruction(
		amount,
		c.signer.PublicKey(),
		to,
	).Build()

	// Create the transaction
	recent, err := c.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to get recent blockhash: %v", err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		recent.Value.Blockhash,
		solana.TransactionPayer(c.signer.PublicKey()),
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to create transaction: %v", err)
	}

	// Sign and send the transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(c.signer.PublicKey()) {
			return &c.signer
		}
		return nil
	})
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to sign transaction: %v", err)
	}

	sig, err := confirm.SendAndConfirmTransaction(
		ctx,
		c.client,
		c.wsClient,
		tx,
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to send and confirm transaction: %v", err)
	}

	return sig, nil
}
