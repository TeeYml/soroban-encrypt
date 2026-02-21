#![cfg(test)]

use super::{AllowlistContract, AllowlistContractClient};
use soroban_sdk::{testutils::Address as _, Address, Bytes, Env};

#[test]
fn test_allowlist_flow() {
    let env = Env::default();
    env.mock_all_auths();

    let contract_id = env.register_contract(None, AllowlistContract);
    let client = AllowlistContractClient::new(&env, &contract_id);

    let admin = Address::generate(&env);
    let name = Bytes::from_slice(&env, b"my_secure_allowlist");

    // Initialize the contract
    client.init(&admin, &name);

    let user_a = Address::generate(&env);
    let user_b = Address::generate(&env);

    // Initial state: neither user is allowed
    assert!(!client.is_allowed(&user_a));
    assert!(!client.is_allowed(&user_b));

    // Admin adds user_a
    client.add(&user_a);
    assert!(client.is_allowed(&user_a));
    assert!(!client.is_allowed(&user_b));

    // approve should succeed for user_a (no panic)
    client.approve(&Bytes::from_slice(&env, b"object_123"), &user_a);

    // Admin removes user_a
    client.remove(&user_a);
    assert!(!client.is_allowed(&user_a));
}

#[test]
#[should_panic(expected = "ENoAccess")]
fn test_approve_denied() {
    let env = Env::default();
    env.mock_all_auths();

    let contract_id = env.register_contract(None, AllowlistContract);
    let client = AllowlistContractClient::new(&env, &contract_id);

    let admin = Address::generate(&env);
    let name = Bytes::from_slice(&env, b"my_secure_allowlist");

    client.init(&admin, &name);

    let user = Address::generate(&env);

    // user is not added, so approve should panic
    client.approve(&Bytes::from_slice(&env, b"object_123"), &user);
}
