use anchor_lang::prelude::*;
use ephemeral_rollups_sdk::anchor::{commit};
use ephemeral_rollups_sdk::ephem::{commit_and_undelegate_accounts};
use crate::instructions::account;

pub fn undelegate_node_account(ctx: Context<UndelegateNodeAccount>, _account: Pubkey) -> Result<()> {

    commit_and_undelegate_accounts(
        &ctx.accounts.receiver,
        vec![&ctx.accounts.node.to_account_info()],
        &ctx.accounts.magic_context,
        &ctx.accounts.magic_program,
    )?;

    Ok(())
}

#[commit]
#[derive(Accounts)]
#[instruction(account: Pubkey)]
pub struct UndelegateNodeAccount<'info> {
    #[account(
        mut,
        seeds=[account.as_ref(), registry.key().as_ref()],
        bump
    )]
    pub node: Account<'info, account::NodeEntry>,

    pub registry: Account<'info, account::Registry>,

    #[account(mut)]
    pub receiver: Signer<'info>,
}