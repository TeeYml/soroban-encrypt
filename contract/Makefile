.PHONY: build test release

build:
	cargo build --target wasm32-unknown-unknown --release

test:
	cargo test

release:
	cargo build --target wasm32-unknown-unknown --release
	wasm-opt -Oz \
		target/wasm32-unknown-unknown/release/soroban_encrypt_contract.wasm \
		-o target/soroban_encrypt_contract_optimised.wasm

clean:
	cargo clean
