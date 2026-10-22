# romana

**An Agent Skill that teaches AI agents to write, correct and review Romanian the way a careful Romanian editor would.**

*Un skill pentru agenții AI (Claude Code, Codex) care scriu, corectează și verifică texte în română: fără calcuri din engleză, fără cuvinte inventate și fără „corecturi” aplicate unor forme corecte.*

Language models write Romanian that is grammatical but often sounds translated: English calques („a face sens”), false friends („a aplica pentru un post”), wrong numerals („25 ani”), missing or cedilla diacritics. Asked to fix a text, they also "correct" forms that were right all along, such as „Dați click” or „orice întârzieri”. `romana` gives the agent the rules, the sources behind them and, just as important, the list of things that are *not* errors.

## What it does

The skill works in three modes, chosen from the request:

| Mode | The agent… |
|---|---|
| **Write** (*redactare*) | writes new Romanian text, or translates into Romanian, following the rules from the start |
| **Correct** (*corectură*) | fixes the errors in a text and leaves everything else untouched |
| **Review** (*verificare*) | reports problems without changing the text, each with its place, level and reason |

It covers four areas: grammar, semantics (word choice, false friends, calques), style and coherence, plus typography, numbers and dates. It adapts to the kind of text: news, marketing, blog, e-mail, user interface, documentation.

Every finding has a **level**, and the level decides what the agent may do:

| Level | Meaning | What the agent does |
|---|---|---|
| norm (*normă*) | contradicts DOOM 3, ÎOOP or the dictionary | corrects it |
| recommendation (*recomandare*) | a preference of the sources, clarity, register | proposes, does not impose |
| disputed (*disputat*) | the sources disagree | mentions both uses, changes nothing |
| avoidable anglicism (*anglicism evitabil*) | a borrowing the dictionary records, where Romanian has a natural equivalent | proposes the equivalent |
| flag (*de semnalat*) | an ambiguity, not an error | asks the author |
| not an error (*ce nu e greșeală*) | accepted variants and pseudo-rules | leaves it alone |

Two principles run through all of it. A correct form "corrected" into another is treated as seriously as a missed error. And the agent never changes the facts: who said what, how certain it is, the numbers, what a quotation says.

## What is inside

```
skills/romana/
├── SKILL.md            the entry point: modes, levels, method, which file to read when
└── references/         29 files, 854 rules
    ├── gramatica.md, semantica.md, coerenta.md, stilistica.md, tipografie.md, …
    ├── ce-nu-e-greseala*.md     99 forms that look wrong but are correct
    ├── limbi-sursa/             calques from English, German and Italian
    └── tipuri/                  conventions for news, marketing, blog, e-mail, UI, docs
```

Each rule has an ID, a level, a wrong → right pair, the reason and the source:

```
### RO-NUM-001 „De” după numeralul care se termină în 20–99 sau 00
- Nivel: normă
- Greșit: „1.500 cazuri noi” → Corect: „1.500 de cazuri noi”
- De ce: Decid ultimele două cifre. …
- Sursă: INJ p. 84–86; DOOM 3 §6.5.3; ÎOOP6 p. 108; …
```

Of the 854 rules, 289 are norms, 439 recommendations, 99 "not an error", 20 flags, 6 disputed and 1 avoidable anglicism. The norms rest on DOOM 3 (the official orthographic dictionary, 2021), ÎOOP 6 (the guide to spelling and punctuation, 2024), the explanatory dictionary and the academic grammar. The skill contains rules and examples written for this project; no text is copied from those works.

The skill is self-contained: it needs no network, no local files and no other tool. A container test (`tools/independenta`) checks this before each release.

## Install

**Claude Code**, straight from a clone of this repository:

```
claude --plugin-dir .
```

or through a local marketplace:

```
claude plugin marketplace add ./
claude plugin install romana@romana
```

The skill appears as `romana:romana`.

**Codex**, through a local marketplace:

```
codex plugin marketplace add .
codex plugin add romana@romana
```

**Other agents** (Cursor, Gemini CLI, Copilot and others): copy the `skills/romana/` folder into the tool's skills directory.

The local installs above were tested with Claude Code 2.x and codex-cli 0.163. Installing directly from GitHub and through skills.sh (`npx skills add doITmagic/romana`) will be tested once the repository is public.

### The `romana` binary (optional dictionary tool)

The skill and the binary work together or separately. The skill needs no binary; the binary (CLI and MCP server) needs no skill and can be used from a terminal or from any MCP client.

You do not need Go. On Linux and macOS:

```
curl -fsSL https://raw.githubusercontent.com/doITmagic/romana/main/install.sh | sh
```

The script picks the right archive from GitHub Releases, checks it against `SHA256SUMS` and installs `romana` into `~/.local/bin` (`ROMANA_DIR` changes the directory, `ROMANA_VERSIUNE` picks a version). It tells you if the directory is not in `PATH`. On Windows, download the `.zip` for your architecture from the Releases page and put `romana.exe` in a directory on `PATH`; the plugin hook does not run on Windows. With Go installed, `go install github.com/doITmagic/romana/cmd/romana@latest` works too.

On first use, `romana cauta`, `romana verifica` and `romana mcp` download the dictionary data by themselves from the official dexonline export (GPL v2+), about 380 MB, and import it into a local SQLite database in roughly one minute. The export and the database are never shipped with romana.

- No network: `--offline` or `ROMANA_FARA_RETEA=1` stops every request.
- From a file you downloaded yourself: `romana import <dex-database.sql.gz>`.
- Updates are checked at most once a week; `romana update` checks now.
- Database path: `--db`, then `$ROMANA_DEX`, then `~/.local/share/romana/dex.sqlite`.

The plugin declares the MCP server (`.mcp.json`) and the hook (`hooks/hooks.json`), so they start once the binary is on `PATH`. Both are only for Claude Code; Codex reads neither, so register the server there by hand (below). The hook does not run on Windows, which has no `sh`. Without the plugin, register the server by hand:

```
claude mcp add romana -- romana mcp
```

and for Codex, in `~/.codex/config.toml`:

```
[mcp_servers.romana]
command = "romana"
args = ["mcp"]
```

## How it is tested

`evals/` holds the test cases: 146 public ones, written for this project or taken from CC BY texts, plus a private set of real press and literary texts that stays outside the repository. Every case runs on two engines (Claude and Codex), once with the skill and once without, and is scored mechanically:

- **false positives**: correct text that was changed or flagged;
- **missed**: real errors left in place;
- **forbidden**: forms that must stay untouched, touched anyway.

The groups include clean texts, which must come out unchanged, correct forms that look wrong, corrections, fidelity to a source, writing (scored by a separate judge model), and whether the skill starts only when it should. The skill passes when the agent does better with it than without it. Whether the skill starts is measured on Claude only, because Codex does not report which skills it loaded.

### Results

From the full run on 22 October 2026: 161 cases (146 public, 15 private), one run per case and arm. The saved outputs were then scored again after one fix to the scorer (quotes shorter than three letters) and three test cases that were missing real errors (`go run ./tools/evals/cmd --set tot --renoteaza <brut.json>`). The last skill changes before the release (new examples, two rules) were checked by re-running the 16 cases they could affect: no regressions. Claude is Claude Sonnet, Codex is gpt-6.1-sol, and the writing judge is Claude Opus (median of 3 judgements). Each cell shows the total over the group, **without skill → with skill**; lower is better.

| Group | What a good answer does | Claude | Codex |
|---|---|---|---|
| Correct forms that look wrong (17) | leaves them alone | false positives 9 → **0** | 11 → **0** |
| Clean texts (34) | leaves the text unchanged | false positives 1 → **0** | 10 → **6** |
| Corrections (14) | fixes the real errors | missed 14 → **1**, false positives 0 → 2 | missed 2 → **1**, false positives 22 → **4** |
| Review of a text in context (18) | reports real problems only | false positives 21 → **3**, missed 11 → **3** | false positives 27 → **3**, missed 10 → **5** |
| Review of real press texts (6) | reports real problems only | false positives 2 → 4, missed 6 → **1** | false positives 1 → 1, missed 7 → **2** |
| Fidelity and preservation (30) | keeps facts, figures and quotations | all faithful; false positives 5 → **0** | all faithful; false positives 14 → **0** |
| Writing and translation (23) | no calques, coherent text | judge-counted problems 7 → **3** | 1 → 2 |
| Starting the skill (22) | starts only for Romanian text work | 11 of 11 relevant, **0 of 11** unrelated | not measurable |

The skill meets every threshold on Claude. Two weak spots remain: on real press texts Claude still flags a little more than the bare model (4 false positives against 2, on 6 texts), and in corrections it makes 2 unnecessary changes where the bare model makes none, while missing far fewer real errors. On Codex it meets every threshold except writing, where the counts are too small to tell: 21 of the 23 cases came out the same, one better and one worse. An earlier run gave 6 → 0 on the same group.

The full tables, per engine and per group, are in [`evals/RAPORT.md`](evals/RAPORT.md).

**An example.** The case `hipercorectie-09` asks for a correction of a text that is already correct:

> Datele se sincronizează în timp real. Dați click pe pictogramă, apoi dați refresh listei. Nu au fost găsite rezultate pentru căutarea dvs.

Without the skill, both engines "fixed" „click” into „clic”, and Codex also rewrote „dați refresh listei”. With the skill, both returned the text unchanged: DOOM 3 records „click” and „refresh”, and the sentence is grammatical.

**What these numbers are not.** They come from this project's own test set, with one run per case, and they hold for the models listed above. There is no official benchmark yet for an agent that writes Romanian by the norms and avoids hypercorrection. A run on RONACC, the academic Romanian grammatical error correction corpus, is planned for 0.1.1.

### Cost in tokens

The skill makes the agent read rules, so it costs tokens. Measured on 23 October 2026 on 37 cases (correct forms that look wrong, clean press texts, corrections), one run per case, average input tokens per run:

| | Without skill | Skill v0.1.0 | Skill after the split (0.1.1) |
|---|---|---|---|
| Claude | 4,453 | 103,172 | **86,839** (−16%) |
| Codex | 18,155 | 129,001 | 132,521 |

A large part of the input is read from the prompt cache, which is billed at a lower rate: about 35,000 tokens per run on Claude and 96,000 on Codex. Output grows much less: on Claude, from 677 to 1,872 tokens.

What has been done to keep the cost down:
- The agent loads only `SKILL.md` up front and opens a reference file only when the text needs it. The table in `SKILL.md` says which file covers which problem.
- The four largest reference files were split by topic, so a file the agent opens is at most about 6,000 words instead of 12,000. The validator now refuses a `SKILL.md` above 6,000 words or a reference file above 6,500.
- The test runner reports the cost of every run (`evals/RAPORT.md`, table „Consum”), so any change that makes the skill more expensive shows up.

The split helped Claude, which reads whole files. It did not change Codex, which reads files in chunks with shell commands. The next steps planned: a compact index of the "not an error" forms, which are read before every correction (about 17,000 tokens today), and a shorter `SKILL.md`.

Run the tests yourself with:

```
go run ./tools/evals/cmd --motor amandoua --ambele --raport evals/RAPORT.md
```

## Development

The tools are in Go, under `tools/`:

```
go run ./tools/validare/cmd skills/romana   # strict format check of the skill
go run ./tools/paza/cmd .                   # blocks files that must not enter the repository
go test ./tools/...
bash tools/independenta/ruleaza.sh static   # the skill works on its own, offline
```

The pre-commit hook in `.githooks/` runs the guard and the validator. See [CONTRIBUTING.md](CONTRIBUTING.md) for the rules on sources and examples.

## Roadmap

Phase 2 adds a dictionary tool (a CLI and an MCP server) built on the official dexonline database export, so the agent can look a word up instead of guessing. Each user builds the database locally from the export; neither the export nor the database is distributed, and the tool makes no requests to dexonline.ro. The skill keeps working without it.

## License

MIT, see [LICENSE](LICENSE). Attributions are in [NOTICE](NOTICE); the sources and licenses of the test cases are in [`evals/SURSE.md`](evals/SURSE.md).
