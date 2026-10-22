# Traducerea din engleză

Fișierul acesta te ajută în două situații: traduci sau adaptezi un text din engleză în română, ori primești un text românesc care sună tradus din engleză. Regulile de limbă ale textului final sunt cele din restul fișierelor. Aici găsești metoda de traducere și cazurile care nu au alt loc.

Abrevierile din câmpul „Sursă”: DEX ’98 = *Dicționarul explicativ al limbii române*, ediția a II-a, 1998; Mic ghid CE = *Mic ghid de redactare clară*, Comisia Europeană, 2013 (ghidul prezintă recomandări, nu reguli); ghid gov-md = *Ghid de redactare lingvistică și tehnoredactare a actelor normative*, Cancelaria de Stat a Republicii Moldova, 2026 (recomandare, fără caracter normativ). Detaliile sunt în [surse.md](../surse.md).

## Cum folosești fișierul

- **La traducere sau adaptare** parcurgi metoda de mai jos, apoi regulile de la sfârșit.
- **La corectura unui text care sună tradus** cauți semnele: topica străină, substantive abstracte în loc de verbe, conectori din alt registru, expresii traduse cuvânt cu cuvânt. Nu corectezi o frază doar pentru că seamănă cu structura engleză. Ordinea subiect, verb, complement e normală în română (vezi RO-NUG-061 din [ce-nu-e-greseala.md](../ce-nu-e-greseala.md)).
- **În orice caz** fidelitatea față de original e prima condiție. Ce păstrezi și ce nu ai voie să adaugi e în [fidelitate.md](../fidelitate.md).

## Metoda de traducere

Metoda are cinci pași.

1. **Citești tot textul înainte să traduci.** Ai nevoie să știi cine vorbește, cui, cu ce scop și în ce registru. Fără asta alegi cuvinte, nu sensuri.
2. **Spui sensul în proză liberă.** Pentru fiecare frază sau paragraf îți formulezi, în două-trei propoziții simple, ce afirmă originalul: cine face ce, cât de sigur e, când, în ce condiții. Încă nu te gândești la formă. Dacă nu poți spune sensul fără să te uiți la cuvintele englezești, nu l-ai înțeles încă.
3. **Scrii afirmațiile direct în română.** Pornești de la ce ai notat la pasul doi, nu de la structura frazei englezești și nu schimbi cuvinte prin sinonime pe același schelet. Alegi subiectul, verbul și complementele ca un autor român care ar scrie lucrul acesta de la zero.
4. **Compari cu originalul.** Pentru fiecare frază românească te întrebi ce afirmație din original îi corespunde. Pentru fiecare afirmație din original te întrebi unde e în text. O afirmație nouă e un fapt adăugat, una lipsă e un fapt șters.
5. **Întrebările finale, pentru fiecare paragraf.** Ar spune un autor român lucrul acesta în felul acesta? S-a pierdut vreo nuanță? O adaptare nu e un rezumat: scurtezi doar dacă ți se cere.

Exemplu, cu originalul „We take data privacy seriously and are committed to ensuring that your information remains protected.”

Slab:

> Luăm confidențialitatea datelor în serios și suntem angajați să asigurăm că informațiile dumneavoastră rămân protejate.

Bun:

> Ținem la confidențialitatea datelor dumneavoastră și ne angajăm să le protejăm.

Prima variantă păstrează structura engleză: „a fi angajat să asigure că”, „informațiile rămân protejate”. A doua spune același lucru cu verbe românești, iar angajamentul („ne angajăm”) și protecția au rămas.

### Când există deja o versiune

Dacă ți se dă o traducere existentă, o revizuiești, nu o iei de la capăt. Păstrezi ce e corect și firesc, schimbi doar ce spune altceva decât originalul sau sună tradus și nu atingi alegerile de terminologie ale autorului, dacă sunt consecvente. O traducere refăcută de la zero introduce abateri noi acolo unde prima era bună.

### Sensul din frază, nu echivalentul din dicționar

Un cuvânt englezesc are în dicționarul bilingv mai multe echivalente românești. Echivalentul din capul listei nu e mereu cel din frază. Alegi după ce spune fraza, nu după cuvântul care seamănă. Când fraza nu arată domeniul, alegi formularea adevărată în orice domeniu, nu termenul unui domeniu anume: „output” al unei echipe e „câtă treabă am făcut” sau „rezultatele”, nu „producția”, care numește bunurile fabricate. Prietenii falși cunoscuți sunt în [prieteni-falsi-si-calcuri.md](../prieteni-falsi-si-calcuri.md); aici e regula generală, în RO-CAL-042.

### Unde sunt celelalte reguli

Regulile generale contra românei englezite sunt în fișierele lor, nu aici:

- subiectul nearticulat, substantivele alăturate, gerunziul fără subiect propriu, pronumele care reia genul englezesc, imperativul negativ: [gramatica.md](../gramatica.md) (RO-GRA-149, RO-GRA-150, RO-GRA-159, RO-GRA-066) și [coerenta.md](../coerenta.md) (RO-COE-026);
- prietenii falși și calcurile frazeologice („a face sens”, „a lua loc”): [prieteni-falsi-si-calcuri.md](../prieteni-falsi-si-calcuri.md);
- stilul nominal, pasivul lung, ritmul de fraze scurte și paralele, conectorii din alt registru, titlul care nu se poate citi cu voce tare: [stilistica.md](../stilistica.md);
- majusculele din titluri (fără Title Case), ghilimelele, linia de pauză: [tipografie.md](../tipografie.md);
- „billion”, separatorul de mii, „de” după numeral: [numere-date-unitati.md](../numere-date-unitati.md);
- „Please…” devine „Te rugăm”, „X de succes”, pluralul în mesajele de interfață: [tipuri/interfata.md](../tipuri/interfata.md);
- cuvintele englezești lăsate în text, articolul și pluralul lor: [imprumuturi.md](../imprumuturi.md).

## Reguli

### RO-CAL-501 Traducerea automată literală se reformulează, nu se repară gramatical
- Nivel: de semnalat
- Greșit: „să repete lucru în repository” → Corect: „indicații pe care le pot aplica din nou”
- De ce: Un text tradus automat, cuvânt cu cuvânt, poate să nu aibă o greșeală de gramatică care să se corecteze și totuși să nu spună nimic. Aici originalul vorbea despre indicații care pot fi aplicate de mai multe ori (engl. „repeatable guidance”), iar traducerea a lipit cuvinte. Nu pui articolul sau acordul lipsă, ci cauți ce voia originalul și scrii asta. Dacă originalul nu e la îndemână și sensul nu se poate ghici, semnalezi fraza și întrebi.
- Sursă: recomandare a skillului, fără normă

### RO-CAL-502 Idiomurile se traduc după sens
- Nivel: recomandare
- Greșit: „Examenul a fost o bucată de prăjitură.” → Corect: „Examenul a fost floare la ureche.”
- De ce: Engl. „a piece of cake” nu are în română imaginea prăjiturii; echivalentul consacrat e „floare la ureche”, cu sensul „lucru de puțină importanță sau gravitate, foarte ușor de rezolvat” (DEX ’98). Se poate spune și „simplu ca bună ziua”. Când o imagine există și în română (cu aceeași formă), o păstrezi; când nu există, cauți expresia românească cu același sens și același registru. Dacă nu găsești una, spui sensul fără imagine.
- Sursă: DEX ’98 s.v. floare, simplu

### RO-CAL-503 Se traduce sensul, nu cuvintele
- Nivel: recomandare
- Greșit: „Te fac faimos.” → Corect: „Te fac cunoscut.”
- De ce: Într-un text religios, engl. „make You famous” nu înseamnă celebritate lumească, ci a face cunoscut. „Faimos” (DEX ’98: „care a devenit extrem de cunoscut; vestit, celebru”) sună a glorie lumească și schimbă registrul; „cunoscut” sau „a vesti” („a aduce la cunoștință, a face cunoscut”, DEX ’98) merg mai bine. Alegerea finală ține de registrul textului și de formula liturgică a limbii-țintă, nu de o regulă. Nu adaugi cuvinte pe care originalul nu le are, iar majusculele la pronumele care se referă la divinitate rămân cum le are textul.
- Sursă: DEX ’98 s.v. faimos, vesti

### RO-CAL-504 Termenul tradus: originalul în paranteză, o singură dată
- Nivel: recomandare
- Greșit: „echilibrarea încărcării (load balancing) … echilibrarea încărcării (load balancing)” → Corect: „echilibrarea încărcării (load balancing)” la prima apariție, apoi doar „echilibrarea încărcării”
- De ce: Convenție de redactare, nu normă. Când traduci un termen tehnic, la prima apariție dai echivalentul românesc și originalul în paranteză, apoi folosești doar forma românească. Termenul păstrat ca împrumut („SQL”, „React”) nu primește paranteză; la nevoie îl explici la prima apariție. Termenul străin dat în paranteză se scrie cu cursive dacă textul le permite (RO-TIP-060 din [tipografie.md](../tipografie.md)). Titlul unei opere străine rămâne în original, între ghilimele românești, cu traducerea în paranteză la prima mențiune: „The Great Gatsby” (în română, „Marele Gatsby”). Ordinea e invers decât la termeni, pentru că titlul e numele operei, iar traducerea doar îl explică.
- Sursă: convenție de redactare, fără normă

### RO-CAL-505 Claritatea are prioritate față de literalitate
- Nivel: recomandare
- Greșit: „Utilizatorii sunt capabili de a realiza modificarea setărilor prin intermediul meniului.” → Corect: „Utilizatorii pot schimba setările din meniu.”
- De ce: Se traduce sensul, nu structura frazei englezești, fără să se schimbe intenția sau conținutul mesajului. Prima variantă urmează engl. „are able to perform changes to settings via the menu”; a doua spune același lucru cu un verb și un complement. Regula nu te lasă să scurtezi ce contează: dacă originalul face o distincție („pot”, nu „trebuie”), o păstrezi.
- Sursă: Mic ghid CE, introducere (ghidul prezintă recomandări, nu reguli)

### RO-CAL-506 Întrebările retorice rămân întrebări
- Nivel: recomandare
- Corect: „Și cum arată aceste condiții?” (pentru engl. „So what do these conditions look like?”). Nu o transformi în „Condițiile se prezintă după cum urmează.”
- De ce: Originalul pune o întrebare retorică, mijloc expresiv pe care îl are și româna. O transformi în afirmație doar dacă registrul nu admite întrebări retorice (text juridic, administrativ, mesaje de sistem); altfel pierzi tonul autorului. Nici nu o calchiezi cuvânt cu cuvânt.
- Sursă: recomandare a skillului, fără normă

### RO-CAL-042 Sensul din frază, nu primul echivalent din dicționar
- Nivel: recomandare
- Greșit: „Comitetul a rezolvat să amâne decizia.” → Corect: „Comitetul a hotărât să amâne decizia.”
- De ce: Engl. „resolved to postpone” înseamnă că a luat hotărârea de a amâna. „A rezolva” are în română alt sens, „a găsi soluția unei probleme, a pune capăt unei controverse” (DEX ’98), deci fraza spune altceva decât originalul. Echivalentul din dicționarul bilingv a fost alegerea greșită. Verbul românesc îl alegi după ce face subiectul în frază. Dacă nu poți deduce sensul, întrebi, nu ghicești.
- Sursă: DEX ’98 s.v. rezolva; DEX ’96 s.v. hotărî

### RO-CAL-051 Cuvintele de umplutură ale originalului nu se traduc una câte una
- Nivel: recomandare
- Greșit: „Sincer, noi practic doar am vrut să facem într-adevăr lucrurile simple.” → Corect: „Am vrut doar să simplificăm lucrurile.”
- De ce: Cuvinte ca „just”, „really”, „actually”, „basically”, „honestly” ajustează tonul fără să aducă informație. Traduse fiecare cu câte un adverb românesc, fac fraza grea și neutră, cum ar fi: *sincer, practic, doar, într-adevăr*. Le redai după sens (dacă schimbă ceva) sau le omiți, mai ales când se adună. Când cuvântul poartă informație („doar” = „numai”, „actually” = „de fapt, contrar așteptărilor”), îl păstrezi. În citat nu scoți nimic (vezi RO-FID-026 din [fidelitate.md](../fidelitate.md)).
- Sursă: recomandare a skillului, fără normă
