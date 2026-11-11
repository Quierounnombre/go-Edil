#!/usr/bin/env bash
# Summarizes fixloop metrics.
#
#   cmd/fix/metrics.sh [metrics.tsv]
#
# Round 1 counts are what the agent wrote before any feedback: the rules to
# put in its prompt. Rolled back counts show where it fights the gate.
set -euo pipefail

METRICS=${1:-${METRICS:-${XDG_STATE_HOME:-$HOME/.local/state}/edil/metrics.tsv}}
TOP=10
[ -f "$METRICS" ] || { echo "no metrics at $METRICS" >&2; exit 1; }

echo "first-pass findings by rule (round 1):"
awk -F'\t' '$3 == 1 { total[$4] += $5 } END { for (rule in total) print total[rule], rule }' "$METRICS" |
	sort -rn | head -"$TOP"

echo
echo "rounds lost to rollbacks, by rule:"
awk -F'\t' '$6 ~ /^rolled_back/ { total[$4] += 1 } END { for (rule in total) print total[rule], rule }' "$METRICS" |
	sort -rn | head -"$TOP"

echo
echo "rounds by outcome:"
awk -F'\t' '{ key = $1 FS $2 FS $3 FS $6; if (!(key in seen)) { seen[key] = 1; outcome[$6] += 1 } }
	END { for (name in outcome) print outcome[name], name }' "$METRICS" | sort -rn
