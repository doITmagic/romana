# Jurnal de modificări

## Neversionat

- **Fișierele mari din `references/` sparte pe subiecte**, ca agentul să citească mai puțin: `gramatica.md` → plus `gramatica-forme.md` și `gramatica-scriere.md`; `tipografie.md` → plus `majuscule-si-nume.md`; `ce-nu-e-greseala.md` → plus `ce-nu-e-greseala-stil.md`; `prieteni-falsi-si-calcuri.md` → plus `calcuri.md`. ID-urile regulilor rămân aceleași.
- **Validatorul** limitează mărimea: `SKILL.md` până la 6.000 de cuvinte, fiecare fișier din `references/` până la 6.500.
- **Evals:** consumul în tokeni pe rulare (tabelul „Consum”) și pornirea skillului măsurată și pe Codex.
- **Corectat:** v0.1.0 avea de fapt 854 de reguli, dintre care 289 de normă (RO-SEM-247 coborâtă la recomandare și RO-CAL-507 nouă, ambele intrate înaintea tagului); README-ul și intrarea de mai jos spuneau 853 și 290.

## 0.1.0 (2026-10-22)

Prima versiune: skillul și uneltele cu care se verifică.

- **Skillul `romana`.** `SKILL.md` descrie cele trei moduri (redactare, corectură, verificare), cele șase niveluri, metoda pe treceri și profilurile de text. Cele 24 de fișiere din `references/` au 853 de reguli, dintre care 290 de normă, fiecare cu sursă.
- **Ce nu e greșeală:** 99 de forme corecte care par greșite, ca agentul să nu le „corecteze”.
- **Pluginul** pentru Claude Code și Codex: manifestele din rădăcină, din `.claude-plugin/` și din `.agents/`.
- **Evals:** 146 de cazuri publice, rulate cu și fără skill pe Claude și pe Codex, cu notare mecanică și cu un judecător separat pentru redactare. Rezultatele sunt în `evals/RAPORT.md`.
- **Unelte:** validatorul skillului, paza fișierelor, exportul pentru pagina de recenzie și testul de independență în container (modul static).

Rezultatele din `evals/RAPORT.md` (161 de cazuri, Claude și Codex): pragul e atins pe Claude la toate criteriile, iar pe Codex la toate în afară de redactare, unde diferența e de un singur caz. Rămân pentru 0.1.1: verificarea semnalează încă puțin prea mult pe presa reală, iar pe textele curate Codex mai face sugestii de stil; câteva nuanțe ale notării (date reformatate, marcaje pomenite în text); o rulare pe RONACC, corpusul academic de corectură gramaticală pentru română.
