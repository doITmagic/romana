# Testul de independență

Verifică spec §4.1: skill-ul `skills/romana` funcționează singur, pe o mașină fără cărți și fără dinlume.

Imaginea (`Dockerfile`) pornește de la `golang:1.26.0-bookworm` (versiunea din `go.mod`) și conține doar repo-ul, fără `.git`, `.superpowers` sau `evals-privat` (vezi `Dockerfile.dockerignore`). Nu conține credențiale.

## Modul `static` (implicit)

```
tools/independenta/ruleaza.sh          # sau: ruleaza.sh static
```

Construiește imaginea și rulează `verifica-skill.sh` în container, cu `--network none`. Se oprește la primul eșec:

1. `go test ./tools/...`
2. `go run ./tools/validare/cmd skills/romana`
3. `go run ./tools/paza/cmd .`
4. `grep` după `/home/`, `romana-privat`, `dexonline.ro/definitie`, `romana-backup`, `.superpowers` în `skills/` (nu trebuie să găsească nimic)
5. fără legături simbolice în `skills/romana`; trimiterile Markdown din `SKILL.md` și `references/` rămân în interiorul skill-ului, nu sunt absolute și țin spre fișiere existente.

## Modul `evals` (pregătit, nerulat)

```
ANTHROPIC_API_KEY=... tools/independenta/ruleaza.sh evals
```

Refuză să pornească fără `ANTHROPIC_API_KEY`. Cheia trece doar ca variabilă de mediu la rulare, nu intră în imagine.

- Imaginea se construiește cu `--build-arg INSTALEAZA_CLAUDE=1` (implicit 0). Etapa folosește installerul nativ oficial (`curl -fsSL https://claude.ai/install.sh | bash -s <versiune>`, documentat la code.claude.com/docs/en/setup). `CLAUDE_VERSIUNE` implicit e `stable`; de fixat la un număr exact înainte de rularea reală.
- Rețeaua: containerul de evals e pe o rețea Docker `--internal`, fără ieșire directă. Singura cale spre exterior e un container tinyproxy (`tinyproxy.conf`, `filtru`) cu `FilterDefaultDeny Yes`, care permite doar `api.anthropic.com` pe portul 443. Am ales tinyproxy fiindcă e mic, filtrează după nume de gazdă la `CONNECT` și nu cere certificate proprii. Telemetria și auto-actualizarea CLI-ului sunt oprite (`CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1`, `DISABLE_AUTOUPDATER=1`).
- Comanda din container: `go run ./tools/evals/cmd --set public --motor claude --ambele`.

## Ce rămâne de făcut

- Rularea modului `evals`, după ce Răzvan confirmă cazurile din `evals/cazuri` (spec §8.5) și există o cheie API.
- Verificat la prima rulare că CLI-ul nu cere și alte gazde decât `api.anthropic.com` (dacă da, se adaugă în `filtru` doar ce e strict necesar). Pe urmă, versiunea `claude` fixată exact.
- Pragul din Task 11 în container: se verifică din raportul `evals`.
