package registry

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
	"context"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

const (
	ClientEntrySize = 8 + 32 + 32 + 8 + 4 // discriminator + parent + registered + until + limit
	NodeEntrySize = 8 + 32 + 32 + 4 + 253 + 4 + 1  // discriminator + parent + registered + domain length + domain + online + active
)

// Anchor instruction discriminators
var (
	InitRegistryDiscriminator = []byte{131, 22, 4, 103, 24, 94, 163, 239}
	AddClientToRegistryDiscriminator = []byte{198, 64, 62, 101, 62, 204, 69, 108}
	AddNodeToRegistryDiscriminator = []byte{135, 249, 13, 74, 61, 190, 188, 33}
	CheckClientDiscriminator = []byte{56, 122, 178, 30, 199, 2, 243, 22}
	CheckNodeDiscriminator = []byte{62, 101, 38, 142, 134, 79, 122, 116}
	RemoveClientFromRegistryDiscriminator = []byte{32, 83, 79, 126, 155, 239, 104, 60}
	RemoveNodeFromRegistryDiscriminator = []byte{96, 10, 183, 238, 187, 248, 96, 36}
	UpdateNodeOnlineDiscriminator = []byte{35, 22, 232, 250, 60, 30, 62, 83}
	UpdateNodeActiveDiscriminator = []byte{121, 150, 132, 175, 172, 145, 197, 132}
)

// findRegistryPDA finds the PDA for a registry with the given name
func findRegistryPDA(programID solana.PublicKey, authority solana.PublicKey, name string) (solana.PublicKey, uint8, error) {
	return solana.FindProgramAddress(
		[][]byte{
			authority.Bytes(),
			[]byte(name),
		},
		programID,
	)
}

// findRegistryEntryPDA finds the PDA for a registry entry
func findRegistryEntryPDA(programID solana.PublicKey, accountToAdd solana.PublicKey, registry solana.PublicKey) (solana.PublicKey, uint8, error) {
	return solana.FindProgramAddress(
		[][]byte{
			accountToAdd.Bytes(),
			registry.Bytes(),
		},
		programID,
	)
}

// buildInitRegistryInstruction builds the instruction to initialize a new registry
func buildInitRegistryInstruction(
	programID solana.PublicKey,
	authority solana.PublicKey,
	name string,
) (solana.Instruction, solana.PublicKey, error) {
	registryPDA, _, err := findRegistryPDA(programID, authority, name)
	if err != nil {
		return nil, solana.PublicKey{}, fmt.Errorf("failed to find registry PDA: %v", err)
	}

	// Encode the instruction data
	data := new(bytes.Buffer)
	// Write instruction discriminator
	data.Write(InitRegistryDiscriminator)
	// Encode the name string
	nameBytes := []byte(name)
	binary.Write(data, binary.LittleEndian, uint32(len(nameBytes))) // Use uint32 for string length
	data.Write(nameBytes)

	accounts := solana.AccountMetaSlice{
		solana.Meta(registryPDA).WRITE(),
		solana.Meta(authority).SIGNER().WRITE(),
		solana.Meta(solana.SystemProgramID),
	}

	return solana.NewInstruction(
		programID,
		accounts,
		data.Bytes(),
	), registryPDA, nil
}

// buildAddClientToRegistryInstruction builds the instruction to add a client account to the registry
func buildAddClientToRegistryInstruction(
	programID solana.PublicKey,
	authority solana.PublicKey,
	registry solana.PublicKey,
	accountToAdd solana.PublicKey,
	validUntil time.Time,
	limit uint32,
) (solana.Instruction, error) {
	entryPDA, _, err := findRegistryEntryPDA(programID, accountToAdd, registry)
	if err != nil {
		return nil, fmt.Errorf("failed to find entry PDA: %v", err)
	}

	// Encode the instruction data
	data := new(bytes.Buffer)
	// Write instruction discriminator
	data.Write(AddClientToRegistryDiscriminator)
	// Encode account to add
	data.Write(accountToAdd.Bytes())
	// Encode until timestamp
	binary.Write(data, binary.LittleEndian, validUntil.Unix())
	// Encode limit
	binary.Write(data, binary.LittleEndian, limit)

	accounts := solana.AccountMetaSlice{
		solana.Meta(entryPDA).WRITE(),
		solana.Meta(registry),
		solana.Meta(authority).SIGNER().WRITE(),
		solana.Meta(solana.SystemProgramID),
	}

	return solana.NewInstruction(
		programID,
		accounts,
		data.Bytes(),
	), nil
}

// buildAddNodeToRegistryInstruction builds the instruction to add a node account to the registry
func buildAddNodeToRegistryInstruction(
	programID solana.PublicKey,
	authority solana.PublicKey,
	registry solana.PublicKey,
	accountToAdd solana.PublicKey,
	domain string,
) (solana.Instruction, error) {
	entryPDA, _, err := findRegistryEntryPDA(programID, accountToAdd, registry)
	if err != nil {
		return nil, fmt.Errorf("failed to find entry PDA: %v", err)
	}

	// Encode the instruction data
	data := new(bytes.Buffer)
	// Write instruction discriminator
	data.Write(AddNodeToRegistryDiscriminator)
	// Encode account to add
	data.Write(accountToAdd.Bytes())
	// Encode domain string
	binary.Write(data, binary.LittleEndian, uint32(len(domain)))
	data.Write([]byte(domain))

	accounts := solana.AccountMetaSlice{
		solana.Meta(entryPDA).WRITE(),
		solana.Meta(registry),
		solana.Meta(authority).SIGNER().WRITE(),
		solana.Meta(solana.SystemProgramID),
	}

	return solana.NewInstruction(
		programID,
		accounts,
		data.Bytes(),
	), nil
}

// buildRemoveClientFromRegistryInstruction builds the instruction to remove a client account from the registry
func buildRemoveClientFromRegistryInstruction(
	programID solana.PublicKey,
	authority solana.PublicKey,
	registry solana.PublicKey,
	accountToDelete solana.PublicKey,
) (solana.Instruction, error) {
	entryPDA, _, err := findRegistryEntryPDA(programID, accountToDelete, registry)
	if err != nil {
		return nil, fmt.Errorf("failed to find entry PDA: %v", err)
	}

	// Encode the instruction data
	data := new(bytes.Buffer)
	// Write instruction discriminator
	data.Write(RemoveClientFromRegistryDiscriminator)
	// Encode account to delete
	data.Write(accountToDelete.Bytes())

	accounts := solana.AccountMetaSlice{
		solana.Meta(entryPDA).WRITE(),
		solana.Meta(registry),
		solana.Meta(authority).SIGNER().WRITE(),
	}

	return solana.NewInstruction(
		programID,
		accounts,
		data.Bytes(),
	), nil
}

// buildRemoveNodeFromRegistryInstruction builds the instruction to remove a node account from the registry
func buildRemoveNodeFromRegistryInstruction(
	programID solana.PublicKey,
	authority solana.PublicKey,
	registry solana.PublicKey,
	accountToDelete solana.PublicKey,
) (solana.Instruction, error) {
	entryPDA, _, err := findRegistryEntryPDA(programID, accountToDelete, registry)
	if err != nil {
		return nil, fmt.Errorf("failed to find entry PDA: %v", err)
	}

	// Encode the instruction data
	data := new(bytes.Buffer)
	// Write instruction discriminator
	data.Write(RemoveNodeFromRegistryDiscriminator)
	// Encode account to delete
	data.Write(accountToDelete.Bytes())

	accounts := solana.AccountMetaSlice{
		solana.Meta(entryPDA).WRITE(),
		solana.Meta(registry),
		solana.Meta(authority).SIGNER().WRITE(),
	}

	return solana.NewInstruction(
		programID,
		accounts,
		data.Bytes(),
	), nil
}

// buildUpdateNodeOnlineInstruction builds the instruction to update node online status
func buildUpdateNodeOnlineInstruction(
	programID solana.PublicKey,
	authority solana.PublicKey,
	registry solana.PublicKey,
	accountToUpdate solana.PublicKey,
	value int32,
) (solana.Instruction, error) {
	entryPDA, _, err := findRegistryEntryPDA(programID, accountToUpdate, registry)
	if err != nil {
		return nil, fmt.Errorf("failed to find entry PDA: %v", err)
	}

	// Encode the instruction data
	data := new(bytes.Buffer)
	// Write instruction discriminator
	data.Write(UpdateNodeOnlineDiscriminator)
	// Encode account to update
	data.Write(accountToUpdate.Bytes())
	// Encode online value
	binary.Write(data, binary.LittleEndian, value)

	accounts := solana.AccountMetaSlice{
		solana.Meta(entryPDA).WRITE(),
		solana.Meta(registry),
		solana.Meta(authority).SIGNER().WRITE(),
	}

	return solana.NewInstruction(
		programID,
		accounts,
		data.Bytes(),
	), nil
}

// buildUpdateNodeActiveInstruction builds the instruction to update node active status
func buildUpdateNodeActiveInstruction(
	programID solana.PublicKey,
	accountToUpdate solana.PublicKey,
	registry solana.PublicKey,
	authority solana.PublicKey,
	active bool,
) (solana.Instruction, error) {
	// Find the entry PDA
	entryPDA, _, err := findRegistryEntryPDA(programID, accountToUpdate, registry)
	if err != nil {
		return nil, fmt.Errorf("failed to find entry PDA: %v", err)
	}

	// Find the authority node PDA
	authorityNodePDA, _, err := findRegistryEntryPDA(programID, authority, registry)
	if err != nil {
		return nil, fmt.Errorf("failed to find authority node PDA: %v", err)
	}

	// Encode the instruction data
	data := new(bytes.Buffer)
	// Write instruction discriminator
	data.Write(UpdateNodeActiveDiscriminator)
	// Encode account to update
	data.Write(accountToUpdate.Bytes())
	// Encode active value
	if active {
		data.Write([]byte{1})
	} else {
		data.Write([]byte{0})
	}

	accounts := solana.AccountMetaSlice{
		solana.Meta(entryPDA).WRITE(),
		solana.Meta(registry),
		solana.Meta(authorityNodePDA),
		solana.Meta(authority).SIGNER(),
	}

	return solana.NewInstruction(
		programID,
		accounts,
		data.Bytes(),
	), nil
}

// getClientEntry retrieves a client entry account data
func getClientEntry(
	ctx context.Context,
	client *rpc.Client,
	programID solana.PublicKey,
	registry solana.PublicKey,
	accountToCheck solana.PublicKey,
) (*ClientEntry, error) {
	entryPDA, _, err := findRegistryEntryPDA(programID, accountToCheck, registry)
	if err != nil {
		return nil, fmt.Errorf("failed to find entry PDA: %v", err)
	}

	// Get the account info
	accountInfo, err := client.GetAccountInfo(ctx, entryPDA)
	if err != nil {
		return nil, fmt.Errorf("failed to get account info: %v", err)
	}

	if accountInfo == nil || len(accountInfo.Value.Data.GetBinary()) == 0 {
		return nil, nil // Account doesn't exist
	}

	// Parse the account data
	data := accountInfo.Value.Data.GetBinary()
	if len(data) != ClientEntrySize {
		return nil, fmt.Errorf("invalid account data size: expected %d, got %d", ClientEntrySize, len(data))
	}

	// Skip the 8-byte discriminator
	data = data[8:]

	entry := &ClientEntry{
		Parent:    solana.PublicKeyFromBytes(data[:32]),
		Registred: solana.PublicKeyFromBytes(data[32:64]),
		Until:     int64(binary.LittleEndian.Uint64(data[64:72])),
		Limit:     binary.LittleEndian.Uint32(data[72:76]),
	}

	return entry, nil
}

// getNodeEntry retrieves a node entry account data
func getNodeEntry(
	ctx context.Context,
	client *rpc.Client,
	programID solana.PublicKey,
	registry solana.PublicKey,
	accountToCheck solana.PublicKey,
) (*NodeEntry, error) {
	entryPDA, _, err := findRegistryEntryPDA(programID, accountToCheck, registry)
	if err != nil {
		return nil, fmt.Errorf("failed to find entry PDA: %v", err)
	}

	// Get the account info
	accountInfo, err := client.GetAccountInfo(ctx, entryPDA)
	if err != nil {
		return nil, fmt.Errorf("failed to get account info: %v", err)
	}

	if accountInfo == nil || len(accountInfo.Value.Data.GetBinary()) == 0 {
		return nil, nil // Account doesn't exist
	}

	// Parse the account data
	data := accountInfo.Value.Data.GetBinary()
	if len(data) != NodeEntrySize {
		return nil, fmt.Errorf("invalid account data size: expected %d, got %d", NodeEntrySize, len(data))
	}

	// Skip the 8-byte discriminator
	data = data[8:]

	// Read domain string length (4 bytes)
	domainLen := binary.LittleEndian.Uint32(data[64:68])
	
	entry := &NodeEntry{
		Parent:    solana.PublicKeyFromBytes(data[:32]),
		Registred: solana.PublicKeyFromBytes(data[32:64]),
		Domain:    string(data[68:68+domainLen]),
		Online:    int32(binary.LittleEndian.Uint32(data[68+domainLen:72+domainLen])),
		Active:    data[72+domainLen] == 1,
	}

	return entry, nil
}