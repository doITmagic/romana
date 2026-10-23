# Evals

Repetări pe caz și braț: 2. Judecător (redactare): opus, mediana a 3 judecăți, doar pe textul produs. Probleme = încălcări Interzis; Judecător = calcuri și coerență numărate de judecător, raportate separat.

### claude (sonnet)

Valorile sunt „fără skill / cu skill”, mai puțin e mai bine (Fidelitate: cazuri fidele din n_comun). „Skill pornit” numără rulările din brațul cu skill în care skill-ul a fost invocat. „Erori/nemăsurat” numără rulările căzute sau nemăsurabile (trigger pe Codex, redactare fără judecător), pe braț. La trigger, brațul fără skill nu are sens („–”).

| Grupa | n_comun | Probleme | Judecător | Fals pozitive | Lipsă | Fidelitate | Skill pornit | Erori/nemăsurat |
|---|---|---|---|---|---|---|---|---|
| context | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | – | 0 / 0 |
| corectura | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | – | 0 / 0 |
| curat | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | – | 0 / 0 |
| curat-micro | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | – | 0 / 0 |
| cuvinte | 16 | 0 / 0 | 0 / 0 | 0 / 0 | 5.5±0.5 / 3.5±0.5 | 16/16 / 16/16 | 32/32 | 0 / 0 |
| fidelitate | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | – | 0 / 0 |
| hipercorectie | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | – | 0 / 0 |
| pastrare | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | – | 0 / 0 |
| redactare | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | – | 0 / 0 |
| trigger-negativ | 0 | – / 0 | – / 0 | – / 0 | – / 0 | – / 0/0 | – | – / 0 |
| trigger-pozitiv | 0 | – / 0 | – / 0 | – / 0 | – / 0 | – / 0/0 | – | – / 0 |

### codex (gpt-6.1-sol)

Valorile sunt „fără skill / cu skill”, mai puțin e mai bine (Fidelitate: cazuri fidele din n_comun). „Skill pornit” numără rulările din brațul cu skill în care skill-ul a fost invocat. „Erori/nemăsurat” numără rulările căzute sau nemăsurabile (trigger pe Codex, redactare fără judecător), pe braț. La trigger, brațul fără skill nu are sens („–”).

| Grupa | n_comun | Probleme | Judecător | Fals pozitive | Lipsă | Fidelitate | Skill pornit | Erori/nemăsurat |
|---|---|---|---|---|---|---|---|---|
| context | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | – | 0 / 0 |
| corectura | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | – | 0 / 0 |
| curat | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | – | 0 / 0 |
| curat-micro | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | – | 0 / 0 |
| cuvinte | 16 | 0 / 0 | 0 / 0 | 0 / 0 | 5.5±0.5 / 1±1 | 16/16 / 16/16 | 32/32 | 0 / 0 |
| fidelitate | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | – | 0 / 0 |
| hipercorectie | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | – | 0 / 0 |
| pastrare | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | – | 0 / 0 |
| redactare | 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 | 0/0 / 0/0 | – | 0 / 0 |
| trigger-negativ | 0 | – / 0 | – / 0 | – / 0 | – / 0 | – / 0/0 | – | – / 0 |
| trigger-pozitiv | 0 | – / 0 | – / 0 | – / 0 | – / 0 | – / 0/0 | – | – / 0 |

### Consum

Tokeni medii pe rulare, „fără skill / cu skill”, fără judecător. Intrarea include partea citită din cache.

| Motor | Intrare | Ieșire | Rulări |
|---|---|---|---|
| claude | 4.320 / 65.203 | 231 / 1.031 | 32 / 32 |
| codex | 23.960 / 149.860 | 220 / 579 | 32 / 32 |

### MCP (claude)

Valorile sunt „claude cu skill / claude+mcp”, pe cazurile cu rezultat valid în ambele brațe; mai puțin e mai bine. Consumul e media tokenilor pe rulare (intrarea include cache-ul).

| Grupa | n_comun | Probleme | Fals pozitive | Lipsă | Intrare | Ieșire |
|---|---|---|---|---|---|---|
| context | 0 | 0 / 0 | 0 / 0 | 0 / 0 | – / – | – / – |
| corectura | 0 | 0 / 0 | 0 / 0 | 0 / 0 | – / – | – / – |
| curat | 0 | 0 / 0 | 0 / 0 | 0 / 0 | – / – | – / – |
| curat-micro | 0 | 0 / 0 | 0 / 0 | 0 / 0 | – / – | – / – |
| cuvinte | 16 | 0 / 0 | 0 / 0 | 3.5 / 1.5 | 65.203 / 74.938 | 1.031 / 745 |
| fidelitate | 0 | 0 / 0 | 0 / 0 | 0 / 0 | – / – | – / – |
| hipercorectie | 0 | 0 / 0 | 0 / 0 | 0 / 0 | – / – | – / – |
| pastrare | 0 | 0 / 0 | 0 / 0 | 0 / 0 | – / – | – / – |
| redactare | 0 | 0 / 0 | 0 / 0 | 0 / 0 | – / – | – / – |
| trigger-negativ | 0 | 0 / 0 | 0 / 0 | 0 / 0 | – / – | – / – |
| trigger-pozitiv | 0 | 0 / 0 | 0 / 0 | 0 / 0 | – / – | – / – |

Criteriu: cu MCP, fals pozitivele și lipsurile nu cresc față de „cu skill”: **atins** (fals pozitive 0 → 0, lipsă 3.5 → 1.5, 16 cazuri).

### MCP (codex)

Valorile sunt „codex cu skill / codex+mcp”, pe cazurile cu rezultat valid în ambele brațe; mai puțin e mai bine. Consumul e media tokenilor pe rulare (intrarea include cache-ul).

| Grupa | n_comun | Probleme | Fals pozitive | Lipsă | Intrare | Ieșire |
|---|---|---|---|---|---|---|
| context | 0 | 0 / 0 | 0 / 0 | 0 / 0 | – / – | – / – |
| corectura | 0 | 0 / 0 | 0 / 0 | 0 / 0 | – / – | – / – |
| curat | 0 | 0 / 0 | 0 / 0 | 0 / 0 | – / – | – / – |
| curat-micro | 0 | 0 / 0 | 0 / 0 | 0 / 0 | – / – | – / – |
| cuvinte | 16 | 0 / 0 | 0 / 0 | 1 / 0 | 149.860 / 141.969 | 579 / 540 |
| fidelitate | 0 | 0 / 0 | 0 / 0 | 0 / 0 | – / – | – / – |
| hipercorectie | 0 | 0 / 0 | 0 / 0 | 0 / 0 | – / – | – / – |
| pastrare | 0 | 0 / 0 | 0 / 0 | 0 / 0 | – / – | – / – |
| redactare | 0 | 0 / 0 | 0 / 0 | 0 / 0 | – / – | – / – |
| trigger-negativ | 0 | 0 / 0 | 0 / 0 | 0 / 0 | – / – | – / – |
| trigger-pozitiv | 0 | 0 / 0 | 0 / 0 | 0 / 0 | – / – | – / – |

Criteriu: cu MCP, fals pozitivele și lipsurile nu cresc față de „cu skill”: **atins** (fals pozitive 0 → 0, lipsă 1 → 0, 16 cazuri).

### Patru brațe (claude)

Valorile sunt „fără skill / doar dicționar / doar skill / skill + dicționar”, pe cazurile cu rezultat valid în toate brațele; mai puțin e mai bine. „Skill pornit” numără rulările în care skillul a fost invocat, „doar skill / skill + dicționar”; „–” înseamnă un braț fără rulări.

| Grupa | n_comun | Fals pozitive | Lipsă | Skill pornit |
|---|---|---|---|---|
| cuvinte | 16 | 0 / 0 / 0 / 0 | 5.5 / 3.5 / 3.5 / 1.5 | 32/32 / 32/32 |

### Patru brațe (codex)

Valorile sunt „fără skill / doar dicționar / doar skill / skill + dicționar”, pe cazurile cu rezultat valid în toate brațele; mai puțin e mai bine. „Skill pornit” numără rulările în care skillul a fost invocat, „doar skill / skill + dicționar”; „–” înseamnă un braț fără rulări.

| Grupa | n_comun | Fals pozitive | Lipsă | Skill pornit |
|---|---|---|---|---|
| cuvinte | 16 | 0 / 0 / 0 / 0 | 5.5 / 4 / 1 / 0 | 32/32 / 32/32 |

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

