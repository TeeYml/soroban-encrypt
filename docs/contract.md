# Contract Deployment & Upgrade

## Deploying

```bash
./client-bin admin deploy \
  --wasm target/soroban_encrypt_contract_optimised.wasm \
  --secret SADMIN_SECRET
```

## Initialising

```bash
./client-bin admin init \
  --contract CCONTRACT_ID \
  --admin GADMIN_ADDRESS \
  --secret SADMIN_SECRET
```

## Upgrading

1. Compile and optimise the new WASM:
   ```bash
   make -C contract release
   ```
2. Upload and upgrade:
   ```bash
   ./client-bin admin deploy --wasm contract/target/soroban_encrypt_contract_optimised.wasm --secret SADMIN_SECRET
   ```
   The contract emits a `ContractUpgraded` event containing the new WASM hash.

## Multi-Admin

Add a second admin:
```bash
./client-bin admin add-admin --contract CXXX --address GNEW_ADMIN --secret SADMIN_SECRET
```
