# FEM Supply and Fee Policy

## Keep the networks separate

| Network | Chain ID | Genesis supply |
|---|---:|---|
| Legacy Replit devnet | 23124 | 4,000,000 FEM in its four simulation allocations |
| Mainnet draft | 23124 | 25,000,000,000 FEM across the approved allocation categories |
| Railway testnet | 23126 | Disposable test balances; not the mainnet supply |

The legacy devnet allocation is not the mainnet tokenomics. The Railway
testnet must use a different chain ID and test-only validator keys. No testnet
genesis has been finalized yet.

## Approved 25B mainnet allocation draft

| Category | FEM | Share |
|---|---:|---:|
| Treasury | 8,000,000,000 | 32% |
| Early users | 300,000,000 | 1.2% |
| LP/DEX | 900,000,000 | 3.6% |
| Validator rewards | 4,500,000,000 | 18% |
| Ecosystem/developers | 3,000,000,000 | 12% |
| Community incentives | 2,000,000,000 | 8% |
| Team/founders | 1,500,000,000 | 6% |
| Partnerships | 1,000,000,000 | 4% |
| Future development/reserve | 3,800,000,000 | 15.2% |
| **Total** | **25,000,000,000** | **100%** |

The addresses are present in `genesis-mainnet-template.json`, but that file is
a draft: its bootnodes are empty and it is not a launch-ready network
configuration. Reconcile every address and amount before producing a final
mainnet genesis. The 4.5B validator-rewards allocation is a funded allocation,
not an automatic block-reward schedule.

## Gas fees and issuance

- Eligible FEM gas fees are split 50/50 between the configured FEM Treasury
  (`0x38B9bbECA3bB1437edeBA2c6298024e1c1C8E997`) and the block proposer.
- For an odd wei amount, Treasury receives the floor of half and the proposer
  receives the remaining wei.
- The 1 gwei minimum is enforced in the FEM transaction path and server
  transaction-pool configuration for chain IDs 23124 and 23126.
- London/EIP-1559 remains disabled in the current genesis drafts. The existing
  base-fee/burn path has not been changed or activated by the fee split.
- IBFT does not automatically mint the 4.5B reward allocation. A future reward
  schedule would need a separately designed and tested protocol change.
- The user plans wFEM as a separate wrapped asset after deployment. Its supply,
  bridge, custody, and security design are not set here.

## Before using token balances

1. Confirm the final 25B addresses and amounts sum exactly to 25B FEM.
2. Confirm each mainnet recipient controls the intended address and has an
   appropriate custody process; use a multisig for treasury funds.
3. Keep the testnet balances and validator keys independent from mainnet.
4. Verify the final genesis allocation and chain ID from a reproducible
   genesis artifact before any network launch.
