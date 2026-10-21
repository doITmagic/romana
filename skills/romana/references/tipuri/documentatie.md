# Documentație

Profilul acesta se aplică documentației tehnice, ghidurilor, tutorialelor și articolelor de ajutor. Ce îl deosebește: termenii tehnici trebuie să rămână aceiași de la un capăt la altul, pașii sunt scriși ca instrucțiuni, iar codul, comenzile și numele de fișiere nu se ating. Cititorul vine cu o sarcină, nu cu timp de citit; textul îl ajută să o ducă la capăt. Normele nu se schimbă.

Abrevierile din câmpul „Sursă”: DOOM 3 = *Dicționarul ortografic, ortoepic și morfologic al limbii române*, ediția a treia; DEX '98 = *Dicționar explicativ al limbii române*, 1998; DEX '96 = *Dicționar explicativ al limbii române*, 1996; INJ = *Corectitudine lingvistică în actele judecătorești*, Institutul Național al Justiției, Chișinău, 2020; ghid gov-md = *Ghid de redactare lingvistică și tehnoredactare a actelor normative*, Cancelaria de Stat a Republicii Moldova, 2026; Mic ghid CE = *Mic ghid de redactare clară*, Comisia Europeană, 2013.

## Cum citești perechile

La o recomandare, forma din „Greșit” e forma de evitat, nu neapărat o greșeală de limbă. Dacă forma e corectă (de pildă „matrice”, „implementare”, „sare”), regula o spune în „De ce” și nu o corectezi, ci o propui. Doar nivelul „normă” cere corectare.

## Ce se schimbă față de limba literară standard

- **Registrul și adresarea.** Neutru, fără oralitate. Adresarea („tu” sau „dumneavoastră”) se alege o dată pentru întreaga documentație (RO-COE-019 din [coerenta.md](../coerenta.md); alegerea, în [registru-si-adresare.md](../registru-si-adresare.md)). Documentația poate evita adresarea cu totul și descrie ce face produsul.
- **Pașii.** Au același mod verbal de la primul la ultimul (toți la imperativ, toți la infinitiv sau toți la forma impersonală), iar condiția vine înaintea acțiunii (RO-COE-020 și RO-COE-023 din [coerenta.md](../coerenta.md)). Un pas, o acțiune.
- **Termenii.** Fiecare termen tehnic se definește o dată și se folosește la fel peste tot. Un concept, un termen (RO-COE-006 din [coerenta.md](../coerenta.md)).
- **Titlurile.** Nominale, scurte, cu majusculă doar la primul cuvânt și la numele proprii.
- **Cod și nume.** Codul, comenzile, numele de fișiere, de variabile și de funcții rămân neatinse, inclusiv în exemple și în citate. Diacriticele și ghilimelele românești nu intră în cod.
- **Vocabularul.** Cuvântul românesc obișnuit unde există și cititorul îl folosește; termenul tehnic consacrat rămâne. Verbele de interfață care se pot confunda se aleg precis („apăsați”, „selectați”).
- **Ce nu se schimbă.** Fidelitatea față de sursă: nu adaugi pași, avertismente sau valori care nu sunt în sursă (vezi [fidelitate.md](../fidelitate.md)).

## Decizia per termen

Pentru fiecare termen tehnic care apare în documentație, hotărăști o dată dintre trei variante: îl păstrezi în engleză, îl păstrezi cu o explicație la prima apariție sau îl traduci. Notezi decizia într-un glosar scurt al documentației și o respecți peste tot, inclusiv în titluri și în exemple. Aceeași noțiune tradusă o dată „implementare” și lăsată altă dată „deploy” e o rupere de coerență (RO-COE-006). Forma împrumuturilor (articol, plural, cratimă) se alege după DOOM 3 (vezi [imprumuturi.md](../imprumuturi.md)).

Doi termeni se leagă de decizie: nomenclatura (SQL, React, nume de produse) se păstrează, iar conceptele nu se traduc cuvânt cu cuvânt.

### RO-STI-818 Nomenclatura se păstrează, conceptele nu se traduc cuvânt cu cuvânt
- Nivel: recomandare
- Greșit: „Serverul face echilibrarea încărcării între instanțe.” → Corect: „Serverul distribuie cererile între instanțe (load balancing).”
- De ce: Nomenclatura tehnică (numele limbajelor, ale bibliotecilor, ale produselor) se păstrează. O traducere literală („load balancing” → „echilibrarea încărcării”) se folosește doar dacă termenul e consacrat în domeniu; „load balancing” nu e în dicționare, iar echivalentul literal nu e un termen fixat. Oricare ai alege, îl ții consecvent. Dacă publicul larg nu cunoaște termenul englezesc, îl explici la prima apariție.
- Sursă: Mic ghid CE p. 11–13 (Rec. 9); DOOM 3 (lipsa intrării „load balancing”)

### RO-STI-817 Termenul IT primește echivalentul românesc firesc, când există
- Nivel: recomandare
- Greșit: „Autentificarea se face cu OTP, fără password.” → Corect: „Autentificarea se face cu cod de unică folosință, fără parolă.”
- De ce: Termenii informatici au de regulă un echivalent românesc obișnuit, iar documentația îl folosește consecvent: „cod de unică folosință” pentru „OTP”, „fără parolă” pentru „passwordless”. Formulările sunt făcute din cuvinte românești atestate (DOOM 3 are „unic” și „folosință”), dar îmbinările nu sunt fixate în dicționare. Sigla „OTP” rămâne după prima explicație, dacă documentația o folosește. E o recomandare de uz, nu o normă.
- Sursă: convenția de localizare Mozilla pentru limba română (uz, nu normă; nereverificat local)

## Termeni de programare

### RO-STI-805 „Limbaj de programare”
- Nivel: recomandare
- Greșit: „Python este o limbă de programare ușor de învățat.” → Corect: „Python este un limbaj de programare ușor de învățat.”
- De ce: Termenul consacrat este „limbaj de programare”: DEX '98 dă pentru „limbaj” sensul informatic „sistem de caractere și simboluri folosit în programare”. „Limbă” înseamnă în dicționar limba unei comunități, deci „limbă de programare” nu contrazice dicționarul, dar nu e îmbinarea folosită în domeniu. Pentru limbile naturale rămâne „limbă”.
- Sursă: DEX '98 s.v. limbaj

### RO-STI-806 „Interpretor”, nu „interpret”
- Nivel: recomandare
- Greșit: „Interpretul Python citește fișierul rând cu rând.” → Corect: „Interpretorul Python citește fișierul rând cu rând.”
- De ce: Pentru programul care execută cod, forma folosită în domeniu este „interpretor”. „Interpret” are în DEX '98 sensurile de persoană (traducător, artist, comentator). „Interpretor” este atestat doar în Ortografic și în MDN '00, nu în DOOM 3 sau în DEX '98, de aceea propui forma, nu o impui.
- Sursă: Ortografic și MDN '00 s.v. interpretor; DEX '98 s.v. interpret

### RO-STI-807 „Indice” pentru poziția unui element
- Nivel: recomandare
- Greșit: „Elementul de la index 3 este ultimul.” → Corect: „Elementul de la indicele 3 este ultimul.”
- De ce: DEX '98 dă pentru „indice” sensul de număr, literă sau simbol care precizează valoarea sau înțelesul unui element, iar numărul de ordine al unui element dintr-un tablou se încadrează aici. În același dicționar „indice” (sensul 5) trimite și la „index”, deci cele două sunt înrudite. „Index” este tot în DOOM 3, cu alt sens („listă”, deci „indexul unei cărți” sau „indexul unei baze de date”), și rămâne corect acolo. Nu corectezi „index” în aceste sensuri.
- Sursă: DEX '98 s.v. indice, index; DOOM 3 s.v. index

### RO-STI-808 „Programare”, nu „codare”
- Nivel: recomandare
- Greșit: „Cursul predă coding pentru copii.” → Corect: „Cursul predă programare pentru copii.”
- De ce: „Coding” nu e în dicționare. Pentru activitatea de a scrie programe, DEX '98 are „programare” („elaborarea programului în vederea rezolvării unei probleme la calculator”). „Codare” înseamnă în dicționar transformarea unui text într-un cod, deci în sens tehnic corect se poate spune „codarea caracterelor”, dar nu desemnează scrierea de programe. Pentru limbajul de zi cu zi în educație, „codare” nu se condamnă, doar se precizează.
- Sursă: DEX '98 s.v. programare, codare

### RO-STI-809 „Statement”: „afirmație” sau „instrucțiune”, după context
- Nivel: recomandare
- Greșit: „Fiecare statement din funcție se execută pe rând.” → Corect: „Fiecare instrucțiune din funcție se execută pe rând.”
- De ce: În cod, „statement” se redă prin „instrucțiune”, care înseamnă în general „indicație dată cuiva” (DEX '98) și are același sens imperativ. În logică, „statement” e „afirmație”. Cele două nu se confundă. DOOM 3 are „statement” ca intrare nouă, deci forma englezească nu e greșită, doar mai puțin precisă.
- Sursă: DEX '98 s.v. instrucțiune; DOOM 3 s.v. statement

### RO-STI-811 „Deploy”
- Nivel: recomandare
- Greșit: „Facem deploy pe server vineri.” → Corect: „Instalăm aplicația pe server vineri.”
- De ce: „Deploy” (formă dedusă: lipsește din dicționare, la fel „a deploya”) are în română mai multe echivalente, după sens: „a instala pe server”, „a pune în producție” (doar dacă e vorba de mediul de producție), „a lansa” sau „a publica”. „A publica” se potrivește doar pentru lansarea publică: DEX '98 îl definește prin „a face cunoscut tuturor”, deci nu merge pentru un mediu intern sau de test. „A implementa” este corect (RO-NUG-075 din [ce-nu-e-greseala.md](../ce-nu-e-greseala.md)) și nu se corectează, dar poate fi confundat cu „to implement” (a scrie funcționalitatea), de aceea nu e prima alegere.
- Sursă: DEX '98 s.v. publica; DEX '96 s.v. implementa

### RO-STI-812 „Tablou” sau „vector” pentru „array”
- Nivel: recomandare
- Greșit: „Funcția returnează o matrice cu numele utilizatorilor.” → Corect: „Funcția returnează un tablou cu numele utilizatorilor.”
- De ce: Pentru structura unidimensională „array”, informatica școlară și universitară folosește „tablou” sau „vector”. „Matrice” este în DEX '98 un tablou dreptunghiular cu linii și coloane, adică bidimensional, deci se potrivește când structura chiar are două dimensiuni. Sursele se contrazic, de aceea nu corectezi forma aleasă de autor, ci doar o propui și o ții consecventă.
- Sursă: DEX '98 s.v. matrice, tablou, vector

### RO-STI-813 „Paranteze pătrate” sau „paranteze drepte”
- Nivel: recomandare
- Greșit: „Scrieți lista între square brackets.” → Corect: „Scrieți lista între paranteze pătrate.”
- De ce: „Paranteză pătrată” este atestat în ghidul INJ, iar „parantezele drepte” apar în același ghid și în ghidul gov-md. Ambele denumiri sunt în uz și niciuna nu se marchează drept greșeală: alegi una și o ții în tot textul.
- Sursă: INJ cap. 17; ghid gov-md §9.6

### RO-STI-814 „Acolade”
- Nivel: recomandare
- Greșit: „Blocul de cod stă între curly braces.” → Corect: „Blocul de cod stă între acolade.”
- De ce: „Acolade” este termenul românesc corect și ajunge singur. „Paranteze acolade” se poate scurta, dar nu e pleonasm: DEX '98 definește acolada ca semn grafic în formă de arc, nu ca fel de paranteză, deci sensul nu se repetă.
- Sursă: DEX '98 s.v. acoladă

### RO-STI-815 „Paranteze unghiulare”
- Nivel: recomandare
- Greșit: „Etichetele se scriu între angle brackets.” → Corect: „Etichetele se scriu între paranteze unghiulare.”
- De ce: „Paranteze unghiulare” este denumirea românească a semnelor < și >. „Unghiular” este în DOOM 3 și în DEX '98 („care formează un unghi”); se mai spune și „paranteze ascuțite”. Alegi o denumire și o ții.
- Sursă: DOOM 3 s.v. unghiular; DEX '98 s.v. unghiular (sinonimul: nereverificat local)

## Securitate și termeni definiți

### RO-STI-810 „Autentificare” și „autorizare” sunt noțiuni diferite
- Nivel: normă
- Greșit: „Autentificarea stabilește ce acțiuni are voie să facă utilizatorul.” → Corect: „Autorizarea stabilește ce acțiuni are voie să facă utilizatorul.”
- De ce: Autentificarea verifică identitatea („a autentifica” înseamnă „a face ca un act să devină autentic”, DEX '98), iar autorizarea acordă drepturi („a autoriza” înseamnă „a da cuiva dreptul de a face ceva”, DEX '98). Cuvintele nu sunt interschimbabile, pentru că schimbarea lor schimbă ce afirmă documentația despre sistem.
- Sursă: DOOM 3 și DEX '98 s.v. autentifica, autoriza

### RO-STI-804 Hash nu este cifrare
- Nivel: normă
- Greșit: „Parolele sunt cifrate cu scrypt.” → Corect: „Parolele sunt stocate ca hash scrypt.”
- De ce: „A cifra” înseamnă „a aplica unui text un cifru” (DEX '98), operație care se poate inversa cu o cheie. Hash-ul nu se inversează, deci „cifrat” spune despre sistem altceva decât se întâmplă. „Hash” și „scrypt” sunt forme deduse (lipsesc din dicționare), de aceea la prima apariție le explici. Corectezi doar când contextul arată limpede că parola e stocată ca hash (nu poate fi recuperată, se verifică prin comparație). Dacă nu e limpede ce face sistemul, întrebi autorul.
- Sursă: DEX '98 s.v. cifra

### RO-STI-803 „Salt” în criptografie
- Nivel: recomandare
- Greșit: „Parola se amestecă cu o sare aleatorie înainte de hash.” → Corect: „Parola se amestecă cu un salt aleatoriu înainte de hash.”
- De ce: Dicționarele nu înregistrează sensul criptografic nici pentru „salt”, nici pentru „sare”; „salt” apare în DOOM 3 și în DEX doar cu sensul de săritură. „Sare” nu e greșeală de limbă, ci o traducere literală fără acoperire în dicționar, de aceea propui ca termenul să rămână „salt” (formă dedusă, ca împrumut din engleză) și se explică la prima apariție: „o valoare aleatorie adăugată parolei înainte de hash”.
- Sursă: DOOM 3 s.v. salt; DEX '98 s.v. salt, sare

## Termeni din limbajul „agentic”

### RO-STI-801 „Livrare agentică”
- Nivel: recomandare
- Greșit: „Echipa folosește livrare agentică.” → Corect: „Echipa livrează software cu ajutorul agenților.”
- De ce: „Agentic” nu e în dicționarele consultate, nici în formă flexionată. „Livrare” cere un complement: DEX '98 o leagă de mărfuri și servicii, deci „livrare agentică”, singură, e obscură și calchiată după engl. „agentic delivery”. Formularea cu „agenți” este descriptivă (formă dedusă), nu un termen consacrat; la prima apariție explici ce înseamnă.
- Sursă: DEX '98 s.v. livrare

### RO-STI-802 „Coding agentic”
- Nivel: recomandare
- Greșit: „Lucrăm în coding agentic.” → Corect: „Programăm cu ajutorul agenților.”
- De ce: „Coding” și „agentic” lipsesc amândouă din dicționare, deci „coding agentic” e o construcție hibridă cu două cuvinte englezești. „Programare cu ajutorul agenților” este o formulare descriptivă (formă dedusă); la prima apariție explici ce înseamnă.
- Sursă: DEX '98 s.v. programare

## Titluri și etichete de prezentare

### RO-STI-816 Titlurile de slide sunt nominale
- Nivel: recomandare
- Greșit: „Overview” → Corect: „Prezentare generală”
- De ce: Titlurile de slide se scriu nominal: „Overview” devine „Prezentare generală”, „Features” devine „Funcționalități”. „Funcționalitate” este corect (DOOM 3, DEX '98). Pe slide rămân puține cuvinte, nu propoziții întregi. E o convenție de redactare, fără normă.
- Sursă: DOOM 3 s.v. prezentare, funcționalitate; convenție de redactare, fără normă
