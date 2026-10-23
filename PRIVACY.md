# Privacy

romana is made by DO IT MAGIC SOFTWARE SRL (Pitești, Romania, CUI 34288367). It has three parts, and only one of them receives anything from you: the public MCP server at `mcp.doitmagic.ro`.

## In short

| | Skill | Public MCP server | `romana` binary |
|---|---|---|---|
| **Personal data collected** | none | the text you send to `verifica_greseli` and the query you send to `cauta_regula`, plus the technical data of any web request (IP address, time, client name) | none |
| **Purpose** | | to answer that request, and nothing else | |
| **Recipients** | | nobody; the server runs on a Hetzner Online GmbH machine in Germany | |
| **Retention** | | texts and queries: none, they stay in memory only while the answer is computed; web access log: see below | |
| **Your control** | | send only the text you want checked; use the skill or the binary instead, which send nothing | |

No accounts, no tracking, no analytics, no advertising, and nothing is sold or shared.

## The skill

The skill is a set of instructions and reference files that the AI assistant reads. It makes no network requests and stores nothing. The text you give the assistant to write, correct or check is handled by the app you use it in (Claude, ChatGPT, Codex or another agent), under that provider's privacy policy.

## The public MCP server (`mcp.doitmagic.ro`)

When your AI assistant calls one of the server's tools, it sends the server the text to check (`verifica_greseli`) or the words to look up (`cauta_regula`). The server:

- **processes them in memory** to compute the answer and keeps nothing: no copy of the text, no database, no log of requests;
- **does not use them** for anything else, such as training, statistics or profiling;
- **runs in a container on a server rented from Hetzner Online GmbH** in Germany, which stores no content of the requests.

The web proxy in front of the server, shared with the other sites on the same machine, keeps the usual access log of every web request: IP address, time, the address requested (`/mcp`), the response code and the client name (User-Agent). It does not contain the texts or the queries. The log is used only to keep the server secure and running.

The server limits how many requests it accepts from the same address; for that, it keeps the IP address in memory for a few minutes.

## The `romana` binary (optional, installed separately)

The binary gives the assistant the dictionary tools and, in Claude Code, the hook that checks the files the agent writes. Everything it does runs on your computer.

- **Its one network request is to dexonline.ro.** It downloads the official dexonline export (`https://dexonline.ro/static/download/dex-database.sql.gz`, GPL v2+) when the local database is missing, and checks for a newer export at most once a week. The request is a plain download with the User-Agent `romana/<version> (+https://github.com/doITmagic/romana)`. It carries no text and no identifier of yours; like any download, dexonline.ro sees your IP address. `--offline` or `ROMANA_FARA_RETEA=1` turns it off.
- **Local files.** The database goes to `~/.local/share/romana/dex.sqlite` (or the path in `$ROMANA_DEX`). The downloaded export and the update state go to `~/.cache/romana/`. They stay there until you delete them.
- **Nothing leaves your computer.** The MCP tools (`cauta_cuvant`, `verifica_text`) and the hook only read: the words or the text the assistant passes them, and the file the agent has just written. They answer the assistant locally and send nothing anywhere.
- **The install script.** `install.sh` downloads the binary from this repository's GitHub Releases and checks it against `SHA256SUMS`. GitHub's privacy policy applies to that download.

## Contact

Questions, requests or concerns about your data: open an issue at <https://github.com/doITmagic/romana/issues>, or use the contact page at <https://doitmagic.ro/>.

---

*Pe scurt, în română:* skillul nu trimite nimic, iar binarul opțional descarcă doar exportul oficial de pe dexonline.ro și lucrează local. Serverul public `mcp.doitmagic.ro` primește doar textul sau întrebarea pe care asistentul i le trimite, le procesează în memorie ca să răspundă și nu le păstrează. Proxy-ul web din fața lui ține jurnalul obișnuit de acces (adresa IP, ora, adresa cerută), fără texte.
