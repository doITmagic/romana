# romana

Skill pentru agenți AI (Claude Code, Codex) care scriu, corectează și verifică texte în română. În faza 2 se adaugă o unealtă de dicționar (CLI și server MCP) construită peste exportul oficial dexonline.

## Structura

- `skills/romana/`: skillul. Conține `SKILL.md` și `references/` (24 de fișiere, ~850 de reguli `RO-<PIL>-<NNN>`).
- `evals/`: cazurile de test publice (`cazuri/*.json`) și sursele lor (`SURSE.md`).
- `tools/`: unelte în Go.
  - `paza`: ce nu are voie să intre în repo.
  - `validare`: formatul skillului.
  - `recenzie`: exportul pentru pagina de recenzie.
  - `evals`: rularea cazurilor și notarea lor.
  - `independenta`: testul în container.
- `docs/superpowers/`: designul și planul fazei 1. `docs/STARE.md` e starea curentă.

## Comenzi

```
go run ./tools/validare/cmd skills/romana   # validatorul strict; trebuie să treacă
go run ./tools/paza/cmd .                   # paza
go test ./tools/...
go run ./tools/evals/cmd -h                 # evals: --motor claude|codex|amandoua, --ambele, --caz <id>
bash tools/independenta/ruleaza.sh static   # testul de independență, fără rețea
```

Hook-ul de pre-commit din `.githooks/` rulează paza și validatorul. Nu-l ocoli cu `--no-verify`.

## Reguli ferme

- **Nicio cerere automată către dexonline.ro**, nici din skill, nici de la agenți. Condițiile site-ului o interzic.
  - Singura excepție, decisă de Răzvan pe 23.10.2026: binarul `romana` descarcă exportul oficial (`https://dexonline.ro/static/download/dex-database.sql.gz`, licență GPL v2+) când baza lipsește și verifică actualizări cel mult o dată pe săptămână, cu User-Agent `romana/<versiune>`. Nimic altceva de pe site.
  - Testele și CI-ul nu fac cereri reale; descărcarea se testează cu un server local.
- **Licențe.**
  - Repo-ul e MIT.
  - Nu intră text copiat din cărți, din DOOM 3, din ÎOOP sau din export, nici liste extrase din export.
  - Tiparele de text generat preluate din alte skilluri se rescriu de la zero.
  - Setul public de evals are doar texte scrise de noi și texte CC BY.
  - Presa, literatura și manualele intră doar în setul privat, care stă local, în afara repo-ului.
- **Skillul e independent.** Nu cere nicio cale locală, niciun URL de urmat și nicio unealtă. Verificarea e `tools/independenta`.
- **Niveluri.** Autoritatea e tabelul din `SKILL.md`. Hipercorectitudinea, adică o formă corectă condamnată, e la fel de gravă ca o greșeală lăsată să treacă. Deciziile proprietarului despre niveluri sunt în registrul local (vezi `CLAUDE.local.md`).
- **Diacritice:** ș și ț cu virgulă. Ghilimele: „…”, cu «…» pentru citatul din citat. Apostrof: ’.
- **Corecturile se fac în skill, nu în cazurile de test**, când un eval pică.
- **Nimic pe GitHub fără cererea explicită a proprietarului.** Înainte de primul push public, documentele interne (`docs/STARE.md`, `docs/superpowers/`) și istoria trebuie curățate de căi locale și de numele surselor private.
- **Commit-uri** în română, cu prefix (`feat(skill):`, `fix(evals):` …).
