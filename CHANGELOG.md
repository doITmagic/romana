# Jurnal de modificări

## 0.1.0 (2026-10-22)

Prima versiune: skillul și uneltele cu care se verifică.

- **Skillul `romana`.** `SKILL.md` descrie cele trei moduri (redactare, corectură, verificare), cele șase niveluri, metoda pe treceri și profilurile de text. Cele 24 de fișiere din `references/` au 853 de reguli, dintre care 290 de normă, fiecare cu sursă.
- **Ce nu e greșeală:** 99 de forme corecte care par greșite, ca agentul să nu le „corecteze”.
- **Pluginul** pentru Claude Code și Codex: manifestele din rădăcină, din `.claude-plugin/` și din `.agents/`.
- **Evals:** 146 de cazuri publice, rulate cu și fără skill pe Claude și pe Codex, cu notare mecanică și cu un judecător separat pentru redactare. Rezultatele sunt în `evals/RAPORT.md`.
- **Unelte:** validatorul skillului, paza fișierelor, exportul pentru pagina de recenzie și testul de independență în container (modul static).

Rezultatele din `evals/RAPORT.md` (161 de cazuri, Claude și Codex): pragul e atins pe Claude la toate criteriile, iar pe Codex la toate în afară de redactare, unde diferența e de un singur caz. Rămân pentru 0.1.1: verificarea semnalează încă puțin prea mult pe presa reală, iar pe textele curate Codex mai face sugestii de stil; câteva nuanțe ale notării (date reformatate, marcaje pomenite în text); o rulare pe RONACC, corpusul academic de corectură gramaticală pentru română.
