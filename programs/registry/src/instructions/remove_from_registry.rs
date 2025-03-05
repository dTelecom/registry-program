use anchor_lang::prelude::*;
use crate::instructions::account::{Registry, ClientEntry, NodeEntry};

pub fn remove_client_from_registry(
    _ctx: Context<RemoveClientFromRegistry>,
    _account_to_delete: Pubkey,
) -> Result<()> {
    Ok(())
}

pub fn remove_node_from_registry(
    _ctx: Context<RemoveNodeFromRegistry>,
    _account_to_delete: Pubkey,
) -> Result<()> {
    Ok(())
}

#[derive(Accounts)]
#[instruction(account_to_delete: Pubkey)]
pub struct RemoveClientFromRegistry<'info> {
    #[account(
        mut,
        close = authority,
        seeds=[account_to_delete.as_ref(), registry.key().as_ref()],
        bump,
    )]
    entry: Account<'info, ClientEntry>,

    registry: Account<'info, Registry>,

    #[account(
        mut,
        constraint = registry.authority == authority.key()
    )]
    authority: Signer<'info>,
}

#[derive(Accounts)]
#[instruction(account_to_delete: Pubkey)]
pub struct RemoveNodeFromRegistry<'info> {
    #[account(
        mut,
        close = authority,
        seeds=[account_to_delete.as_ref(), registry.key().as_ref()],
        bump,
    )]
    entry: Account<'info, NodeEntry>,

    registry: Account<'info, Registry>,

    #[account(
        mut,
        constraint = registry.authority == authority.key()
    )]
    authority: Signer<'info>,
}