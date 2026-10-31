#!/usr/bin/env bash
# Fem Chain — 4-validator IBFT devnet
# Node 1's JSON-RPC is exposed on :8080 (publicly reachable via Replit).
# Gas price floor: 1 gwei (1000000000 wei)

set -euo pipefail

BINARY="./polygon-edge"
GENESIS="./genesis.json"

# Rebuild if the binary is missing or older than any Go source/module file.
if [ ! -x "$BINARY" ] ||
  [ go.mod -nt "$BINARY" ] ||
  [ go.sum -nt "$BINARY" ] ||
  find . -path './.git' -prune -o -name '*.go' -newer "$BINARY" -print -quit | grep -q .; then
  echo "Polygon Edge source changed; building the node binary..."
  go build -o "$BINARY" \
    -ldflags="-X 'github.com/0xPolygon/polygon-edge/versioning.Version=v1.0.0' \
              -X 'github.com/0xPolygon/polygon-edge/versioning.Branch=main'" main.go
fi

# Never kill unrelated Polygon Edge processes. Refuse to start duplicate local
# validators and let the operator stop this project's nodes explicitly.
for cmdline_path in /proc/[0-9]*/cmdline; do
  [ -r "$cmdline_path" ] || continue
  cmdline="$(tr '\0' ' ' < "$cmdline_path" 2>/dev/null || true)"
  case "$cmdline" in
    *"./polygon-edge server --data-dir ./fem-chain-"*)
      echo "A Fem Chain validator is already running. Stop this devnet before starting it again." >&2
      exit 1
      ;;
  esac
done

echo "=============================="
echo "  Starting Fem Chain Devnet"
echo "  Chain ID : 23124"
echo "  Token    : FEM"
echo "  Consensus: IBFT (4 validators)"
echo "=============================="

NODE_PIDS=()
cleanup() {
  local pid
  for pid in "${NODE_PIDS[@]}"; do
    kill -TERM "$pid" 2>/dev/null || true
  done
}
trap cleanup EXIT
trap 'exit 0' INT TERM

# Node 1 — JSON-RPC on :8080 (public)
echo "[1/4] Starting validator 1 (JSON-RPC :8080)..."
"$BINARY" server \
  --data-dir ./fem-chain-1 \
  --chain "$GENESIS" \
  --grpc-address 127.0.0.1:10000 \
  --libp2p :30301 \
  --jsonrpc 0.0.0.0:8080 \
  --seal \
  --price-limit 1000000000 \
  --log-level INFO &
NODE1_PID=$!
NODE_PIDS+=("$NODE1_PID")
echo "  PID: $NODE1_PID"

sleep 3

# Node 2
echo "[2/4] Starting validator 2..."
"$BINARY" server \
  --data-dir ./fem-chain-2 \
  --chain "$GENESIS" \
  --grpc-address 127.0.0.1:10002 \
  --libp2p :30302 \
  --jsonrpc 127.0.0.1:10202 \
  --seal \
  --price-limit 1000000000 \
  --log-level WARN &
NODE2_PID=$!
NODE_PIDS+=("$NODE2_PID")
echo "  PID: $NODE2_PID"

# Node 3
echo "[3/4] Starting validator 3..."
"$BINARY" server \
  --data-dir ./fem-chain-3 \
  --chain "$GENESIS" \
  --grpc-address 127.0.0.1:10003 \
  --libp2p :30303 \
  --jsonrpc 127.0.0.1:10203 \
  --seal \
  --price-limit 1000000000 \
  --log-level WARN &
NODE3_PID=$!
NODE_PIDS+=("$NODE3_PID")
echo "  PID: $NODE3_PID"

# Node 4
echo "[4/4] Starting validator 4..."
"$BINARY" server \
  --data-dir ./fem-chain-4 \
  --chain "$GENESIS" \
  --grpc-address 127.0.0.1:10004 \
  --libp2p :30304 \
  --jsonrpc 127.0.0.1:10204 \
  --seal \
  --price-limit 1000000000 \
  --log-level WARN &
NODE4_PID=$!
NODE_PIDS+=("$NODE4_PID")
echo "  PID: $NODE4_PID"

echo ""
echo "=============================="
echo "  Fem Chain validator processes started"
echo ""
echo "  JSON-RPC : http://0.0.0.0:8080"
echo "  Chain ID : 23124"
echo "  Symbol   : FEM"
echo ""
echo "  MetaMask settings:"
echo "    Network name : Fem"
echo "    RPC URL      : https://<your-replit-url>-8080.replit.dev"
echo "    Chain ID     : 23124"
echo "    Currency sym : FEM"
echo "=============================="

# Keep alive — wait for node 1
wait "$NODE1_PID"
