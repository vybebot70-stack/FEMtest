# FEM Railway Testnet

**Status:** configuration and instructions only. No Railway services have been
created or deployed.

## Network separation

| Network | Chain ID | Intended use |
|---|---:|---|
| FEM mainnet draft | 23124 | Reserved for mainnet |
| FEM Railway testnet | 23126 | Disposable integration testing |
| Existing Replit devnet | 23124 | Legacy local simulation; keep its genesis and data separate |

Do not use mainnet allocation addresses or production validator keys on the
testnet. Testnet balances have no mainnet value. The testnet must use a valid
IBFT/BLS genesis generated from four fresh test-only validator key sets.

## Before creating Railway services

1. Confirm the Railway project can run four continuously running services,
   attach one persistent volume to each service, and allow private TCP traffic
   between those services.
2. Create four separate validator services from this repository. Attach a
   different persistent volume to each at `/data`.
3. For each service, temporarily set its start command to:

   ```sh
   ./polygon-edge secrets init --data-dir /data
   ```

   Run this once per empty volume. Do not copy the tracked Replit validator
   keys into Railway.
4. Obtain only the public validator address, BLS public key, and libp2p node ID
   from each service:

   ```sh
   ./polygon-edge secrets output --data-dir /data
   ```

   Keep all private key files inside their respective persistent volumes.
5. Generate a shared genesis file with those four public validator entries,
   chain ID `23126`, IBFT PoA/BLS, and a 2-second block time. Use the Polygon
   Edge version built by this repository and verify the generated validator
   set against the four public outputs before starting any validator.
6. Configure one advertised Railway private DNS name per validator. Use those
   names and the corresponding libp2p node IDs to form bootnodes such as:

   ```text
   /dns4/<validator-private-hostname>/tcp/30303/p2p/<libp2p-peer-id>
   ```

   Set the comma-separated list of reachable multiaddresses as
   `FEM_BOOTNODES` on each validator. The server's `--bootnodes` option
   overrides any loopback bootnodes in the genesis file.

## Per-service environment

Set these variables on each Railway validator service:

| Variable | Example or default | Purpose |
|---|---|---|
| `FEM_VALIDATOR_NAME` | `validator-1` | Human-readable log label |
| `FEM_CHAIN_FILE` | `./genesis-railway-testnet.json` | Shared testnet genesis in the image |
| `FEM_DATA_DIR` | `/data` | Persistent volume mount |
| `FEM_P2P_DNS_NAME` | `/dns4/<private-hostname>` | Address this node advertises |
| `FEM_BOOTNODES` | Comma-separated P2P multiaddresses | Private peer discovery |
| `FEM_LIBP2P_ADDRESS` | `0.0.0.0:30303` | P2P listen address |
| `FEM_JSONRPC_ADDRESS` | `127.0.0.1:8545` | Private JSON-RPC bind |
| `FEM_GRPC_ADDRESS` | `127.0.0.1:10000` | Private gRPC bind |

The start script fails closed unless the genesis uses chain ID `23126`, the
volume is writable, the three required local key files exist, and P2P DNS and
bootnodes are configured. It enforces the 1 gwei gas-price floor. JSON-RPC and
gRPC bind to loopback by default; do not make validator RPC ports public.

The `Dockerfile` exposes only P2P port `30303/tcp` as image metadata. Railway
network access still needs to be configured for validator-to-validator
connectivity. Do not add a public JSON-RPC domain to validator services. If
external wallet or explorer access is needed, create a separate RPC service
with explicit authentication, rate limits, and monitoring.

## Readiness checks

Before considering the testnet online:

- Confirm all four services use different volumes and distinct test-only keys.
- Confirm every genesis validator public key matches exactly one service's
  private key.
- Confirm every bootnode hostname resolves from the other services and peers
  connect in both directions.
- Query `eth_chainId` and verify `0x5a56` (decimal `23126`) on every RPC node.
- Verify blocks are produced at the configured interval and all four
  validators remain in the IBFT set.
- Verify transactions below `1000000000` wei are rejected and eligible gas
  fees split equally between the FEM Treasury and proposer, with odd wei
  remainder assigned to the proposer.
- Confirm London/EIP-1559 remains disabled and do not infer base-fee burn
  behavior from the genesis `baseFee` field.
- Test restart recovery from each persistent volume and record a tested backup
  and restore procedure before relying on the network.

## Important limits

- This repository does not yet contain a bootable Railway testnet genesis.
  Fresh validator public keys and Railway private P2P addresses are required
  before one can be generated.
- This guide does not deploy, publish, or expose a service.
- Railway plan/network limits must be checked in the Railway account before
  provisioning. The repository cannot verify the selected plan's networking
  or storage allowance.
- The fee-policy code recognizes chain IDs `23124` and `23126`; London remains
  unchanged and disabled for the current FEM genesis drafts.
