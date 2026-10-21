#!/usr/bin/env bash
# Testul de independență. Folosire: ruleaza.sh [static|evals]
set -euo pipefail
MOD=${1:-static}
AICI=$(cd "$(dirname "$0")" && pwd)
REPO=$(cd "$AICI/../.." && pwd)
IMG=romana-independenta
cd "$REPO"

construieste() {  # $1 = 0|1 (instalează claude)
  docker build -q -f tools/independenta/Dockerfile \
    --build-arg INSTALEAZA_CLAUDE="$1" -t "$IMG:$2" .
}

case "$MOD" in
static)
  construieste 0 static >/dev/null
  docker run --rm --network none "$IMG:static" bash tools/independenta/verifica-skill.sh
  ;;
evals)
  : "${ANTHROPIC_API_KEY:?ANTHROPIC_API_KEY lipsește; setează-l în mediu, nu e copiat în imagine}"
  construieste 1 evals >/dev/null
  RET=romana-indep-int; PROXY=romana-indep-proxy
  cleanup() { docker rm -f "$PROXY" >/dev/null 2>&1 || true; docker network rm "$RET" >/dev/null 2>&1 || true; }
  trap cleanup EXIT
  # rețea internă (fără ieșire directă); doar proxy-ul are și o rețea externă
  docker network create --internal "$RET" >/dev/null
  docker build -q -t romana-indep-proxy -f - tools/independenta >/dev/null <<'DF'
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends tinyproxy && rm -rf /var/lib/apt/lists/*
COPY tinyproxy.conf /etc/tinyproxy/tinyproxy.conf
COPY filtru /etc/tinyproxy/filtru
CMD ["tinyproxy","-d","-c","/etc/tinyproxy/tinyproxy.conf"]
DF
  docker create --name "$PROXY" --network bridge romana-indep-proxy >/dev/null
  docker network connect "$RET" "$PROXY"
  docker start "$PROXY" >/dev/null
  P="http://$PROXY:8888"
  # cheia trece doar ca variabilă de mediu la rulare (-e fără valoare)
  docker run --rm --network "$RET" -e ANTHROPIC_API_KEY \
    -e HTTPS_PROXY="$P" -e HTTP_PROXY="$P" -e NO_PROXY= \
    -e CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1 -e DISABLE_AUTOUPDATER=1 \
    "$IMG:evals" go run ./tools/evals/cmd --set public --motor claude --ambele
  ;;
*)
  echo "folosire: $0 [static|evals]" >&2; exit 2 ;;
esac
