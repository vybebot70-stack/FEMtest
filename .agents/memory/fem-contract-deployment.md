---
name: FEM contract deployment
description: User-stated policy for who may deploy EVM contracts on FEM.
---

Contract deployment on FEM should be permissionless, like Ethereum: any funded address can deploy a contract by submitting a signed transaction and paying gas. Do not enable a contract-deployer allowlist unless the user changes this policy.

**Why:** The user selected deployment by any funded address and said it should work like Ethereum.

**How to apply:** Preserve permissionless contract deployment in FEM genesis and runtime settings.
