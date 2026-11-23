#!/usr/bin/env bash
# Regression test for fixloop.sh: scripted agents walk every outcome.
#
#   EDIL=/path/to/edil cmd/fix/fixloop_test.sh
set -euo pipefail

EDIL=${EDIL:?set EDIL to the edil binary}
LOOP=$(readlink -f "$(dirname "$0")/fixloop.sh")
WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT
failures=0

BAD_MAIN='package main\n\nconst maxRetries = 3\n\nfunc main() {\n}\n'
CLEAN_MAIN='package main\n\nfunc main() {\n}\n'
DEFINES='package main\n\nconst MAX_RETRIES = 3\n'
FIX="printf '$CLEAN_MAIN' >main.go; printf '$DEFINES' >defines.go"

# repo MAIN_SOURCE: a fresh git repo with one commit.
repo() {
	rm -rf "$WORK/repo" && mkdir "$WORK/repo" && cd "$WORK/repo"
	git init -q && git config user.email test@edil && git config user.name test
	printf 'module demo\n\ngo 1.26\n' >go.mod
	printf "$1" >main.go
	printf 'package main\n\nfunc helper() {\n}\n' >other.go
	git add -A && git commit -qm init
}

# expect NAME MAIN_SOURCE AGENT WANT_EXIT WANT_OUTCOME
expect() {
	repo "$2"
	rm -f "$WORK/metrics.tsv"
	got=0
	METRICS=$WORK/metrics.tsv ROUND_TIMEOUT=2s AGENT="cat >/dev/null; $3" bash "$LOOP" 2 >/dev/null 2>&1 || got=$?
	outcome=$(tail -n 1 "$WORK/metrics.tsv" 2>/dev/null | cut -f6 || true)
	dirty=$(git status --porcelain)
	# A build artifact committed by the loop counts as dirt too.
	git ls-files --error-unmatch demo >/dev/null 2>&1 && dirty="binary demo committed"
	if [ "$got" -eq "$4" ] && [ "$outcome" = "$5" ] && [ -z "$dirty" ]; then
		echo "ok   $1"
		return
	fi
	echo "FAIL $1: exit $got (want $4), outcome '$outcome' (want '$5'), dirty tree: ${dirty:-no}"
	failures=$((failures + 1))
}

expect "fixes and commits" "$BAD_MAIN" "$FIX" 0 committed
expect "already clean" "$CLEAN_MAIN" "true" 0 ""
expect "out of scope" "$BAD_MAIN" "$FIX; echo >>other.go" 1 rolled_back_scope
expect "breaks the build" "$BAD_MAIN" "printf 'package main\n' >main.go" 1 rolled_back_build
expect "changes nothing" "$BAD_MAIN" "true" 1 no_change
expect "agent crashes" "$BAD_MAIN" "echo >>main.go; exit 3" 1 agent_failed
expect "agent hangs" "$BAD_MAIN" "sleep 30" 1 agent_timeout
expect "does not compile" 'package main\n\nfunc main() {\n' "true" 1 ""

[ "$failures" -eq 0 ] || { echo "$failures failed" >&2; exit 1; }
echo "all passed"
