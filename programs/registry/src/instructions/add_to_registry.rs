use anchor_lang::prelude::*;
use crate::instructions::account;

pub fn add_client_to_registry(
    ctx: Context<AddClientToRegistry>,
    account_to_add: Pubkey,
    until: i64,
    limit: u32,
) -> Result<()> {
    let entry = &mut ctx.accounts.entry;
    entry.parent = *ctx.accounts.registry.to_account_info().key;
    entry.registred = account_to_add;
    entry.until = until;
    entry.limit = limit;
    Ok(())
}

#[derive(Accounts)]
#[instruction(account_to_add: Pubkey)]
pub struct AddClientToRegistry<'info> {
    #[account(
        init,
        payer=authority,
        space=account::ClientEntry::SIZE,
        seeds=[account_to_add.as_ref(), registry.key().as_ref()],
        bump
    )]
    entry: Account<'info, account::ClientEntry>,

    registry: Account<'info, account::Registry>,

    #[account(
        mut,
        constraint = registry.authority == authority.key()
    )]
    authority: Signer<'info>,

    system_program: Program<'info, System>,
}