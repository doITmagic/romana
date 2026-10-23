# Privacy

romana collects no personal data. DO IT MAGIC SOFTWARE, which makes it, receives nothing from it: no texts, no usage data, no telemetry, no accounts.

## The skill

The skill is a set of instructions and reference files that the AI assistant reads. It makes no network requests and stores nothing. The text you give the assistant to write, correct or check is handled by the app you use it in (Claude, ChatGPT, Codex or another agent), under that provider's privacy policy.

## The `romana` binary (optional, installed separately)

The binary gives the assistant the dictionary tools and, in Claude Code, the hook that checks the files the agent writes. Everything it does runs on your computer.

- **Its one network request is to dexonline.ro.** It downloads the official dexonline export (`https://dexonline.ro/static/download/dex-database.sql.gz`, GPL v2+) when the local database is missing, and checks for a newer export at most once a week. The request is a plain download with the User-Agent `romana/<version> (+https://github.com/doITmagic/romana)`. It carries no text and no identifier of yours; like any download, dexonline.ro sees your IP address.
- **Local files.** The database goes to `~/.local/share/romana/dex.sqlite` (or the path in `$ROMANA_DEX`). The downloaded export and the update state go to `~/.cache/romana/`. They stay there until you delete them.
- **Nothing leaves your computer.** The MCP tools (`cauta_cuvant`, `verifica_text`) and the hook only read: the words or the text the assistant passes them, and the file the agent has just written. They answer the assistant locally and send nothing anywhere.
- **The install script.** `install.sh` downloads the binary from this repository's GitHub Releases and checks it against `SHA256SUMS`. GitHub's privacy policy applies to that download.

## Contact

Questions or concerns: open an issue at <https://github.com/doITmagic/romana/issues>, or use the contact page at <https://doitmagic.ro/>.

---

*Pe scurt, în română:* romana nu colectează date. Skillul nu face nicio cerere în rețea. Binarul opțional descarcă doar exportul oficial de pe dexonline.ro, cel mult o dată pe săptămână, și ține baza pe calculatorul tău. Uneltele lui și hook-ul lucrează local și nu trimit nimic nicăieri.
