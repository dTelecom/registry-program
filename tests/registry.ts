import * as anchor from "@coral-xyz/anchor";
import { Program } from "@coral-xyz/anchor";
import { Registry } from "../target/types/registry";
import { PublicKey, Keypair, SystemProgram } from "@solana/web3.js";
import { expect } from "chai";

describe("registry", () => {
  const provider = anchor.AnchorProvider.env();
  anchor.setProvider(provider);

  const program = anchor.workspace.Registry as Program<Registry>;
  const authority = provider.wallet;
  const registryName = "test_registry";
  const clientAccount = Keypair.generate();
  const nodeAccount = Keypair.generate();

  let registryPDA: PublicKey;
  let clientEntryPDA: PublicKey;
  let nodeEntryPDA: PublicKey;

  before(async () => {
    // Find PDAs
    [registryPDA] = await PublicKey.findProgramAddress(
      [authority.publicKey.toBuffer(), Buffer.from(registryName)],
      program.programId
    );

    [clientEntryPDA] = await PublicKey.findProgramAddress(
      [clientAccount.publicKey.toBuffer(), registryPDA.toBuffer()],
      program.programId
    );

    [nodeEntryPDA] = await PublicKey.findProgramAddress(
      [nodeAccount.publicKey.toBuffer(), registryPDA.toBuffer()],
      program.programId
    );
  });

  it("Creates a new registry", async () => {
    await program.methods
      .initRegistry(registryName)
      .accounts({
        registry: registryPDA,
        authority: authority.publicKey,
        systemProgram: SystemProgram.programId,
      })
      .rpc();

    const registry = await program.account.registry.fetch(registryPDA);
    expect(registry.authority.toString()).to.equal(authority.publicKey.toString());
    expect(registry.name).to.equal(registryName);
  });

  it("Adds a client to the registry", async () => {
    const now = Math.floor(Date.now() / 1000);
    const validUntil = now + 30 * 24 * 60 * 60; // 30 days from now
    const limit = 1000;

    await program.methods
      .addClientToRegistry(clientAccount.publicKey, new anchor.BN(validUntil), limit)
      .accounts({
        entry: clientEntryPDA,
        registry: registryPDA,
        authority: authority.publicKey,
        systemProgram: SystemProgram.programId,
      })
      .rpc();

    const entry = await program.account.clientEntry.fetch(clientEntryPDA);
    expect(entry.parent.toString()).to.equal(registryPDA.toString());
    expect(entry.registred.toString()).to.equal(clientAccount.publicKey.toString());
    expect(entry.until.toNumber()).to.equal(validUntil);
    expect(entry.limit).to.equal(limit);
  });

  it("Checks client registration", async () => {
    const info = await program.methods
      .checkClient(clientAccount.publicKey)
      .accounts({
        entry: clientEntryPDA,
        registry: registryPDA,
      })
      .view();

    expect(info.until.toNumber()).to.be.greaterThan(Math.floor(Date.now() / 1000));
    expect(info.limit).to.equal(1000);
  });

  it("Removes client from registry", async () => {
    await program.methods
      .removeClientFromRegistry(clientAccount.publicKey)
      .accounts({
        entry: clientEntryPDA,
        registry: registryPDA,
        authority: authority.publicKey,
      })
      .rpc();

    // Verify the account is closed
    const accountInfo = await provider.connection.getAccountInfo(clientEntryPDA);
    expect(accountInfo).to.be.null;
  });

  it("Adds a node to the registry", async () => {
    const domain = "example.com";

    await program.methods
      .addNodeToRegistry(nodeAccount.publicKey, domain)
      .accounts({
        entry: nodeEntryPDA,
        registry: registryPDA,
        authority: authority.publicKey,
        systemProgram: SystemProgram.programId,
      })
      .rpc();

    const entry = await program.account.nodeEntry.fetch(nodeEntryPDA);
    expect(entry.parent.toString()).to.equal(registryPDA.toString());
    expect(entry.registred.toString()).to.equal(nodeAccount.publicKey.toString());
    expect(entry.domain).to.equal(domain);
    expect(entry.online).to.equal(0);
    expect(entry.active).to.equal(false);
  });

  it("Fails to add node with domain name exceeding 253 characters", async () => {
    const longDomain = "a".repeat(254);
    const newNodeAccount = Keypair.generate();
    const [newNodeEntryPDA] = await PublicKey.findProgramAddress(
      [newNodeAccount.publicKey.toBuffer(), registryPDA.toBuffer()],
      program.programId
    );

    try {
      await program.methods
        .addNodeToRegistry(newNodeAccount.publicKey, longDomain)
        .accounts({
          entry: newNodeEntryPDA,
          registry: registryPDA,
          authority: authority.publicKey,
          systemProgram: SystemProgram.programId,
        })
        .rpc();
      expect.fail("Expected error was not thrown");
    } catch (error) {
      expect(error.toString()).to.include("DomainTooLong");
    }
  });

  it("Successfully adds node with maximum length domain name (253 characters)", async () => {
    const maxDomain = "a".repeat(253);
    const newNodeAccount = Keypair.generate();
    const [newNodeEntryPDA] = await PublicKey.findProgramAddress(
      [newNodeAccount.publicKey.toBuffer(), registryPDA.toBuffer()],
      program.programId
    );

    await program.methods
      .addNodeToRegistry(newNodeAccount.publicKey, maxDomain)
      .accounts({
        entry: newNodeEntryPDA,
        registry: registryPDA,
        authority: authority.publicKey,
        systemProgram: SystemProgram.programId,
      })
      .rpc();

    const entry = await program.account.nodeEntry.fetch(newNodeEntryPDA);
    expect(entry.domain).to.equal(maxDomain);
    expect(entry.domain.length).to.equal(253);
  });

  it("Updates node online status", async () => {
    // First, we need to airdrop some SOL to the node account to make it a valid signer
    const signature = await provider.connection.requestAirdrop(
      nodeAccount.publicKey,
      1000000000 // 1 SOL
    );
    await provider.connection.confirmTransaction(signature);

    // Try to update with a valid value
    const onlineValue = 42;
    await program.methods
      .updateNodeOnline(nodeAccount.publicKey, onlineValue)
      .accounts({
        entry: nodeEntryPDA,
        registry: registryPDA,
        authority: nodeAccount.publicKey,
      })
      .signers([nodeAccount])
      .rpc();

    let entry = await program.account.nodeEntry.fetch(nodeEntryPDA);
    expect(entry.online).to.equal(onlineValue);

    // Try to update with another valid value
    const newOnlineValue = 100;
    await program.methods
      .updateNodeOnline(nodeAccount.publicKey, newOnlineValue)
      .accounts({
        entry: nodeEntryPDA,
        registry: registryPDA,
        authority: nodeAccount.publicKey,
      })
      .signers([nodeAccount])
      .rpc();

    entry = await program.account.nodeEntry.fetch(nodeEntryPDA);
    expect(entry.online).to.equal(newOnlineValue);
  });

  it("Fails to update node online status with negative value", async () => {
    try {
      await program.methods
        .updateNodeOnline(nodeAccount.publicKey, -1)
        .accounts({
          entry: nodeEntryPDA,
          registry: registryPDA,
          authority: nodeAccount.publicKey,
        })
        .signers([nodeAccount])
        .rpc();
      expect.fail("Expected error was not thrown");
    } catch (error) {
      expect(error.toString()).to.include("InvalidOnlineValue");
    }
  });

  it("Updates node active status by node owner", async () => {
    // First, we need to airdrop some SOL to the node account if not already done
    try {
      const signature = await provider.connection.requestAirdrop(
        nodeAccount.publicKey,
        1000000000 // 1 SOL
      );
      await provider.connection.confirmTransaction(signature);
    } catch (error) {
      // Ignore error if account already has SOL
    }

    await program.methods
      .updateNodeActive(nodeAccount.publicKey, true)
      .accounts({
        entry: nodeEntryPDA,
        registry: registryPDA,
        authorityNode: nodeEntryPDA,
        authority: nodeAccount.publicKey,
      })
      .signers([nodeAccount])
      .rpc();

    const entry = await program.account.nodeEntry.fetch(nodeEntryPDA);
    expect(entry.active).to.equal(true);
  });

  it("Updates node active status by node in same registry", async () => {
    // Create another node in the same registry
    const otherNode = Keypair.generate();
    const [otherNodeEntryPDA] = await PublicKey.findProgramAddress(
      [otherNode.publicKey.toBuffer(), registryPDA.toBuffer()],
      program.programId
    );

    // Airdrop SOL to the other node
    const signature = await provider.connection.requestAirdrop(
      otherNode.publicKey,
      1000000000 // 1 SOL
    );
    await provider.connection.confirmTransaction(signature);

    // Add the other node to the registry
    await program.methods
      .addNodeToRegistry(otherNode.publicKey, "other.example.com")
      .accounts({
        entry: otherNodeEntryPDA,
        registry: registryPDA,
        authority: authority.publicKey,
        systemProgram: SystemProgram.programId,
      })
      .rpc();

    // First node updates other node's active status
    await program.methods
      .updateNodeActive(otherNode.publicKey, true)
      .accounts({
        entry: otherNodeEntryPDA,
        registry: registryPDA,
        authorityNode: nodeEntryPDA,
        authority: nodeAccount.publicKey,
      })
      .signers([nodeAccount])
      .rpc();

    const entry = await program.account.nodeEntry.fetch(otherNodeEntryPDA);
    expect(entry.active).to.equal(true);
  });

  it("Updates node active status by node from different registry must fail", async () => {
    // Create a different registry
    const otherRegistryName = "other_registry";
    const [otherRegistryPDA] = await PublicKey.findProgramAddress(
      [authority.publicKey.toBuffer(), Buffer.from(otherRegistryName)],
      program.programId
    );

    // Create the other registry
    await program.methods
      .initRegistry(otherRegistryName)
      .accounts({
        registry: otherRegistryPDA,
        authority: authority.publicKey,
        systemProgram: SystemProgram.programId,
      })
      .rpc();

    // Create a node in the other registry
    const otherNode = Keypair.generate();
    const [otherNodeEntryPDA] = await PublicKey.findProgramAddress(
      [otherNode.publicKey.toBuffer(), otherRegistryPDA.toBuffer()],
      program.programId
    );

    // Airdrop SOL to the other node
    const signature = await provider.connection.requestAirdrop(
      otherNode.publicKey,
      1000000000 // 1 SOL
    );
    await provider.connection.confirmTransaction(signature);

    // Add the other node to the other registry
    await program.methods
      .addNodeToRegistry(otherNode.publicKey, "other.example.com")
      .accounts({
        entry: otherNodeEntryPDA,
        registry: otherRegistryPDA,
        authority: authority.publicKey,
        systemProgram: SystemProgram.programId,
      })
      .rpc();

    // Try to update node in first registry from node in second registry - should fail
    try {
      await program.methods
        .updateNodeActive(nodeAccount.publicKey, false)
        .accounts({
          entry: nodeEntryPDA,
          registry: registryPDA,
          authorityNode: otherNodeEntryPDA,
          authority: otherNode.publicKey,
        })
        .signers([otherNode])
        .rpc();
      expect.fail("Expected error was not thrown");
    } catch (error) {
      // The error will be about the PDA derivation failing because we're using a node from a different registry
      expect(error.toString()).to.include("AnchorError caused by account: author");
    }

    // Verify the node's status didn't change
    const entry = await program.account.nodeEntry.fetch(nodeEntryPDA);
    expect(entry.active).to.equal(true);
  });

  it("Checks node registration includes active status", async () => {
    const info = await program.methods
      .checkNode(nodeAccount.publicKey)
      .accounts({
        entry: nodeEntryPDA,
        registry: registryPDA,
      })
      .view();

    expect(info.domain).to.equal("example.com");
    expect(info.active).to.equal(true);
  });

  it("Removes node from registry", async () => {
    await program.methods
      .removeNodeFromRegistry(nodeAccount.publicKey)
      .accounts({
        entry: nodeEntryPDA,
        registry: registryPDA,
        authority: authority.publicKey,
      })
      .rpc();

    // Verify the account is closed
    const accountInfo = await provider.connection.getAccountInfo(nodeEntryPDA);
    expect(accountInfo).to.be.null;
  });
});
