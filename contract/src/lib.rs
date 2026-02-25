#![no_std]
use soroban_sdk::{contract, contractimpl, contracttype, Address, Bytes, Env};

#[contract]
pub struct AllowlistContract;

#[derive(Clone)]
#[contracttype]
pub enum DataKey {
    Admin,
    AllowedUser(Address),
    Name,
}

#[contractimpl]
impl AllowlistContract {
    /// Sets up the allowlist admin and name.
    pub fn init(env: Env, admin: Address, name: Bytes) {
        if env.storage().instance().has(&DataKey::Admin) {
            panic!("Already initialized");
        }
        env.storage().instance().set(&DataKey::Admin, &admin);
        env.storage().instance().set(&DataKey::Name, &name);
    }

    /// Adds a user (requires admin auth).
    pub fn add(env: Env, account: Address) {
        let admin: Address = env.storage().instance().get(&DataKey::Admin).expect("Not initialized");
        admin.require_auth();
        env.storage().persistent().set(&DataKey::AllowedUser(account), &true);
    }

    /// Removes a user (requires admin auth).
    pub fn remove(env: Env, account: Address) {
        let admin: Address = env.storage().instance().get(&DataKey::Admin).expect("Not initialized");
        admin.require_auth();
        env.storage().persistent().remove(&DataKey::AllowedUser(account));
    }

    /// Returns true if the account is allowed.
    pub fn is_allowed(env: Env, account: Address) -> bool {
        env.storage().persistent().has(&DataKey::AllowedUser(account))
    }

    /// Approve checks if a caller has access; panics with "ENoAccess" if not.
    pub fn approve(env: Env, _id: Bytes, caller: Address) {
        if !Self::is_allowed(env.clone(), caller) {
            panic!("ENoAccess");
        }
    }
}

#[cfg(test)]
mod test;
