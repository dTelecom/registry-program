#![allow(unexpected_cfgs)]

use anchor_lang::prelude::*;

pub mod instructions;
use instructions::{
    init_registry::*, 
    add_to_registry::*,
    add_node::*,
    delegate_node::*,
    undelegate_node::*,
    check_registred::*,
    remove_from_registry::*,
    update_node_online::*,
    update_node_active::*,
    account::*,
};

declare_id!("E2FcHsC9STeB6FEtxBKGAwMTX7cbfYMyjSHKs4QbBAmh");

#[error_code]
pub enum ErrorCode {
    #[msg("Online value must be non-negative")]
    InvalidOnlineValue,
    #[msg("Domain name must be 253 characters or less")]
    DomainTooLong,
}

#[program]
pub mod registry {
    use super::*;

    pub fn init_registry(ctx: Context<InitRegistry>, name: String) -> Result<()> {
        crate::instructions::init_registry::init_registry(ctx, name)
    }

    pub fn add_client_to_registry(
        ctx: Context<AddClientToRegistry>,
        account_to_add: Pubkey,
        until: i64,
        limit: u32,
    ) -> Result<()> {
        crate::instructions::add_to_registry::add_client_to_registry(ctx, account_to_add, until, limit)
    }

    pub fn add_node_to_registry(
        ctx: Context<AddNodeToRegistry>,
        account_to_add: Pubkey,
        domain: String,
    ) -> Result<()> {
        crate::instructions::add_node::add_node_to_registry(ctx, account_to_add, domain)
    }

    pub fn delegate_node_account(
        ctx: Context<DelegateNodeAccount>,
        account: Pubkey,
    ) -> Result<()> {
        crate::instructions::delegate_node::delegate_node_account(ctx, account)
    }

    pub fn undelegate_node_acount(
        ctx: Context<UndelegateNodeAccount>,
        account: Pubkey,
    ) -> Result<()> {
        crate::instructions::undelegate_node::undelegate_node_account(ctx, account)
    }

    pub fn check_client(
        ctx: Context<CheckClient>,
        account_to_check: Pubkey,
    ) -> Result<ClientInfo> {
        crate::instructions::check_registred::check_client(ctx, account_to_check)
    }

    pub fn check_node(
        ctx: Context<CheckNode>,
        account_to_check: Pubkey,
    ) -> Result<NodeInfo> {
        crate::instructions::check_registred::check_node(ctx, account_to_check)
    }

    pub fn remove_client_from_registry(
        ctx: Context<RemoveClientFromRegistry>,
        account_to_delete: Pubkey,
    ) -> Result<()> {
        crate::instructions::remove_from_registry::remove_client_from_registry(ctx, account_to_delete)
    }

    pub fn remove_node_from_registry(
        ctx: Context<RemoveNodeFromRegistry>,
        account_to_delete: Pubkey,
    ) -> Result<()> {
        crate::instructions::remove_from_registry::remove_node_from_registry(ctx, account_to_delete)
    }

    pub fn update_node_online(
        ctx: Context<UpdateNodeOnline>,
        account_to_update: Pubkey,
        online: i32,
    ) -> Result<()> {
        crate::instructions::update_node_online::update_node_online(ctx, account_to_update, online)
    }

    pub fn update_node_active(
        ctx: Context<UpdateNodeActive>,
        account_to_update: Pubkey,
        active: bool,
    ) -> Result<()> {
        crate::instructions::update_node_active::update_node_active(ctx, account_to_update, active)
    }
}