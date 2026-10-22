#!/bin/sh
# Testul pentru install.sh: arhive false, servite de un server local.
# Nu face nicio cerere în afara mașinii locale.
set -eu

AICI=$(cd "$(dirname "$0")" && pwd)
TMP=$(mktemp -d)
PID=""
curata() {
  if [ -n "$PID" ]; then kill "$PID" 2>/dev/null || true; fi
  rm -rf "$TMP"
}
trap curata EXIT INT TERM

TAG=v9.9.9
case $(uname -s) in
  Linux) OS=linux ;;
  Darwin) OS=darwin ;;
  *) echo "SARIT: sistem neacceptat"; exit 0 ;;
esac
case $(uname -m) in
  x86_64|amd64) ARH=amd64 ;;
  aarch64|arm64) ARH=arm64 ;;
  *) echo "SARIT: arhitectură neacceptată"; exit 0 ;;
esac
command -v python3 >/dev/null 2>&1 || { echo "SARIT: lipsește python3"; exit 0; }

suma() {
  if command -v sha256sum >/dev/null 2>&1; then sha256sum "$@"; else shasum -a 256 "$@"; fi
}

# Arhiva falsă, în structura unei versiuni de pe GitHub Releases.
SERV="$TMP/serv"
REL="$SERV/download/$TAG"
mkdir -p "$REL" "$TMP/pachet"
printf '#!/bin/sh\necho "romana fals"\n' > "$TMP/pachet/romana"
chmod +x "$TMP/pachet/romana"
ARHIVA="romana_${TAG}_${OS}_${ARH}.tar.gz"
tar -C "$TMP/pachet" -czf "$REL/$ARHIVA" romana
(cd "$REL" && suma "$ARHIVA" > SHA256SUMS)

PORT=$(python3 -c 'import socket;s=socket.socket();s.bind(("127.0.0.1",0));print(s.getsockname()[1])')
python3 -m http.server "$PORT" --bind 127.0.0.1 --directory "$SERV" >/dev/null 2>&1 &
PID=$!
sleep 1

PICA=0
verifica() {
  nume=$1
  shift
  if "$@"; then echo "ok   $nume"; else echo "PICĂ $nume"; PICA=1; fi
}
instaleaza() { # director, versiune, fișier de ieșire
  ROMANA_URL_BAZA="http://127.0.0.1:$PORT" ROMANA_VERSIUNE="$2" ROMANA_DIR="$1" \
    sh "$AICI/install.sh" >"$3" 2>&1
}

# 1. Instalare reușită.
DIR1="$TMP/bin1"
instaleaza "$DIR1" "$TAG" "$TMP/ies1" || { cat "$TMP/ies1"; PICA=1; }
verifica "binarul ajunge în directorul țintă" test -x "$DIR1/romana"
verifica "binarul instalat rulează" test "$("$DIR1/romana")" = "romana fals"
verifica "avertizează că directorul nu e în PATH" grep -q "PATH" "$TMP/ies1"

# 2. Suma greșită oprește instalarea.
sed 's/^./0/' "$REL/SHA256SUMS" > "$REL/SHA256SUMS.nou"
mv "$REL/SHA256SUMS.nou" "$REL/SHA256SUMS"
DIR2="$TMP/bin2"
if instaleaza "$DIR2" "$TAG" "$TMP/ies2"; then
  echo "PICĂ suma greșită trebuia să oprească instalarea"; PICA=1
else
  echo "ok   suma greșită oprește instalarea"
fi
verifica "nu rămâne niciun binar după suma greșită" test ! -e "$DIR2/romana"
verifica "mesajul numește suma de control" grep -qi "suma de control" "$TMP/ies2"

# 3. Versiune inexistentă.
DIR3="$TMP/bin3"
if instaleaza "$DIR3" v0.0.1 "$TMP/ies3"; then
  echo "PICĂ versiunea inexistentă trebuia să pice"; PICA=1
else
  echo "ok   versiunea inexistentă oprește instalarea"
fi
verifica "nu rămâne niciun binar după versiunea inexistentă" test ! -e "$DIR3/romana"

exit "$PICA"
