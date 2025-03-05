use anchor_lang::prelude::*;
use crate::instructions::account::*;

#[derive(Accounts)]
#[instruction(account_to_update: Pubkey)]
pub struct UpdateNodeActive<'info> {
    #[account(
        mut,
        seeds = [account_to_update.as_ref(), registry.key().as_ref()],
        bump,
        constraint = entry.parent == registry.key() // Must be in this registry
    )]
    pub entry: Account<'info, NodeEntry>,

    pub registry: Account<'info, Registry>,

    #[account(
        seeds = [authority.key().as_ref(), registry.key().as_ref()],
        bump,
        constraint = authority_node.parent == registry.key() && // Must be in this registry
                    authority_node.registred == authority.key() // Must be the authority's node
    )]
    pub authority_node: Account<'info, NodeEntry>,

    pub authority: Signer<'info>,
}

pub fn update_node_active(ctx: Context<UpdateNodeActive>, _account_to_update: Pubkey, active: bool) -> Result<()> {
    let entry = &mut ctx.accounts.entry;
    entry.active = active;
    Ok(())
} 