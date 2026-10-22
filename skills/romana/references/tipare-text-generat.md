# Tiparele textului generat

Fișierul acesta te ajută să recunoști și să corectezi tiparele care fac un text să sune generat automat: laude fără fapte, replici care anunță ce urmează să spună textul în loc să spună, deschideri și încheieri care nu zic nimic, verbe pompoase puse în locul lui „a fi”. Descrierile și exemplele sunt scrise pentru acest skill, pe subiecte românești obișnuite (un comunicat de primărie, o descriere de produs, un articol de blog).

Un lucru contează mai mult decât restul: **tiparele nu sunt greșeli de limbă.** Fiecare cuvânt și fiecare construcție de aici există în dicționare și în gramatici și, luate separat, sunt corecte. Un semnal izolat nu se raportează. Ce contează e grupul: același tic repetat, sau mai multe tipare diferite strânse într-un pasaj. De aceea toate regulile de aici sunt recomandări, gravitatea e de regulă P2, și se aplică doar când semnalele se adună; excepție fac urmele conversației (RO-STI-207, RO-STI-209), care se scot și când apar o singură dată.

Abrevierile din câmpul „Sursă”: DSL = A. Bidu-Vrănceanu et al., *Dicționar general de științe. Științe ale limbii*, 1997; Mic ghid CE = *Mic ghid de redactare clară*, Comisia Europeană, 2013; ghid gov-md = *Ghid de redactare lingvistică și tehnoredactare a actelor normative*, Cancelaria de Stat a Republicii Moldova, 2026; INJ = *Corectitudine lingvistică în actele judecătorești*, Institutul Național al Justiției, Chișinău, 2020; DOOM 3 = *Dicționarul ortografic, ortoepic și morfologic al limbii române*, ediția a III-a; DEX ’98 = *Dicționarul explicativ al limbii române*, ediția a II-a, 1998. „Tipar de stil” înseamnă că tiparul nu are o normă în spate: e o observație de stil, care se semnalează doar în grup.

## Cum folosești fișierul

- **La redactare**, nu pornești de la lista asta. Pornești de la ce vrei să spui, cu faptele pe care le ai (vezi SKILL.md). Lista o folosești abia la recitirea ciornei, ca să găsești locurile unde textul sună generat.
- **În corectură**, propui rescrierea unui loc, separat de text, doar când tiparul face parte dintr-un grup și când poți păstra faptele sursei (vezi [fidelitate.md](fidelitate.md)). Un tipar izolat îl lași.
- **În verificare**, raportezi grupul, nu fiecare apariție: numești tiparele găsite, citezi formele exact și spui că textul „seamănă probabil cu un text generat”. Nu spui niciodată „e scris de un model”. Densitatea e un indiciu, nu o dovadă.
- **Ce pui în loc.** Un tipar se scoate înlocuind formula cu faptul pe care o ținea în loc. Dacă sursa nu are faptul, nu-l inventezi: scurtezi sau lași „[de completat: …]”. Unele corecturi din regulile de mai jos sunt scrise cu un asemenea semn tocmai pentru că faptul lipsește din exemplu.
- **Urmele conversației** (replicile adresate celui care a cerut textul, un disclaimer de model) nu sunt tipare de stil, ci resturi din dialogul cu autorul. Se scot din textul livrat, chiar dacă apar o singură dată (RO-STI-207, RO-STI-209).

## Cum se numără

Fără statistici: nu avem o sursă pentru cât de des apar tiparele în textele generate, așa că fișierul nu dă cifre de acest fel. Pragurile de mai jos sunt orientative și servesc la a decide când merită să intervii.

- **Un grup** înseamnă același tipar de cel puțin trei ori în același text, sau cel puțin trei tipare diferite într-un singur pasaj de câteva paragrafe, din cel puțin două dintre categoriile de mai jos.
- **Vocabularul** (RO-STI-213, RO-STI-225) se numără după aceeași regulă: un cuvânt mare într-un text nu e semnal; câteva, strânse într-un paragraf, da.
- **Expresiile-cadru** (RO-STI-214) se numără mai strict: aceeași expresie de două ori, sau trei expresii diferite.
- **Rescrierea completă** a unui text se justifică doar când aproape fiecare paragraf are câte un tipar și, pe deasupra, ritmul e uniform: fraze de aceeași lungime, construite la fel. Altfel corectezi punctual, locul cu locul.
- Un text fără niciunul dintre aceste tipare nu e dovada că a fost scris de om, iar un text cu ele nu e dovada contrară. Nu lași și nu introduci greșeli ca să pară textul scris de om (vezi RO-NUG-098 în [ce-nu-e-greseala-stil.md](ce-nu-e-greseala-stil.md)).

## Tiparele după cât țin de limbă

Unele tipare se văd mai ales în româna de birou sau de presă: verbul-paravan cu substantiv verbal, copula ocolită, lanțul de genitive, locuțiunile lungi în loc de o prepoziție. Altele apar în orice limbă: antiteza falsă, deschiderea cu întrebare retorică, lauda fără fapte, finalul care nu afirmă nimic. Cele din a doua categorie le recunoști mai ușor; cele din prima le trece cititorul român cu vederea, pentru că arată a limbaj oficial.

## Tabelul de recunoaștere

| Tipar | Cum arată | Ce faci | Regula |
|---|---|---|---|
| Replică adresată autorului | „Iată varianta cerută:” | scoți | RO-STI-207 |
| Disclaimer de model | „Ca model de limbaj, nu am acces…” | scoți din text, spui autorului separat | RO-STI-209 |
| Anunțarea intenției | „În rândurile următoare vă prezentăm…” | intri direct | RO-STI-208 |
| Copulă ocolită | „se prezintă ca o soluție” | „e o soluție” | RO-STI-212 |
| Verb-paravan | „asigură realizarea instalării” | verbul direct, dacă subiectul face acțiunea | RO-STI-227 |
| Atribuire vagă | „Specialiștii recomandă…” | numești sursa sau renunți | RO-STI-231 |
| Importanță declarată | „Nu trebuie uitat că…” | spui faptul | RO-STI-228 |
| Antiteză falsă | „Nu e doar o băutură, ci un ritual” | spui ce este | RO-STI-233 |
| Final fără afirmație | „Rămâneți aproape.” | te oprești la ultimul fapt | RO-STI-210 |

## Urmele conversației

### RO-STI-207 Replici adresate celui care a cerut textul
- Nivel: recomandare
- Greșit: „Sigur, iată varianta cerută: Cantina primăriei se închide vineri la ora 14.” → Corect: „Cantina primăriei se închide vineri la ora 14.”
- De ce: Replici ca „Iată varianta cerută”, „Sper că ți-a fost de folos!” sau „Spune-mi dacă vrei alt ton” se adresează celui care a comandat textul, nu cititorului lui. Se scot din textul livrat, chiar dacă apar o singură dată. Ca parte a unei fraze obișnuite („un răspuns sigur”) cuvintele rămân.
- Sursă: tipar de stil

### RO-STI-209 Disclaimerul de model în textul livrat
- Nivel: recomandare
- Greșit: „Ca model de limbaj, nu am acces la tarifele de azi, deci cifrele de mai jos pot fi depășite.” → Corect: „[de completat: tarifele în vigoare]”
- De ce: Un disclaimer de model nu are ce căuta într-un text publicat. Dar nu doar îl ștergi: dacă limita e reală (lipsesc date actuale), autorul trebuie să afle, într-o notă separată, în afara textului, iar locul rămâne marcat cu „[de completat: …]”. Ștergerea simplă ar ascunde o lacună.
- Sursă: tipar de stil; regulă de fidelitate (vezi [fidelitate.md](fidelitate.md))

### RO-STI-208 Anunțarea intenției
- Nivel: recomandare
- Greșit: „În rândurile următoare vă prezentăm principalele lucrări începute de primărie în această toamnă.” → Corect: „Primăria a început în această toamnă [de completat: lucrările, începând cu prima].”
- De ce: Frazele de felul „Să vedem pe rând cum stau lucrurile”, „Înainte de a intra în detalii, să clarificăm contextul” anunță ce va face textul în loc să facă. Nu sunt greșeli: conjunctivul de îndemn e o construcție românească firească (vezi RO-NUG-114 în [ce-nu-e-greseala.md](ce-nu-e-greseala.md)). Problema e doar funcția de anunț, deci intri direct în subiect. Un cuprins sau o frază de orientare într-un text foarte lung e altceva și rămâne.
- Sursă: tipar de stil

## Construcții

### RO-STI-212 Copula ocolită
- Nivel: recomandare
- Greșit: „Brânza de capră reprezintă un aliment bogat în calciu.” → Corect: „Brânza de capră e un aliment bogat în calciu.”
- De ce: Un verb mai solemn ține locul lui „a fi” sau al unui verb simplu: „reprezintă”, „se prezintă ca”, „se constituie ca”, „vine în sprijinul” (în loc de „ajută”). Nu e greșeală de limbă (DEX ’98 atestă sensul „a fi”) și nu se înlocuiește unde are sens propriu: „a reprezenta pe cineva”, „reprezintă 60% din total”, „constituie infracțiune” în drept.
- Sursă: DEX ’98 s.v. reprezenta, constitui; Mic ghid CE p. 8 (Rec. 6)

### RO-STI-227 Verbele-paravan
- Nivel: recomandare
- Greșit: „Tehnicianul asigură realizarea instalării.” → Corect: „Tehnicianul face instalarea.”
- De ce: Un verb gol urmat de un substantiv verbal ascunde acțiunea: „asigură realizarea”, „contribuie la derularea”, „ajută la îmbunătățirea”. Înlocuirea cu verbul direct se face doar dacă subiectul chiar face acțiunea. „Permite” nu înseamnă „face”: „Programul permite exportul datelor” nu e același lucru cu „Programul exportă datele”, deci acolo verbul rămâne. Altfel înlocuirea umflă afirmația și încalcă fidelitatea (vezi RO-FID-021 în [fidelitate.md](fidelitate.md)).
- Sursă: Mic ghid CE p. 8 (Rec. 6); tipar de stil

### RO-STI-219 Gerunziul de la final de frază
- Nivel: recomandare
- Greșit: „Muzeul a primit 5.000 de vizitatori, demonstrând astfel interesul publicului.” → Corect: „Muzeul a primit 5.000 de vizitatori.”
- De ce: Gerunziul pus la sfârșitul frazei ca să adauge o concluzie („…, demonstrând astfel…”, „…, reflectând o tendință”) comentează propriul text în loc să aducă un fapt. Verbele în sine sunt uzuale și nu se semnalează izolat; se semnalează tiparul. Dacă gerunziul spune ceva care e în sursă, îl faci propoziție de sine stătătoare.
- Sursă: DSL, s.v. GERUNZIU; tipar de stil

### RO-STI-235 Pseudo-explicația
- Nivel: recomandare
- Greșit: „Livrarea întârzie. Explicația este simplă: depozitul e închis.” → Corect: „Livrarea întârzie pentru că depozitul e închis.”
- De ce: Formule ca „Explicația este simplă:”, „Totul se reduce la faptul că” sau „Acest lucru se datorează faptului că” pun o ramă cauzală în fața unei simple afirmații. Nu e greșeală („a se datora” e în DOOM 3), ci prolixitate. Când relația cauzală e reală, o păstrezi, dar o comprimi: „X, pentru că Y”.
- Sursă: DOOM 3 s.v. datora; tipar de stil

### RO-STI-233 Antiteza falsă „nu X, ci Y”
- Nivel: recomandare
- Greșit: „Cafeaua noastră nu e doar o băutură, ci un ritual.” → Corect: „Cafeaua noastră [de completat: ce o deosebește, concret].”
- De ce: „Nu…, ci” e construcția românească firească pentru corectarea unei presupuneri reale („Nu e un lux, e o necesitate” rămâne când cititorul chiar crede că e un lux). E tic doar când nimeni nu credea X și X e inventat ca să fie respins. Testul „taie negația” aplicat orbește pierde sens.
- Sursă: DSL, s.v. ANTITEZĂ; tipar de stil

### RO-STI-234 „Nu doar…, ci și…” pentru două calități obișnuite
- Nivel: recomandare
- Greșit: „Un ghid nu doar scurt, ci și clar.” → Corect: „Un ghid scurt și clar.”
- De ce: Corelativul „nu doar / nu numai…, ci și” e normat (ghid gov-md îl folosește și îi fixează virgula) și rămâne când al doilea termen contrazice o așteptare a cititorului sau marchează o gradare reală. Pentru două calități obișnuite, fără gradare, o coordonare simplă cu „și” spune același lucru, mai scurt.
- Sursă: ghid gov-md; tipar de stil

### RO-STI-217 Concesia falsă
- Nivel: recomandare
- Greșit: „Chiar dacă drumul nu a fost lipsit de obstacole, rezultatele sunt remarcabile.” → Corect: „Rezultatele sunt remarcabile, deși [de completat: ce a mers prost, concret].”
- De ce: O rezervă generică, urmată imediat de o concluzie pozitivă care o anulează, nu spune nimic despre obstacole. Concesia în sine e o construcție gramaticală (DSL, s.v. CONCESIV); tiparul e lipsa conținutului. Numești concret compromisul real sau renunți la rezervă.
- Sursă: DSL, s.v. CONCESIV; tipar de stil

### RO-STI-230 Prudența excesivă
- Nivel: recomandare
- Greșit: „S-ar putea spune că, într-o oarecare măsură, prețurile au crescut.” → Corect: „Prețurile au crescut [de completat: cu cât, la ce produse, după ce sursă].”
- De ce: Atenuările puse în lanț („s-ar putea spune”, „într-o oarecare măsură”, „în anumite contexte”, fără să se spună care) nu arată grad de certitudine, ci evitare. Gradul de certitudine din sursă se păstrează (vezi RO-FID-021 în [fidelitate.md](fidelitate.md)); scoți doar ambalajul, nu transformi atenuarea în afirmație categorică. „Ar putea” rămâne corect când sursa e numită (RO-NUG-070).
- Sursă: tipar de stil

### RO-STI-231 Atribuirea vagă
- Nivel: recomandare
- Greșit: „Specialiștii recomandă copii de rezervă lunare.” → Corect: „[de completat: cine recomandă] copii de rezervă lunare.”
- De ce: „Specialiștii recomandă…”, „Se știe că…”, „Cercetările recente au arătat…”, „Mulți utilizatori spun…” ascund sursa. Numești sursa sau renunți la afirmație. Nu ștergi doar atribuirea, lăsând afirmația fără sprijin, și nu inventezi sursa: dacă lipsește, semnalezi golul.
- Sursă: tipar de stil

### RO-STI-204 „Se recomandă” fără agent
- Nivel: recomandare
- Greșit: „Se recomandă actualizarea.” → Corect: „[de completat: cine recomandă] recomandă actualizarea.”
- De ce: Numești recomandantul doar dacă îl găsești în text sau în sursă (de pildă, textul vine dintr-un manual al producătorului și poți scrie „Producătorul recomandă actualizarea”). Dacă nu-l găsești, semnalezi lipsa atribuirii și nu-l inventezi. Pasivul reflexiv rămâne corect când agentul e evident (vezi RO-NUG-085).
- Sursă: Mic ghid CE p. 10 (Rec. 8); tipar de stil

### RO-STI-203 „Uite care e treaba:”
- Nivel: recomandare
- Greșit: „Uite care e treaba: stocul e aproape epuizat.” → Corect: „Problema e că stocul e aproape epuizat.”
- De ce: Formula traduce o introducere conversațională din engleza de afaceri și coboară registrul („uite care e treaba” e o formulă conversațională, iar „chestia” e familiară; vezi RO-STI-428 în [registru-si-adresare.md](registru-si-adresare.md)). Intri direct în subiect sau folosești o tranziție care spune ceva: „De fapt, …”, „Problema e că…”.
- Sursă: DEX ’98 s.v. chestiune

## Vocabular și laudă

### RO-STI-213 Vocabular de efect, folosit în serie
- Nivel: recomandare
- Greșit: „Aplicația oferă o experiență fluidă într-un ecosistem digital integrat.” → Corect: „Aplicația [de completat: ce face, pe ce dispozitive].”
- De ce: Cuvinte ca „ecosistem”, „experiență fluidă”, „a debloca potențialul” sunt corecte, dar puse în serie, fără fapte, ocupă locul informației. Nici cuvintele mai vechi nu sunt calcuri recente (vezi RO-NUG-072 în [ce-nu-e-greseala.md](ce-nu-e-greseala.md)); doar aglomerarea lor e semnal. Când le înlocuiești, pui cifra, faptul sau cuvântul simplu, nu un alt cuvânt la fel de gol. Fiecare cuvânt pe care îl declari semnal îl verifici întâi în dicționar.
- Sursă: tipar de stil; Mic ghid CE p. 9 (Rec. 7)

### RO-STI-225 Intensificatoare goale
- Nivel: recomandare
- Greșit: „Un aspect capital, major și determinant pentru buna funcționare a orașului.” → Corect: „[de completat: ce depinde de acest aspect]”
- De ce: „Capital”, „major”, „determinant” (la fel „fundamental”, „vital”) sunt adjective corecte. Devin semnal doar ca intensificator fără justificare, mai ales strânse în același pasaj. Sensurile tehnice rămân („drepturi fundamentale”, „semne vitale”, „Banca Centrală”).
- Sursă: DOOM 3 s.v. capital, major, determinant; DEX ’98

### RO-STI-205 Adjective-laudă goale
- Nivel: recomandare
- Greșit: „soluție robustă” → Corect: „soluție care rezistă la [de completat: ce anume]”
- De ce: „Robust”, „cuprinzător”, „de ultimă generație” sunt laude fără conținut: nu spun ce rezistă, ce acoperă, ce tehnologie e. Spui proprietatea măsurabilă sau tehnologia pe nume. Sensurile tehnice rămân („statistică robustă”).
- Sursă: DOOM 3 s.v. robust, cuprinzător; DEX ’98; Mic ghid CE p. 9 (Rec. 7)

### RO-STI-220 Limbaj de reclamă pe vocea unui text informativ
- Nivel: recomandare
- Greșit: „Pensiunea te îmbie cu priveliști de vis și o atmosferă de poveste.” → Corect: „Din camerele pensiunii se vede [de completat: ce se vede].”
- De ce: „Te îmbie cu”, „de vis”, „de poveste”, „o experiență de neuitat” sunt laude fără conținut verificabil, puse pe vocea unui text care ar trebui să informeze. În profilul marketing toleranța e mai mare, dar faptele rămân cele din sursă.
- Sursă: tipar de stil; DEX ’98

### RO-STI-232 Inflația retorică
- Nivel: recomandare
- Greșit: „Piața centrală e o simfonie de culori și arome.” → Corect: „În piața centrală se vând [de completat: ce se vinde].”
- De ce: Imaginea pompoasă ține locul unei descrieri. Metafora e o figură legitimă într-un text literar sau într-un reportaj asumat; într-un text informativ, imaginea pompoasă sună artificială. Semnalezi epitetul de umplutură și densitatea lui, nu cuvintele.
- Sursă: DSL, s.v. HIPERBOLĂ; tipar de stil

### RO-STI-224 Acumulările de adjective
- Nivel: recomandare
- Greșit: „Un telefon sigur, rapid, intuitiv și elegant.” → Corect: „Un telefon [de completat: singurul adjectiv susținut de fapte].”
- De ce: Un șir de adjective evaluative generice, pus pe vocea textului, nu spune nimic măsurabil. Păstrezi doar adjectivul pe care îl susțin faptele. Acumularea e o figură legitimă, folosită mai ales în descriere (DSL); aici se semnalează, nu se taie automat.
- Sursă: DSL, s.v. ACUMULARE; tipar de stil

### RO-STI-214 Expresii-cadru goale
- Nivel: recomandare
- Greșit: „La ora actuală, pe fondul evoluțiilor din domeniu, cofetăria a mutat comenzile online.” → Corect: „Cofetăria a mutat comenzile online.”
- De ce: „La ora actuală”, „pe fondul evoluțiilor din domeniu”, „în vremurile pe care le trăim” sunt cadraje care nu adaugă informație. Le scoți când nu trimit la nimic și le păstrezi când numesc chiar conceptul (o strategie de digitalizare, un plan de mobilitate urbană). Pragul e mai strict decât la vocabular: aceeași expresie de două ori sau trei expresii diferite.
- Sursă: Mic ghid CE p. 11 (Rec. 9, jargonul); tipar de stil

## Importanță declarată

### RO-STI-228 Importanța declarată, nu demonstrată
- Nivel: recomandare
- Greșit: „Un aspect de care merită să ținem cont este că parcarea se închide la ora 22.” → Corect: „Parcarea se închide la ora 22.”
- De ce: Formule ca „Nu trebuie uitat că…”, „Vom vedea că aceasta contează enorm” anunță că ceva contează, în loc să arate de ce. Scoți formula; dacă importanța e ea însăși o informație, spui de ce contează faptul. Dacă formula marca un contrast, rămâne doar conectorul („Totuși”). Nu sunt greșeli de limbă.
- Sursă: tipar de stil

### RO-STI-201 Formulele impersonale de evidențiere
- Nivel: recomandare
- Greșit: „Se cuvine precizat că ședința a fost amânată.” → Corect: „Ședința a fost amânată.”
- De ce: Când numești formula într-o observație, o numești exact. „Se cuvine precizat” are cliticul „se” și o construcție impersonală; „vrednic de reținut este că” e adjectiv plus supin, fără „se”. Niciuna nu e pasiv reflexiv în sens strict: în pasivul reflexiv subiectul real e obiectul acțiunii („Cartea se citește cu plăcere”, DSL). Toate sunt corecte în limbă; devin semnal doar când se adună și când nu aduc nimic. Formula rămâne dacă introduce o constatare asumată, cu sursă.
- Sursă: DSL, s.v. PASIV; tipar de stil

### RO-STI-202 „Joacă un rol vital”
- Nivel: recomandare
- Greșit: „Apa joacă un rol vital în funcționarea unui oraș.” → Corect: „Apa [de completat: ce anume face] în funcționarea unui oraș.”
- De ce: „A juca un rol” e o expresie corectă și consacrată (DEX ’98), nu un calc. Se semnalează doar formula goală, când nu spune ce anume face subiectul. Spui ce face: ce transportă, ce reglează, ce schimbă.
- Sursă: DEX ’98 s.v. rol; tipar de stil

### RO-STI-218 Importanța umflată
- Nivel: recomandare
- Greșit: „Inaugurarea a fost un moment definitoriu pentru oraș.” → Corect: „Inaugurarea a schimbat [de completat: ce s-a schimbat, cu sursa].”
- De ce: „Un moment definitoriu”, „o pagină de istorie”, „va rămâne în memoria orașului” țin locul faptului printr-o etichetă de importanță. Expresiile sunt corecte; se semnalează doar când sunt goale. Rămân în citat și acolo unde fraza spune în continuare ce s-a schimbat.
- Sursă: DOOM 3; DEX ’98

### RO-STI-221 Formule care pretind profunzime
- Nivel: recomandare
- Greșit: „La final de zi, ce contează e cât costă.” → Corect: „Contează cât costă.”
- De ce: „La final de zi” (calc după engl. „at the end of the day”), „Adevărul e că…”, „Să fim sinceri:” pretind o profunzime pe care fraza nu o aduce. Le reformulezi direct, nu le ștergi mecanic: se păstrează când chiar opun o idee alteia sau când rezumă cu adevărat.
- Sursă: tipar de stil

### RO-STI-222 Formulele aforistice
- Nivel: recomandare
- Greșit: „Un site bun este o conversație care nu se termină.” → Corect: „Site-ul [de completat: ce permite cititorului să facă].”
- De ce: Formula are aerul unui adevăr general, dar nu spune nimic verificabil despre obiect. Aforismul e o figură de stil legitimă (DSL), mai ales în text literar sau în eseu; nu se semnalează în citate, în texte literare sau acolo unde e ales cu intenție.
- Sursă: DSL, s.v. AFORISM; tipar de stil

## Deschideri și încheieri

### RO-STI-229 Deschideri care amână subiectul
- Nivel: recomandare
- Greșit: „Nu e un secret pentru nimeni că grădinăritul devine tot mai popular.” → Corect: „Grădinăritul devine tot mai popular [de completat: faptul, cu cifra sau sursa].”
- De ce: Deschideri ca „Nu e un secret pentru nimeni că…”, „De când lumea și pământul…” amână subiectul. Prima frază spune faptul. Dacă faptul nu e de notorietate, îl atribui.
- Sursă: tipar de stil

### RO-STI-223 Deschiderea fals sinceră
- Nivel: recomandare
- Greșit: „Răspunsul scurt? Nu.” → Corect: „Nu.”
- De ce: Întrebarea cu răspuns imediat („Verdictul?”, „Adevărul?”) e un tic conversațional fără loc în textul informativ. Nu se semnalează automat orice „sincer” sau „sincer vorbind” folosit ca adverb de atitudine al autorului: ghidul INJ îl folosește în propriul text.
- Sursă: INJ; tipar de stil

### RO-STI-216 Deschiderea cu întrebare retorică de șablon
- Nivel: recomandare
- Greșit: „Te-ai întrebat vreodată cât timp pierzi cu facturile? Aplicația le plătește automat, la scadență.” → Corect: „Aplicația plătește facturile automat, la scadență.”
- De ce: Întrebarea retorică de șablon, care amână teza, se semnalează și se rescrie pornind direct de la teză. Interogația retorică e însă o figură de stil recunoscută (DSL), deci nu se semnalează în corpul unui text de opinie și nici când citezi o întrebare a sursei.
- Sursă: DSL, s.v. INTEROGAȚIE; tipar de stil

### RO-STI-210 Finalul care nu afirmă nimic
- Nivel: recomandare
- Greșit: „Detaliile le vom afla la momentul potrivit, iar până atunci rămâneți aproape.” → Corect: „[ultimul fapt din sursă sau pasul următor concret]”
- De ce: Încheieri ca „Drumul continuă”, „Rămâneți aproape”, „Detaliile vor veni la momentul potrivit” nu spun nimic verificabil. Dacă sursa are un pas următor concret (o dată, un loc, o acțiune), acela încheie textul; altfel textul se oprește la ultimul fapt.
- Sursă: tipar de stil

### RO-STI-211 Formule goale de rezumat și de concluzie
- Nivel: recomandare
- Greșit: „Așadar, în concluzie, parcarea nouă are 80 de locuri.” → Corect: „Parcarea nouă are 80 de locuri.”
- De ce: „Pe scurt, putem spune că…”, „Așadar, în concluzie…” amână afirmația. Concluzia se spune direct și se atribuie dacă nu e a autorului. Problema e de stil, nu de vocabular: „a concluziona” e corect (DOOM 3). Rezumatul se taie întreg dacă nu aduce nimic nou.
- Sursă: DOOM 3 s.v. concluziona; DEX ’98

## Ritm și structură

### RO-STI-206 Textul nu se taie în staccato
- Nivel: recomandare
- Greșit: „Lansarea a avut loc ieri. Mulți clienți. Fără probleme.” → Corect: „Lansarea a avut loc ieri, cu mulți clienți și fără probleme.”
- De ce: Frazele obișnuite nu se transformă în propoziții scurte și dramatice doar pentru ritm. Lungimea o variezi prin conținutul frazelor, nu prin tăiere, și la împărțire păstrezi conectorii (Mic ghid CE cere, când împarți o frază lungă, conectorii adecvați, ca documentul să rămână coerent). Se semnalează doar când apar în serie, fără funcție expresivă.
- Sursă: Mic ghid CE p. 6 (Rec. 4); tipar de stil

### RO-STI-215 Listele fără verb
- Nivel: recomandare
- Greșit: „• Calitate • Încredere • Experiență” → Corect: „Firma [de completat: de când există, ce face, cu ce garanții].”
- De ce: Se semnalează listele de puncte formate doar din substantive abstracte sau sloganuri, fără verb, cifră sau referent concret. Nu se semnalează listele fără verb cu elemente concrete (nume, instituții, date, sume, pași, categorii definite), care sunt normale într-un act normativ, într-o casetă de fapte sau într-un program.
- Sursă: tipar de stil; Mic ghid CE p. 9 (Rec. 7)

### RO-STI-226 Formă mecanică: triade, bold, emoji
- Nivel: recomandare
- Greșit: „Un serviciu **rapid**, **sigur** și **simplu** 🚀” → Corect: „Un serviciu [de completat: ce îl face rapid, sigur sau simplu].”
- De ce: Se semnalează, ca indicii și nu ca greșeli de limbă: emoji decorativ într-un text informativ, cuvinte îngroșate puse mecanic în mijlocul frazelor și gruparea forțată în câte trei elemente când conținutul are două sau patru. Fiecare caz se judecă în context: un emoji într-un mesaj prietenos, boldul într-un termen definit și o enumerare firească de trei lucruri nu sunt tipare.
- Sursă: tipar de stil; DSL, s.v. ENUMERARE

### RO-STI-236 Regula densității
- Nivel: recomandare
- Greșit: „Textul conține un singur «esențial» și este raportat ca text generat.” → Corect: „Un singur «esențial» nu se raportează; se raportează grupul.”
- De ce: Detecția caută grupuri de semnale, adică densitatea unui tic, nu apariția lui izolată. Regula se aplică și la conectori și la adjectivele-semnal: „esențial”, „crucial”, „semnificativ” sunt cuvinte literare din DEX și DOOM 3. Chiar peste prag, verdictul spune „probabil”, nu „e scris de un model”.
- Sursă: DOOM 3; DEX ’98; tipar de stil

## Vezi și

- [coerenta.md](coerenta.md): aglomerarea de conectori (RO-COE-016) și testul tranziției.
- [stilistica.md](stilistica.md): limba de lemn, densitatea substantivelor abstracte, politețea tradusă.
- [fidelitate.md](fidelitate.md): ce nu ai voie să adaugi când scoți un tipar.
- [registru-si-adresare.md](registru-si-adresare.md): registrul familiar care nu face un text să pară scris de om.
- [ce-nu-e-greseala.md](ce-nu-e-greseala.md): cuvintele și construcțiile corecte pe care nu le numești tipare.
