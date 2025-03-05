use anchor_lang::prelude::*;
use crate::instructions::account;
use crate::ErrorCode;

pub fn add_node_to_registry(
    ctx: Context<AddNodeToRegistry>,
    account_to_add: Pubkey,
    domain: String,
) -> Result<()> {
    require!(domain.len() <= 253, ErrorCode::DomainTooLong);
    let entry = &mut ctx.accounts.entry;
    entry.parent = *ctx.accounts.registry.to_account_info().key;
    entry.registred = account_to_add;
    entry.domain = domain;
    Ok(())
}

#[derive(Accounts)]
#[instruction(account_to_add: Pubkey)]
pub struct AddNodeToRegistry<'info> {
    #[account(
        init,
        payer=authority,
        space=account::NodeEntry::SIZE,
        seeds=[account_to_add.as_ref(), registry.key().as_ref()],
        bump
    )]
    entry: Account<'info, account::NodeEntry>,

    registry: Account<'info, account::Registry>,

    #[account(
        mut,
        constraint = registry.authority == authority.key()
    )]
    authority: Signer<'info>,

    system_program: Program<'info, System>,
} 