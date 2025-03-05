use anchor_lang::prelude::*;
use crate::instructions::account::{Registry, NodeEntry};
use crate::ErrorCode;

pub fn update_node_online(
    ctx: Context<UpdateNodeOnline>,
    _account_to_update: Pubkey,
    value: i32,
) -> Result<()> {
    require!(value >= 0, ErrorCode::InvalidOnlineValue);
    let entry = &mut ctx.accounts.entry;
    entry.online = value;
    Ok(())
}

#[derive(Accounts)]
#[instruction(account_to_update: Pubkey)]
pub struct UpdateNodeOnline<'info> {
    #[account(
        mut,
        seeds=[account_to_update.as_ref(), registry.key().as_ref()],
        bump,
        constraint = entry.registred == authority.key()
    )]
    entry: Account<'info, NodeEntry>,
    registry: Account<'info, Registry>,
    authority: Signer<'info>,
} 