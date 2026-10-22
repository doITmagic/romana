# Evals

Repetări pe caz și braț: 1. Judecător (redactare): opus, mediana a 3 judecăți, doar pe textul produs. Probleme = încălcări Interzis; Judecător = calcuri și coerență numărate de judecător, raportate separat.

### claude (sonnet)

Valorile sunt „fără skill / cu skill”, mai puțin e mai bine (Fidelitate: cazuri fidele din n_comun). „Skill pornit” numără rulările din brațul cu skill în care skill-ul a fost invocat. „Erori/nemăsurat” numără rulările căzute sau nemăsurabile (trigger pe Codex, redactare fără judecător), pe braț. La trigger, brațul fără skill nu are sens („–”).

| Grupa | n_comun | Probleme | Judecător | Fals pozitive | Lipsă | Fidelitate | Skill pornit | Erori/nemăsurat |
|---|---|---|---|---|---|---|---|---|
| context | 18 | 0 / 0 | 0 / 0 | 21 / 3 | 11 / 3 | 18/18 / 18/18 | 18/18 | 0 / 0 |
| corectura | 14 | 0 / 0 | 0 / 0 | 0 / 2 | 14 / 1 | 14/14 / 14/14 | 14/14 | 0 / 0 |
| curat | 31 | 0 / 0 | 0 / 0 | 1 / 0 | 0 / 0 | 31/31 / 31/31 | 31/31 | 0 / 0 |
| curat-micro | 3 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 3/3 / 3/3 | 3/3 | 0 / 0 |
| fidelitate | 17 | 0 / 0 | 0 / 2 | 1 / 0 | 0 / 0 | 17/17 / 17/17 | 17/17 | 0 / 0 |
| hipercorectie | 17 | 7 / 0 | 0 / 0 | 9 / 0 | 0 / 0 | 17/17 / 17/17 | 17/17 | 0 / 0 |
| pastrare | 13 | 0 / 0 | 0 / 0 | 5 / 0 | 1 / 0 | 13/13 / 13/13 | 13/13 | 0 / 0 |
| redactare | 20 | 0 / 0 | 7 / 1 | 0 / 0 | 0 / 0 | 20/20 / 20/20 | 20/20 | 0 / 0 |
| trigger-negativ | 11 | – / 0 | – / 0 | – / 0 | – / 0 | – / 11/11 | 0/11 | – / 0 |
| trigger-pozitiv | 11 | – / 0 | – / 0 | – / 0 | – / 0 | – / 11/11 | 11/11 | – / 0 |
| verificare | 6 | 0 / 0 | 0 / 0 | 2 / 4 | 6 / 1 | 6/6 / 6/6 | 6/6 | 0 / 0 |

### codex (gpt-6.1-sol)

Valorile sunt „fără skill / cu skill”, mai puțin e mai bine (Fidelitate: cazuri fidele din n_comun). „Skill pornit” numără rulările din brațul cu skill în care skill-ul a fost invocat. „Erori/nemăsurat” numără rulările căzute sau nemăsurabile (trigger pe Codex, redactare fără judecător), pe braț. La trigger, brațul fără skill nu are sens („–”).

| Grupa | n_comun | Probleme | Judecător | Fals pozitive | Lipsă | Fidelitate | Skill pornit | Erori/nemăsurat |
|---|---|---|---|---|---|---|---|---|
| context | 18 | 4 / 0 | 0 / 0 | 27 / 3 | 10 / 5 | 18/18 / 18/18 | – | 0 / 0 |
| corectura | 14 | 0 / 0 | 0 / 0 | 22 / 4 | 2 / 1 | 14/14 / 14/14 | – | 0 / 0 |
| curat | 31 | 0 / 0 | 0 / 0 | 10 / 6 | 0 / 0 | 31/31 / 31/31 | – | 0 / 0 |
| curat-micro | 3 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 3/3 / 3/3 | – | 0 / 0 |
| fidelitate | 17 | 0 / 0 | 1 / 0 | 0 / 0 | 0 / 0 | 17/17 / 17/17 | – | 0 / 0 |
| hipercorectie | 17 | 5 / 0 | 0 / 0 | 11 / 0 | 0 / 0 | 17/17 / 17/17 | – | 0 / 0 |
| pastrare | 13 | 0 / 0 | 0 / 0 | 14 / 0 | 2 / 0 | 13/13 / 13/13 | – | 0 / 0 |
| redactare | 20 | 0 / 0 | 0 / 2 | 0 / 0 | 0 / 0 | 20/20 / 20/20 | – | 0 / 0 |
| trigger-negativ | 0 | – / 0 | – / 0 | – / 0 | – / 0 | – / 0/0 | – | – / 11 |
| trigger-pozitiv | 0 | – / 0 | – / 0 | – / 0 | – / 0 | – / 0/0 | – | – / 11 |
| verificare | 6 | 0 / 0 | 0 / 0 | 1 / 1 | 7 / 2 | 6/6 / 6/6 | – | 0 / 0 |

### Prag de reușită

| Motor | Criteriu | Verdict | Detaliu |
|---|---|---|---|
| claude | Redactare: probleme (judecător) scad cu cel puțin 50% | atins | fără skill 7, cu skill 3, 23 cazuri; skill mai bun/egal/mai rău: 7/14/2 |
| claude | Redactare: încălcări Interzis cu skill ≤ fără skill | atins | fără skill 0, cu skill 0 |
| claude | Fals pozitive cu skill ≤ fără skill: curat | atins | fără skill 1, cu skill 0, 31 cazuri |
| claude | Fals pozitive cu skill ≤ fără skill: curat-micro | atins | fără skill 0, cu skill 0, 3 cazuri |
| claude | Fals pozitive cu skill ≤ fără skill: hipercorectie | atins | fără skill 9, cu skill 0, 17 cazuri |
| claude | Fidelitate adevărată cu skill pe fidelitate-* și pastrare-* | atins | 30 cazuri, toate fidele |
| codex | Redactare: probleme (judecător) scad cu cel puțin 50% | neatins | fără skill 1, cu skill 2, 23 cazuri; skill mai bun/egal/mai rău: 1/21/1 |
| codex | Redactare: încălcări Interzis cu skill ≤ fără skill | atins | fără skill 0, cu skill 0 |
| codex | Fals pozitive cu skill ≤ fără skill: curat | atins | fără skill 10, cu skill 6, 31 cazuri |
| codex | Fals pozitive cu skill ≤ fără skill: curat-micro | atins | fără skill 0, cu skill 0, 3 cazuri |
| codex | Fals pozitive cu skill ≤ fără skill: hipercorectie | atins | fără skill 11, cu skill 0, 17 cazuri |
| codex | Fidelitate adevărată cu skill pe fidelitate-* și pastrare-* | atins | 30 cazuri, toate fidele |

