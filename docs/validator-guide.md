# FEM Validator Operations

## Status

This document is not a production launch authorization. The current Replit
workflow is a local four-process simulation. The separate Railway testnet is
planned for chain ID `23126`; see `docs/railway-testnet.md`. Mainnet remains a
draft on chain ID `23124`.

Do not use the former generic launch commands from this guide: they used a
loopback bootnode, exposed JSON-RPC and gRPC broadly, and included placeholder
release and genesis information.

## Key requirements

- Use a distinct persistent data volume for each validator.
- Generate fresh keys on the host or volume that will operate the validator.
- Keep validator, BLS, and libp2p private keys out of source control, container
  build contexts, chat, and logs.
- The validator public key in genesis must match that node's private key.
- Do not reuse the Replit simulation validator keys for Railway testnet or
  mainnet. Existing Replit validator key files were present in repository
  history; treat them as exposed and rotate before using those identities
  outside the local simulation. Removing them from future commits does not
  remove them from Git history.
- Use `/dns4/<hostname>/tcp/<port>/p2p/<peer-id>` bootnode multiaddresses for
  private DNS peers. Never use `127.0.0.1` for validators on separate hosts.

## Network access

- Allow only the P2P listener needed for peer communication.
- Keep validator JSON-RPC and gRPC bound to loopback or private service
  networking. Do not publish these ports directly.
- If public JSON-RPC is required, place it on a separate RPC service with
  authentication or network restrictions, rate limits, request limits, and
  monitoring.
- Test peer discovery, quorum, process recovery, storage recovery, and
  coordinated upgrades on the testnet before considering mainnet.

See `docs/railway-testnet.md` for the Railway testnet configuration and
`docs/mainnet-launch-checklist.md` for mainnet readiness gates.
