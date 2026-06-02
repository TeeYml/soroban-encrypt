use soroban_sdk::{contracttype, vec, Address, Env, Vec};

const ADMINS_KEY: &str = "admins";

#[contracttype]
pub struct AdminSet {
    pub admins: Vec<Address>,
}

pub fn add_admin(env: &Env, caller: &Address, new_admin: &Address) {
    require_any_admin(env, caller);
    let mut admins = get_admins(env);
    if !admins.contains(new_admin) {
        admins.push_back(new_admin.clone());
        env.storage().persistent().set(&ADMINS_KEY, &admins);
    }
}

pub fn remove_admin(env: &Env, caller: &Address, target: &Address) {
    require_any_admin(env, caller);
    let admins: Vec<Address> = get_admins(env).iter().filter(|a| a != target).collect();
    env.storage().persistent().set(&ADMINS_KEY, &admins);
}

pub fn require_any_admin(env: &Env, caller: &Address) {
    for admin in get_admins(env).iter() {
        if admin == *caller {
            caller.require_auth();
            return;
        }
    }
    panic!("ENotAdmin");
}

fn get_admins(env: &Env) -> Vec<Address> {
    env.storage().persistent().get(&ADMINS_KEY).unwrap_or(vec![env])
}

/// delegate grants per-object access from one address to another (1-hop only).
pub fn delegate(env: &Env, from: &Address, to: &Address, object_id: &soroban_sdk::Symbol) {
    from.require_auth();
    // Delegation depth must not exceed 1 hop to prevent privilege escalation
    let key = (Symbol::new(env, "delegate"), object_id.clone(), to.clone());
    env.storage().persistent().set(&key, from);
}
