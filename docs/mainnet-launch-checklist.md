# FEM Mainnet Readiness Checklist

**Status: not launch-ready.** This checklist describes gates, not a completed
deployment. No mainnet has been deployed or published.

## Chain identity and genesis

- [x] Keep mainnet chain ID at `23124`.
- [x] Reserve `23126` for the separate Railway testnet.
- [ ] Reconcile the 25B FEM allocation and recipient addresses against the
  approved tokenomics in `docs/token-distribution.md`.
- [ ] Generate and independently verify final genesis extraData from the
  production validator public keys.
- [ ] Configure and test public P2P bootnodes. The current mainnet template has
  an empty `bootnodes` list.
- [ ] Confirm London/EIP-1559 fork settings and burn behavior from a complete,
  tested design before changing them. Current drafts leave London disabled.
- [ ] Generate and publish a reproducible genesis hash only after the above
  checks pass.

## Validators and keys

- [ ] Use four independent production validator operators and hosts.
- [ ] Generate fresh production BLS, validator, and P2P keys on operator-owned
  hosts. Do not reuse Replit simulation keys or Railway testnet keys.
- [ ] Confirm each private key matches exactly one public validator entry.
- [ ] Store backups securely and test recovery without placing private keys in
  source control, container images, chat, or logs.
- [ ] Verify IBFT quorum, validator replacement, network partitions, and
  coordinated upgrades on a separate testnet first.

## Fees, RPC, and application compatibility

- [ ] Run the 50/50 Treasury/proposer fee policy tests successfully on every
  validator build before rollout. Odd wei remainder goes to the proposer.
- [ ] Confirm the 1 gwei transaction-price floor on all validators.
- [ ] Keep validator JSON-RPC and gRPC private. If public RPC is required,
  deploy a separate protected RPC layer with rate limits and monitoring.
- [ ] Test wallet signing, chain ID display, explorer indexing, contract
  deployment, and the target DEX on the testnet.
- [ ] Confirm no London/EIP-1559 or base-fee/burn behavior changes are included
  without a separate tested configuration.

## Build, operations, and release

- [ ] Resolve dependency firewall/security findings through compatible
  dependency updates; do not bypass the package firewall.
- [ ] Run successful unit tests, integration tests, and reproducible builds.
- [ ] Upgrade all validators consistently at an agreed protocol activation
  point; do not mix fee-policy binaries.
- [ ] Configure alerting for block production, peer count, validator quorum,
  RPC errors, disk capacity, and process restarts.
- [ ] Test persistent-storage snapshots, backup restoration, and documented
  recovery procedures.
- [ ] Complete an external security review and testnet soak period.
- [ ] Obtain explicit operator approval for the final genesis, validator set,
  RPC exposure, and mainnet launch.

## Current project state

- The Replit workflow is a local four-process simulation, not four independent
  production validators.
- The Railway runner is being prepared for chain ID `23126`; no Railway service
  has been deployed.
- The mainnet draft remains chain ID `23124`, with a 25B allocation template
  and no public bootnodes.
- The fee-policy source is present, but dependency downloads have blocked
  successful Go test execution in this environment.
