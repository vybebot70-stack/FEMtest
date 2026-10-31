#!/usr/bin/env bash
set -euo pipefail

stopped=0
for cmdline_path in /proc/[0-9]*/cmdline; do
  [ -r "$cmdline_path" ] || continue
  pid="${cmdline_path#/proc/}"
  pid="${pid%/cmdline}"
  cmdline="$(tr '\0' ' ' < "$cmdline_path" 2>/dev/null || true)"

  case "$cmdline" in
    *"./polygon-edge server --data-dir ./fem-chain-1 "*|\
    *"./polygon-edge server --data-dir ./fem-chain-2 "*|\
    *"./polygon-edge server --data-dir ./fem-chain-3 "*|\
    *"./polygon-edge server --data-dir ./fem-chain-4 "*)
      echo "Sending SIGTERM to Fem Chain validator PID $pid"
      kill -TERM "$pid" 2>/dev/null || true
      stopped=1
      ;;
  esac
done

if [ "$stopped" -eq 0 ]; then
  echo "No validators for this Fem Chain devnet were found."
fi
