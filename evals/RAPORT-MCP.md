# Evals

Repetări pe caz și braț: 2. Judecător (redactare): opus, mediana a 3 judecăți, doar pe textul produs. Probleme = încălcări Interzis; Judecător = calcuri și coerență numărate de judecător, raportate separat.

### claude (sonnet)

Valorile sunt „fără skill / cu skill”, mai puțin e mai bine (Fidelitate: cazuri fidele din n_comun). „Skill pornit” numără rulările din brațul cu skill în care skill-ul a fost invocat. „Erori/nemăsurat” numără rulările căzute sau nemăsurabile (trigger pe Codex, redactare fără judecător), pe braț. La trigger, brațul fără skill nu are sens („–”).

| Grupa | n_comun | Probleme | Judecător | Fals pozitive | Lipsă | Fidelitate | Skill pornit | Erori/nemăsurat |
|---|---|---|---|---|---|---|---|---|
| context | 12 | 0 / 0 | 0 / 0 | 14±2 / 7.5±2.5 | 2.5±0.5 / 0 | 12/12 / 12/12 | 24/24 | 0 / 0 |
| corectura | 13 | 0 / 0 | 0 / 0 | 0 / 0 | 12 / 1±1 | 13/13 / 13/13 | 26/26 | 0 / 0 |
| cuvinte | 4 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 4/4 / 4/4 | 8/8 | 0 / 0 |
| hipercorectie | 17 | 7±1 / 0 | 0 / 0 | 9±1 / 0 | 0 / 0 | 17/17 / 17/17 | 34/34 | 0 / 0 |

### codex (gpt-6.1-sol)

Valorile sunt „fără skill / cu skill”, mai puțin e mai bine (Fidelitate: cazuri fidele din n_comun). „Skill pornit” numără rulările din brațul cu skill în care skill-ul a fost invocat. „Erori/nemăsurat” numără rulările căzute sau nemăsurabile (trigger pe Codex, redactare fără judecător), pe braț. La trigger, brațul fără skill nu are sens („–”).

| Grupa | n_comun | Probleme | Judecător | Fals pozitive | Lipsă | Fidelitate | Skill pornit | Erori/nemăsurat |
|---|---|---|---|---|---|---|---|---|
| context | 12 | 0 / 0 | 0 / 0 | 16.5±0.5 / 0 | 4±1 / 2 | 12/12 / 12/12 | 24/24 | 0 / 0 |
| corectura | 13 | 0 / 0 | 0 / 0 | 8.5±0.5 / 1.5±0.5 | 2±1 / 0.5±0.5 | 13/13 / 13/13 | 26/26 | 0 / 0 |
| cuvinte | 4 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 4/4 / 4/4 | 8/8 | 0 / 0 |
| hipercorectie | 17 | 5 / 0 | 0 / 0 | 11 / 0 | 0 / 0 | 17/17 / 17/17 | 34/34 | 0 / 0 |

### Consum

Tokeni medii pe rulare, „fără skill / cu skill”, fără judecător. Intrarea include partea citită din cache.

| Motor | Intrare | Ieșire | Rulări |
|---|---|---|---|
| claude | 4.369 / 77.698 | 564 / 1.444 | 92 / 92 |
| codex | 19.048 / 162.816 | 212 / 666 | 92 / 92 |

### MCP (claude)

Valorile sunt „claude cu skill / claude+mcp”, pe cazurile cu rezultat valid în ambele brațe; mai puțin e mai bine. Consumul e media tokenilor pe rulare (intrarea include cache-ul).

| Grupa | n_comun | Probleme | Fals pozitive | Lipsă | Intrare | Ieșire |
|---|---|---|---|---|---|---|
| context | 12 | 0 / 0 | 7.5 / 0 | 0 / 0 | 41.480 / 76.704 | 966 / 992 |
| corectura | 13 | 0 / 0 | 0 / 0 | 1 / 1 | 89.808 / 107.580 | 2.138 / 1.714 |
| cuvinte | 4 | 0 / 0 | 0 / 0 | 0 / 0 | 92.328 / 77.250 | 1.339 / 814 |
| hipercorectie | 17 | 0 / 0 | 0 / 0 | 0 / 0 | 90.560 / 95.800 | 1.277 / 1.176 |

Criteriu: cu MCP, fals pozitivele și lipsurile nu cresc față de „cu skill”: **atins** (fals pozitive 7.5 → 0, lipsă 1 → 1, 46 cazuri).

### MCP (codex)

Valorile sunt „codex cu skill / codex+mcp”, pe cazurile cu rezultat valid în ambele brațe; mai puțin e mai bine. Consumul e media tokenilor pe rulare (intrarea include cache-ul).

| Grupa | n_comun | Probleme | Fals pozitive | Lipsă | Intrare | Ieșire |
|---|---|---|---|---|---|---|
| context | 12 | 0 / 0 | 0 / 0.5 | 2 / 2 | 173.536 / 145.455 | 663 / 586 |
| corectura | 13 | 0 / 0 | 1.5 / 1 | 0.5 / 1 | 177.548 / 150.210 | 841 / 824 |
| cuvinte | 4 | 0 / 0 | 0 / 0 | 0 / 0 | 146.714 / 155.927 | 548 / 604 |
| hipercorectie | 17 | 0 / 0 | 0 / 0 | 0 / 0 | 147.772 / 138.149 | 561 / 563 |

Criteriu: cu MCP, fals pozitivele și lipsurile nu cresc față de „cu skill”: **neatins** (fals pozitive 1.5 → 1.5, lipsă 2.5 → 3, 46 cazuri).

### Patru brațe (claude)

Valorile sunt „fără skill / doar dicționar / doar skill / skill + dicționar”, pe cazurile cu rezultat valid în toate brațele; mai puțin e mai bine. „Skill pornit” numără rulările în care skillul a fost invocat, „doar skill / skill + dicționar”; „–” înseamnă un braț fără rulări.

| Grupa | n_comun | Fals pozitive | Lipsă | Skill pornit |
|---|---|---|---|---|
| context | 12 | 14 / 2 / 7.5 / 0 | 2.5 / 3.5 / 0 / 0 | 24/24 / 24/24 |
| corectura | 13 | 0 / 0 / 0 / 0 | 12 / 9 / 1 / 1 | 26/26 / 26/26 |
| cuvinte | 4 | 0 / 0 / 0 / 0 | 0 / 0 / 0 / 0 | 8/8 / 8/8 |
| hipercorectie | 17 | 9 / 11.5 / 0 / 0 | 0 / 0 / 0 / 0 | 34/34 / 34/34 |

### Patru brațe (codex)

Valorile sunt „fără skill / doar dicționar / doar skill / skill + dicționar”, pe cazurile cu rezultat valid în toate brațele; mai puțin e mai bine. „Skill pornit” numără rulările în care skillul a fost invocat, „doar skill / skill + dicționar”; „–” înseamnă un braț fără rulări.

| Grupa | n_comun | Fals pozitive | Lipsă | Skill pornit |
|---|---|---|---|---|
| context | 12 | 16.5 / 16.5 / 0 / 0.5 | 4 / 4 / 2 / 2 | 24/24 / 23/23 |
| corectura | 13 | 8.5 / 13 / 1.5 / 1 | 2 / 2.5 / 0.5 / 1 | 26/26 / 26/26 |
| cuvinte | 4 | 0 / 0 / 0 / 0 | 0 / 0 / 0 / 0 | 8/8 / 8/8 |
| hipercorectie | 17 | 11 / 11.5 / 0 / 0 | 0 / 0 / 0 / 0 | 34/34 / 34/34 |

### Prag de reușită

| Motor | Criteriu | Verdict | Detaliu |
|---|---|---|---|
| claude | Redactare: probleme (judecător) scad cu cel puțin 50% | nemăsurabil | 0 cazuri comune, minim 20 |
| claude | Redactare: încălcări Interzis cu skill ≤ fără skill | nemăsurabil | 0 cazuri comune, minim 20 |
| claude | Fals pozitive cu skill ≤ fără skill: curat | nemăsurabil | 0 cazuri comune, minim 5 |
| claude | Fals pozitive cu skill ≤ fără skill: curat-micro | nemăsurabil | 0 cazuri comune, minim 3 |
| claude | Fals pozitive cu skill ≤ fără skill: hipercorectie | atins | fără skill 9, cu skill 0, 17 cazuri |
| claude | Fidelitate adevărată cu skill pe fidelitate-* și pastrare-* | nemăsurabil | 0 cazuri, minim 5 |
| codex | Redactare: probleme (judecător) scad cu cel puțin 50% | nemăsurabil | 0 cazuri comune, minim 20 |
| codex | Redactare: încălcări Interzis cu skill ≤ fără skill | nemăsurabil | 0 cazuri comune, minim 20 |
| codex | Fals pozitive cu skill ≤ fără skill: curat | nemăsurabil | 0 cazuri comune, minim 5 |
| codex | Fals pozitive cu skill ≤ fără skill: curat-micro | nemăsurabil | 0 cazuri comune, minim 3 |
| codex | Fals pozitive cu skill ≤ fără skill: hipercorectie | atins | fără skill 11, cu skill 0, 17 cazuri |
| codex | Fidelitate adevărată cu skill pe fidelitate-* și pastrare-* | nemăsurabil | 0 cazuri, minim 5 |

