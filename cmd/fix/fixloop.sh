#!/usr/bin/env bash
# Runs an agent until Edil reports nothing.
#
#   AGENT='claude -p' EDIL=/path/to/edil cmd/fix/fixloop.sh [max_rounds]
#
# AGENT is a shell command that reads the prompt on stdin and edits files in
# the current directory. Each round may only touch files that have findings,
# plus defines.go and type.go next to them; anything else is rolled back.
set -euo pipefail

EDIL=${EDIL:?set EDIL to the edil binary}
AGENT=${AGENT:?set AGENT to a command that reads the prompt on stdin}
MAX_ROUNDS=${1:-5}
FAILURE=1
# One line per rule per round: time, repo, round, rule, findings, outcome.
# Lives outside the repo so it never dirties the working tree.
METRICS=${METRICS:-${XDG_STATE_HOME:-$HOME/.local/state}/edil/metrics.tsv}
mkdir -p "$(dirname "$METRICS")"

command -v jq >/dev/null || { echo "jq is required" >&2; exit "$FAILURE"; }
if [ -n "$(git status --porcelain)" ]; then
	echo "working tree must be clean" >&2
	exit "$FAILURE"
fi

findings() {
	go vet -json -vettool="$EDIL" ./... 2>&1 | grep -v '^#' | jq -r --arg root "$PWD/" '
		to_entries[] | .value | to_entries[] | .key as $rule | .value[]
		| "\(.posn | ltrimstr($root))\t\($rule)\t\(.message)"'
}

allowed_files() {
	cut -f1 <<<"$1" | sed -E 's/:[0-9]+:[0-9]+$//' | sort -u | while read -r file; do
		dir=$(dirname "$file")
		echo "$file"
		echo "$dir/defines.go" | sed 's#^\./##'
		echo "$dir/type.go" | sed 's#^\./##'
	done | sort -u
}

changed_files() {
	{ git diff --name-only; git ls-files --others --exclude-standard; } | sort -u
}

rollback() {
	git reset --hard -q
	git clean -fdq
}

record() {
	cut -f2 <<<"$current" | sort | uniq -c | while read -r count rule; do
		printf '%s\t%s\t%s\t%s\t%s\t%s\n' "$(date -u +%FT%TZ)" "$(basename "$PWD")" "$round" "$rule" "$count" "$1"
	done >>"$METRICS"
}

prompt_for() {
	cat <<PROMPT
Fix every finding below. Rules are in CONVENTIONS.md.
Edit only the listed files, plus defines.go and type.go in the same directory.
Never add suppression comments. Keep behavior identical.

$1
$2
PROMPT
}

feedback=""
for round in $(seq 1 "$MAX_ROUNDS"); do
	current=$(findings)
	if [ -z "$current" ]; then
		echo "clean after $((round - 1)) rounds"
		exit 0
	fi
	echo "round $round: $(wc -l <<<"$current") findings"
	allowed=$(allowed_files "$current")
	prompt_for "$(tr '\t' ' ' <<<"$current")" "$feedback" | bash -c "$AGENT" >/dev/null || true
	feedback=""
	if [ -z "$(changed_files)" ]; then
		feedback="Your previous attempt changed nothing; the findings are still there."
		echo "no changes"
		record no_change
		continue
	fi
	outside=$(comm -23 <(changed_files) <(echo "$allowed") || true)
	if [ -n "$outside" ]; then
		feedback="Your previous attempt was rolled back: it touched files with no findings: $(tr '\n' ' ' <<<"$outside")"
		echo "$feedback"
		rollback
		record rolled_back_scope
		continue
	fi
	if ! build_log=$(go build ./... 2>&1 && go test ./... 2>&1); then
		feedback="Your previous attempt was rolled back: build or tests failed:
$build_log"
		echo "build/test failed, rolled back"
		rollback
		record rolled_back_build
		continue
	fi
	git add -A
	git commit -q -m "edil fixloop round $round"
	record committed
done

if [ -z "$(findings)" ]; then
	echo "clean after $MAX_ROUNDS rounds"
	exit 0
fi
echo "still failing after $MAX_ROUNDS rounds" >&2
exit "$FAILURE"
