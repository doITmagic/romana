# Jurnal de modificări

## 0.2.4 (2026-10-23)

- `TERMS.md`, condițiile de folosire (licența MIT, fără garanție, fără serviciu și fără cont), cu legătura în ambele manifeste; catalogul OpenAI o cere.
- Descrierea pentru catalogul OpenAI nu mai numește alte aplicații de AI, iar categoria e „Communication”.
- `romana-plugin.zip` conține și `TERMS.md`.

## 0.2.3 (2026-10-23)

- Pictograma nouă, pentru cataloage: marca „romana” de pe banner, pe care un corector o face de mână „română”, cu „^”, cu căciula și cu nota „corectă” deasupra. Sursa HTML e alături de imagine.

## 0.2.2 (2026-10-23)

- Pregătit pentru cataloagele oficiale (directorul Claude și catalogul de pluginuri OpenAI): pictograma pluginului, `PRIVACY.md` și datele listării OpenAI în `plugin.json` (`extensions.com.openai`), cu descrierea și în română.
- `romana-plugin.zip` conține și pictograma, și `PRIVACY.md`.
- Bannerul din README e în sintaxă Markdown, iar brieful lui nu mai numește imaginea în cod; validatorul directorului Claude le semnala pe amândouă.

## 0.2.1 (2026-10-23)

- **Skillul pornește singur și în claude.ai.** Descrierea spune acum explicit să fie folosit la orice cerere de scris, corectat, verificat sau tradus în română, oricât de scurtă. Într-o conversație nouă din claude.ai, la o corectură scurtă, skillul nu pornea, iar modelul „corecta” forme corecte („unul din”, „prețuri scumpe”, „click”). Testele de pornire, pe Claude: 22 din 22 la cazurile în care trebuie să pornească, 0 din 22 la cele în care nu.
- **Instalarea direct din GitHub, fără clonare:** `claude plugin marketplace add doITmagic/romana` și `codex plugin marketplace add doITmagic/romana`. În claude.ai, prin „Add marketplace” → `doITmagic/romana`; în ChatGPT, prin arhiva `romana-plugin.zip`.
- **Release-ul are și `romana-plugin.zip`**, cu skillul, fără serverul MCP și hook-ul, pentru aplicațiile care încarcă un plugin ca arhivă.

## 0.2.0 (2026-10-23)

A doua fază: unealta de dicționar. Skillul și binarul merg și împreună, și separat.

### Binarul `romana`

- **Comenzile:**
  - `cauta <cuvânt>`: de la orice formă flexionată la leme, cu partea de vorbire, formele cu etichetă, dicționarele care atestă cuvântul și definițiile pe care licența le permite;
  - `verifica [fișier]`: cuvintele care nu sunt în dicționar și greșelile sigure;
  - `mcp`: server MCP pe stdio, cu `cauta_cuvant` și `verifica_text`;
  - `import [fișier]` și `update`: baza;
  - `verifica --hook`: pentru hook-ul PostToolUse din Claude Code.

  Toate au `--json` și `--offline`.
- **Datele** vin din exportul oficial dexonline (GPL v2+). Exportul și baza nu se distribuie: fiecare utilizator își construiește baza local, într-un minut.
  - La prima folosire, binarul descarcă singur exportul (~360 MB).
  - Verifică actualizări cel mult o dată pe săptămână, cu cerere condiționată, într-un proces separat. După o descărcare eșuată face o pauză de o oră, iar `romana import` o pornește imediat.
  - Afișează doar definițiile surselor care permit asta (DEX ’98, DEX ’96). Celelalte apar doar ca nume, în lista de surse.
- **DOOM 3:** câmpul `doom3` spune dacă dicționarul ortografic al Academiei atestă lema. E un indiciu, nu o dovadă.
  - Când DOOM 3 nu atestă o lemă, dar definiția ei trimite la una atestată („DEGIVRA vb. I v. dejivra.”), `nota` spune ce se recomandă.
  - Dacă forma e atestată ca alt cuvânt (un omonim), nota cere judecata după sens.
- **Greșelile sigure** (linterul):
  - ș și ț cu sedilă (RO-TIP-027);
  - grafia veche cu „î” în interior și „sînt” (RO-TIP-128);
  - „vroiam” (RO-GRA-067);
  - spațiile înaintea punctuației și spațiul dublu.

  Nu se ating:
  - codul, URL-urile, HTML-ul și front matter-ul;
  - formulele și tabelele;
  - citatele;
  - numele proprii, cuvintele străine și compusele cu prefixoide („videoînregistrare”).
- **Un text scris complet fără diacritice** primește o singură notă (`fara_diacritice`), nu câte o observație pentru fiecare cuvânt.
- **`verifica_text`** spune în fiecare răspuns ce nu verifică: acordul, sensul, coerența, punctuația și stilul. Un raport gol nu înseamnă un text corect.
- **Distribuirea:**
  - binare gata compilate pentru Linux, macOS și Windows (amd64 și arm64), cu `SHA256SUMS`;
  - `install.sh`, fără Go;
  - sau `go install github.com/doITmagic/romana/cmd/romana@latest`.
- **Pluginul pentru Claude Code** declară serverul MCP (`.mcp.json`) și hook-ul (`hooks/hooks.json`). Uneltele se declară doar de citire, așa că le cheamă și Codex fără aprobare, chiar într-un sandbox doar de citire.

### Skillul

- **s6** („Caută cuvântul, nu-l presupune”) numește uneltele și ordinea lor: MCP, apoi CLI, apoi metoda de până acum. Mai spune:
  - un raport gol nu scutește de treceri;
  - ce face agentul cu `nota`.

  Fără unelte, skillul merge ca înainte; testul de independență trece.
- **Skillul și dicționarul se încarcă împreună.** Descrierea skillului spune acum că se folosește și când agentul are unelte de dicționar, iar instrucțiunile serverului MCP îi cer agentului să încarce întâi un skill pentru texte în română. Înainte, cu dicționarul conectat, Claude nu încărca skillul în 19 din 84 de rulări; acum îl încarcă în toate.
- **Fișierele mari din `references/` sparte pe subiecte**, ca agentul să citească mai puțin: `gramatica.md` → plus `gramatica-forme.md` și `gramatica-scriere.md`; `tipografie.md` → plus `majuscule-si-nume.md`; `ce-nu-e-greseala.md` → plus `ce-nu-e-greseala-stil.md`; `prieteni-falsi-si-calcuri.md` → plus `calcuri.md`. ID-urile regulilor rămân aceleași.

### Evals și unelte

- **Evals:**
  - patru brațe: fără skill, doar dicționar, doar skill, skill și dicționar (`--ambele --mcp`, pe Claude și pe Codex), cu tabelul „Patru brațe”;
  - 20 de cazuri noi pe cuvinte (`cuvinte-*`): variante, forme greșite și cuvinte pe care dicționarul nu le are, cu așteptarea nouă „neverificat”;
  - rezultatul final, pe 42 de cazuri, în medie pe rulare: cu skill și dicționar, Claude nu mai semnalează nicio formă corectă (7,5 doar cu skill) și ratează o greșeală, ca doar cu skill; Codex e la egalitate cu skillul singur la formele corecte (1,5) și ratează 3 greșeli, față de 2,5. Pe cuvinte, dicționarul scade greșelile scăpate de la 3,5 la 1,5 pe Claude și de la 1 la 0 pe Codex ([`evals/RAPORT-MCP.md`](evals/RAPORT-MCP.md), [`evals/RAPORT-CUVINTE.md`](evals/RAPORT-CUVINTE.md));
  - consumul în tokeni pe rulare (tabelul „Consum”);
  - pornirea skillului, măsurată și pe Codex.
- **Validatorul** limitează mărimea: `SKILL.md` până la 6.000 de cuvinte, fiecare fișier din `references/` până la 6.500.
- **Paza** blochează fișierele de bază de date și arhivele (`.sqlite`, `.db`, `.sql`, `.gz`).
- **Repo-ul devine public.** Documentele interne (starea, designul, planurile) nu mai sunt în repo. Runner-ul de evals nu mai are căi locale implicite: setul privat se ia din `ROMANA_EVALS_PRIVAT` (implicit `~/.local/share/romana/evals-privat`), iar `flux.json` din `ROMANA_FLUX`.
- **Corectat:** v0.1.0 avea de fapt 854 de reguli, dintre care 289 de normă (RO-SEM-247 coborâtă la recomandare și RO-CAL-507 nouă, ambele intrate înaintea tagului); README-ul și intrarea de mai jos spuneau 853 și 290.

## 0.1.0 (2026-10-22)

Prima versiune: skillul și uneltele cu care se verifică.

- **Skillul `romana`.** `SKILL.md` descrie cele trei moduri (redactare, corectură, verificare), cele șase niveluri, metoda pe treceri și profilurile de text. Cele 24 de fișiere din `references/` au 853 de reguli, dintre care 290 de normă, fiecare cu sursă.
- **Ce nu e greșeală:** 99 de forme corecte care par greșite, ca agentul să nu le „corecteze”.
- **Pluginul** pentru Claude Code și Codex: manifestele din rădăcină, din `.claude-plugin/` și din `.agents/`.
- **Evals:** 146 de cazuri publice, rulate cu și fără skill pe Claude și pe Codex, cu notare mecanică și cu un judecător separat pentru redactare. Rezultatele sunt în `evals/RAPORT.md`.
- **Unelte:** validatorul skillului, paza fișierelor, exportul pentru pagina de recenzie și testul de independență în container (modul static).

Rezultatele din `evals/RAPORT.md` (161 de cazuri, Claude și Codex): pragul e atins pe Claude la toate criteriile, iar pe Codex la toate în afară de redactare, unde diferența e de un singur caz. Rămân pentru 0.1.1: verificarea semnalează încă puțin prea mult pe presa reală, iar pe textele curate Codex mai face sugestii de stil; câteva nuanțe ale notării (date reformatate, marcaje pomenite în text); o rulare pe RONACC, corpusul academic de corectură gramaticală pentru română.
