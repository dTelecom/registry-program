use anchor_lang::prelude::*;
use crate::instructions::account;

pub fn init_registry(
    ctx: Context<InitRegistry>,
    name: String // user-provided name of the registry
) -> Result<()> {
    let registry = &mut ctx.accounts.registry;
    // add the account that created this registry as the registry admin
    registry.authority = *ctx.accounts.authority.signer_key().unwrap();
    registry.name = name;
    Ok(())
}

#[derive(Accounts)]
#[instruction(name: String)]
pub struct InitRegistry<'info> {
    #[account(
        init,
        payer=authority,
        space=account::Registry::SIZE,
        // use the authority and registry name as seeds for the PDA
        seeds=[authority.key().as_ref(), name.as_bytes()],
        bump
    )]
    registry: Account<'info, account::Registry>,

    #[account(mut)]
    authority: Signer<'info>, // the account sending (and signing) this transaction
    system_program: Program<'info, System>,
}