# Evals

Repetări pe caz și braț: 1. Judecător (redactare): opus, mediana a 3 judecăți, doar pe textul produs. Probleme = încălcări Interzis; Judecător = calcuri și coerență numărate de judecător, raportate separat.

### claude (sonnet)

Valorile sunt „fără skill / cu skill”, mai puțin e mai bine (Fidelitate: cazuri fidele din n_comun). „Skill pornit” numără rulările din brațul cu skill în care skill-ul a fost invocat. „Erori/nemăsurat” numără rulările căzute sau nemăsurabile (trigger pe Codex, redactare fără judecător), pe braț. La trigger, brațul fără skill nu are sens („–”).

| Grupa | n_comun | Probleme | Judecător | Fals pozitive | Lipsă | Fidelitate | Skill pornit | Erori/nemăsurat |
|---|---|---|---|---|---|---|---|---|
| context | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | 12/12 | 0 / 0 |
| corectura | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | 13/13 | 0 / 0 |
| hipercorectie | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | 17/17 | 0 / 0 |

### codex (gpt-6.1-sol)

Valorile sunt „fără skill / cu skill”, mai puțin e mai bine (Fidelitate: cazuri fidele din n_comun). „Skill pornit” numără rulările din brațul cu skill în care skill-ul a fost invocat. „Erori/nemăsurat” numără rulările căzute sau nemăsurabile (trigger pe Codex, redactare fără judecător), pe braț. La trigger, brațul fără skill nu are sens („–”).

| Grupa | n_comun | Probleme | Judecător | Fals pozitive | Lipsă | Fidelitate | Skill pornit | Erori/nemăsurat |
|---|---|---|---|---|---|---|---|---|
| context | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | 12/12 | 0 / 0 |
| corectura | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | 13/13 | 0 / 0 |
| hipercorectie | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | 17/17 | 0 / 0 |

### Consum

Tokeni medii pe rulare, „fără skill / cu skill”, fără judecător. Intrarea include partea citită din cache.

| Motor | Intrare | Ieșire | Rulări |
|---|---|---|---|
| claude | – / 60.495 | – / 1.417 | 0 / 42 |
| codex | – / 121.366 | – / 517 | 0 / 42 |

### MCP (claude)

Valorile sunt „claude cu skill / claude+mcp”, pe cazurile cu rezultat valid în ambele brațe; mai puțin e mai bine. Consumul e media tokenilor pe rulare (intrarea include cache-ul).

| Grupa | n_comun | Probleme | Fals pozitive | Lipsă | Intrare | Ieșire |
|---|---|---|---|---|---|---|
| context | 12 | 0 / 0 | 13 / 11 | 1 / 2 | 36.727 / 37.817 | 982 / 931 |
| corectura | 13 | 0 / 0 | 0 / 0 | 2 / 2 | 65.979 / 97.467 | 2.066 / 1.847 |
| hipercorectie | 17 | 1 / 1 | 1 / 1 | 0 / 0 | 73.077 / 90.770 | 1.228 / 1.145 |

Criteriu: cu MCP, fals pozitivele și lipsurile nu cresc față de „cu skill”: **neatins** (fals pozitive 14 → 12, lipsă 3 → 4, 42 cazuri).

### MCP (codex)

Valorile sunt „codex cu skill / codex+mcp”, pe cazurile cu rezultat valid în ambele brațe; mai puțin e mai bine. Consumul e media tokenilor pe rulare (intrarea include cache-ul).

| Grupa | n_comun | Probleme | Fals pozitive | Lipsă | Intrare | Ieșire |
|---|---|---|---|---|---|---|
| context | 12 | 0 / 0 | 0 / 0 | 2 / 2 | 127.709 / 134.022 | 493 / 528 |
| corectura | 13 | 0 / 0 | 2 / 1 | 0 / 1 | 126.233 / 133.013 | 682 / 674 |
| hipercorectie | 17 | 0 / 0 | 0 / 0 | 0 / 0 | 113.167 / 122.402 | 408 / 418 |

Criteriu: cu MCP, fals pozitivele și lipsurile nu cresc față de „cu skill”: **neatins** (fals pozitive 2 → 1, lipsă 2 → 3, 42 cazuri).

### Prag de reușită

| Motor | Criteriu | Verdict | Detaliu |
|---|---|---|---|
| claude | Redactare: probleme (judecător) scad cu cel puțin 50% | nemăsurabil | 0 cazuri comune, minim 20 |
| claude | Redactare: încălcări Interzis cu skill ≤ fără skill | nemăsurabil | 0 cazuri comune, minim 20 |
| claude | Fals pozitive cu skill ≤ fără skill: curat | nemăsurabil | 0 cazuri comune, minim 5 |
| claude | Fals pozitive cu skill ≤ fără skill: curat-micro | nemăsurabil | 0 cazuri comune, minim 3 |
| claude | Fals pozitive cu skill ≤ fără skill: hipercorectie | nemăsurabil | 0 cazuri comune, minim 5 |
| claude | Fidelitate adevărată cu skill pe fidelitate-* și pastrare-* | nemăsurabil | 0 cazuri, minim 5 |
| codex | Redactare: probleme (judecător) scad cu cel puțin 50% | nemăsurabil | 0 cazuri comune, minim 20 |
| codex | Redactare: încălcări Interzis cu skill ≤ fără skill | nemăsurabil | 0 cazuri comune, minim 20 |
| codex | Fals pozitive cu skill ≤ fără skill: curat | nemăsurabil | 0 cazuri comune, minim 5 |
| codex | Fals pozitive cu skill ≤ fără skill: curat-micro | nemăsurabil | 0 cazuri comune, minim 3 |
| codex | Fals pozitive cu skill ≤ fără skill: hipercorectie | nemăsurabil | 0 cazuri comune, minim 5 |
| codex | Fidelitate adevărată cu skill pe fidelitate-* și pastrare-* | nemăsurabil | 0 cazuri, minim 5 |

