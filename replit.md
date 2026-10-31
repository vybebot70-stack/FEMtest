# Fem Chain — Polygon Edge Devnet

A 4-validator IBFT (Proof of Authority) Ethereum-compatible devnet built on [Polygon Edge](https://github.com/0xPolygon/polygon-edge) v1.0.0.

## Stack
- **Language**: Go 1.21
- **Framework**: Polygon Edge (fork of 0xPolygon/polygon-edge)
- **Consensus**: IBFT PoA (4 validators)

## How to run

The workflow **"Fem Chain Devnet"** starts the chain. If the `polygon-edge` binary
is missing (such as after a fresh import), `start-fem-chain.sh` builds it first.

To build the binary manually:
```bash
go build -o polygon-edge \
  -ldflags="-X 'github.com/0xPolygon/polygon-edge/versioning.Version=v1.0.0' \
            -X 'github.com/0xPolygon/polygon-edge/versioning.Branch=main'" main.go
```

To start the devnet manually:
```bash
bash start-fem-chain.sh
```

## Network details

| Property       | Value                          |
|----------------|-------------------------------|
| Network name   | Fem                            |
| Chain ID       | 23124                          |
| Currency       | FEM                            |
| JSON-RPC       | `http://0.0.0.0:8080`         |
| Block time     | ~2 seconds                     |
| Gas price      | **1 gwei** minimum (1,000,000,000 wei) |
| Gas limit/block| 10,000,000                     |
| Consensus      | IBFT PoA                       |

## MetaMask connection

1. Open MetaMask → Add Network → Add a network manually
2. Fill in:
   - **Network name**: Fem
   - **RPC URL**: `https://<your-replit-url>.replit.dev` (port 80 → maps to 8080)
   - **Chain ID**: `23124`
   - **Currency symbol**: `FEM`

## Pre-funded validator accounts

| Address | Balance |
|---------|---------|
| `0x246Bd0e7Da038524F3341a9dfB874173b00BB437` | 1,000,000 FEM |
| `0x939B951C6CFdB409BDD2f67Fe94b9Fe94597273c` | 1,000,000 FEM |
| `0x95C06222EFc6463B8882cEFad9bfb45025772FAD` | 1,000,000 FEM |
| `0xe2C84f1a57151b1a22e0fbfb92542881Ac95e8d7` | 1,000,000 FEM |

> Private keys for these accounts are stored in the `fem-chain-*/consensus/` directories.

## Validator nodes

| Node | gRPC    | libp2p | JSON-RPC |
|------|---------|--------|----------|
| 1    | :10000  | :30301 | :8080 (public) |
| 2    | :10002  | :30302 | :10202 |
| 3    | :10003  | :30303 | :10203 |
| 4    | :10004  | :30304 | :10204 |

## Gas fees

The validator nodes enforce a minimum transaction price of **1 gwei**. The genesis `baseFee` is also set to 1 gwei, and `burnContract` is null.

## FEM mainnet target

The production target is chain ID `23124`, name `fem`, EVM-compatible IBFT PoA, 2-second blocks, four independent validators at launch, and a 1 gwei minimum price. The approved 25-billion native FEM allocation is: 8 billion Treasury, 300 million early users, 900 million LP/DEX, 4.5 billion validator rewards, 3 billion ecosystem/developers, 2 billion community incentives, 1.5 billion team/founders, 1 billion partnerships, and 3.8 billion future development/reserve. Recipient addresses and four matching validator EVM/BLS public-key pairs have been supplied. `genesis-mainnet-template.json` records these allocations and validator identities, but its bootnode list is empty until production P2P multiaddresses are provided.

Contract deployment is permissionless, like Ethereum: any funded account can deploy a contract by submitting a signed transaction and paying gas. No contract-deployer allowlist is configured.

The user intends FEM to run as a production mainnet, but the current Replit chain remains a simulation. Keep its existing validator signers for simulation only; do not reuse them as production keys. Production IBFT validators must operate independently and retain their private keys on their own hosts. The 4.5-billion validator-rewards wallet is only a funded allocation; it does not automatically distribute rewards. The current `genesis.json` and running devnet are unchanged. `genesis-mainnet-template.json` is not launch-ready until bootnodes and production node operations are configured. wFEM is a separate wrapped ERC-20 planned after the FEM chain is deployed.

## Mainnet readiness

| Item | Status |
|------|--------|
| Replit devnet simulation | ✅ Live |
| Mainnet genesis template | ✅ 25B allocation + 4 validator public-key records |
| Public mainnet bootnodes | 🔲 Need validator P2P multiaddresses |
| Mainnet launched | ❌ Not launched |
| Gas price floor | ✅ 1 gwei |
| Public JSON-RPC (dev) | ✅ Port 80 |
| Replit production publish | 🔲 Click Publish |
| Chainlist.org submission | 🔲 See `docs/chainlist-submission/` |
| External validators | 🔲 See `docs/validator-guide.md` |
| Block explorer | 🔲 Not yet deployed |
| Token listing (CoinGecko) | 🔲 Needs explorer + community |

See `docs/mainnet-launch-checklist.md` for the full step-by-step plan.

## User preferences
