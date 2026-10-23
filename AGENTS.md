# AGENTS.md

Instructions for AI coding agents. People should start with [README.md](README.md) and [CONTRIBUTING.md](CONTRIBUTING.md).

## What this repository is

- `skills/romana/`: an Agent Skill for writing, correcting and reviewing Romanian. Entry point: `SKILL.md`; rules in `references/`.
- `cmd/romana/`: the `romana` binary (CLI, MCP server, Claude Code hook). Packages:
  - `dex/`: dictionary;
  - `baza/`: download and update of the export;
  - `lint/`: sure errors;
  - `verifica/`: text check;
  - `mcp/`: MCP server (stdio with the dictionary; the public set without it, in `public.go`);
  - `reguli/`: the skill's rules, embedded in the binary (`reguli.go` at the root), searched by ID or words.
- `action.yml`: the GitHub Action (`uses: doITmagic/romana@v0`): installs the binary from Releases and runs `romana verifica --sigure --github`; no dictionary, no network beyond GitHub.
- `cmd/romana-public/`: the public MCP server at https://mcp.doitmagic.ro/mcp (Streamable HTTP, Echo v5): `verifica_greseli` and `cauta_regula`, never the dictionary; `GET /` is a small embedded HTML page that links to doitmagic.ro (with a looser CSP for that route only), plus `/icon.png`, `/favicon.ico` and `/.well-known/security.txt` (the dexonline export must not be redistributed). `deploy/mcp/` runs it in Docker behind nginx-proxy.
- `evals/`: public test cases; `tools/`: validator, guard, test runner, independence test.

## Commands

Run these before every commit; the pre-commit hook runs the first two again.

```bash
go run ./tools/validare/cmd skills/romana   # skill format; must pass
go run ./tools/paza/cmd .                   # files that must not enter the repository
gofmt -l . && go vet ./...                  # also GOOS=windows and GOOS=darwin go vet ./...
go test -race ./...
bash tools/independenta/ruleaza.sh static   # the skill works on its own, offline
```

Any `romana` command you run yourself: add `ROMANA_FARA_RETEA=1` (or `--offline`) and, if you need the real database, `--db ~/.local/share/romana/dex.sqlite`.

## Always

- **Romanian text:**
  - ș and ț with a comma below, never with a cedilla;
  - quotation marks „…”, with «…» inside a quotation;
  - apostrophe ’.
- **Code comments and user-facing messages** of the binary are in Romanian; README and AGENTS.md are in English.
- **Tests first** (TDD). A test must fail without the fix.
- **Commit messages** in Romanian, with a prefix: `feat(skill):`, `fix(lint):`, `docs(readme):` …
- **Dependencies:** the standard library and `modernc.org/sqlite`, nothing else.

## Ask first

- Changing the level of a rule (normă, recomandare, …): the owner decides levels.
- Removing or renaming a field in the JSON returned by `romana cauta --json`, `romana verifica --json` or the MCP tools. Adding a field is fine while the version is 0.x.
- Adding a dependency, or a new network request of any kind.

## Never

- **dexonline.ro:** never send a request there. The only exception is the binary downloading the official export (`https://dexonline.ro/static/download/dex-database.sql.gz`) when the database is missing, and checking it at most once a week. Tests use a local `httptest` server.
- **Licensed material:**
  - never commit the export or the database (`.sql`, `.gz`, `.sqlite`, `.db`);
  - never copy text from DOOM 3, ÎOOP, the dictionaries or the export into the repository: examples are written from scratch.
- **Evals:** never change a test case to make an eval pass. The fix goes into the skill.
- **Hooks:** never use `git commit --no-verify`.
- **GitHub:** never push without the owner's explicit request.
- **Hypercorrection:** never condemn a correct form. It counts as much as a missed error.

## Using `romana` as a tool

If you are an agent that uses the dictionary rather than develops it:

1. **What to look up:** call `cauta_cuvant` (or `romana cauta <word> --json`) for any rare word, neologism, derived word or inflected form you are not sure of, before you condemn or correct it.
2. **How to read the answer:**
   - `doom3: true` means DOOM 3 records the lemma. It is a hint, not proof that every inflected form is correct.
   - A `nota` on a lemma says which form the norm prefers ("degivra" → "dejivra", hence "dejivrat"). Follow it, with the reason it gives.
   - `exista: false` or `necunoscut` means the dictionary does not have the form. It is not the same as wrong: put the word on a "not verified" list and leave it unchanged.
3. **`verifica_text`** reports only the errors that are wrong in any context, plus unknown words. Its `limite` field is there for a reason: an empty report does not mean the text is correct. Agreement, meaning, coherence, punctuation and style are still yours to check.
