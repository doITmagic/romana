#!/bin/sh
# Instalează binarul „romana” de pe GitHub Releases, fără Go.
#
# Variabile:
#   ROMANA_VERSIUNE   versiunea de instalat (de exemplu v0.2.0); implicit, ultima
#   ROMANA_DIR        directorul țintă; implicit $HOME/.local/bin
#   ROMANA_URL_BAZA   adresa de bază a release-urilor (folosită de teste)
#
# Scriptul nu face nicio cerere către dexonline.ro. Baza de date o descarcă
# binarul însuși, la prima pornire.
set -eu

URL_BAZA=${ROMANA_URL_BAZA:-https://github.com/doITmagic/romana/releases}
DIR=${ROMANA_DIR:-$HOME/.local/bin}

eroare() {
  echo "romana: $*" >&2
  exit 1
}

# Sistemul și arhitectura.
case $(uname -s) in
  Linux) OS=linux ;;
  Darwin) OS=darwin ;;
  *) eroare "sistem neacceptat: $(uname -s). Pe Windows descarcă arhiva .zip din $URL_BAZA" ;;
esac
case $(uname -m) in
  x86_64|amd64) ARH=amd64 ;;
  aarch64|arm64) ARH=arm64 ;;
  *) eroare "arhitectură neacceptată: $(uname -m)" ;;
esac

# Descărcarea: curl sau wget.
if command -v curl >/dev/null 2>&1; then
  descarca() { curl -fsSL -o "$2" "$1"; }
  adresa_finala() { curl -fsSL -o /dev/null -w '%{url_effective}' "$1"; }
elif command -v wget >/dev/null 2>&1; then
  descarca() { wget -q -O "$2" "$1"; }
  adresa_finala() { wget -q --spider -S "$1" 2>&1 | sed -n 's/^ *[Ll]ocation: *//p' | tail -n 1; }
else
  eroare "lipsește curl sau wget"
fi

# Suma de control: sha256sum sau shasum.
if command -v sha256sum >/dev/null 2>&1; then
  suma_fisier() { sha256sum "$1" | cut -d ' ' -f 1; }
elif command -v shasum >/dev/null 2>&1; then
  suma_fisier() { shasum -a 256 "$1" | cut -d ' ' -f 1; }
else
  eroare "lipsește sha256sum sau shasum, deci nu pot verifica suma de control"
fi

# Versiunea: cea cerută sau ultima.
TAG=${ROMANA_VERSIUNE:-}
if [ -z "$TAG" ]; then
  TAG=$(adresa_finala "$URL_BAZA/latest" | sed 's|.*/||') || TAG=""
  case $TAG in
    v[0-9]*) ;;
    *) eroare "nu am putut afla ultima versiune; setează ROMANA_VERSIUNE" ;;
  esac
fi

ARHIVA="romana_${TAG}_${OS}_${ARH}.tar.gz"
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT INT TERM

echo "romana: descarc $ARHIVA"
descarca "$URL_BAZA/download/$TAG/$ARHIVA" "$TMP/$ARHIVA" || eroare "nu am putut descărca $ARHIVA"
descarca "$URL_BAZA/download/$TAG/SHA256SUMS" "$TMP/SHA256SUMS" || eroare "nu am putut descărca SHA256SUMS"

# Verifică suma înainte de orice dezarhivare.
SUMA_ASTEPTATA=$(awk -v f="$ARHIVA" '$2 == f || $2 == "*" f { print $1; exit }' "$TMP/SHA256SUMS")
[ -n "$SUMA_ASTEPTATA" ] || eroare "$ARHIVA nu apare în SHA256SUMS"
SUMA_CALCULATA=$(suma_fisier "$TMP/$ARHIVA")
if [ "$SUMA_ASTEPTATA" != "$SUMA_CALCULATA" ]; then
  eroare "suma de control nu se potrivește pentru $ARHIVA (așteptată $SUMA_ASTEPTATA, găsită $SUMA_CALCULATA); nu instalez nimic"
fi

mkdir -p "$TMP/extras"
tar -xzf "$TMP/$ARHIVA" -C "$TMP/extras" || eroare "arhiva nu se poate dezarhiva"
[ -f "$TMP/extras/romana" ] || eroare "arhiva nu conține binarul „romana”"

mkdir -p "$DIR"
cp "$TMP/extras/romana" "$DIR/romana.nou"
chmod 755 "$DIR/romana.nou"
mv "$DIR/romana.nou" "$DIR/romana"
echo "romana: instalat $DIR/romana ($TAG)"

case ":$PATH:" in
  *":$DIR:"*) ;;
  *) echo "romana: $DIR nu e în PATH. Adaugă în profilul shell-ului: export PATH=\"$DIR:\$PATH\"" ;;
esac
