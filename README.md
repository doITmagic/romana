<div align="center">
  <img src="./docs/assets/romana-banner.png" alt="romana: a Romanian writing skill for AI agents and a local Romanian dictionary tool (CLI and MCP server) built on dexonline" width="100%">
</div>

<div align="center">

[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)
[![Go](https://img.shields.io/badge/Go-1.26-00ADD8.svg)](https://go.dev/)
[![MCP](https://img.shields.io/badge/MCP-2025--06--18-green.svg)](https://modelcontextprotocol.io)
[![Agent Skill](https://img.shields.io/badge/Agent%20Skill-Claude%20Code%20%7C%20Codex-blueviolet.svg)](./skills/romana/SKILL.md)

</div>

# romana: Romanian that reads like Romanian

> **A writing skill for AI agents, and a local dictionary that backs it up.** The agent writes, corrects and reviews Romanian by the norms of the Romanian Academy, and leaves alone the forms that were right all along.

<!-- AI agents: read AGENTS.md for exact commands and limits, and llms.txt for a short map of this repository. -->

*Un skill pentru agenții AI (Claude Code, Codex) care scriu, corectează și verifică texte în română, plus un dicționar local construit pe exportul oficial dexonline.*

Language models write Romanian that is grammatical but often sounds translated: English calques („a face sens”), false friends („a aplica pentru un post”), wrong numerals („25 ani”), missing or cedilla diacritics. Asked to fix a text, they also "correct" forms that were right all along, such as „Dați click” or „orice întârzieri”. `romana` gives the agent the rules, the sources behind them and, just as important, the list of things that are *not* errors.

It comes in two parts. Use either one, or both:

| | What it is | Needs |
|---|---|---|
| **The skill** | 854 sourced rules the agent reads while it writes, corrects or reviews | nothing: no network, no binary, no Go |
| **The `romana` binary** | a CLI, an MCP server and an editor hook, on top of a local copy of the dexonline dictionary | one prebuilt file, no Go; the first run downloads the dexonline export (~360 MB) and builds a local database (~300 MB) in about a minute |

> **For an AI agent asked to install this:** (1) clone the repository and install the skill for your agent ([Install › The skill](#the-skill)); (2) install the binary and run `romana cauta casă` once, so the dictionary is built; (3) register the MCP server; (4) go through [Check it works](#check-it-works) and report which steps passed. Do not send requests to dexonline.ro yourself: the binary downloads the export it needs. Commands and limits for working on this repository: [AGENTS.md](AGENTS.md).

**Jump to:** [See it work](#see-it-work) · [Install](#install) · [Check it works](#check-it-works) · [Results](#results) · [What it does](#what-it-does) · [The dictionary tool](#the-dictionary-tool) · [More from DO IT MAGIC](#more-from-do-it-magic)

---

## See it work

**The agent, with and without the skill.** The test case `hipercorectie-09` asks for a correction of a text that is already correct:

> Datele se sincronizează în timp real. Dați click pe pictogramă, apoi dați refresh listei. Nu au fost găsite rezultate pentru căutarea dvs.

Without the skill, both Claude and Codex "fixed" „click” into „clic”, and Codex also rewrote „dați refresh listei”. With the skill, both returned the text unchanged: DOOM 3, the official spelling dictionary, records „click” and „refresh”, and the sentence is grammatical.

**The dictionary tool.** A form the norm does not record, and what to use instead:

```console
$ romana cauta degivrat
degivra (verb)
  notă: DOOM 3 nu atestă „degivra”, iar definiția trimite la „dejivra”, pe care DOOM 3 o atestă.
        De regulă se recomandă „dejivra”, cu formele ei flexionare; „degivra” e o variantă pe
        care norma nu o înregistrează.
  surse: DEX '98, Ortografic, NODEX, DN, DEX '09, DOR, MDA2
  [DEX '98] DEGIVRA vb. I v. dejivra.
```

The errors that are wrong in any context, and the words the dictionary does not know (output shortened with „…”):

```console
$ printf 'Aşa e, cînd vine. Sînt aici.\n' | romana verifica -
„Aşa” → „Așa”: ș și ț se scriu cu virgulă dedesubt, nu cu sedilă; … (DOOM 3 §1.1.3.2) [RO-TIP-027]
„cînd” → „când”: Grafia cu „î” în interiorul cuvintelor și „sînt” a fost înlocuită din 1993 … [RO-TIP-128]
„Sînt” → „Sunt”: … [RO-TIP-128]
Verificarea acoperă doar existența fiecărui cuvânt în dicționar și greșelile sigure, …

$ echo "Aplicatia trimite raportul." | romana verifica -
necunoscut: „Aplicatia” – poate: aplicația
```

The tool speaks Romanian, because so do the texts it checks. Every report ends with what it does *not* check (agreement, meaning, coherence, punctuation, style), so an empty report is never read as "the text is correct".

---

## Install

### The skill

Start from a clone of this repository:

```bash
git clone https://github.com/doITmagic/romana.git
cd romana
```

**Claude Code**, installed permanently through a local marketplace:

```bash
claude plugin marketplace add ./
claude plugin install romana@romana
```

or for one session only: `claude --plugin-dir .`. The skill appears as `romana:romana`. The plugin also brings the MCP server and the hook: the MCP server starts once `romana` is on `PATH`; the hook also finds it in `~/.local/bin` or `~/go/bin`.

**Codex**, through a local marketplace:

```bash
codex plugin marketplace add .
codex plugin add romana@romana
```

**Other agents** (Cursor, Gemini CLI, Copilot and others): copy the `skills/romana/` folder into the directory where your tool looks for Agent Skills (see its documentation). We test Claude Code and Codex.

### The `romana` binary (optional)

You do not need Go. On Linux and macOS:

```bash
curl -fsSL https://raw.githubusercontent.com/doITmagic/romana/main/install.sh | sh
```

The script picks the archive for your system from GitHub Releases, checks it against `SHA256SUMS`, and installs `romana` into `~/.local/bin`. `ROMANA_DIR` changes the directory and `ROMANA_VERSIUNE` picks a version. It tells you if the directory is not on `PATH`.

<details>
<summary><b>Windows, with Go, or from source</b></summary>

- **Windows:** download the `.zip` for your architecture from the Releases page, put `romana.exe` in a folder and add that folder to `PATH`. The CLI and the MCP server work (use the entry under [Other MCP clients](#the-mcp-server)); the plugin hook does not, because Windows has no `sh`.
- **With Go:** `go install github.com/doITmagic/romana/cmd/romana@latest`.
- **From source** (also before the first release is published), in the clone: `go build -o ~/.local/bin/romana ./cmd/romana`. A binary built from source reports its version as `romana dev`; the release binaries report the tag.

</details>

Then build the dictionary once, before you register the MCP server:

```bash
romana cauta casă
```

The first run downloads the official dexonline export (GPL v2+, about 360 MB) and builds a local SQLite database in about a minute; later lookups answer at once. Nothing is downloaded again until the export changes.

### The MCP server

With the Claude Code plugin above, there is nothing to do. Otherwise:

**Claude Code:**

```bash
claude mcp add romana -- romana mcp
```

**Codex**, in `~/.codex/config.toml`:

```toml
[mcp_servers.romana]
command = "romana"
args = ["mcp"]
```

<details>
<summary><b>Other MCP clients</b></summary>

Clients that read an `mcpServers` block (Claude Desktop, Cursor and others) take the same entry:

```json
{
  "mcpServers": {
    "romana": { "command": "romana", "args": ["mcp"] }
  }
}
```

We test Claude Code and Codex. Other clients should work, since the server speaks plain MCP over stdio, but they are not in our test runs.

</details>

---

## Check it works

1. **The binary and the dictionary:**

   ```bash
   romana --version
   romana cauta casă       # lemma, forms, sources, definitions
   ```

2. **The skill is loaded.** In Claude Code, type `/romana`: the skill `romana:romana` shows up in the list. In Codex, ask the agent which skills it has. If the skill is missing, repeat the install step for your agent before going on.

3. **The skill works.** Ask the agent, in any language:

   > Corectează textul: „Dați click pe buton. Avem 25 ani de experiență.”

   A working skill keeps „Dați click” (DOOM 3 records it) and changes „25 ani” into „25 de ani” (rule RO-NUM-001), with the level and the reason.

4. **The MCP server.** In Claude Code, `/mcp` lists `romana` with two tools, `cauta_cuvant` and `verifica_text`; in Codex, `codex mcp list` shows `romana`. Then ask the agent: *"Use the cauta_cuvant tool on «degivrat»."* The answer must mention the note that recommends „dejivra”.

5. **The hook** (Claude Code, with the plugin). It runs the offline linter, without the dictionary, on `.md`, `.txt` and `.html` files right after the agent writes them, when the text is Romanian. Ask the agent to create a `.md` file with its file-writing tool and exactly this content: *„Aşa e, cînd vine echipa, iar raportul e gata pentru că sunt aici.”* Right after the write, the agent receives the list of sure errors and fixes them. If nothing happens, check `command -v romana` or `ls ~/.local/bin/romana`: the hook stays silent when it cannot find the binary.

---

## Results

### The skill, with and without

From the full run on 22 October 2026: 161 cases (146 public, 15 private), one run per case and arm. Claude is Claude Sonnet, Codex is gpt-6.1-sol, and the writing judge is Claude Opus (median of 3 judgements). Each cell shows the total over the group, **without skill → with skill**; lower is better.

| Group | What a good answer does | Claude | Codex |
|---|---|---|---|
| Correct forms that look wrong (17) | leaves them alone | false positives 9 → **0** | 11 → **0** |
| Clean texts (34) | leaves the text unchanged | false positives 1 → **0** | 10 → **6** |
| Corrections (14) | fixes the real errors | missed 14 → **1**, false positives 0 → 2 | missed 2 → **1**, false positives 22 → **4** |
| Review of a text in context (18) | reports real problems only | false positives 21 → **3**, missed 11 → **3** | false positives 27 → **3**, missed 10 → **5** |
| Review of real press texts (6) | reports real problems only | false positives 2 → 4, missed 6 → **1** | false positives 1 → 1, missed 7 → **2** |
| Fidelity and preservation (30) | keeps facts, figures and quotations | all faithful; false positives 6 → **0** | all faithful; false positives 14 → **0** |
| Writing, scored by a judge model (23: 20 writing, 3 fidelity cases) | no calques, coherent text | judge-counted problems 7 → **3** | 1 → 2 |
| Starting the skill (22) | starts only for Romanian text work | 11 of 11 relevant, **0 of 11** unrelated | not measurable |

The skill meets every threshold on Claude. Two weak spots remain: on real press texts Claude still flags a little more than the bare model (4 false positives against 2, on 6 texts), and in corrections it makes 2 unnecessary changes where the bare model makes none, while missing far fewer real errors. On Codex it meets every threshold except writing, where the counts are too small to tell: 21 of the 23 cases came out the same, one better and one worse.

### The skill, with and without the dictionary tool

<!-- MCP-RESULTS -->

### What these numbers are not

They come from this project's own test set, with one or two runs per case, and they hold for the models listed above. There is no official benchmark yet for an agent that writes Romanian by the norms and avoids hypercorrection. A run on RONACC, the academic corpus for Romanian grammatical error correction, is in progress; its reference pairs are checked against DOOM 3 first, because some of them are stricter than the norm.

The full tables, per engine and per group, are in [`evals/RAPORT.md`](evals/RAPORT.md) and [`evals/RAPORT-MCP.md`](evals/RAPORT-MCP.md).

---

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

---

## The dictionary tool

### Commands

| Command | What it does |
|---|---|
| `romana cauta <word>` | from any inflected form to its lemmas: part of speech, inflected forms, the dictionaries that record it, the definitions whose license allows display, DOOM 3 attestation |
| `romana verifica [file\|-]` | the words not in the dictionary (with suggestions that differ only in diacritics or a hyphen) and the errors that are wrong in any context |
| `romana mcp` | the MCP server, over stdio |
| `romana verifica --hook` | the Claude Code `PostToolUse` hook: runs only the linter (no database, no network) on `.md`, `.txt` and `.html` files right after the agent writes them, and gives the agent the sure errors |
| `romana import [file]` | builds the database from an export you downloaded, or downloads it |
| `romana update` | checks now whether the export changed |

`cauta` and `verifica` take `--json`; every command takes `--db <path>` and `--offline`, after the command name.

**Exit codes:** 0 when the command ran, findings or not; 1 on an error (missing file, no database); 2 on wrong usage. To fail a build on findings, use `romana verifica --json` and check that `lint` is not empty.

### MCP tools

Both tools are declared read-only (`readOnlyHint`), so Codex calls them without asking for approval, even in a read-only sandbox.

| Tool | Parameter | Returns |
|---|---|---|
| `cauta_cuvant` | `cuvant` (string, required) | lemmas with `doom3`, `surse`, `definitii`, `forme`; `trimite_la` and `nota` when DOOM 3 does not record the lemma but points elsewhere; `sugestii` when the word is unknown |
| `verifica_text` | `text` (string, required, up to 200,000 characters) | `necunoscute` (with byte positions, `sugestii`, `neverificat` and `motiv`), `lint` (the sure errors, each with a proposal and a reason), `limite` (what is not checked); at most 200 items per list, with `trunchiat` when cut |

How to read the answers:
- **`doom3: true`** means DOOM 3 records the lemma. It is a strong hint, not proof that every inflected form is right: the form table also holds old and non-recommended variants.
- **A `nota`** says, in words, which form the norm prefers ("degivra" → "dejivra"). When the same form is recorded by DOOM 3 as another word, the note asks you to decide by meaning.
- **`necunoscut`** means the dictionary does not have the form. It is not the same as wrong: new words, names and terms are often missing. Names, acronyms, code and words in other alphabets come back as `neverificat`, with the reason.

### Data, network and licenses

- **The source** is the official dexonline export, under GPL v2+. Neither the export nor the database is shipped with `romana`: each user builds the database locally.
- **Only freely redistributable text is shown.** Definitions are displayed only from the dictionaries whose license allows it (DEX ’98 and DEX ’96). The others, DOOM 3 among them, appear by name in the list of sources.
- **The only request** the binary makes is to the export itself: when the database is missing, and at most once a week to check for a newer export, with a conditional request. The check runs as a separate background process, so a lookup never waits for it; its log is `~/.cache/romana/actualizare.log`. A running MCP server switches to the new database by itself. After a failed download, it waits an hour before trying again; `romana import` starts it at once.
- **No network at all:** `--offline` or `ROMANA_FARA_RETEA=1`.
- **Database path:** `--db`, then `$ROMANA_DEX`, then `~/.local/share/romana/dex.sqlite`.

---

## What it is not

- **Not a full grammar checker.** The binary catches only what is wrong in any context (cedilla, the old spelling with „î”, „vroiam”, spaces). Agreement, meaning and style are the skill's job, done by the agent, and they need judgement.
- **Not an authority over the author.** Findings come with a level. Only norms are corrected; everything else is proposed or flagged.
- **Not a copy of the sources.** The skill holds rules and examples written for this project; no text is copied from DOOM 3, ÎOOP or the dictionaries.
- **Not a dexonline client.** It makes no requests to dexonline.ro beyond downloading the official export.
- **Not for every kind of text.** Poetry, fiction and translations *out of* Romanian are outside its scope.

---

## What is inside

```
skills/romana/
├── SKILL.md            the entry point: modes, levels, method, which file to read when
└── references/         29 files, 854 rules
    ├── gramatica*.md, semantica.md, coerenta.md, stilistica.md, tipografie.md, …
    ├── ce-nu-e-greseala*.md     93 of the 99 forms that look wrong but are correct (the rest sit next to the rules they concern)
    ├── limbi-sursa/             calques from English, German and Italian
    └── tipuri/                  conventions for news, marketing, blog, e-mail, UI, docs
cmd/romana/            the binary
dex/ baza/ lint/ verifica/ mcp/    dictionary, download, sure errors, text check, MCP server
evals/                 the public test cases and their sources
tools/                 validator, guard, test runner, independence test
```

Each rule has an ID, a level, a wrong → right pair, the reason and the source:

```
### RO-NUM-001 „De” după numeralul care se termină în 20–99 sau 00
- Nivel: normă
- Greșit: „1.500 cazuri noi” → Corect: „1.500 de cazuri noi”
- De ce: Decid ultimele două cifre. …
- Sursă: INJ p. 84–86; DOOM 3 §6.5.3; ÎOOP6 p. 108; …
```

Of the 854 rules, 289 are norms, 439 recommendations, 99 "not an error", 20 flags, 6 disputed and 1 avoidable anglicism. The norms rest on DOOM 3 (the official orthographic dictionary, 2021), ÎOOP 6 (the guide to spelling and punctuation, 2024), the explanatory dictionary and the academic grammar.

The skill is self-contained: it needs no network, no local files and no other tool. A container test (`tools/independenta`) checks this before each release.

---

## How it is tested

`evals/` holds the test cases: 146 public ones, written for this project or taken from CC BY texts, plus a private set of real press and literary texts that stays outside the repository. Almost every case runs on two engines (Claude and Codex), with and without the skill; the 22 cases about when the skill starts run only with the skill. For the dictionary tool, cases run with the skill alone and with the skill plus the MCP server. Scoring is mechanical:

- **false positives:** correct text that was changed or flagged;
- **missed:** real errors left in place;
- **forbidden:** forms that must stay untouched, touched anyway.

The groups include clean texts, which must come out unchanged, correct forms that look wrong, corrections, fidelity to a source, writing (scored by a separate judge model), and whether the skill starts only when it should. The skill passes when the agent does better with it than without it; the dictionary tool passes when it does not make the agent worse. Whether the skill starts is measured on Claude only, because Codex does not report which skills it loaded.

```bash
go run ./tools/evals/cmd --motor amandoua --ambele --raport evals/RAPORT.md     # skill vs no skill
go run ./tools/evals/cmd --motor amandoua --cu-skill --mcp --raport evals/RAPORT-MCP.md   # + dictionary tool
```

The binary itself has its own tests (`go test -race ./...`), including runs against the real database when it is present, and it was reviewed in six rounds before 0.2.0, each finding checked by a second reviewer that tried to refute it.

---

## Cost in tokens

The skill makes the agent read rules, so it costs tokens. Measured on 23 October 2026 on 37 cases (correct forms that look wrong, clean press texts, corrections), one run per case, average input tokens per run:

| | Without skill | Skill v0.1.0 | Skill after the split |
|---|---|---|---|
| Claude | 4,453 | 103,172 | **86,839** (−16%) |
| Codex | 18,155 | 129,001 | 132,521 |

A large part of the input is read from the prompt cache, which is billed at a lower rate: about 35,000 tokens per run on Claude and 96,000 on Codex. Output grows much less: on Claude, from 677 to 1,872 tokens.

What keeps the cost down:
- The agent loads only `SKILL.md` up front and opens a reference file only when the text needs it.
- The four largest reference files were split by topic, so a file the agent opens is at most about 6,000 words. The validator refuses a `SKILL.md` above 6,000 words or a reference file above 6,500.
- The test runner reports the cost of every run, so a change that makes the skill more expensive shows up.

The split helped Claude, which reads whole files. It did not change Codex, which reads files in chunks with shell commands.

---

## Stability

- **Versions:** `romana` is at 0.x. Every finding cites its rule ID (`RO-<PIL>-<NNN>`), and changes to rules, levels included, are listed in the changelog. The JSON of the CLI and the MCP tools may still gain fields.
- **Tested with:** Claude Code 2.x, codex-cli 0.163, Go 1.26. Installing straight from GitHub and through skills.sh (`npx skills add doITmagic/romana`) is tested once the repository is public.
- **Changes:** [CHANGELOG.md](CHANGELOG.md).

---

## Development

```bash
go run ./tools/validare/cmd skills/romana   # strict format check of the skill
go run ./tools/paza/cmd .                   # blocks files that must not enter the repository
go test -race ./...                         # with the real database, if present, the dex tests also run on it
bash tools/independenta/ruleaza.sh static   # the skill works on its own, offline
```

The pre-commit hook in `.githooks/` runs the guard and the validator. [CONTRIBUTING.md](CONTRIBUTING.md) has the rules on sources and examples; [AGENTS.md](AGENTS.md) has the commands and limits for coding agents.

---

## Roadmap

- a measurement on RONACC, the academic corpus for Romanian grammatical error correction, with the reference pairs checked against DOOM 3 first;
- a compact index of the "not an error" forms, so the agent reads less before every correction;
- more rules for the gaps the tests found („dejivra”, „dintr-odată”).

---

## More from DO IT MAGIC

`romana` is built by [DO IT MAGIC SOFTWARE](https://doitmagic.ro/), a software studio in Pitești, Romania, working on AI agents, Go and local-first tools.

- **[RagCode MCP](https://github.com/doITmagic/rag-code-mcp)**: semantic code search for AI assistants, 100% local, written in Go. The story behind it: [case study on doitmagic.ro](https://doitmagic.ro/blog/studiu-de-caz-ragcode-mcp/).
- **[ai-agent-skills](https://github.com/doITmagic/ai-agent-skills)**: a community registry of agent skill packs for RagCode.
- **[trello-mcp-complete](https://github.com/doITmagic/trello-mcp-complete)**: an MCP server for Trello boards, lists and cards.
- **[Speak2Me](https://doitmagic.ro/software/)**: Romanian voice dictation, offline and fully local. How we chose the model: [Romanian speech-to-text, tested offline on CPU](https://doitmagic.ro/blog/speech-to-text-romana-offline/).
- **Working with AI agents:** [programming with AI agents](https://doitmagic.ro/blog/programare-cu-agenti-ai-40-la-suta-mai-putin-cod/), and [from vibe code to production](https://doitmagic.ro/servicii/vibe-code-productie/) when an app built with AI needs to be checked and shipped.
- **Technical SEO for Romanian sites:** [site audit, structured data, sitemap and indexing for Google and Bing](https://doitmagic.ro/servicii/seo-geo/).

All projects: [doitmagic.ro/proiecte](https://doitmagic.ro/proiecte/).

---

## License

MIT, see [LICENSE](LICENSE). Attributions are in [NOTICE](NOTICE); the sources and licenses of the test cases are in [`evals/SURSE.md`](evals/SURSE.md).

The dictionary data is not part of this repository: it comes from the official [dexonline](https://dexonline.ro/) export, under GPL v2+, and each user downloads and builds it locally.

<div align="center">

**Built in Romania by [DO IT MAGIC SOFTWARE](https://doitmagic.ro/)** · [Star `romana` on GitHub](https://github.com/doITmagic/romana) if it helps your agents write better Romanian.

</div>
