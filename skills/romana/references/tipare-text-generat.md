# Tiparele textului generat

Fișierul acesta te ajută să recunoști și să corectezi tiparele care fac un text să sune generat automat: laude fără fapte, formule care anunță ce urmează să spună textul în loc să spună, deschideri și încheieri care nu zic nimic, construcții în care un verb pompos ține locul lui „a fi”. Tiparele sunt scrise de noi, pe româna de azi, cu exemple noi.

Un lucru contează mai mult decât restul: **tiparele nu sunt greșeli de limbă.** Fiecare cuvânt și fiecare construcție de aici există în dicționare și în gramatici și, luate separat, sunt corecte. Un semnal izolat nu se raportează. Ce contează e grupul: același tic repetat, sau mai multe tipare diferite strânse într-un pasaj. De aceea toate regulile de aici sunt recomandări, iar gravitatea e de regulă P2, și se aplică doar când semnalele se adună.

Abrevierile din câmpul „Sursă”: DSL = A. Bidu-Vrănceanu et al., *Dicționar general de științe. Științe ale limbii*, 1997; Mic ghid CE = *Mic ghid de redactare clară*, Comisia Europeană, 2013; ghid gov-md = *Ghid de redactare lingvistică și tehnoredactare a actelor normative*, Cancelaria de Stat a Republicii Moldova, 2026; INJ = *Corectitudine lingvistică în actele judecătorești*, Institutul Național al Justiției, Chișinău, 2020; DOOM 3 = *Dicționarul ortografic, ortoepic și morfologic al limbii române*, ediția a III-a; DEX '98 = *Dicționarul explicativ al limbii române*, ediția a II-a, 1998. „Tipar de stil” înseamnă că tiparul nu are o normă în spate: e o observație de stil, care se semnalează doar în grup.

## Cum folosești fișierul

- **La redactare**, nu pornești de la lista asta. Pornești de la ce vrei să spui, cu faptele pe care le ai (vezi SKILL.md). Lista o folosești abia la recitirea ciornei, ca să găsești locurile unde textul sună generat.
- **În corectură**, rescrii un loc doar când tiparul face parte dintr-un grup și când poți păstra faptele sursei (vezi [fidelitate.md](fidelitate.md)). Un tipar izolat îl lași.
- **În verificare**, raportezi grupul, nu fiecare apariție: numești tiparele găsite, citezi formele exact și spui că textul „seamănă probabil cu un text generat”. Nu spui niciodată „e scris de un model”. Densitatea e un indiciu, nu o dovadă.
- **Ce pui în loc.** Un tipar se scoate înlocuind formula cu faptul pe care o ținea în loc. Dacă sursa nu are faptul, nu-l inventezi: scurtezi sau lași „[de completat: …]”. Unele corecturi din aceste reguli sunt scrise cu un asemenea semn tocmai pentru că faptul lipsește din exemplu.
- **Urmele de chatbot** („Desigur!”, „Iată textul revizuit:”, un disclaimer de model) nu sunt tipare de stil, ci resturi din conversația cu cel care a cerut textul. Se scot din textul livrat, chiar dacă apar o singură dată (RO-STI-207, RO-STI-209).

## Cum se numără

Fără statistici: nu există aici cifre despre cât de des apar tiparele în textele generate, pentru că nu avem o sursă pentru ele. Pragurile de mai jos sunt praguri de lucru ale acestui skill, nu măsurători.

- **Un grup** înseamnă același tipar de cel puțin trei ori în același text, sau cel puțin trei tipare diferite într-un singur pasaj de câteva paragrafe, din cel puțin două dintre categoriile de mai jos.
- **Vocabularul** (RO-STI-213, RO-STI-225) se numără după aceeași regulă: un „esențial” într-un text nu e semnal; trei sau patru cuvinte din listă într-un paragraf, da.
- **Expresiile-cadru** (RO-STI-214) se numără mai strict: aceeași expresie de două ori, sau trei expresii diferite.
- **Rescrierea completă** a unui text se justifică doar când tiparele sunt multe (cel puțin cinci), vin din cel puțin trei categorii și textul are și un ritm uniform, cu fraze de aceeași lungime și de aceeași construcție. Altfel corectezi punctual, locul cu locul.
- Un text fără niciunul dintre aceste tipare nu e dovada că a fost scris de om, iar un text cu ele nu e dovada contrară. Nu lași și nu introduci greșeli ca să pară textul scris de om (vezi RO-NUG-098 în [ce-nu-e-greseala.md](ce-nu-e-greseala.md)).

## Tiparele după cât țin de limbă

Unele tipare sunt specifice românei, adică se văd în felul în care româna se scrie la birou sau în presă: verbul-paravan cu substantiv verbal („permite optimizarea”), copula ocolită („reprezintă o soluție”), lanțul de genitive, „prin intermediul”, „în contextul”, gerunziul de la final de frază. Altele sunt universale și apar în orice limbă: antiteza falsă, regula lui trei, deschiderea cu întrebare retorică, importanța declarată, concluzia care repetă corpul textului. Cele din a doua categorie le recunoști mai ușor; cele din prima le trece cititorul român cu vederea, pentru că arată a limbaj oficial.

## Tabelul de recunoaștere

| Tipar | Cum arată | Ce faci | Regula |
|---|---|---|---|
| Formulă de chatbot | „Desigur!”, „Cu plăcere!” | scoți | RO-STI-207 |
| Anunțarea intenției | „Haideți să explorăm…” | intri direct | RO-STI-208 |
| Disclaimer de model | „Întrucât nu dețin informații recente…” | scoți din text, spui autorului separat | RO-STI-209 |
| Deschidere-clișeu | „În lumea de astăzi…” | prima frază spune faptul | RO-STI-229 |
| Importanță declarată | „Este important de menționat că…” | spui faptul | RO-STI-228 |
| Verb-paravan | „permite optimizarea” | verbul direct, dacă subiectul face acțiunea | RO-STI-227 |
| Copulă ocolită | „reprezintă o soluție” | „e o soluție” | RO-STI-212 |
| Gerunziu de final | „…, subliniind importanța” | tai sau faci propoziție | RO-STI-219 |
| Atribuire vagă | „Studiile arată…” | numești sursa sau renunți | RO-STI-231 |
| Prudență excesivă | „Unii ar putea considera…” | numești cine, sau afirmi ce poți susține | RO-STI-230 |
| Antiteză falsă | „Nu e doar X, ci Y” | afirmi direct Y | RO-STI-233 |
| Final nefalsificabil | „Vremuri interesante ne așteaptă.” | te oprești la ultimul fapt | RO-STI-210 |

## Importanța declarată, nu arătată

### RO-STI-228 Importanța declarată, nu demonstrată
- Nivel: recomandare
- Greșit: „Este important de menționat că termenul de livrare a fost scurtat.” → Corect: „Termenul de livrare a fost scurtat.”
- De ce: Formule ca „Merită menționat faptul că…”, „Trebuie subliniat că…”, „Este esențial să înțelegem…”, „Interesant este că…”, „o lucrare de o importanță deosebită” anunță că ceva contează, în loc să arate de ce. Scoți formula; dacă importanța e ea însăși o informație, spui de ce contează faptul. Dacă formula marca un contrast, rămâne doar conectorul („Totuși”). Nu sunt greșeli de limbă.
- Sursă: tipar de stil

### RO-STI-201 Formulele impersonale de evidențiere
- Nivel: recomandare
- Greșit: „Se poate observa că vânzările au scăzut.” → Corect: „Vânzările au scăzut.”
- De ce: Când numești formula într-o observație, o numești exact. „Se poate observa” și „se constată” sunt pasive reflexive, cu cliticul „se” și un verb. „Trebuie menționat că” e o construcție impersonală cu participiu, fără „se”, iar „este de remarcat că” e copulă plus supin. Toate sunt corecte în limbă; devin semnal doar când se adună și când nu aduc nimic.
- Sursă: DSL, s.v. PASIV; tipar de stil

### RO-STI-202 „Joacă un rol esențial”
- Nivel: recomandare
- Greșit: „Cercetarea joacă un rol esențial în procesul decizional.” → Corect: „Cercetarea [de completat: ce anume face] în procesul decizional.”
- De ce: „A juca un rol” e o expresie corectă și consacrată (DEX '98 o atestă), nu un calc. Se semnalează doar formula goală „joacă un rol esențial, crucial, important”, când nu spune ce anume face subiectul. Spui ce face: ce hotărăște, ce măsoară, ce schimbă.
- Sursă: DEX '98 s.v. rol; tipar de stil

### RO-STI-218 Importanța umflată
- Nivel: recomandare
- Greșit: „Evenimentul a marcat un punct de cotitură.” → Corect: „Evenimentul [de completat: ce s-a schimbat, cu sursa].”
- De ce: „O amprentă de neșters”, „simbolizând o moștenire de durată”, „marcând un punct de cotitură” țin locul faptului printr-o etichetă de importanță. „Moment de cotitură” e corect (DEX '98); se semnalează doar când e gol. Rămâne în citat și acolo unde fraza spune în continuare ce s-a schimbat.
- Sursă: DOOM 3; DEX '98 s.v. cotitură

### RO-STI-221 Tropii de autoritate
- Nivel: recomandare
- Greșit: „Adevărata întrebare este dacă proiectul merită continuat.” → Corect: „Merită continuat proiectul?”
- De ce: „Adevărata întrebare este…”, „ceea ce contează cu adevărat”, „în esență” pretind o profunzime pe care fraza nu o aduce. Le reformulezi direct, nu le ștergi mecanic: se păstrează când chiar opun o întrebare reală alteia sau când rezumă cu adevărat.
- Sursă: tipar de stil

### RO-STI-222 Formulele aforistice
- Nivel: recomandare
- Greșit: „Claritatea este limba respectului.” → Corect: „[de completat: ce anume face textul clar și pentru cine]”
- De ce: Formula are aerul unui adevăr general, dar nu spune nimic verificabil. Aforismul e o figură de stil legitimă (DSL), mai ales în text literar sau în eseu; nu se semnalează în citate, în texte literare sau acolo unde e ales cu intenție.
- Sursă: DSL, s.v. AFORISM; tipar de stil

### RO-STI-205 Adjective-laudă goale
- Nivel: recomandare
- Greșit: „soluție robustă” → Corect: „soluție care rezistă la [de completat: ce anume]”
- De ce: „Robust”, „cuprinzător”, „de ultimă generație” laudă fără conținut: nu spun ce rezistă, ce acoperă, ce tehnologie e. Spui proprietatea măsurabilă sau tehnologia pe nume. Sensurile tehnice rămân („statistică robustă”).
- Sursă: DOOM 3 s.v. robust, cuprinzător; DEX '98; Mic ghid CE p. 9 (Rec. 7)

### RO-STI-220 Limbaj promoțional pe vocea unui text informativ
- Nivel: recomandare
- Greșit: „Orașul se mândrește cu un muzeu.” → Corect: „Orașul are un muzeu.”
- De ce: „Se mândrește cu” (despre locuri), „cuibărit în inima…”, „vibrant”, „bogat” folosit evaluativ, „care îți taie respirația”, „o destinație de neratat”, „angajamentul față de [valoare abstractă]” sunt laude fără conținut verificabil. În profilul marketing toleranța e mai mare, dar faptele rămân cele din sursă.
- Sursă: DEX '98; tipar de stil

### RO-STI-225 Intensificatoare goale
- Nivel: recomandare
- Greșit: „Rolul acestui factor este vital, central și decisiv.” → Corect: „[de completat: ce depinde de acest factor]”
- De ce: „Fundamental”, „vital”, „central”, „decisiv”, „primordial” sunt adjective corecte (toate în DOOM 3). Devin semnal doar ca intensificator fără justificare, mai ales strânse în același pasaj. Sensurile tehnice rămân („drepturi fundamentale”, „semne vitale”, „Banca Centrală”).
- Sursă: DOOM 3; DEX '98

### RO-STI-232 Inflația retorică
- Nivel: recomandare
- Greșit: „Catedrala se ridică ca o mărturie a credinței locuitorilor.” → Corect: „Catedrala arată credința locuitorilor.”
- De ce: Pompa moderată, copiată din engleză, sună exagerat în română: „bogat patrimoniu cultural”, „se ridică ca o mărturie a”. „E o mărturie a” e corect (DEX '98: mărturie = dovadă); semnalezi epitetul de umplutură și densitatea, nu cuvântul. „Patrimoniu cultural” ca termen rămâne.
- Sursă: DOOM 3; DEX '98 s.v. mărturie; DSL, s.v. HIPERBOLĂ

### RO-STI-213 Vocabular de model, semnalat doar la densitate
- Nivel: recomandare
- Greșit: „Datele se trimit prin intermediul aplicației.” → Corect: „Datele se trimit prin aplicație.”
- De ce: Cuvinte ca „esențial”, „crucial”, „inovator”, „transformator”, „holistic”, „semnificativ” (fără cifre), „provocare” (pentru „dificultate”), „oportunitate”, „abordare”, „impact”, „prin intermediul”, „în contextul”, „din perspectiva”, „numeroase”, „divers” sunt corecte și se află în dicționare („holistic” e înregistrat în DEX '96, nu în DOOM 3; „inovativ” e în DOOM 3). Nu sunt calcuri recente (vezi RO-NUG-072 în [ce-nu-e-greseala.md](ce-nu-e-greseala.md)). Doar aglomerarea lor, fără conținut, e semnal. Când le înlocuiești, pui cifra, faptul sau cuvântul simplu („prin”, „multe”, „efect”), nu un alt cuvânt la fel de gol.
- Sursă: DOOM 3 s.v. esențial, crucial, abordare, impact, intermediu; DEX '98; DEX '96 s.v. holistic

### RO-STI-214 Expresii-cadru goale
- Nivel: recomandare
- Greșit: „În contextul actual, magazinul închide mai devreme.” → Corect: „Acum magazinul închide mai devreme.”
- De ce: „În lumea de astăzi”, „în era digitală”, „bune practici”, „provocări și oportunități”, „transformare digitală”, „dezvoltare durabilă” scoasă din context, triada „beneficii, oportunități, provocări” sau „valorificarea sinergiilor” sunt cadraje care nu adaugă informație. Le scoți când nu trimit la nimic și le păstrezi când numesc chiar conceptul (o strategie de digitalizare, un obiectiv de dezvoltare durabilă). Pragul e mai strict decât la vocabular: aceeași expresie de două ori sau trei expresii diferite.
- Sursă: Mic ghid CE p. 11 (Rec. 9, jargonul); tipar de stil

### RO-STI-224 Acumulările de adjective
- Nivel: recomandare
- Greșit: „O platformă modernă, inovatoare, versatilă și cuprinzătoare.” → Corect: „O platformă [de completat: singurul adjectiv susținut de fapte].”
- De ce: Un șir de adjective evaluative generice, pus pe vocea textului, nu spune nimic măsurabil. Păstrezi doar adjectivul pe care îl susțin faptele. Acumularea e o figură legitimă, folosită mai ales în descriere (DSL); aici se semnalează, nu se taie automat.
- Sursă: DSL, s.v. ACUMULARE; INJ (adjectivul evaluativ)

## Deschideri, încheieri și formule adresate cititorului

### RO-STI-207 Formule de chatbot la începutul răspunsului
- Nivel: recomandare
- Greșit: „Desigur! Iată descrierea: Aparatul se pornește de la buton.” → Corect: „Aparatul se pornește de la buton.”
- De ce: „Desigur!”, „Cu siguranță!”, „Absolut!”, „Cu plăcere!”, „Mă bucur să vă ajut.” sunt replici adresate celui care a cerut textul, nu conținut. Se scot din textul livrat, chiar dacă apar o singură dată. Ca adverbe în frază („absolut necesar”) rămân.
- Sursă: tipar de stil; DOOM 3 s.v. desigur, absolut

### RO-STI-208 Anunțarea intenției
- Nivel: recomandare
- Greșit: „Haideți să explorăm împreună cum funcționează rețeaua.” → Corect: „[prima frază cu conținut despre rețea]”
- De ce: „Permiteți-mi să vă ghidez prin…”, „Să analizăm pas cu pas…”, „Pentru a răspunde la întrebarea dumneavoastră…”, „În cele ce urmează vom prezenta…” anunță ce urmează să facă textul în loc să facă. Nu sunt greșeli: „haide, haideți” și „a explora” sunt corecte, iar conjunctivul de îndemn e o construcție românească (vezi RO-NUG-114 în [ce-nu-e-greseala.md](ce-nu-e-greseala.md)). Problema e doar funcția de anunț, deci intri direct în subiect.
- Sursă: DOOM 3 s.v. haide, explora; tipar de stil

### RO-STI-209 Disclaimerul de model în textul livrat
- Nivel: recomandare
- Greșit: „Întrucât nu dețin informații recente, textul se oprește la anul 2023.” → Corect: „[de completat: cifrele din anul curent]”
- De ce: Un disclaimer de model lingvistic nu are ce căuta într-un text publicat. Dar nu doar îl ștergi: dacă limita e reală (lipsesc date actuale), autorul trebuie să afle, într-o notă separată, în afara textului, și locul rămâne marcat cu „[de completat: …]”. Ștergerea simplă ar ascunde o lacună.
- Sursă: tipar de stil; regulă de fidelitate (vezi [fidelitate.md](fidelitate.md))

### RO-STI-229 Deschideri-clișeu
- Nivel: recomandare
- Greșit: „Trăim într-o epocă în care informația circulă repede.” → Corect: „[prima frază spune faptul despre care e textul]”
- De ce: „În lumea de astăzi…”, „Trăim într-o epocă în care…”, „Este bine cunoscut faptul că…”, „Fie că ești X sau Y…” amână subiectul. Prima frază spune faptul. Dacă faptul nu e de notorietate, îl atribui; dacă X și Y primesc de fapt aceleași informații, formula „fie că” cade.
- Sursă: tipar de stil; DSL, s.v. CLIȘEU (nereverificat local)

### RO-STI-223 Deschiderea fals sinceră
- Nivel: recomandare
- Greșit: „Verdictul? Da, merită.” → Corect: „Da, merită.”
- De ce: Întrebarea cu răspuns imediat („Sincer? Depinde.”, „Pe scurt? Da.”) e un tic conversațional fără loc în textul informativ. Nu se semnalează automat orice „sincer” sau „sincer vorbind” folosit ca adverb de atitudine al autorului: ghidul INJ îl folosește în propriul text.
- Sursă: INJ; tipar de stil

### RO-STI-216 Deschiderea cu întrebare retorică de șablon
- Nivel: recomandare
- Greșit: „Ce ar fi dacă facturile s-ar plăti singure? Aplicația plătește facturile automat, la scadență.” → Corect: „Aplicația plătește facturile automat, la scadență.”
- De ce: Întrebarea retorică de șablon, care amână teza, se semnalează și se rescrie pornind direct de la teză. Interogația retorică e însă o figură de stil recunoscută (DSL), deci nu se semnalează în corpul unui text de opinie și nici când citezi o întrebare a sursei.
- Sursă: DSL, s.v. INTEROGAȚIE; tipar de stil

### RO-STI-210 Finalul nefalsificabil sau futurologic
- Nivel: recomandare
- Greșit: „Vremuri interesante ne așteaptă.” → Corect: „[ultimul fapt din sursă sau pasul următor concret]”
- De ce: „Suntem la începutul unui nou capitol”, „ar putea deveni unul dintre cele mai importante subiecte ale deceniului” nu spun nimic verificabil. Dacă sursa are un pas următor concret, acela încheie textul; altfel textul se oprește la ultimul fapt.
- Sursă: tipar de stil

### RO-STI-211 Formule goale de rezumat și de concluzie
- Nivel: recomandare
- Greșit: „Se poate concluziona că programul a redus abandonul școlar.” → Corect: „Programul a redus abandonul școlar.”
- De ce: „Se poate concluziona că X”, „În ansamblu, se poate afirma că X”, „Rezumând cele de mai sus…” amână afirmația. Concluzia se spune direct și se atribuie dacă nu e a autorului. Problema e de stil, nu de vocabular: „a concluziona” e corect (DOOM 3). Rezumatul se taie întreg dacă nu aduce nimic nou.
- Sursă: DOOM 3 s.v. concluziona; DEX '98

## Construcții

### RO-STI-212 Copula ocolită și perifrazele verbale
- Nivel: recomandare
- Greșit: „Platforma reprezintă o soluție.” → Corect: „Platforma e o soluție.”
- De ce: Verbul pompos ține locul lui „a fi”: „reprezintă o soluție”, „constituie un element”, „vine în întâmpinarea nevoilor” (în loc de „răspunde nevoilor”). Nu e greșeală de limbă (DEX '98 atestă sensul „a fi”) și nu se înlocuiește unde are sens propriu: „a reprezenta pe cineva”, „reprezintă 60% din total”, „constituie infracțiune” în drept.
- Sursă: DEX '98 s.v. reprezenta, constitui; Mic ghid CE p. 8 (Rec. 6)

### RO-STI-227 Verbele-paravan
- Nivel: recomandare
- Greșit: „Programul permite optimizarea procesului.” → Corect: „Programul optimizează procesul.”
- De ce: Un verb gol plus un substantiv verbal ascunde acțiunea: „permite optimizarea”, „ajută la îmbunătățirea”, „contribuie la creșterea”, „facilitează înțelegerea”. Înlocuirea („optimizează”, „îmbunătățește”, „mărește”, „explică”) se face doar dacă subiectul chiar face acțiunea; „permite” spune că face posibil, nu că face. Altfel înlocuirea umflă afirmația și încalcă fidelitatea (vezi [fidelitate.md](fidelitate.md)).
- Sursă: Mic ghid CE p. 8 (Rec. 6); tipar de stil

### RO-STI-219 Gerunziul de la final de frază
- Nivel: recomandare
- Greșit: „Compania a lansat aplicația, subliniind importanța digitalizării.” → Corect: „Compania a lansat aplicația.”
- De ce: Gerunziul de la final de frază („…, subliniind / evidențiind / reflectând…”) e un tic de text generat. „A sublinia” și „a prezenta” sunt verbe uzuale și nu se semnalează ca vocabular izolat, ci doar în tipare: acest gerunziu, comentariul despre propriul text („merită subliniat”) și umflarea importanței („subliniază importanța sa”). Dacă gerunziul spune ceva care e în sursă, îl faci propoziție de sine stătătoare.
- Sursă: DSL, s.v. GERUNZIU; DEX '98 s.v. sublinia

### RO-STI-235 Pseudo-explicația
- Nivel: recomandare
- Greșit: „Livrarea întârzie. Motivul este că depozitul e închis.” → Corect: „Livrarea întârzie pentru că depozitul e închis.”
- De ce: Formulele „Motivul este că…” și „Acest lucru se datorează faptului că…” puse în fața unei simple afirmații aduc o ramă cauzală în plus. Nu e greșeală („a se datora” e în DOOM 3), ci prolixitate. Când relația cauzală e reală, o păstrezi, dar o comprimi: „X, pentru că Y”.
- Sursă: DOOM 3 s.v. datora; tipar de stil

### RO-STI-233 Antiteza falsă „nu X, ci Y”
- Nivel: recomandare
- Greșit: „Nu e doar o aplicație, ci o filozofie.” → Corect: „Aplicația [de completat: ce este, în termeni concreți].”
- De ce: „Nu…, ci” e construcția românească firească pentru corectarea unei presupuneri reale („Nu e un lux, e o necesitate” rămâne când cititorul chiar crede că e un lux). E tic doar când nimeni nu credea X și X e inventat ca să fie respins. Testul „taie negația” aplicat orbește pierde sens.
- Sursă: DSL, s.v. ANTITEZĂ; tipar de stil

### RO-STI-234 „Nu doar…, ci și…” pentru două calități obișnuite
- Nivel: recomandare
- Greșit: „Nu doar rapid, ci și simplu.” → Corect: „Rapid și simplu.”
- De ce: Corelativul „nu doar / nu numai…, ci și” e normat (ghid gov-md îl folosește și îi fixează virgula) și rămâne când al doilea termen contrazice o așteptare a cititorului sau marchează o gradare reală. Pentru două calități obișnuite, fără gradare, o coordonare simplă cu „și” spune același lucru, mai scurt.
- Sursă: ghid gov-md; tipar de stil

### RO-STI-217 Concesia falsă
- Nivel: recomandare
- Greșit: „Deși există unele provocări, proiectul este un mare succes.” → Corect: „Proiectul a fost finalizat, dar [de completat: compromisul real].”
- De ce: O rezervă generică („deși există unele provocări”), urmată imediat de o concluzie pozitivă care o anulează. Concesia în sine e o construcție gramaticală (DSL, s.v. CONCESIV); tiparul e lipsa conținutului. Numești concret compromisul real sau renunți la rezervă.
- Sursă: DSL, s.v. CONCESIV; tipar de stil

### RO-STI-230 Prudența excesivă
- Nivel: recomandare
- Greșit: „Unii ar putea considera că proiectul e riscant.” → Corect: „[de completat: cine consideră] că proiectul e riscant.”
- De ce: „Ar putea fi argumentat că…”, „Unii ar putea considera…”, „În anumite contexte…” (fără să se spună care) nu arată grad de certitudine, ci evitare. Gradul de certitudine din sursă se păstrează (vezi RO-FID-021 în [fidelitate.md](fidelitate.md)); scoți doar ambalajul. Nu transformi atenuarea în afirmație categorică. „Ar putea” rămâne corect când sursa e numită (RO-NUG-070).
- Sursă: tipar de stil

### RO-STI-231 Atribuirea vagă
- Nivel: recomandare
- Greșit: „Studiile arată că munca de acasă crește productivitatea.” → Corect: „[de completat: studiul și anul] arată că munca de acasă crește productivitatea.”
- De ce: „Studiile arată…”, „Experții consideră…”, „literatura de specialitate sugerează”, „revistele consacrate arată” ascund sursa. Numești sursa sau renunți la afirmație. Nu ștergi doar atribuirea, lăsând afirmația fără sprijin, și nu inventezi sursa: dacă lipsește, semnalezi golul.
- Sursă: tipar de stil

### RO-STI-204 „Se recomandă” fără agent
- Nivel: recomandare
- Greșit: „Se recomandă actualizarea.” → Corect: „Manualul producătorului recomandă actualizarea.”
- De ce: În exemplu, sursa din care vine fraza este manualul producătorului, deci recomandantul se poate numi. Numești recomandantul doar dacă îl găsești în text sau în sursă („raportul recomandă”, „medicii recomandă”). Într-o știre, dacă nu-l găsești, semnalezi lipsa atribuirii și nu-l inventezi. Pasivul reflexiv rămâne corect când agentul e evident (vezi RO-NUG-085).
- Sursă: Mic ghid CE p. 10 (Rec. 8); tipar de stil

### RO-STI-203 „Iată chestia:”
- Nivel: recomandare
- Greșit: „Iată chestia: stocul e aproape epuizat.” → Corect: „Problema e că stocul e aproape epuizat.”
- De ce: Formula traduce cuvânt cu cuvânt engl. „Here's the thing:” și coboară registrul („chestie” e familiar, vezi RO-STI-428 în [registru-si-adresare.md](registru-si-adresare.md)). Intri direct în subiect sau folosești o tranziție care spune ceva: „De fapt, …”, „Problema e că…”.
- Sursă: DOOM 3 s.v. chestie; DEX '98 s.v. chestiune

## Ritm și structură

### RO-STI-206 Textul nu se taie în staccato
- Nivel: recomandare
- Greșit: „Lansarea a avut loc ieri. Mulți clienți. Fără probleme.” → Corect: „Lansarea a avut loc ieri, cu mulți clienți și fără probleme.”
- De ce: Frazele obișnuite nu se transformă în propoziții scurte și dramatice doar pentru ritm. Lungimea o variezi prin conținutul frazelor, nu prin tăiere, și la împărțire păstrezi conectorii (Mic ghid CE cere, când împarți o frază lungă, conectorii adecvați, ca documentul să rămână coerent). Se semnalează doar când apar în serie, fără funcție expresivă.
- Sursă: Mic ghid CE p. 6 (Rec. 4); tipar de stil

### RO-STI-215 Listele fără verb
- Nivel: recomandare
- Greșit: „• Incluziune socială • Acces la educație • Parteneriate durabile” → Corect: „Proiectul [de completat: ce face, pentru cine, cu ce rezultat].”
- De ce: Se semnalează listele de puncte formate doar din substantive abstracte sau sloganuri, fără verb, cifră sau referent concret. Nu se semnalează listele fără verb cu elemente concrete (nume, instituții, date, sume, pași, categorii definite), care sunt normale într-un act normativ, într-o casetă de fapte sau într-un program.
- Sursă: tipar de stil; Mic ghid CE p. 9 (Rec. 7)

### RO-STI-226 Catalog de semnale care se raportează doar în grup
- Nivel: recomandare
- Greșit: „De la ghid până la formular, site-ul are totul.” → Corect: „Site-ul [de completat: ce conține, concret].”
- De ce: Se semnalează, ca indicii și nu ca greșeli de limbă: secțiunile formulaice („Provocări și perspective”, „În ciuda acestor provocări”); gruparea forțată în câte trei; „de la X la Y” când X și Y nu formează o scală reală; boldul pus mecanic în mijlocul frazelor; listele cu „Titlu:” la început; emoji în texte neinformale; atribuirile vagi; importanța sau notorietatea umflate fără date. Intervalele false și rotirea sinonimelor intră doar în grup. Fiecare caz se judecă în context.
- Sursă: tipar de stil; ENUMERARE în DSL

### RO-STI-236 Regula densității
- Nivel: recomandare
- Greșit: „Textul conține un singur «esențial» și este raportat ca text generat.” → Corect: „Un singur «esențial» nu se raportează; se raportează grupul.”
- De ce: Detecția caută grupuri de semnale, adică densitatea unui tic, nu apariția lui izolată. Regula se aplică și la conectori și la adjectivele-semnal: „esențial”, „crucial”, „semnificativ” sunt cuvinte literare din DEX și DOOM 3. Chiar peste prag, verdictul spune „probabil”, nu „e scris de un model”.
- Sursă: DOOM 3; DEX '98; tipar de stil

## Vezi și

- [coerenta.md](coerenta.md): aglomerarea de conectori (RO-COE-016) și testul tranziției.
- [stilistica.md](stilistica.md): limba de lemn, densitatea substantivelor abstracte, politețea tradusă.
- [fidelitate.md](fidelitate.md): ce nu ai voie să adaugi când scoți un tipar.
- [registru-si-adresare.md](registru-si-adresare.md): registrul familiar care nu face un text să pară scris de om.
- [ce-nu-e-greseala.md](ce-nu-e-greseala.md): cuvintele și construcțiile corecte pe care nu le numești tipare.
