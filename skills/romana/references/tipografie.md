# Tipografia

Fișierul privește felul în care se scrie textul, nu ce spune: diacriticele, majusculele, virgula și celelalte semne de punctuație, ghilimelele, linia de pauză și cratima, apostroful, spațierea, cuvintele străine și numele proprii străine. Numerele, datele, orele, sumele și unitățile sunt în [numere-date-unitati.md](numere-date-unitati.md). Formulele de chatbot și celelalte tipare ale textului generat sunt în [tipare-text-generat.md](tipare-text-generat.md).

Tipografia e ultima trecere din ordinea de lucru (vezi SKILL.md). O verifici după ce textul e gata ca sens, structură și cuvinte, pentru că o virgulă corectată într-o frază care se rescrie oricum e muncă pierdută.

## Cum folosești fișierul

- **La redactare** scrii direct formele corecte; regulile de aici le folosești la recitirea ciornei.
- **În corectură** corectezi ce e normă și propui ce e recomandare. Un text corect rămâne neschimbat.
- **În verificare** raportezi fiecare abatere cu nivelul ei. Unde sursele se contrazic (nivelul „disputat”), nu alegi tu o variantă: ceri doar ca alegerea să fie aceeași în tot textul.
- **Citatul direct** păstrează cuvintele sursei. În citat corectezi doar tipografia: ghilimelele, diacriticele lipsă sau cu sedilă, spațiile. Grafia veche cu „î” („sînt”, „cînd”) rămâne.
- **Înainte să corectezi, verifici lista „ce nu e greșeală”** din [ce-nu-e-greseala.md](ce-nu-e-greseala.md): linia de pauză, ghilimelele «…», apostroful, punctul-și-virgula, „de ex.”, bara oblică, „totuși” cu sau fără virgule sunt forme corecte.

Abrevierile din câmpul „Sursă”: DOOM 3 = *Dicționarul ortografic, ortoepic și morfologic al limbii române*, ed. a III-a, Academia Română, 2021; ÎOOP6 = *Îndreptar ortografic, ortoepic, morfologic și de punctuație*, ed. a VI-a, Academia Română, 2024; ghid gov-md = *Ghid de redactare lingvistică și tehnoredactare a actelor normative*, Cancelaria de Stat a Republicii Moldova, 2026; Ghid UE = *Ghidul de redactare interinstituțional* al Oficiului pentru Publicații al Uniunii Europene, partea în limba română (convenții ale actelor UE, nu normă generală); IER = Institutul European din România, *Ghid stilistic de traducere în limba română pentru uzul traducătorilor acquis-ului comunitar*, ed. a V-a, 2008 (practică de traducere juridică, nu normă pentru presă); INJ = *Corectitudine lingvistică în actele judecătorești*, Institutul Național al Justiției, Chișinău, 2020; B. Croitor et al. = B. Croitor et al., *Ești cool și dacă vorbești corect*, Institutul de Lingvistică „Iorgu Iordan – Al. Rosetti”, 2014; Legea 290/2018 = legea privind denumirile localităților din România. Mențiunea „(nereverificat local)” arată că trimiterea la capitol, paragraf sau pagină n-a putut fi verificată pe text local.

## 1. Diacriticele

### Metoda: judeci textul, nu cuvintele izolate

Româna are cinci litere cu semne: ă, â, î, ș, ț. Ș și ț se scriu cu virgulă dedesubt (Ș U+0218, ș U+0219, Ț U+021A, ț U+021B). Literele cu sedilă (U+015E, U+015F, U+0162, U+0163) aparțin altor alfabete și sunt greșite în cuvintele românești. Majusculele au și ele diacritice: ȘTIRI, ȚARA, ÎNCEPUT.

Pașii:

1. Mai întâi vezi dacă textul e în română. Un text englezesc, o adresă web sau un nume de fișier nu se „diacritizează”. Cuvintele funcționale (și, în, este, sunt, dar, sau, care, din) te ajută să recunoști româna.
2. Cauți sedilele, la litere mici și la majuscule: ele se pot înlocui sigur cu ș și ț, dar numai în cuvintele românești.
3. Dacă diacriticele lipsesc de peste tot, nu repari cuvânt cu cuvânt. Semnalezi lipsa pentru tot textul (RO-TIP-119). Restaurarea automată e nesigură, pentru că multe forme fără semne sunt ambigue (RO-TIP-124).
4. Dacă lipsesc doar la unele cuvinte, hotărăști după sens: „sa” poate fi „șa”, dar și „sa” din „casa sa”.

Slab: o frază ca „Știri din București: țara se pregătește de iarnă”, în care ș, ț și Ș sunt scrise cu sedilă (U+015F, U+0163, U+015E). Pe ecran arată aproape la fel, dar codurile sunt altele.

Bun: aceeași frază, cu ș, ț și Ș cu virgulă (U+0219, U+021B, U+0218).

### RO-TIP-027 Diacriticele sunt obligatorii, cu virgulă la ș și ț
- Nivel: normă
- Greșit: „ştiinţă, BUCUREŞTI” → Corect: „știință, BUCUREȘTI”
- De ce: Semnele diacritice sunt obligatorii în română, inclusiv pe majuscule, iar ș și ț se scriu cu virgulă, nu cu sedilă. Caracterele cu sedilă par la fel pe ecran, dar au alt cod și strică căutarea, sortarea și potrivirea cuvintelor. Sedila se corectează în toate cuvintele românești, mici sau mari, și în titluri.
- Sursă: ghid gov-md §2.1.3 și §2.1.5; DOOM 3 §1.1.3.2

### RO-TIP-029 Textul în altă limbă și adresele web nu se „diacritizează”
- Nivel: normă
- Corect: „The report covers Romania in 2025.” Aici „Romania” nu primește diacritice, la fel ca în adresa unei pagini.
- De ce: Diacriticele românești se scriu în cuvinte românești. Un titlu englezesc, un citat într-o altă limbă, o adresă de e-mail, un URL sau un nume de fișier rămân cum sunt. Un verificator care marchează „Romania” dintr-un text englezesc ca „lipsă de diacritice” greșește.
- Sursă: ghid gov-md §2.3.1–§2.3.3; DOOM 3 §1.1.3.2 (cuvintele și numele străine își păstrează grafia)
- Vezi și: RO-TIP-071

### RO-TIP-030 Lipsa diacriticelor nu se repară după listă
- Nivel: normă
- Corect: „Casa sa e lângă râu.” Aici „sa” e corect, fiindcă e adjectivul posesiv.
- De ce: Multe cuvinte fără diacritice sunt ele însele corecte: „sa” (posesiv), „tine” (pronume), „si”, „para” (verbul „a para”). Nu înlocuiești mecanic după o listă; citești fraza și hotărăști după sens, cu proba de înlocuire (pui forma la alt caz sau la alt număr).
- Sursă: DOOM 3 §1.1.3.2; B. Croitor et al., 2014

### RO-TIP-116 Diacriticele sunt obligatorii și în titluri, meta și profiluri
- Nivel: normă
- Greșit: „Cum alegi o saltea potrivita pentru copii” → Corect: „Cum alegi o saltea potrivită pentru copii”
- De ce: Diacriticele se scriu în orice text românesc, inclusiv în titluri, în descrierile de pagină, în titlul H1 și în profilurile publice; lipsa semnelor poate schimba sensul („pături” și „paturi” sunt cuvinte diferite), așa că și un titlu fără diacritice are nevoie de ele. Fraza „Cele mai bune paturi pentru copii” e însă corectă fără nicio diacritică lipsă; nu o „repari” fără un context care cere „pături” (vezi RO-TIP-030). Rămân fără diacritice doar identificatorii tehnici: slug-ul din adresa paginii (de exemplu „cele-mai-bune-paturi-copii”), adresele de e-mail și numele de fișiere. Nu corectezi nici siglele consacrate care se scriu fără semne (ANIF, SNSPA) și nici numele proprii pe care purtătorii le scriu deliberat fără ele (de exemplu Hasdeu, nu Hașdeu, dacă așa își scrie familia numele).
- Sursă: ghid gov-md §2.1.4, §2.1.5 și §2.3; INJ p. 93–94; DOOM 3 §1.1.3.2

### RO-TIP-119 Lipsa totală a diacriticelor se semnalează
- Nivel: de semnalat
- Corect: „Observație pentru autor: textul are peste 300 de cuvinte și nicio diacritică; verifică tastatura sau exportul.”
- De ce: Un text românesc de peste 20–40 de cuvinte fără nicio diacritică e suspect: de obicei vine de la o tastatură sau de la un export care a pierdut semnele. Semnalezi lipsa pentru tot textul, o dată, și nu o repari cuvânt cu cuvânt. Fragmentele scurte nu se marchează automat, pentru că pot fi corecte fără diacritice.
- Sursă: ghid gov-md §2.1.4; INJ p. 93–94

### RO-TIP-124 Diacriticele nu se restaurează automat
- Nivel: normă
- Greșit: „Importanta deciziei a crescut.” → Corect: „Importanța deciziei a crescut.”
- De ce: Formele fără diacritice sunt adesea ambigue: „importanta” poate fi „importantă”, „importanța” sau „importanță”; „tara” poate fi „țară”, „țara” sau „tară”. Din această cauză lipsa semnelor se semnalează, iar corectura mecanică produce erori noi. Restaurezi un cuvânt numai când contextul alege o singură formă: în exemplul de mai sus, „deciziei” cere un substantiv articulat, deci „importanța”.
- Sursă: INJ p. 93–94; DEX ’98, s.v. țară, importanță

## 2. Majusculele

### Metoda: întâi felul numelui, apoi regula

Româna scrie cu literă mică substantivele comune, funcțiile, titlurile, lunile, zilele, popoarele și limbile, și cu majusculă numele proprii și, în stilul oficial, anumite denumiri. Majuscula „de politețe” sau „de importanță”, pusă după modelul englezesc, e cea mai frecventă greșeală. Ordinea de lucru:

1. Te întrebi dacă e un nume propriu (o instituție anume, un loc, un eveniment unic) sau un substantiv comun (un guvern, un oraș, o funcție).
2. La un nume oficial complet scrii cu majusculă fiecare cuvânt plin și lași mici cuvintele ajutătoare (prepoziții, conjuncții, articole): Ministerul Afacerilor Externe. Compartimentele din interiorul unei instituții iau majusculă doar la primul cuvânt: „Direcția de studii economice” (RO-TIP-009).
3. La un titlu de operă, de lege, de document, de publicație sau de emisiune scrii cu majusculă doar primul cuvânt și numele proprii din interior.
4. Textul citat și numele purtate de oameni se scriu cum sunt.

Slab (majuscule englezești):

> Guvernul Aprobă Noul Buget De Stat. Ministrul Finanțelor, Ion Popescu, A Declarat Luni Că Se Va Aplica Din Martie.

Bun:

> Guvernul aprobă noul buget de stat. Ministrul finanțelor, Ion Popescu, a declarat luni că se va aplica din martie.

### 2.1 Funcții și titluri

### RO-TIP-001 Funcțiile au literă mică în textul neoficial
- Nivel: normă
- Greșit: „a declarat Ministrul de Externe” → Corect: „a declarat ministrul de externe”
- De ce: Numele funcțiilor (președintele, premierul, prim-ministrul, ministrul, primarul, senatorul) se scriu cu literă mică în registrul neoficial, deci și în presă, oricât de înaltă ar fi funcția și chiar când stau înaintea numelui: „președintele Ion Popescu”. Majuscula la început de propoziție rămâne, desigur, obligatorie.
- Sursă: DOOM 3 §3.2 II.2.2–2.3; ghid gov-md §4.2.1.1

### RO-TIP-002 Domeniul funcției rămâne cu literă mică
- Nivel: normă
- Greșit: „Ministrul Finanțelor a declarat” → Corect: „ministrul finanțelor a declarat”
- De ce: Domeniul care întregește numele funcției (ministrul finanțelor, ministrul de externe) nu face parte dintr-un nume propriu și nu primește majusculă. Instituția, în schimb, o primește: „Ministerul Finanțelor”.
- Sursă: DOOM 3 §3.2 II.2.2–2.3; ghid gov-md §4.2.1.1
- Vezi și: RO-TIP-001

### RO-TIP-003 Majuscula la funcții: doar în stilul oficial sau când funcția denumește instituția
- Nivel: normă
- Corect: „Președintele Republicii semnează decretul.” (stil oficial, funcția celui mai înalt demnitar)
- De ce: Majuscula la o funcție apare în stilul oficial reverențios (acte, protocol), în unele funcții înalte din acte și când funcția denumește chiar o instituție. Într-un articol de presă sau pe o pagină de site se scrie cu literă mică. Într-un e-mail sau într-o scrisoare, DOOM 3 admite majuscula când autorul marchează o atitudine ceremonioasă ori reverențioasă: n-o corectezi. Dacă textul e un act sau o scrisoare protocolară, urmezi convenția lui și ții aceeași alegere în tot textul.
- Sursă: DOOM 3 §3.2 II.2.2; ghid gov-md §4.1.5.3
- Vezi și: RO-TIP-001

### RO-TIP-004 Titlurile științifice și onorifice au literă mică
- Nivel: normă
- Greșit: „Profesorul Dumitrescu și Academicianul Ionescu” → Corect: „profesorul Dumitrescu și academicianul Ionescu”
- De ce: Titlurile științifice și onorifice (profesorul, doctorul, academicianul, doctor honoris causa) sunt substantive comune. Numele purtătorului primește majusculă, titlul nu.
- Sursă: ghid gov-md §4.2.1.11; DOOM 3 §3.2 II.2.12

### RO-TIP-022 „Domnia-sa” și „dumnealui” au literă mică
- Nivel: normă
- Greșit: „Am vorbit cu Dumnealui ieri.” → Corect: „Am vorbit cu dumnealui ieri.”
- De ce: Pronumele de politețe formate recent, din seria „domnia-sa”, „domnia-voastră”, „domniile-lor”, „dumnealui”, se scriu cu literă mică. Scrierea cu cratimă („domnia-sa”) e recomandată de DOOM 3, iar „domnia sa”, fără cratimă, nu se corectează. Locuțiunea protocolară „Domnia Sa”, cu majusculă și fără cratimă, e corectă în stilul oficial și nu se atinge (vezi RO-TIP-023). Cuvântul „domnia-sa” e în dicționar.
- Sursă: DOOM 3 §3.1 I.2.22, §3.2 II.4 și §6.7.2; ghid gov-md §4.1.5.11

### RO-TIP-023 Formulele protocolare vechi au majusculă
- Nivel: normă
- Corect: „Excelența Sa, Majestatea Sa, Sanctitatea Sa”
- De ce: Locuțiunile vechi de adresare, pentru domnitori, boieri, demnitari și capete încoronate, se scriu cu majusculă la toate componentele și fără cratimă. Nu le amesteci cu „domnia-sa”, care e un pronume modern (RO-TIP-022).
- Sursă: ghid gov-md §4.1.5.11; DOOM 3 §3.1 I.2.22

### 2.2 Instituții

### RO-TIP-005 Numele oficiale iau majusculă la toate cuvintele pline
- Nivel: normă
- Greșit: „ministerul afacerilor externe” → Corect: „Ministerul Afacerilor Externe”
- De ce: Denumirea oficială a unei instituții se scrie cu majusculă la fiecare cuvânt plin și cu literă mică la cuvintele ajutătoare: Camera Deputaților, Curtea Constituțională, Parlamentul European, Organizația Națiunilor Unite.
- Sursă: ghid gov-md §4.1.5.1; DOOM 3 §3.1 I.2.5

### RO-TIP-006 Forma eliptică a unei instituții anume păstrează majuscula
- Nivel: normă
- Corect: „Guvernul a decis să amâne măsura.”
- De ce: Când cuvântul scurt desemnează instituția anume („Guvernul”, „Parlamentul”, „Camera”, „secretar de stat la Externe”), se scrie cu majusculă, ca și numele complet.
- Sursă: ghid gov-md §4.1.5.1; DOOM 3 §3.1 I.2.5
- Vezi și: RO-TIP-005

### RO-TIP-007 Termenul generic sau neoficial are literă mică
- Nivel: normă
- Greșit: „Guvernul Român a anunțat” → Corect: „guvernul român a anunțat”
- De ce: „Guvernul României” e denumirea instituției, iar „guvernul român” e un termen generic, care se scrie cu literă mică. La fel „guvernele statelor membre”, „am fost la primărie”, „un guvern interimar”.
- Sursă: ghid gov-md §4.2.1.6; DOOM 3 §3.2 II.2.3 și §3.1 I.2.5
- Vezi și: RO-TIP-005

### RO-TIP-008 Elipsa sau termenul generic: hotărăște sensul
- Nivel: normă
- Corect: „Guvernul a adoptat ordonanța.” (instituția anume) și „Guvernele statelor membre au reacționat.” (termenul generic)
- De ce: Alegerea între majusculă și literă mică nu se face după formă, ci după sens: dacă cuvântul trimite la o instituție anume, e nume propriu; dacă numește orice guvern, e substantiv comun. În caz de îndoială, citești fraza și întrebi autorul.
- Sursă: ghid gov-md §4.2.1.6; DOOM 3 §3.2 II.2.3 și §3.1 I.2.5

### RO-TIP-009 Compartimentele unei instituții: doar primul cuvânt cu majusculă
- Nivel: normă
- Greșit: „Direcția de Studii Economice” → Corect: „Direcția de studii economice”
- De ce: Organele de conducere și subdiviziunile din interiorul unei instituții (direcție, secție, serviciu, comisie permanentă) se scriu cu majusculă doar la primul cuvânt: „Direcția de studii economice”. DOOM 3 §3.1 I.3.1 admite și varianta „Direcția Studii economice” (cu două majuscule): nu o corectezi. Forma cu o singură majusculă este cea de bază; „Consiliul de administrație” e dat tot acolo. Excepție: în actele UE, direcțiile generale ale organismelor europene iau majusculă la fiecare cuvânt plin: „Direcția Generală pentru Agricultură și Dezvoltare Rurală” (ÎOOP6 p. 82). Instituția de sine stătătoare, cu personalitate juridică, rămâne la RO-TIP-005.
- Sursă: ghid gov-md §4.1.7.3; DOOM 3 §3.1 I.3.1; ÎOOP6 p. 82

### RO-TIP-085 Instituțiile străine traduse iau majusculă la cuvintele pline
- Nivel: normă
- Greșit: „departamentul de stat american” → Corect: „Departamentul de Stat american”
- De ce: Când o instituție străină are un nume tradus în română, numele tradus se scrie ca orice nume oficial: majusculă la cuvintele pline („Departamentul de Stat”, „Congresul”). Adjectivul care arată țara rămâne cu literă mică.
- Sursă: DOOM 3 §3.1 I.2.5; ghid gov-md §15.2.4
- Vezi și: RO-TIP-086

### RO-TIP-086 Funcțiile străine au literă mică, ca cele românești
- Nivel: normă
- Greșit: „Ieri a ajuns la Berlin Președintele Macron.” → Corect: „Ieri a ajuns la Berlin președintele Macron.”
- De ce: Numele funcției rămâne substantiv comun și se scrie cu literă mică, în orice țară ar fi exercitată. Numele persoanei se scrie cu majusculă. Prenumele la prima mențiune e o preferință de stil la redactare; în corectură nu se adaugă, pentru că ar fi un conținut nou (vezi fidelitate.md).
- Sursă: DOOM 3 §3.2 II.2.2; ghid gov-md §4.2.1.1
- Vezi și: RO-TIP-085

### 2.3 Titluri de opere, de documente și de pagini

### RO-TIP-010 În titluri, doar primul cuvânt și numele proprii au majusculă
- Nivel: normă
- Greșit: „Guvernul Aprobă Noul Buget De Stat” → Corect: „Guvernul aprobă noul buget de stat”
- De ce: Titlurile de opere, legi, documente, publicații, emisiuni, programe și articole se scriu cu majusculă la primul cuvânt și la numele proprii din interior. Majuscula la fiecare cuvânt (Title Case) imită engleza și e greșită. Excepțiile sunt la RO-TIP-127. La fel se scriu titlurile paginilor și ale secțiunilor, butoanele și meniurile. Un titlu tradus nu păstrează majusculele engleze („Cele mai recente articole”, nu „Cele mai Recente Articole”). Numele unei tehnici sau al unui concept nu e nume propriu: „chain of thought” se scrie cu litere mici. Numele de produs se scrie cum îl înregistrează titularul (iPhone, eBay).
- Sursă: ghid gov-md §4.1.5.10 și §4.1.7.1; DOOM 3 §3.1 I.3.3

### RO-TIP-011 Sărbătorile naționale și unele documente de importanță națională fac excepție
- Nivel: normă
- Corect: „Ziua Națională a României”
- De ce: Sărbătorile naționale oficiale și denumirile unor documente sau evenimente de importanță națională se scriu cu majusculă la toate cuvintele pline. Lista exactă a excepțiilor e la RO-TIP-127. Nu transformi o excepție în regulă: un titlu obișnuit rămâne cu literă mică.
- Sursă: ghid gov-md §4.1.5.7 și §4.1.7.1; DOOM 3 §3.1 I.2.19 și I.3.3

### RO-TIP-127 Excepțiile de la regula cu un singur cuvânt cu majusculă
- Nivel: normă
- Corect: „Festivalul Internațional de Teatru de la Sibiu; Planul Național de Redresare și Reziliență; Casa Poporului; Ziua Europei”
- De ce: Majuscula la fiecare cuvânt plin se păstrează, după DOOM 3 §3.1, pentru patru categorii: denumirile de manifestări (festivaluri, conferințe, saloane), denumirile planurilor de importanță națională, denumirile sărbătorilor oficiale și denumirile de clădiri publice. Fără aceste excepții ai corecta forme corecte. Verifici categoria înainte să schimbi un titlu: „Planul național de dezvoltare pentru anii 2025–2027” (titlu obișnuit de plan) nu se confundă cu „Planul Național de Redresare și Reziliență” (plan național cu denumire consacrată).
- Sursă: DOOM 3 §3.1 I.2.5, I.2.6, I.2.11 și I.2.19; ghid gov-md §4.1.5.5, §4.1.5.7, §4.1.5.8 și §4.1.7.1

### 2.4 Calendar, popoare, limbi, discipline

### RO-TIP-012 Lunile, zilele, popoarele, limbile și disciplinele: literă mică
- Nivel: normă
- Greșit: „pe 5 Martie, Luni, un Francez, limba Română” → Corect: „pe 5 martie, luni, un francez, limba română”
- De ce: Lunile, zilele săptămânii, numele popoarelor, ale limbilor și ale disciplinelor de studiu sunt substantive comune; engleza le scrie cu majusculă, româna nu. Numele de oameni derivate din ele rămân cu majusculă acolo unde sunt nume proprii (Marțian).
- Sursă: ghid gov-md §4.2.1.10; DOOM 3 §3.2 II.2.11

### RO-TIP-013 Sărbătorile au majusculă, datele obișnuite nu
- Nivel: normă
- Greșit: „de 1 decembrie, la mulți ani” → Corect: „de 1 Decembrie, la mulți ani”
- De ce: Când data denumește o sărbătoare (1 Decembrie, 9 Mai, Întâi Mai), luna se scrie cu majusculă. Când data e doar un moment din calendar („ședința din 1 decembrie”), luna rămâne cu literă mică.
- Sursă: ghid gov-md §4.1.5.7; DOOM 3 §3.1 I.2.19 și §3.2 II.2.11

### 2.5 Geografie

### RO-TIP-014 Termenul generic dinaintea unui nume geografic: literă mică
- Nivel: normă
- Greșit: „Fluviul Dunărea, pe Strada Lipscani” → Corect: „fluviul Dunărea, pe strada Lipscani”
- De ce: Cuvântul generic (fluviul, județul, strada, munții) rămâne cu literă mică atunci când nu face parte din nume: „județul Cluj”, „munții Carpați”, „bulevardul Grigore Vieru”. În indicarea adresei, pe hărți și pe plăcuțe indicatoare termenul generic se scrie însă cu majusculă (RO-TIP-015).
- Sursă: ghid gov-md §4.2.1.7; DOOM 3 §3.2 II.2.7

### RO-TIP-015 Termenul generic care face parte din nume are majusculă
- Nivel: normă
- Corect: „Marea Neagră, Țara Românească, Baia Mare, Drumul Taberei”
- De ce: Dacă termenul generic e urmat de un determinant care îl individualizează sau și-a pierdut sensul propriu, face parte din numele propriu și se scrie cu majusculă. Pe hărți, pe plăcuțe indicatoare și în adrese, aceste substantive se scriu cu majusculă.
- Sursă: ghid gov-md §4.2.1.7; DOOM 3 §3.1 I.2.1 și §3.2 II.2.7

### RO-TIP-016 Punctele cardinale: literă mică ca direcție, majusculă ca regiune
- Nivel: normă
- Greșit: „Europa de est” → Corect: „Europa de Est”
- De ce: Punctul cardinal care arată o direcție se scrie cu literă mică („spre est”, „în sudul țării”). Cel care denumește o regiune sau face parte dintr-un nume se scrie cu majusculă: Orientul Mijlociu, Europa de Est, America de Nord, Polul Sud.
- Sursă: ghid gov-md §4.1.6; DOOM 3 §3.1 I.2.18 și §3.2 II.2.5

### RO-TIP-017 Cuvintele ajutătoare din nume rămân mici
- Nivel: normă
- Greșit: „Statele Unite Ale Americii” → Corect: „Statele Unite ale Americii”
- De ce: Prepozițiile, articolele și conjuncțiile din interiorul unui nume propriu compus rămân cu literă mică: „Regatul Unit al Marii Britanii și Irlandei de Nord”, „Ștefan cel Mare”, „Ursula von der Leyen”.
- Sursă: ghid gov-md §4.2.1.13; DOOM 3 §3.2 II.3.2

### 2.6 Războaie, epoci, evenimente

### RO-TIP-018 Războaiele unice sau de anvergură au majusculă la toate componentele
- Nivel: normă
- Greșit: „al doilea război mondial, Primul război mondial” → Corect: „Al Doilea Război Mondial, Primul Război Mondial”
- De ce: Războaiele cu denumire proprie consacrată se scriu cu majusculă la toate componentele, inclusiv la cuvântul ajutător de la început: Al Doilea Război Mondial, Războiul de Independență, Războiul Rece. Scrierea războaielor mondiale așa e o noutate a DOOM 3.
- Sursă: ghid gov-md §4.1.5.6; DOOM 3 §3.1 I.2.10 și „Ce e nou în DOOM 3” §2.2.5

### RO-TIP-019 Marile epoci și evenimente au majusculă
- Nivel: normă
- Greșit: „evul mediu” → Corect: „Evul Mediu”
- De ce: Marile epoci istorice și evenimentele unice se scriu cu majusculă la toate cuvintele pline: Antichitatea, Evul Mediu, Renașterea, Marea Unire, Revoluția Franceză.
- Sursă: ghid gov-md §4.1.5.6; DOOM 3 §3.1 I.2.9
- Vezi și: RO-TIP-018

### RO-TIP-020 Adjectivul adăugat rămâne cu literă mică
- Nivel: normă
- Greșit: „Evul Mediu Românesc” → Corect: „Evul Mediu românesc”
- De ce: Adjectivul care precizează o epocă sau un eveniment nu face parte din denumire și nu primește majusculă.
- Sursă: ghid gov-md §4.1.7.1; DOOM 3 §3.1 I.3.3

### RO-TIP-021 Perioadele vagi și războaiele multiple: literă mică
- Nivel: normă
- Greșit: „Epoca Modernă, Războaiele Balcanice” → Corect: „epoca modernă, războaiele balcanice”
- De ce: Perioadele delimitate vag (epoca modernă, perioada postbelică) și războaiele fără nume propriu unic (războaiele balcanice, războiul din Afganistan) sunt denumiri descriptive și se scriu cu literă mică.
- Sursă: ghid gov-md §4.2.1.4 și §4.2.1.5; DOOM 3 §3.2 II.2.6.3–2.6.4

### 2.7 După două-puncte, în citat, în liste

### RO-TIP-024 După două-puncte se continuă cu literă mică
- Nivel: normă
- Greșit: „Motivul e simplu: Nu am timp.” → Corect: „Motivul e simplu: nu am timp.”
- De ce: Când după două-puncte urmează o enumerare, o explicație, o concluzie sau o consecință, textul continuă cu literă mică, pentru că fraza nu s-a încheiat. Excepțiile sunt vorbirea directă și citatul din RO-TIP-025.
- Sursă: ghid gov-md §4.1.3 și §9.4.5; DOOM 3 §3.2 II.1.4

### RO-TIP-025 Vorbirea directă și citatul-frază încep cu majusculă
- Nivel: normă
- Corect: „A spus: «Bugetul nu ajunge.»” În ghilimelele românești: A spus: „Bugetul nu ajunge.”
- De ce: Vorbirea directă, fie între ghilimele, fie introdusă cu linie de dialog, și citatul redat ca enunț de sine stătător încep cu majusculă, chiar după două-puncte. Citatul care e doar un fragment integrat în frază intră la RO-TIP-026.
- Sursă: ghid gov-md §4.1.3 și §9.8.6; DOOM 3 §3.1 I.1.2.3–1.2.4
- Vezi și: RO-TIP-024, RO-TIP-026

### RO-TIP-026 Citatul integrat în frază începe cu literă mică
- Nivel: normă
- Corect: „A vorbit despre «un buget care nu ajunge», de fapt.” Fără două-puncte înainte, fragmentul citat începe cu literă mică.
- De ce: Când citatul continuă propoziția autorului, fără două-puncte, începe cu literă mică și punctul final stă după ghilimelele de închidere (RO-TIP-104). Majuscula ar înlocui o relație sintactică cu un început de frază.
- Sursă: ghid gov-md §4.2.1.12 și §9.8.5; DOOM 3 §3.2 II.1.5
- Vezi și: RO-TIP-024, RO-TIP-025

### RO-TIP-125 Elementele unei liste nu cer majusculă
- Nivel: recomandare
- Greșit: „Ai nevoie de: Cont, Parolă, Adresă de e-mail.” → Corect: „Ai nevoie de: cont, parolă, adresă de e-mail.”
- De ce: Elementele unei enumerări introduse prin două-puncte și despărțite prin virgulă sau punct-și-virgulă încep cu literă mică. Majuscula la începutul fiecărui element e o convenție de casă, nu o normă; o impui numai dacă ți-o cere aplicația. Elementele care sunt fraze întregi, așezate pe rânduri separate, pot începe cu majusculă dacă toate au aceeași formă.
- Sursă: ghid gov-md §9.3.2, §9.4.5 și §4.1.11

### 2.8 Majuscule de accent

### RO-TIP-122 Fără majuscule de accent în text
- Nivel: recomandare
- Greșit: „Oferta e destinată EXCLUSIV clienților noi.” → Corect: „Oferta e destinată exclusiv clienților noi.”
- De ce: Majusculele integrale folosite pentru accent sunt o alegere de stil, nu de ortografie. În textul curent accentul se dă prin formulare sau, la nevoie, prin cursive. Fac excepție siglele, titlurile sau afișele scrise integral cu majuscule și formulele consacrate din acte (HOTĂRĂȘTE, DISPUNE). INJ admite majuscula folosită deliberat, în scop stilistic, pentru a scoate un cuvânt în evidență; un text în care autorul o folosește așa nu se corectează, doar se propune formularea.
- Sursă: INJ p. 92

## 3. Virgula

### Metoda: întâi felul legăturii, apoi virgula

Virgula nu se pune „unde se respiră”. Se pune după legătura sintactică dintre părți. Pașii:

1. Găsești verbul fiecărei propoziții și vezi ce fel de legătură are fiecare propoziție sau fiecare parte de propoziție cu restul.
2. Elementele **obligatorii** (subiectul, predicatul, complementele cerute de verb, completivele și subiectivele puse după regentă) nu se despart prin virgulă de verb.
3. Elementele **adăugate** (apoziția explicativă, incidenta, construcția gerunzială de la început, concesiva) se izolează prin virgulă sau, când conțin virgule, prin linii de pauză.
4. La coordonare: „și” și „sau” simple nu cer virgulă; „dar”, „iar”, „ci”, „însă” o cer; perechile repetate („și… și”, „fie… fie”, „nici… nici”) o cer între termeni.
5. La relative, virgula arată sensul: relativa **explicativă** (adaugă o informație) se desparte, cea **determinativă** (spune care anume) nu.

Slab:

> Directorul a anunțat, că termenul se prelungește iar angajații care au lucrat în weekend, primesc un spor.

Bun:

> Directorul a anunțat că termenul se prelungește, iar angajații care au lucrat în weekend primesc un spor.

Virgula de după „anunțat” e greșită: completiva pusă după regentă nu se desparte. Lipsește virgula dinaintea lui „iar”. Virgula de după „weekend” desparte subiectul de predicat. Relativa „care au lucrat în weekend” e determinativă (arată care angajați primesc sporul) și nu se izolează.

### RO-TIP-036 Fără virgulă între subiect și predicat
- Nivel: normă
- Greșit: „Echipa de proiect, a livrat raportul la timp.” → Corect: „Echipa de proiect a livrat raportul la timp.”
- De ce: Subiectul și predicatul nu se despart prin virgulă, oricât de lung ar fi subiectul și chiar dacă e o propoziție: „Cine dorește să facă renegocierea să meargă la Bruxelles”. Excepția e o apoziție sau o incidentă închisă între două virgule.
- Sursă: ghid gov-md §9.2.1.2; ÎOOP6 p. 59 (Virgula, ix) și p. 61–62 (Greșeli de utilizare, i)

### RO-TIP-037 Fără virgulă înaintea complementului sau a circumstanțialului obligatoriu
- Nivel: normă
- Greșit: „Primarul a semnat, contractul cu firma.” → Corect: „Primarul a semnat contractul cu firma.”
- De ce: Complementul direct sau indirect așezat după verb, numele predicativ și, de regulă, circumstanțialul cerut de verb nu se despart prin virgulă. Nici complementul de agent nu se desparte de verb. O virgulă pusă aici are doar rol emfatic.
- Sursă: ghid gov-md §9.2.1.3; ÎOOP6 p. 60 (Virgula, ix)
- Vezi și: RO-TIP-038

### RO-TIP-038 Fără virgulă înaintea completivei, a subiectivei și a predicativei
- Nivel: normă
- Greșit: „Directorul a anunțat, că termenul se prelungește.” → Corect: „Directorul a anunțat că termenul se prelungește.”
- De ce: Completiva directă sau indirectă pusă după regentă, subiectiva și predicativa nu se despart prin virgulă: „Am observat că funcționează”, „Problema este cine răspunde”. Detaliile pentru „că”, „dacă”, „când” sunt la RO-TIP-113.
- Sursă: ghid gov-md §9.2.2.3 și §9.2.2.6; ÎOOP6 p. 59–60 (Virgula, ix)
- Vezi și: RO-TIP-037

### RO-TIP-039 Fără virgulă înaintea atributivei determinative
- Nivel: normă
- Greșit: „Cel, care a greșit, plătește amenda.” → Corect: „Cel care a greșit plătește amenda.”
- De ce: Atributiva determinativă spune care anume e persoana sau lucrul și nu se izolează prin virgulă. Sensul se schimbă dacă o izolezi (vezi RO-TIP-051).
- Sursă: ghid gov-md §9.2.2.5; ÎOOP6 p. 59–60 (Virgula, ix)

### RO-TIP-040 Fără virgulă înaintea lui „și” și „sau” simple
- Nivel: normă
- Greșit: „A citit raportul, și a plecat.” → Corect: „A citit raportul și a plecat.”
- De ce: Părțile de propoziție și propozițiile legate prin „și” sau „sau” nu se despart prin virgulă. Excepțiile: perechile repetate („Au răspuns și clienții, și furnizorii”, „fie… fie”, „nici… nici”), unde virgula e obligatorie, și „și” cu sensul „iar”, urmat de negație („ghidurile publicate, și nu cele vechi”). Pentru „sau” disjunctiv, cazurile sunt la RO-TIP-114.
- Sursă: ghid gov-md §9.2.1.9, §9.2.1.11 (N.B.), §9.2.1.15 și §9.2.2.1; Ghid UE §10.1.2.4–5 și §10.1.2.13 (nereverificat local); ÎOOP6 p. 57–58 (Virgula, ii–iii)

### RO-TIP-041 „Desigur că”, „bineînțeles că”, „sigur că”: fără virgulă
- Nivel: recomandare
- Greșit: „Desigur, că măsura va fi contestată.” → Corect: „Desigur că măsura va fi contestată.”
- De ce: Când „desigur”, „bineînțeles” sau „sigur” introduc o subiectivă cu „că”, virgula dintre ele nu se pune, pentru că subiectiva nu se desparte de regentă. ÎOOP6 nu tratează expres „desigur că”; regula decurge din principiul că subiectiva nu se desparte de regentă. Virgula apare doar când adverbul ține loc de propoziție: „Desigur, măsura va fi contestată.”
- Sursă: ÎOOP6 p. 59–60 (Virgula, ix: subiectiva nu se desparte de regentă); „desigur că” nu e tratat expres

### RO-TIP-042 Fără virgulă înainte de „etc.” și „ș.a.m.d.”
- Nivel: recomandare
- Greșit: „mere, pere, etc.” → Corect: „mere, pere etc.”
- De ce: „Etc.” înlocuiește ultimul element al enumerării și se leagă de el fără virgulă. Dacă propoziția se termină cu „etc.”, punctul abrevierii este și punctul final al frazei.
- Sursă: ghid gov-md §9.2.1.17; Ghid UE §10.1.2.11 (nereverificat local)

### RO-TIP-043 „Însă” și „deci” din interiorul propoziției, fără virgule
- Nivel: recomandare
- Corect: „Guvernul însă a amânat votul.” (preferat de ghiduri); „Guvernul, însă, a amânat votul.” nu se corectează.
- De ce: Ghidul gov-md și Ghidul UE preferă „însă” și „deci” din interiorul propoziției fără virgule: „Măsura constituie deci un avantaj”. ÎOOP6 lasă virgula la alegere la „totuși” și o admite facultativ în cazuri apropiate (vezi RO-NUG-119), așa că forma cu virgule nu se corectează ca greșeală. Cere doar consecvență în tot textul.
- Sursă: ghid gov-md §9.2.1.16; ÎOOP6 p. 60–61 (Virgula, xi)

### RO-TIP-044 „Așadar” și „prin urmare” din interior stau între virgule
- Nivel: recomandare
- Corect: „Guvernul, așadar, a amânat votul.” (preferat de ghidul gov-md); „Guvernul așadar a amânat votul.” nu se corectează.
- De ce: Ghidul gov-md încadrează între virgule „așadar” și „prin urmare” așezate în interiorul propoziției; alte ghiduri (vezi RO-TIP-043) le scriu fără virgule. Preferințele nu sunt norme, deci ceri doar consecvență. La începutul propoziției nu e nevoie de a doua virgulă: „Așadar, votul a fost amânat.” ÎOOP6 lasă facultativă virgula la cuvinte de acest fel („totuși”, p. 60–61), deci forma fără virgule nu e greșeală.
- Sursă: ghid gov-md §9.2.1.13; ÎOOP6 p. 59 (Virgula, vi) și p. 60–61 (Virgula, xi)

### RO-TIP-045 Virgulă înaintea lui „dar”, „iar”, „ci”, „însă”
- Nivel: normă
- Greșit: „Ministerul a promis fonduri dar nu le-a alocat.” → Corect: „Ministerul a promis fonduri, dar nu le-a alocat.”
- De ce: Conjuncțiile adversative „dar”, „iar”, „ci”, „însă” aflate la începutul unei propoziții sau al unei părți coordonate sunt precedate de virgulă. „Însă” din interiorul propoziției urmează RO-TIP-043.
- Sursă: ghid gov-md §9.2.1.11 și §9.2.2.2; Ghid UE §10.1.2.7 (nereverificat local); ÎOOP6 p. 57–58 (Virgula, ii)

### RO-TIP-046 Virgulă înaintea lui „precum și”, „cât și”, „dar și”
- Nivel: recomandare
- Greșit: „atât statul de origine cât și statul solicitat” → Corect: „atât statul de origine, cât și statul solicitat”
- De ce: Locuțiunile „precum și”, „cât și” (după „atât”) și „dar și” sunt precedate de virgulă. Structura pe care o introduc nu e urmată de virgulă dacă face parte din subiect: „Atât directorul, cât și angajații au semnat.”
- Sursă: ghid gov-md §9.2.1.10; Ghid UE §10.1.2.6 (nereverificat local)

### RO-TIP-047 Apoziția explicativă se izolează
- Nivel: normă
- Greșit: „Ion Popescu fostul ministru a declarat” → Corect: „Ion Popescu, fostul ministru, a declarat”
- De ce: Apoziția explicativă se închide între virgule când fraza continuă după ea (sau între linii de pauză, dacă are virgule în interior; vezi RO-TIP-062). Dacă apoziția stă la sfârșitul propoziției, ajunge o singură virgulă înainte.
- Sursă: ghid gov-md §9.2.1.4.5; ÎOOP6 p. 59–60 (Virgula, viii și x)

### RO-TIP-048 Virgulă după construcția gerunzială sau participială de la început
- Nivel: normă
- Greșit: „Luând în considerare datele Guvernul a decis.” → Corect: „Luând în considerare datele, Guvernul a decis.”
- De ce: Construcția cu gerunziu sau cu participiu așezată la începutul frazei se desparte de restul prin virgulă. Dacă stă după regentă, regula se aplică după legătura cu verbul (vezi RO-TIP-049).
- Sursă: ghid gov-md §9.2.1.4.2; ÎOOP6 p. 58 (Virgula, v) și Ghid UE §10.1.2.2 (nereverificat local)

### RO-TIP-049 Circumstanțialul intercalat sau mutat în față
- Nivel: recomandare
- Greșit: „Echipa la cererea clientului a amânat livrarea.” → Corect: „Echipa, la cererea clientului, a amânat livrarea.”
- De ce: Circumstanțialul pus între subiect și predicat și cel mutat în fața frazei pentru accent se izolează prin virgule. Gradul de legătură cu verbul și topica hotărăsc: circumstanțialele scurte și strâns legate de verb pot rămâne fără virgule, deci propui, nu impui.
- Sursă: ghid gov-md §9.2.1.4.3, §9.2.1.4.4 și §9.2.1.6
- Vezi și: RO-TIP-048

### RO-TIP-050 Concesiva, incidentele și vocativul se izolează
- Nivel: normă
- Greșit: „Deși plouă echipa a plecat.” → Corect: „Deși plouă, echipa a plecat.”
- De ce: Subordonata concesivă se desparte prin virgulă, mai ales când stă înaintea regentei (ÎOOP6 p. 59); la circumstanțialele așezate după regentă virgula e adesea facultativă. La fel se izolează vocativul („Maria, vino încoace”), interjecția și propozițiile intercalate între părțile regentei.
- Sursă: ghid gov-md §9.2.1.5, §9.2.2.7 și §9.2.2.13; ÎOOP6 p. 58–59 (Virgula, iv, vi, viii)

### RO-TIP-051 Atributiva explicativă are virgulă, cea determinativă nu
- Nivel: normă
- Corect: „Deputații, care au votat, au plecat” (toți au votat și toți au plecat) este altceva decât „Deputații care au votat au plecat” (au plecat doar cei care au votat).
- De ce: Virgula schimbă sensul. Când contextul arată limpede una dintre lecturi și virgula o contrazice, e greșeală de punctuație; când textul rămâne ambiguu, semnalezi cazul și întrebi autorul (vezi RO-TIP-112).
- Sursă: ghid gov-md §9.2.2.5; Ghid UE §10.1.2.17 (nereverificat local); ÎOOP6 p. 59–60 (Virgula, viii–ix)

### RO-TIP-052 Cauzala se desparte, în afară de cauza unică
- Nivel: recomandare
- Greșit: „Decizia a fost anulată deoarece procedura nu a fost respectată.” → Corect: „Decizia a fost anulată, deoarece procedura nu a fost respectată.”
- De ce: Cauzala se desparte de regentă prin virgulă oriunde ar sta. Autorul poate arăta că e singura cauză, caz în care virgula dispare („Decizia a fost anulată întrucât procedura nu a fost respectată”); de aceea scoaterea virgulei nu se corectează automat.
- Sursă: ghid gov-md §9.2.2.9; Ghid UE §10.1.2.20 (nereverificat local)

### RO-TIP-053 Condiționala: virgula depinde de legătura cu regenta
- Nivel: recomandare
- Corect: „Apăsați «Continuă», dacă datele sunt corecte.” (condiția e o precizare) și „Nu intri fără permis dacă nu ai însoțitor.” (singura condiție)
- De ce: Condiționala se desparte prin virgulă, după regulă, oriunde ar sta; când arată singura condiție a acțiunii nu se desparte. Nu scoți o virgulă corectă (vezi RO-NUG-115).
- Sursă: ghid gov-md §9.2.2.10; Ghid UE §10.1.2.21 (nereverificat local)
- Vezi și: RO-TIP-052

### RO-TIP-112 Virgula înaintea lui „care” depinde de sens
- Nivel: de semnalat
- Corect: „Colegul, care a venit ieri, a plecat” și „Colegul care a venit ieri a plecat” sunt amândouă corecte, cu sensuri diferite.
- De ce: Virgula înainte de „care” arată felul relativei: explicativă (informație în plus, cu virgule) sau determinativă (identifică, fără virgulă). Nu e o greșeală de scoaterea sau de pus virgula, ci o alegere care schimbă sensul. Semnalezi doar când lectura rămâne ambiguă și întrebi ce a vrut autorul.
- Sursă: ghid gov-md §9.2.2.5

### RO-TIP-113 Virgula nu e obligatorie înainte de „care”, „dacă”, „când”, „că”
- Nivel: normă
- Greșit: „Notificați, când este gata.” → Corect: „Notificați când este gata.”
- De ce: Subordonata pusă după regentă nu primește automat virgulă. Nu se despart: completiva („Am observat că funcționează”), relativa determinativă („Fișierul pe care l-ați selectat”) și temporala pusă după regentă („Notificați când este gata”). Virgula se pune la cauzală (de regulă), la condițională (de regulă), la consecutivă și la subordonata intercalată care întrerupe regenta („Tu, dacă nu ești nevoit, nu faci nimic”). Relativa determinativă nu se desparte nici când stă în interior („Fata care râde este sora lui”). A cere virgulă înaintea oricărei subordonate e hipercorectitudine.
- Sursă: ghid gov-md §9.2.2.3, §9.2.2.5, §9.2.2.6 și §9.2.2.8; ÎOOP6 p. 59–61 (Virgula, viii, ix, xi)

### RO-TIP-114 Fără virgulă înaintea lui „sau” disjunctiv simplu
- Nivel: normă
- Greșit: „apasă pe legătură, sau lipește-o” → Corect: „apasă pe legătură sau lipește-o”
- De ce: „Sau” simplu, între părți de propoziție sau între propoziții coordonate, nu cere virgulă. Virgula e corectă în trei cazuri: la perechea „sau… sau”, la „sau” cu sens explicativ („adică”), pus între virgule, și când înaintea lui se închide o incidentă sau o apoziție.
- Sursă: ghid gov-md §9.2.1.9 și §9.2.2.1; INJ p. 67–70; ÎOOP6 p. 58 (Virgula, ii–iii)

### RO-TIP-115 Virgulă înaintea lui „nu” de opoziție
- Nivel: recomandare
- Greșit: „Scrie ș nu s.” → Corect: „Scrie ș, nu s.”
- De ce: Când „nu” opune doi termeni fără conjuncție („X, nu Y”), coordonarea adversativă eliptică cere virgulă. La fel înaintea lui „și nu” cu sens adversativ („ghidurile publicate, și nu cele vechi”). Nu se pune virgulă când „nu” neagă verbul: „Fișierul nu există.”
- Sursă: ghid gov-md §9.2.1.11 (N.B.) și §9.2.1.7

## 4. Ghilimelele

### Metoda: ce se citează, ce se pune între ghilimele

Ghilimelele românești sunt „…” (deschiderea jos, U+201E; închiderea sus, U+201D). Citatul din citat stă între «…». Al treilea nivel revine la „…”, în schema „… «… „…” …» …”. Nu amesteci sistemele: ghilimelele drepte ("…"), cele englezești (U+201C la deschidere), combinația „…“ (închidere cu U+201C) și ghilimelele cu apostrofuri ‘…’ pentru citat nu sunt forma românească.

Pașii:

1. Verifici deschiderea și închiderea fiecărui citat: jos–sus, cu aceeași pereche.
2. Punctul final: vezi mai jos (citat integrat sau frază completă după două-puncte).
3. Semnul întrebării și semnul exclamării care aparțin citatului stau înăuntru.
4. Între ghilimele și textul citat nu se lasă spațiu.
5. Citatul direct nu se corectează lingvistic, doar tipografic.

Slab:

> Directorul a spus: "Nu plec". Apoi a adăugat: “Nu mai e nimic de discutat“.

Bun:

> Directorul a spus: „Nu plec.” Apoi a adăugat: „Nu mai e nimic de discutat.”

Atenție: ghilimelele curbe nu sunt un semn de text generat. „…” sunt chiar ghilimelele românești; cele drepte nu sunt forma recomandată într-un text îngrijit (vezi RO-NUG-095).

### RO-TIP-054 Ghilimelele românești: „…” și «…»
- Nivel: recomandare
- Greșit: „"Nu plec", a spus el.” → Corect: „„Nu plec”, a spus el.”
- De ce: Citatul stă între „…”, iar citatul din citat între «…». Exemplu cu trei niveluri: „El a spus: «Colega mea a întrebat: „Cine vine?”»”. Ghilimelele drepte se înlocuiesc cu cele românești; ÎOOP6 formulează asta ca recomandare (p. 43).
- Sursă: ghid gov-md §9.8.1–9.8.3; ÎOOP6 p. 43 (Ghilimelele)

### RO-TIP-055 Ghilimelele englezești și combinațiile mixte nu sunt forma recomandată
- Nivel: recomandare
- Greșit: „„Nu plec“, a spus el.” → Corect: „„Nu plec”, a spus el.”
- De ce: În textul românesc nu sunt forma recomandată ghilimelele drepte, cele englezești cu deschidere U+201C și combinația „…“ (cu închidere U+201C, specifică germanei). Închiderea corectă e ” (U+201D). ÎOOP6 descrie doar ghilimelele românești „ ” și le recomandă (p. 43); nu menționează cele englezești sau germane. Înainte să corectezi, verifici fiecare cap al citatului: problema nu e că ghilimelele ar fi „curbe”, ci că perechea nu e jos–sus.
- Sursă: ghid gov-md §9.8.1; ÎOOP6 p. 43 (Ghilimelele)

### RO-TIP-056 Pentru citatul din citat, nu apostrofurile ‘…’
- Nivel: recomandare
- Greșit: „Antrenorul mi-a spus: ‘Lovește puternic’” → Corect: „Antrenorul mi-a spus: «Lovește puternic»”
- De ce: Citatul din citat se pune între «…». Apostrofurile ‘…’, propuse de unele ghiduri de presă, nu sunt forma recomandată pentru citat (ÎOOP6 p. 43 descrie « » pentru citatul din citat). Regula privește doar citatele: ghilimelele simple la explicarea unui cuvânt (deadline ‘termen-limită’) sunt admise (vezi RO-NUG-120).
- Sursă: ghid gov-md §9.8.2; ÎOOP6 p. 43 (Ghilimelele: « » pentru citat în citat; ‘ ’ doar pentru glosare)

### RO-TIP-057 Semnul întrebării sau al exclamării care aparține citatului stă înăuntru
- Nivel: normă
- Greșit: „El a întrebat: "Ce să fie"?” → Corect: „El a întrebat: „Ce să fie?””
- De ce: Semnul întrebării sau al exclamării care ține de citat se pune înaintea ghilimelelor de închidere și după el nu se mai pune punct. Dacă semnul aparține propoziției autorului („Ai citit „Moromeții”?”), stă după ghilimele.
- Sursă: ghid gov-md §9.1.2 și §9.4.4; ÎOOP6 p. 44 (Ghilimelele, i) și p. 56 (Semnul întrebării)

### RO-TIP-058 Punctul după citatul-frază introdus prin două-puncte
- Nivel: disputat
- Corect: „Directorul a spus: „Nu plec.”” (punctul înăuntru) și „Directorul a spus: „Nu plec”.” (punctul după ghilimele).
- De ce: Ghidul gov-md și Ghidul UE pun punctul înăuntru când citatul e o frază completă introdusă prin două-puncte; alte ghiduri de presă îl pun după ghilimele. Nu corectezi niciuna; ceri consecvență în tot textul.
- Sursă: ghid gov-md §9.1.2 și §9.8.6; Ghid UE §10.1 (nereverificat local)

### RO-TIP-059 Linia de dialog nu se combină cu ghilimelele
- Nivel: normă
- Greșit: „– "Bună ziua!"” → Corect: „– Bună ziua!”
- De ce: Replica introdusă cu linie de dialog nu mai primește ghilimele. Una dintre cele două convenții de redare a vorbirii directe se alege pentru tot textul.
- Sursă: ÎOOP6 p. 44 (Ghilimelele, greșeli i) și p. 45 (Linia de dialog, greșeli i)

### RO-TIP-104 Citatul integrat în frază: punctul după ghilimele
- Nivel: normă
- Greșit: „A vorbit despre "un buget care nu ajunge."” → Corect: „A vorbit despre „un buget care nu ajunge”.”
- De ce: Când citatul continuă fraza autorului, fără două-puncte, punctul stă după ghilimelele de închidere. Citatul-frază completă introdus prin două-puncte e cazul disputat de la RO-TIP-058.
- Sursă: ghid gov-md §9.1.2 și §9.8.5; ÎOOP6 p. 51 (Punctul, ii)

## 5. Linia de pauză, cratima și intervalele

### Metoda: trei semne diferite, trei funcții

- **Cratima** (-) leagă: cuvinte compuse (prim-ministru), articolul sau desinența de un cuvânt străin (site-ul), formanții ordinalelor (al XX-lea). Nu are spații în jur.
- **Linia de pauză** (–) separă: delimitează incize și apoziții, mai ales cele cu virgule în interior, și marchează un predicat omis. Se scrie cu spațiu de ambele părți.
- **Intervalele** (anii 2020-2021, orele 8-16) se scriu, după DOOM 3, cu cratimă fără spații. Linia de pauză cu spații e admisă, dar mai lungă.

Verifici un semn din mijlocul frazei astfel: dacă are spații în jur, trebuie să fie linie de pauză, nu cratimă. Dacă nu are spații, trebuie să lege două cuvinte sau două numere.

Slab:

> Ministrul - aflat la Bruxelles - a declarat că magazinul e deschis între 8-21.

Bun:

> Ministrul – aflat la Bruxelles – a declarat că magazinul e deschis între 8 și 21.

### RO-TIP-031 Numele semnelor de punctuație se scriu cu cratimă
- Nivel: recomandare
- Greșit: „punct și virgulă” → Corect: „punct-și-virgulă”
- De ce: Ca nume de semn, „două-puncte” și „punct-și-virgulă” sunt substantive compuse scrise cu cratimă, forma pe care o dă DOOM 3. Forma cu spații nu e o greșeală în vorbirea curentă; într-un text îngrijit alegi forma cu cratimă și o păstrezi.
- Sursă: DOOM 3, s.v. două-puncte, punct-și-virgulă; ghid gov-md §9.3 și §9.4

### RO-TIP-062 Linia de pauză este semn de punctuație românesc
- Nivel: normă
- Corect: „Toate partidele – inclusiv cele din opoziție, care au boicotat ședința – au semnat.”
- De ce: Linia de pauză („–”) delimitează incize și apoziții explicative, mai ales pe cele care conțin ele însele virgule, sau marchează un predicat omis. Regula după care „nu se folosește în română” e falsă. Se scrie cu spațiu de ambele părți.
- Sursă: ghid gov-md §9.9.3–9.9.5; ÎOOP6 p. 45–46 (Linia de pauză, funcții i–iii)

### RO-TIP-063 Cratima nu ține loc de linie de pauză
- Nivel: normă
- Greșit: „Ministrul - aflat la Bruxelles - a declarat” → Corect: „Ministrul – aflat la Bruxelles – a declarat”
- De ce: Cratima cu spații în jur, în mijlocul frazei, e o linie de pauză scrisă greșit: nu „-” cu spații. Linia de pauză se scrie „–”, cu spații; cratima, fără spații, leagă cuvinte sau numere. Excepție: liniuța de la începutul rândului, ca marcaj de listă.
- Sursă: ghid gov-md §9.9.1 și §9.11.9; DOOM 3 §1.2.5
- Vezi și: RO-TIP-062

### RO-TIP-064 Virgula necesară se pune după a doua linie
- Nivel: normă
- Greșit: „Dacă ajungeți la timp, – și sperăm că veți ajunge – vă așteptăm la intrare.” → Corect: „Dacă ajungeți la timp – și sperăm că veți ajunge –, vă așteptăm la intrare.”
- De ce: Când incidenta dintre linii urmează unei construcții după care fraza cere oricum virgulă, virgula se pune după a doua linie, nu înaintea primei. Dacă fraza n-ar fi cerut virgulă, nu se pune nici după a doua linie.
- Sursă: ghid gov-md §9.9.4; ÎOOP6 p. 47 (Linia de pauză, greșeli ii–iii)

### RO-TIP-065 Intervalele: cratimă fără spații
- Nivel: recomandare
- Corect: „anii 2020-2021, orele 8-16, articolele 12-14, perioada 1-23 august, două-trei zile”
- De ce: Pentru limite de timp, de spațiu și pentru aproximări DOOM 3 și ghidul gov-md preferă cratima fără spații; linia de pauză cu spații e admisă de DOOM 3, dar socotită mai lungă fără rost. Ghidul gov-md recomandă linia de pauză doar între două părți (RO-TIP-066) sau când cratima poate fi confundată cu un nume compus. Propui cratima, fără spații, și ceri consecvență; linia de pauză nu se corectează ca greșeală.
- Sursă: ghid gov-md §9.9.7 și §9.11.8; DOOM 3 §1.2.5 II.2–II.6; ÎOOP6 p. 41 (Cratima, funcții ii–iii) și p. 46 (Linia de pauză, iv)

### RO-TIP-066 Relația dintre două părți: linia de pauză sau cratima
- Nivel: disputat
- Corect: „Acordul România – Franța” și „Paris-Londra”.
- De ce: Pentru relația dintre două părți (întâlnirea NATO – Rusia, acordul România – Franța) DOOM 3 și ghidul gov-md preferă linia de pauză, mai ales dacă un element are cratimă internă (distanța Iași – Cluj-Napoca). ÎOOP6 folosește cratima pentru cuvinte comune („raportul semnificant-semnificat”), dar cere linia de pauză când elementele sunt nume proprii sau conțin cratimă (p. 41, p. 46). Nu corectezi; ceri aceeași alegere în tot textul.
- Sursă: ghid gov-md §9.9.7 și §9.11.8; DOOM 3 §1.2.5 II.7; ÎOOP6 p. 41 (Cratima, funcții i, greșeli ii) și p. 46 (Linia de pauză, iv)

### RO-TIP-067 Intervalul: cu prepoziții sau cu semn, nu amândouă
- Nivel: normă
- Greșit: „deschis între 10-18” → Corect: „deschis între 10 și 18”
- De ce: Cratima sau linia de pauză înseamnă „între” sau „dintre”; scrisă după „între”, spune de două ori același lucru. Alegi: „între 10 și 18” sau „10-18”. La fel „distanța Ungheni-Chișinău” sau „distanța dintre Ungheni și Chișinău”, nu „distanța dintre Ungheni-Chișinău”.
- Sursă: ghid gov-md §9.11.8 (N.B.); ÎOOP6 p. 41 (Cratima, greșeli i)

### RO-TIP-121 Titlul cu subtitlu: linie de pauză cu spații
- Nivel: recomandare
- Greșit: „Partea I - Arhitectura” → Corect: „Partea I – Arhitectura”
- De ce: Între titlu și subtitlu, ca și între părțile unui titlu, linia de pauză se scrie cu spațiu de ambele părți; cratima în locul ei e greșită. Linia lungă (—) nu e greșită normativ, dar se uniformizează la „–”, ca să nu amesteci două semne în același text. Folosirea liniei de pauză se dozează: un text în care apare în aproape fiecare frază are un tic de stil (vezi RO-NUG-030).
- Sursă: ghid gov-md §9.9.1 și §9.9.7

## 6. Apostroful, spațiile și celelalte semne

### Metoda: ce nu are spațiu, ce are

- Virgula, punctul, punctul-și-virgula, două-puncte, semnul întrebării și semnul exclamării se lipesc de cuvântul dinainte și sunt urmate de spațiu.
- Între ghilimele sau paranteze și textul din interior nu se lasă spațiu.
- Cratima, bara oblică și punctul din interiorul unei abrevieri nu au spații: „a.c.”, „km/h”, „lt.-maj.”.
- Excepții la „spațiu după semn”: virgula zecimală (3,5), ora (12:30), numărul cu puncte (1.500), adresele web și de e-mail, punctul din abrevieri, semnul care e urmat de alt semn (ghilimele, paranteză închisă).

Verificări rapide (mecanice, dar cu ochii pe excepții):

1. spații duble;
2. spațiu înaintea virgulei sau a punctului („cuvânt ,” → „cuvânt,”);
3. cuvânt repetat imediat („de de”, „la la”), care se verifică pe text: unele repetiții sunt corecte („A plecat și și-a luat geanta”);
4. bara oblică din „km/h” se păstrează și nu devine „km sau h”;
5. ghilimelele de închidere sunt ” (U+201D), nu ghilimeaua dreaptă.

### RO-TIP-032 Apostroful ’ în cuvintele românești
- Nivel: normă
- Greșit: „D'ale carnavalului, dup' o oră” → Corect: „D-ale carnavalului, dup-o oră”
- De ce: DOOM 3 §1.2.1 și ÎOOP6 p. 33 cer cratima, nu apostroful singur sau apostroful plus cratimă, când rezultă o silabă nouă: corect „d-ale carnavalului”, nu „*d’ale” și nici „*d’-ale”. Apostroful se folosește la anii prescurtați (’89, anii ’90, ante-’89), la rostirea familiară reprodusă (dom’le, dom’ profesor) și în nume străine (O’Neill, McDonald’s, D’Annunzio); în citate și în titluri vechi se păstrează grafia originală. Forma tipografică corectă a apostrofului este ’ (U+2019), care arată ca o virgulă ridicată; nu o confunda cu semnul prim (′) din matematică și nici cu apostroful drept de la tastatură, pe care un text îngrijit îl înlocuiește cu ’. În documentele oficiale anii se scriu întregi. Apostroful nu se scoate din formele de mai sus (vezi RO-NUG-096).
- Sursă: ghid gov-md §9.13.1–9.13.5; DOOM 3 §1.2.1; ÎOOP6 p. 33 (Cratima și apostroful)

### RO-TIP-068 Spațiul în jurul semnelor de punctuație
- Nivel: normă
- Greșit: „Ministrul a spus , apoi a plecat .” → Corect: „Ministrul a spus, apoi a plecat.”
- De ce: Virgula, punctul, punctul-și-virgula, două-puncte, semnul întrebării, semnul exclamării și punctele de suspensie se lipesc de cuvântul dinainte și sunt urmate de spațiu (sau de alt semn de punctuație): „Gata. Continuă”, nu „Gata.Continuă”. Excepțiile sunt cele din metoda de mai sus: punctul din abrevieri, virgula zecimală, ora, numerele, adresele web și de e-mail, și cazul în care urmează alt semn. Un verificator mecanic care ignoră aceste excepții strică textul.
- Sursă: ghid gov-md §9.1.11, §9.2.1.1, §9.3.3, §9.4.7, §9.10.2 și §9.14; DOOM 3 §1.2.3 pct. 4 și §1.2.6 I

### RO-TIP-069 Fără spații în jurul cratimei, al barei oblice și în interiorul abrevierilor
- Nivel: normă
- Greșit: „km / h, a . c., prim - ministru” → Corect: „km/h, a.c., prim-ministru”
- De ce: Cratima, bara oblică și punctele din interiorul abrevierilor se scriu fără spații. Excepția e cratima-marcaj de listă, la început de rând. Ce compuse se scriu cu cratimă spune RO-GRA-082 din [gramatica.md](gramatica.md); scrierea siglelor (fără puncte între litere), RO-GRA-090 din același fișier.
- Sursă: ghid gov-md §9.1.11, §9.11.9 și §9.12.4; DOOM 3 §1.2.3 pct. 4; ÎOOP6 p. 34 (Bara oblică) și p. 38–39 (Punctul)

### RO-TIP-035 Titlurile, subtitlurile și legendele nu se încheie cu punct
- Nivel: normă
- Greșit: „Guvernul aprobă bugetul.” → Corect: „Guvernul aprobă bugetul”
- De ce: Titlurile (inclusiv de acte și de capitole), subtitlurile și legendele foto nu se încheie cu punct. Semnul întrebării, semnul exclamării și punctele de suspensie, dacă sunt necesare, rămân. Titlul format din mai multe unități independente poate avea punct în interior („Capitolul I. Dispoziții generale. Definiții”).
- Sursă: ghid gov-md §9.1.5; ÎOOP6 p. 52 (Punctul, greșeli i)

## 7. Cursivele pentru cuvinte străine

### RO-TIP-060 Cuvintele și expresiile străine neadaptate se scriu cu cursive
- Nivel: recomandare
- Corect: „Raportul folosește sintagma *intelligence-led policing*, tradusă aici „poliție bazată pe informații”.”
- De ce: Cuvintele și expresiile latinești și cele din alte limbi, neadaptate, se scriu cu cursive; împrumuturile adaptate (ad-hoc, acquis, post-mortem, site) se scriu drept. Dacă textul nu poate avea cursive (text simplu, interfață), lași cuvântul neschimbat, fără să-l pui între ghilimele ca înlocuitor. Denumirile traduse în română se scriu drept.
- Sursă: ghid gov-md §11.2.1.2 și §11.2.1.3; Ghid UE §10.5 (nereverificat local)

### RO-TIP-061 Ghilimele și cursive nu se pun pe același cuvânt
- Nivel: recomandare
- Corect: „Programul Start pentru tineri” (fără ghilimele, cu cursive la denumirea unui program) sau „Programul «Start pentru tineri»” (cu ghilimele, fără cursive).
- De ce: Dubla evidențiere (ghilimele și cursive) nu se recomandă. Alegi un singur mijloc pentru același cuvânt sau aceeași denumire și îl păstrezi în tot textul.
- Sursă: ghid gov-md §11.2.2; Ghid UE §10.1.10.3 (nereverificat local)

## 8. Nume proprii străine și denumiri de state

### Metoda: ce rămâne în original, ce are formă românească

Pașii:

1. Numele din alfabet latin rămân în forma originală, cu diacriticele lor: München, Škoda, Dvořák, Mickiewicz.
2. Locurile cunoscute de multă vreme au forme românești consacrate (exonime), pe care le folosești în textul curent: Londra, Moscova, Viena, Varșovia, Florența, Neapole, Nisa.
3. Numele din alte alfabete (rusă, bulgară, chineză etc.) se transcriu după regulile românești, nu după forma englezească. Numele sârbești scrise cu chirilice se redau în alfabetul latin sârbesc (RO-TIP-079).
4. Statele se numesc în forma oficială actuală, iar pentru formele concurente alegi una și o păstrezi.
5. Derivatele păstrează grafia numelui: shakespearian, darwinism, newtonian.
6. Articularea urmează pronunția: Bordeaux-ul, show-ul, dar Kremlinul (vezi RO-IMP-004 și RO-IMP-005 în [imprumuturi.md](imprumuturi.md)).

Slab:

> Președintele Macron a fost la Munich și a vorbit despre Pyotr Ilyich Tchaikovsky.

Bun:

> Președintele Macron a fost la München și a vorbit despre Piotr Ilici Ceaikovski.

### RO-TIP-070 Majusculele urmează regulile românei, nu ale limbii sursei
- Nivel: normă
- Greșit: „Prime Minister Ion Popescu” → Corect: „prim-ministrul Ion Popescu”
- De ce: Într-un text românesc, funcțiile și substantivele comune se scriu cu literă mică, iar numele oficiale complete cu majusculă la cuvintele pline, după regulile din secțiunea 2, nu după regulile limbii din care vine informația.
- Sursă: ghid gov-md §4.1.7.1 (N.B.); DOOM 3 §3.1 și §3.2 II.2.2

### RO-TIP-071 Numele din alfabet latin: forma originală, cu diacriticele ei
- Nivel: normă
- Greșit: „Munich, Skoda” → Corect: „München, Škoda”
- De ce: Numele proprii din limbi cu alfabet latin se scriu ca în original, cu diacriticele lor: München, Mickiewicz, Petőfi, Márquez, São Tomé. Nu le anglicizezi și nu le simplifici (Skoda, Munchen). Regula se aplică și sedilei: într-un nume turcesc, litera s cu sedilă (U+015F) e corectă, iar sedila se corectează doar în cuvintele românești. Pentru orașele care au o formă românească tradițională (Cracovia, Londra) se folosește aceasta, iar exonimul rămâne corect; vezi RO-TIP-072.
- Sursă: ghid gov-md §3.5.1 și §2.2.1; DOOM 3 §2.5 pct. 1

### RO-TIP-072 Exonimele consacrate se folosesc în textul curent
- Nivel: recomandare
- Greșit: „London, Warszawa, Wien” → Corect: „Londra, Varșovia, Viena”
- De ce: Locurile cunoscute de mult au forme românești curente: Londra, Moscova, Varșovia, Viena, Praga, Florența, Nisa, Salonic, Neapole. Forma originală (London, Wien) e folosită în lucrări de specialitate, pe hărți și în bibliografii; într-un text românesc obișnuit, exonimul e forma firească.
- Sursă: ghid gov-md §3.5.3; DOOM 3 §2.5 pct. 3
- Vezi și: RO-TIP-071

### RO-TIP-073 Derivatele păstrează grafia numelui
- Nivel: normă
- Corect: „versuri shakespeariene”. Derivatele păstrează grafia numelui de la care vin: shakespearian, darwinism, washingtonian.
- De ce: Un derivat de la un nume propriu se scrie cu grafia numelui, nu se „românizează”. Forma „shakespearian” e în dicționar.
- Sursă: ghid gov-md §3.5.4; DEX ’98, s.v. shakespearian

### RO-TIP-074 Toponimul fără normă românească: nicio formă nu e greșeală
- Nivel: ce nu e greșeală
- Corect: „Kyiv” și „Kiev”, folosite consecvent într-un text, nu se corectează una prin alta.
- De ce: Când nu există o normă românească pentru un toponim (două forme concurente, una cerută de statul respectiv, alta tradițională), niciuna nu e greșeală: alegi una și o păstrezi. Semnalezi doar inconsecvența, nu forma.
- Sursă: ghid gov-md §3.5.3 (nereverificat local pentru cazul concret)

### RO-TIP-076 În presă, transcrierea după convențiile românești
- Nivel: normă
- Corect: „Cehov, nu Chekhov”. Transliterarea ISO („Čehov”) se folosește în acte, în bibliografie și în cartografie, nu în presă.
- De ce: Numele din alfabete nelatine se transcriu după regulile românești, nu se preiau din forma englezească. Transliterarea standardizată (ISO 9 pentru chirilică) servește unor scopuri tehnice.
- Sursă: ghid gov-md §3.2 și §3.4.8; DOOM 3 §1.3.1–1.3.2
- Vezi și: RO-TIP-077

### RO-TIP-077 Rusa
- Nivel: normă
- Greșit: „Chekhov, Tchaikovsky, Yeltsin, Navalny, Kyrgyzstan” → Corect: „Cehov, Ceaikovski, Elțân, Navalnâi, Kârgâzstan”
- De ce: Numele rusești se transcriu după tabelul românesc, nu după forma englezească: ч devine c(e/i), ш devine ș, щ devine șc(i), ж devine j, х devine h, ц devine ț, ы devine â în interior (î la sfârșit de cuvânt), -ий devine i, я devine ia la început de cuvânt și după vocală și ea după consoană, к devine k, е rămâne e (nu ie), г rămâne g (gh înaintea lui e și i). Excepții tradiționale: Dostoievski, Turgheniev.
- Sursă: ghid gov-md §3.3 și §3.4; DOOM 3 §1.3.3.1

### RO-TIP-078 Bulgara
- Nivel: normă
- Corect: „ъ → ă, щ → șt”
- De ce: Numele bulgărești transcrise din chirilică redau ъ prin ă și щ prin șt; г și к urmate de ъ se redau gh, ch.
- Sursă: DOOM 3 §1.3.3.2
- Vezi și: RO-TIP-077

### RO-TIP-079 Sârba
- Nivel: normă
- Greșit: „Vucic, Vuciciu” → Corect: „Vučić”
- De ce: Numele sârbești scrise în chirilică se redau în forma alfabetului latin sârbesc (Đ, Ž, Lj, Nj, Ć, Č, Dž, Š), nu după regulile pentru rusă și nu în forma anglicizată.
- Sursă: DOOM 3 §1.3.3.3
- Vezi și: RO-TIP-077

### RO-TIP-080 Chineza, araba, ebraica, greaca
- Nivel: recomandare
- Greșit: „Peking” → Corect: „Beijing”
- De ce: Chineza se redă în Pinyin, cu forma „Beijing” pentru capitală. DOOM 3 trimite de la „Pekin” la „Beijing”, iar „Peking” nu e o formă românească. Pentru arabă, ebraică și greacă se folosesc standardele de transliterare (ISO 233, ISO 259-2, ISO 843) în acte și bibliografie.
- Sursă: DOOM 3, s.v. Beijing, Pekin; DOOM 3 §1.3.4–1.3.7
- Vezi și: RO-TIP-077

### RO-TIP-081 Statele se numesc în forma oficială actuală
- Nivel: normă
- Greșit: „Bielorusia, Swaziland, Coasta de Fildeș” → Corect: „Belarus, Eswatini, Côte d’Ivoire”
- De ce: DOOM 3 §2.5 pct. 2 cere forma oficială actuală pentru numele statelor: Belarus (adjectivul „belarus”; „Bielorusia” e învechit), Côte d’Ivoire („Coasta de Fildeș” se corectează; DOOM 3 o marchează învechit, iar ÎOOP6 dă și el forma oficială), Myanmar (DOOM 3: „fosta Birmanie”) și Eswatini (DOOM 3: „numele actual al fostului stat Swaziland”). Când textul vorbește despre statul dinainte de schimbarea numelui (Swaziland înainte de 2018, Birmania înainte de 1989), numele de atunci nu e greșeală; poți propune paranteza „(azi Eswatini)”. Excepția nu se aplică la Bielorusia și Coasta de Fildeș, pe care DOOM 3 le marchează învechit. Dacă un ghid instituțional cere altă formă, în textul redactat urmezi DOOM 3 și semnalezi diferența. Cehia (nume scurt) și Capul Verde sunt la RO-TIP-082, ca recomandare.
- Sursă: DOOM 3, s.v. Belarus, Myanmar, Eswatini, Swaziland, Côte d’Ivoire, Coasta de Fildeș; DOOM 3 §2.5 pct. 2; ÎOOP6 (s.v. Côte d’Ivoire); ghid gov-md §3.5.2

### RO-TIP-082 Forma completă a statului sau forma scurtă
- Nivel: recomandare
- Corect: „Acordul se încheie cu Republica Italiană.” (statul ca parte la un act) și „Exporturile din Italia au crescut.” (statul ca spațiu geografic sau economic)
- De ce: Forma oficială, completă, se folosește când statul e subiect de drept (parte la un tratat, titular al unui drept); forma scurtă când e vorba de spațiul geografic, economic sau cultural. Într-un text de presă sau de marketing forma scurtă e curentă. Forma scurtă „Cehia” e curentă. La „Capul Verde” DOOM 3 dă numai această formă (§1.2.2, intrare cu norma modificată); „Cabo Verde” nu figurează în DOOM 3, deci o semnalezi și propui „Capul Verde”, fără s-o corectezi ca greșeală.
- Sursă: ghid gov-md §15.2.1 și §15.2.2

### RO-TIP-083 Statul și regiunea nu se confundă
- Nivel: recomandare
- Corect: „guvernul Regatului Unit” și „guvernul Țărilor de Jos”, în texte oficiale; „premierul Marii Britanii” și „guvernul olandez” sunt uzuale în texte neoficiale.
- De ce: Numele statului nu se confundă cu al unei regiuni: Regatul Unit (statul) cuprinde Marea Britanie (insula) și Irlanda de Nord; Țările de Jos (statul, adjectiv „neerlandez”) cuprind regiunea Olanda. În acte și traduceri juridice se folosesc numele statelor. „Olanda” și „olandez” sunt în dicționar și nu se corectează în presă. Pentru Republica Moldova se scrie „din Republica Moldova”, cu forma completă.
- Sursă: Ghid UE, Anexa A5 (note UK1, NL1) (nereverificat local); ghid gov-md §15.1.1

### RO-TIP-084 „Elen” pentru stat, „grec” pentru limbă
- Nivel: recomandare
- Corect: „guvernul elen”, „limba greacă”
- De ce: În practica de traducere instituțională, „elen” se folosește pentru stat (guvernul elen), iar „grec/greacă” pentru limbă. În uzul curent „guvernul grec” e la fel de corect.
- Sursă: IER, Ghid stilistic de traducere, 2008 (nereverificat local)

### RO-TIP-087 Instituția fără echivalent consacrat rămâne în original
- Nivel: recomandare
- Corect: „Bundestag (Parlamentul Federal al Germaniei)” la prima mențiune, apoi „Bundestagul” sau „parlamentul german”.
- De ce: Instituția sau funcția fără echivalent românesc consacrat rămâne în original, cu explicația la prima mențiune, apoi o formă scurtă românească. Instituțiile cu denumire românească consacrată (Parlamentul European) se scriu în română.
- Sursă: ghid gov-md §15.2.4; practica UE (nereverificat local)
- Vezi și: RO-TIP-085

### RO-TIP-088 Sigla puțin cunoscută se desfășoară la prima apariție
- Nivel: recomandare
- Corect: „Agenția Internațională pentru Energie Atomică (AIEA) a publicat raportul. AIEA recomandă...”
- De ce: La prima apariție sigla puțin cunoscută se introduce între paranteze, după denumirea desfășurată, apoi se folosește forma scurtă. Sigla se folosește cu măsură și cu aceeași formă în tot textul. Siglele pe care le știe orice cititor (NATO, UE, ONU) nu se explică. Regula e valabilă în orice gen de text, presa o aplică la fel (RO-STI-619 din [presa.md](tipuri/presa.md)).
- Sursă: ghid gov-md §5.1.3.3 și §5.2.2.1; Mic ghid CE p. 11–13 (Rec. 9)

## 9. Localitățile din România

### RO-TIP-033 Numele localităților se scriu cum le dă legea
- Nivel: recomandare
- Greșit: „Târgu-Mureș” → Corect: „Târgu Mureș”
- De ce: Denumirile oficiale ale localităților sunt cele din anexa Legii 2/1968 privind organizarea administrativă a teritoriului României, așa cum a fost modificată prin Legea 290/2018. Forma din anexă pentru reședința județului Mureș e „Târgu Mureș”, fără cratimă, iar „Târgu-Neamț” are cratimă; „Cluj-Napoca”, „Satu Mare” și „Râmnicu Vâlcea” se scriu și ele ca în anexă. Folosirea locală nu e uniformă (autoritățile locale și presa scriu și cu cratimă), de aceea abaterea se propune, nu se impune. Orașele pot fi articulate în text: Bucureștiul, Clujul.
- Sursă: Legea 2/1968 privind organizarea administrativă a teritoriului României, anexa, în forma dată de Legea 290/2018 (nereverificat local)

## 10. Urmele lăsate de un chatbot

Pe lângă formulele de conversație și celelalte tipare din [tipare-text-generat.md](tipare-text-generat.md), un text copiat dintr-un chatbot poartă adesea urme de formă. Le cauți în ultima trecere și le tratezi ca semnale, nu ca greșeli de normă:

- asteriscuri de formatare rămase în text („**termen**”) sau diezi de titlu într-un text care nu e Markdown;
- emojiuri decorative puse la începutul rândurilor sau al titlurilor;
- liniuțe lungi în aproape fiecare frază, ca tic (linia de pauză în sine e corectă, vezi RO-NUG-030);
- fraza introductivă a răspunsului lăsată în text: „Iată textul revizuit:”, „Sigur! Iată…”;
- ghilimele drepte ("…") puse în locul celor românești; atenție, ghilimelele curbe „…” sunt tocmai cele corecte;
- titluri cu majuscule pe fiecare cuvânt (RO-TIP-010).

Un semn izolat nu dovedește nimic. Grupul de semne se semnalează ca „de semnalat”, cu propunerea de a curăța textul, și nu schimbi fondul. Asta privește doar originea textului: un semn care încalcă și o regulă proprie își păstrează nivelul ei, chiar izolat (titlul cu majuscule pe fiecare cuvânt e normă, RO-TIP-010; fraza introductivă se scoate, RO-STI-207).

