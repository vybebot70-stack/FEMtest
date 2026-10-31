#!/usr/bin/env bash
# Start exactly one Railway FEM testnet validator.
#
# Each Railway validator service must have its own persistent volume mounted
# at /data, fresh test-only Polygon Edge secrets in that volume, and a genesis
# file whose IBFT validator set matches those secrets.

set -euo pipefail
umask 077

BINARY="./polygon-edge"
GENESIS="${FEM_CHAIN_FILE:-./genesis-railway-testnet.json}"
DATA_DIR="${FEM_DATA_DIR:-/data}"
VALIDATOR_NAME="${FEM_VALIDATOR_NAME:-unnamed-validator}"
P2P_ADDRESS="${FEM_LIBP2P_ADDRESS:-0.0.0.0:30303}"
P2P_DNS_NAME="${FEM_P2P_DNS_NAME:-}"
BOOTNODES="${FEM_BOOTNODES:-}"
GRPC_ADDRESS="${FEM_GRPC_ADDRESS:-127.0.0.1:10000}"
JSONRPC_ADDRESS="${FEM_JSONRPC_ADDRESS:-127.0.0.1:8545}"

if [ ! -x "$BINARY" ]; then
  echo "Polygon Edge binary is missing; build the Railway image first." >&2
  exit 1
fi

if [ ! -f "$GENESIS" ]; then
  echo "FEM testnet genesis file not found: $GENESIS" >&2
  echo "Generate it from fresh test-only validator keys with chain ID 23126." >&2
  exit 1
fi

if ! grep -Eq '"chainID"[[:space:]]*:[[:space:]]*23126([[:space:],}]|$)' "$GENESIS"; then
  echo "Refusing to start: the configured genesis must use chain ID 23126." >&2
  exit 1
fi

if ! grep -Eq '"blockTime"[[:space:]]*:[[:space:]]*2000000000([[:space:],}]|$)' "$GENESIS"; then
  echo "Refusing to start: the configured IBFT block period must be 2 seconds." >&2
  exit 1
fi

if ! grep -Eq '"ibft"[[:space:]]*:[[:space:]]*\{' "$GENESIS" ||
  ! grep -Eq '"validator_type"[[:space:]]*:[[:space:]]*"bls"' "$GENESIS"; then
  echo "Refusing to start: the genesis must configure IBFT with BLS validators." >&2
  exit 1
fi

if [ ! -d "$DATA_DIR" ] || [ ! -w "$DATA_DIR" ]; then
  echo "Validator data directory must exist and be writable: $DATA_DIR" >&2
  echo "Attach a persistent Railway volume at this path." >&2
  exit 1
fi

for key_file in \
  "$DATA_DIR/consensus/validator.key" \
  "$DATA_DIR/consensus/validator-bls.key" \
  "$DATA_DIR/libp2p/libp2p.key"; do
  if [ ! -s "$key_file" ]; then
    echo "Required validator key is missing: $key_file" >&2
    echo "Initialize fresh test-only keys in this validator's persistent volume." >&2
    exit 1
  fi
done

if [ -z "$P2P_DNS_NAME" ]; then
  echo "FEM_P2P_DNS_NAME must identify this validator on the Railway private network." >&2
  exit 1
fi

if [ -z "$BOOTNODES" ]; then
  echo "FEM_BOOTNODES must contain at least one reachable P2P multiaddress." >&2
  exit 1
fi

echo "Starting $VALIDATOR_NAME on FEM testnet chain ID 23126."
echo "JSON-RPC and gRPC are bound to loopback by default."

exec "$BINARY" server \
  --data-dir "$DATA_DIR" \
  --chain "$GENESIS" \
  --grpc-address "$GRPC_ADDRESS" \
  --libp2p "$P2P_ADDRESS" \
  --dns "$P2P_DNS_NAME" \
  --bootnodes "$BOOTNODES" \
  --jsonrpc "$JSONRPC_ADDRESS" \
  --seal \
  --price-limit 1000000000 \
  --log-level INFO
