#!/usr/bin/env bash
# Rulează în container, fără rețea. Se oprește la primul eșec.
set -euo pipefail
cd /repo

echo "== go test ./tools/..."
go test ./tools/...

echo "== validare skills/romana"
go run ./tools/validare/cmd skills/romana

echo "== paza"
go run ./tools/paza/cmd .

echo "== căi locale / URL-uri interzise în skills/"
if grep -rnE '/home/|romana-privat|dexonline\.ro/definitie|romana-backup|\.superpowers' skills/; then
  echo "EȘEC: skills/ conține referințe interzise" >&2
  exit 1
fi
echo "nimic găsit"

echo "== legături simbolice și trimiteri în afara skill-ului"
radacina=$(realpath skills/romana)
if [ -n "$(find skills/romana -type l)" ]; then
  find skills/romana -type l >&2
  echo "EȘEC: legături simbolice în skills/romana" >&2
  exit 1
fi
nr=0
while IFS= read -r -d '' f; do
  dir=$(dirname "$f")
  # trimiteri Markdown: ](țintă) sau ](țintă "titlu")
  while IFS= read -r t; do
    t=${t%%#*}                      # fără ancoră
    [ -z "$t" ] && continue
    case "$t" in
      http://*|https://*|mailto:*) continue ;;
      /*) echo "EȘEC: cale absolută în $f: $t" >&2; exit 1 ;;
    esac
    r=$(realpath -m "$dir/$t")
    case "$r/" in
      "$radacina"/*) ;;
      *) echo "EȘEC: $f trimite în afara skill-ului: $t" >&2; exit 1 ;;
    esac
    [ -e "$r" ] || { echo "EȘEC: $f trimite spre un fișier inexistent: $t" >&2; exit 1; }
    nr=$((nr+1))
  done < <(grep -oE '\]\([^)[:space:]]+' "$f" | sed 's/^](//')
done < <(find skills/romana -type f -name '*.md' \( -name SKILL.md -o -path '*/references/*' \) -print0)
echo "$nr trimiteri verificate, toate în interior"
echo "STATIC: OK"
