---
name: FEM and wFEM plan
description: FEM and wFEM timing, approved native FEM allocation, and production mainnet constraints.
---

The user plans to create wFEM as a separate wrapped token after the FEM chain is deployed. The user approved this 25-billion native FEM allocation:

- Treasury: 8 billion (32%)
- Early users: 300 million (1.2%)
- LP/DEX: 900 million (3.6%)
- Validator rewards: 4.5 billion (18%)
- Ecosystem/developers: 3 billion (12%)
- Community incentives: 2 billion (8%)
- Team/founders: 1.5 billion (6%)
- Partnerships: 1 billion (4%)
- Future development/reserve: 3.8 billion (15.2%)

The user supplied destination addresses for the nine allocation categories and four validator EVM addresses, plus a matching BLS public key per validator. The user intends FEM to run as a production mainnet. The existing four-node Replit setup and its signer keys remain simulation-only; do not reuse them for production. Public P2P bootnode multiaddresses are still needed.

**Why:** The user approved the proposed amounts on 2026-10-26.

**How to apply:** Keep allocations within the 25-billion total. Use a separate production genesis so the Replit devnet stays intact. Before launching, add public P2P bootnode multiaddresses and have operators retain their matching private keys on their own hosts. The 4.5B validator-rewards allocation is a funded wallet, not an automatic payout schedule. Do not use the Replit simulation signers for production. Do not assume a wFEM supply or bridge design.

The user specified mainnet chain ID `23124` and a separate Railway testnet with chain ID `23126`, four validators, and a 2-second block interval. Keep the testnet genesis, validator keys, and balances separate from mainnet and the legacy Replit simulation.

**Why:** The user specified the testnet identity and validator count while preparing a separate Railway testnet.

**How to apply:** Keep the 23126 Railway testnet separate and use four fresh test-only validators with a 2-second block interval. Do not reuse mainnet or Replit simulation keys.
