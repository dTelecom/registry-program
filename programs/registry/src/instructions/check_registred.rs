use anchor_lang::prelude::*;
use crate::instructions::account::{Registry, ClientEntry, NodeEntry, ClientInfo, NodeInfo};

#[derive(AnchorSerialize, AnchorDeserialize, Clone, Debug)]
pub enum RegistryInfo {
    Client(ClientInfo),
    Node(NodeInfo),
}

pub fn check_client(
    ctx: Context<CheckClient>,
    _account_to_check: Pubkey,
) -> Result<ClientInfo> {
    Ok(ClientInfo {
        until: ctx.accounts.entry.until,
        limit: ctx.accounts.entry.limit,
    })
}

pub fn check_node(
    ctx: Context<CheckNode>,
    _account_to_check: Pubkey,
) -> Result<NodeInfo> {
    Ok(NodeInfo {
        domain: ctx.accounts.entry.domain.clone(),
        active: ctx.accounts.entry.active,
    })
}

#[derive(Accounts)]
#[instruction(account_to_check: Pubkey)]
pub struct CheckClient<'info> {
    #[account(
        seeds=[account_to_check.as_ref(), registry.key().as_ref()],
        bump,
    )]
    entry: Account<'info, ClientEntry>,
    registry: Account<'info, Registry>,
}

#[derive(Accounts)]
#[instruction(account_to_check: Pubkey)]
pub struct CheckNode<'info> {
    #[account(
        seeds=[account_to_check.as_ref(), registry.key().as_ref()],
        bump,
    )]
    entry: Account<'info, NodeEntry>,
    registry: Account<'info, Registry>,
}