use anchor_lang::prelude::*;

// the "base" registry account upon which all registry entry account addresses are derived
#[account]
pub struct Registry {
    // the account that created this registry
    pub authority: Pubkey,
    // the name of the registry (max 32 chars)
    pub name: String
}

impl Registry {
    pub const SIZE: usize = 8 + 32 + 32;
}

// a PDA derived from the address of the account to add and the base registry
// defined in create_registry::registry
//
// Checking if an account address X is registred in registry Y
// involves checking if a registryEntry exists whose address is derived from X and Y
#[account]
pub struct ClientEntry {
    pub parent: Pubkey,     // Registry this entry belongs to
    pub registred: Pubkey,  // The registered client account
    pub until: i64,         // Unix timestamp when this registration expires
    pub limit: u32,         // Request limit for this client
}

#[account]
pub struct NodeEntry {
    pub parent: Pubkey,     // Registry this entry belongs to
    pub registred: Pubkey,  // The registered node account
    pub domain: String,     // Domain name (max 253 chars)
    pub online: i32,       // Online status counter
    pub active: bool,      // Active status
}

impl ClientEntry {
    pub const SIZE: usize = 8 + 32 + 32 + 8 + 4;
}

impl NodeEntry {
    // 8 (discriminator) + 32 (parent) + 32 (registered) + 4 + 253 (domain string) + 4 (online) + 1 (active)
    pub const SIZE: usize = 8 + 32 + 32 + 4 + 253 + 4 + 1;
}

#[derive(AnchorSerialize, AnchorDeserialize, Clone, Debug)]
pub struct ClientInfo {
    pub until: i64,
    pub limit: u32,
}

#[derive(AnchorSerialize, AnchorDeserialize, Clone, Debug)]
pub struct NodeInfo {
    pub domain: String,
    pub active: bool,
}