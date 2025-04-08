use anchor_lang::prelude::*;
use ephemeral_rollups_sdk::anchor::{delegate};
use ephemeral_rollups_sdk::cpi::DelegateConfig;
use crate::instructions::account;

pub fn delegate_node_account(ctx: Context<DelegateNodeAccount>, account: Pubkey) -> Result<()> {

    ctx.accounts.delegate_node(
        &ctx.accounts.authority,
        &[account.as_ref(), ctx.accounts.registry.key().as_ref()],
        DelegateConfig::default(),
    )?;

    Ok(())
}

#[delegate]
#[derive(Accounts)]
#[instruction(sender: Pubkey)]
pub struct DelegateNodeAccount<'info> {
    /// CHECK The pda to delegate
    #[account(mut, del)]
    pub node: AccountInfo<'info>,

    pub registry: Account<'info, account::Registry>,

    #[account(
        mut,
        constraint = registry.authority == authority.key()
    )]
    pub authority: Signer<'info>,

    
}