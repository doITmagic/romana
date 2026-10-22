# Interfață

Profilul acesta se aplică butoanelor, meniurilor, titlurilor, mesajelor, notificărilor, etichetelor și șirurilor de traducere dintr-o aplicație sau dintr-un site. Textul de interfață e scurt, e citit în trecere, apare adesea ca șir izolat, fără contextul din jur, și primește valori inserate automat (nume, numere, date). De aceea forma lui verbală depinde de tipul de text, valorile inserate nu primesc articol sau desinență, pluralul se scrie cu trei forme, iar etichetele se acordă cu ce descriu. Normele limbii nu se schimbă.

Abrevierile din câmpul „Sursă”: DOOM 3 = *Dicționarul ortografic, ortoepic și morfologic al limbii române*, ediția a treia; DEX ’98 = *Dicționar explicativ al limbii române*, 1998; DSL = A. Bidu-Vrănceanu et al., *Dicționar general de științe. Științe ale limbii*, 1997; Mic ghid CE = *Mic ghid de redactare clară*, Comisia Europeană, 2013; INJ = *Corectitudine lingvistică în actele judecătorești*, Institutul Național al Justiției, Chișinău, 2020; ghid gov-md = *Ghid de redactare lingvistică și tehnoredactare a actelor normative*, Cancelaria de Stat a Republicii Moldova, 2026; Mozilla RO = ghidul de localizare în limba română al Mozilla, care fixează convenții de uz, nu norme ale limbii (nereverificat local); CLDR = regulile de plural pentru română din datele de localizare Unicode CLDR; ICU = formatul de mesaje ICU MessageFormat; gettext = formula de plural GNU gettext.

## Cum citești perechile

La o recomandare, forma din „Greșit” e forma de evitat, nu neapărat o greșeală de limbă. Dacă forma e corectă (RO-INT-013, 022, 025, 057, 058), regula o spune în „De ce” și nu o corectezi, ci o propui. Doar nivelul „normă” cere corectare.

## Ce se schimbă față de limba literară standard

- **Registrul și adresarea.** O decizie pentru tot produsul (secțiunea următoare).
- **Tonul.** Scurt, politicos, fără oralitate. Mesajul spune ce s-a întâmplat și ce poate face utilizatorul; nu povestește și nu dramatizează.
- **Structura.** Forma verbală se alege după tipul de șir (tabelul de mai jos). Nu amesteci formele în aceeași listă de butoane sau de meniuri (RO-COE-020 din [coerenta.md](../coerenta.md)).
- **Vocabularul.** Termenii de interfață au echivalente românești uzuale și rămân aceiași în tot produsul. Mai multe dintre ele sunt convenții de uz, nu norme, și se propun ca recomandări.
- **Valorile inserate.** Un șir cu valoare inserată se scrie ca să fie corect pentru orice valoare: fără articol, fără desinență, fără adjectiv acordat cu valoarea.
- **Ce nu se schimbă.** Diacriticele (ș, ț cu virgulă), ghilimelele „…”, majusculele (doar la primul cuvânt al unei etichete și la nume proprii) și fidelitatea față de șirul sursă (vezi [tipografie.md](../tipografie.md) și [fidelitate.md](../fidelitate.md)).

## Registrul este o decizie pentru tot produsul

Produsul alege o singură dată între „tu” și „dumneavoastră” și le folosește în toate mesajele, în toate limbile aplicației care folosesc aceeași voce și chiar în aceeași frază. Verbul, imperativul și posesivul urmează alegerea: „Salvează” și „contul tău” cu „tu”, „Salvați” și „contul dumneavoastră” cu „dumneavoastră”. Schimbarea adresării la jumătatea produsului e o rupere de coerență (RO-COE-019 din [coerenta.md](../coerenta.md)); criteriile alegerii sunt în [registru-si-adresare.md](../registru-si-adresare.md). Posesivul care traduce engl. „your” se scurtează când posesorul e clar (RO-CAL-105 din [prieteni-falsi-si-calcuri.md](../prieteni-falsi-si-calcuri.md)); „Salvați-vă munca dumneavoastră” e gramatical și cel mult o recomandare de scurtare (RO-SEM-216 din [paronime-si-pleonasme.md](../paronime-si-pleonasme.md)).

### RO-INT-031 Cererile și întrebările către utilizator, la persoana a II-a
- Nivel: recomandare
- Greșit: „Se dorește continuarea?” → Corect: „Dorești să continui?”
- De ce: Întrebările și cererile către utilizator se adresează direct, la persoana a II-a, în registrul ales. Exemplele din ghidul Mozilla RO sunt la „tu” („Dorești să continui?”, „Te rugăm să încerci din nou mai târziu”), deși același ghid cere o adresare „respectuoasă”: produsul care alege „dumneavoastră” scrie „Doriți să continuați?”.
- Sursă: Mozilla RO (uz, nu normă)

### RO-INT-032 Vocea produsului este „noi”
- Nivel: recomandare
- Greșit: „Vă rog să reîncercați mai târziu.” → Corect: „Vă rugăm să reîncercați mai târziu.”
- De ce: Produsul vorbește în numele echipei care îl face, deci cererile politicoase sunt la persoana I plural: „Te rugăm…” sau „Vă rugăm…”, după registru. Sursa engleză „Please” nu se traduce cu „Vă rog”, pentru că aplicația nu e o singură persoană. Cu registrul „tu”: „Te rugăm să încerci din nou mai târziu.”
- Sursă: Mozilla RO (uz, nu normă)

### RO-INT-012 Fără abrevieri și forme contrase în interfață
- Nivel: recomandare
- Greșit: „Contul dvs. nu-i activ.” → Corect: „Contul dumneavoastră nu este activ.”
- De ce: În interfață se evită abrevierile („dvs.”) și formele contrase („nu-i” pentru „nu este”), pentru că produsul vorbește în nume propriu, într-un registru neutru. „Dvs.” și „nu-i” sunt forme literare corecte (ghidul gov-md înregistrează abrevierea „dvs.”), deci e o convenție de stil, nu o normă; cu adresare prin „tu”, „dumneavoastră” dispare. Contragerea „nu-i” se evită doar când „-i” e verbul „a fi”.
- Sursă: Mozilla RO (uz, nu normă); ghid gov-md §5.3.1.6

### RO-INT-040 Fără acord cu genul necunoscut al utilizatorului
- Nivel: recomandare
- Greșit: „Ești sigur că vrei să ștergi acest semn de carte?” → Corect: „Sigur dorești să ștergi acest semn de carte?”
- De ce: Genul utilizatorului nu se cunoaște, iar adjectivul „sigur” din „ești sigur” se acordă cu el. Folosit ca adverb, „sigur” e invariabil („cu siguranță”), deci propoziția rămâne corectă pentru orice utilizator. Masculinul generic nu e o greșeală gramaticală; formularea neutră e o recomandare.
- Sursă: DEX ’98 s.v. sigur

## Forma verbală pe tipuri de șir

| Tip de șir | Formă | Exemplu |
|---|---|---|
| Buton sau comandă | imperativ (sau substantiv verbal), cel ales pentru tot produsul | „Salvează” / „Salvați”, „Deschide”, „Salvare” |
| Buton sau meniu care deschide alt meniu sau un dialog | substantiv verbal, nominal | „Gestionare…” |
| Titlu de fereastră, titlu de pagină, meniu | substantiv | „Setări”, „Istoric”, „Salvarea configurației” |
| Acțiune în curs | reflexiv cu „se” | „Se salvează…” |
| Rezultat finalizat | participiu | „Modificări salvate” |
| Listă goală | „Niciun…” sau „Nu există…” | „Niciun rezultat” |
| Eroare | ce s-a întâmplat și ce urmează | „Fișierul nu a putut fi procesat. Verifică formatul.” |
| Confirmare | întrebare directă | „Sigur dorești să ștergi acest semn de carte?” |
| Stare raportată de sistem | reflexiv sau impersonal | „Nu s-a putut deschide fișierul” |
| Etichetă de stare | adjectiv acordat cu referentul | „Cont activ”, „Notificări activate” |

### RO-INT-019 Butoanele, la imperativ în registrul ales
- Nivel: recomandare
- Greșit: „Salvezi modificările” → Corect: „Salvează modificările”
- De ce: Butoanele și comenzile se scriu la imperativ, în registrul ales („Salvează” cu „tu”, „Salvați” cu „dumneavoastră”; „Deschide”, „Caută fișiere”), sau ca substantiv verbal („Salvare”). Varianta aleasă rămâne la fel pe tot produsul, iar indicativul nu se folosește pentru comenzi.
- Sursă: Mozilla RO (uz, nu normă)

### RO-INT-026 Descrierea nu seamănă cu o comandă
- Nivel: recomandare
- Greșit: „Validează câmpurile” → Corect: „Se validează câmpurile…”
- De ce: La verbele de conjugarea I, persoana a III-a a indicativului prezent coincide cu imperativul persoana a II-a („validează”). Un cititor poate lua drept comandă descrierea unei acțiuni automate. Pentru descrieri folosești substantivul verbal sau reflexivul.
- Sursă: Coteanu, *Gramatica de bază a limbii române*, 1982 (imperativul)

### RO-INT-029 Butonul care deschide alt meniu, la substantiv verbal
- Nivel: recomandare
- Greșit: „Gestionează…” → Corect: „Gestionare…”
- De ce: Ghidul Mozilla RO recomandă forma impersonală, nominală, pentru butoanele care doar deschid un meniu sau un dialog suplimentar (engl. „Manage…”). Imperativul („Gestionează…”) nu e greșeală, iar ghidul trimite la o discuție deschisă despre imperativ și impersonal.
- Sursă: Mozilla RO (uz, nu normă); DOOM 3 s.v. gestiona

### RO-INT-030 Titlurile de fereastră, la substantiv
- Nivel: recomandare
- Greșit: „Salvează configurația” → Corect: „Salvarea configurației”
- De ce: Titlul unei ferestre e nominal, nu o comandă. Comanda stă pe butonul din fereastră.
- Sursă: Mozilla RO (uz, nu normă)

### RO-INT-034 Acțiunea în curs, la reflexiv
- Nivel: recomandare
- Greșit: „Salvat…” → Corect: „Se salvează…”
- De ce: O acțiune încă neterminată se scrie cu reflexivul și elipsă („Se salvează…” sau „Se salvează modificările…”), nu cu participiul. „În curs de salvare…” e admis. În titlul unei ferestre se folosește forma nominală („Salvarea configurației”).
- Sursă: Mozilla RO (uz, nu normă)

### RO-INT-033 Starea raportată de sistem, cu reflexivul
- Nivel: recomandare
- Greșit: „Nu am putut deschide fișierul.” → Corect: „Nu s-a putut deschide fișierul.”
- De ce: Calculatorul nu e o persoană; starea raportată de sistem se formulează cu reflexivul („Nu s-a putut deschide fișierul %s”, „Se primesc date de la %s”). Variantele impersonale („Imposibil de deschis fișierul”) sunt admise.
- Sursă: Mozilla RO (uz, nu normă)

### RO-INT-035 Finalizarea, la participiu
- Nivel: recomandare
- Greșit: „Salvare de succes” → Corect: „Modificări salvate”
- De ce: Mesajul de finalizare se formulează cu participiul: „Modificări salvate”, „Plată efectuată”. „De succes” e corect în sensul „care are succes” („afacere de succes”), deci nu se corectează în afara mesajelor de stare.
- Sursă: DEX ’98 s.v. succes

### RO-INT-036 „Cu succes” ca umplutură
- Nivel: recomandare
- Greșit: „Modificările au fost salvate cu succes.” → Corect: „Modificările au fost salvate.”
- De ce: Dacă mesajul apare după reușită, „cu succes” nu adaugă nimic. E o preferință de concizie: locuțiunea nu e greșită și rămâne firească unde reușita nu reiese din context („a absolvit cu succes”). Despre „Încărcare reușită” vezi RO-NUG-090 din [ce-nu-e-greseala.md](../ce-nu-e-greseala.md).
- Sursă: DEX ’98 s.v. succes

### RO-INT-058 „Operațiune completă”
- Nivel: recomandare
- Greșit: „Operațiune completă” → Corect: „Operațiune finalizată”
- De ce: „Operațiune completă” nu e greșeală: DEX ’98 dă pentru „complet” sinonimele „întreg, desăvârșit, deplin, împlinit”, deci se poate înțelege și „împlinită”. Pentru engl. „operation complete” sensul „terminată” e mai clar cu „finalizată”, „încheiată” sau „Gata”. Pentru verbul „a completa” vezi RO-CAL-058 din [prieteni-falsi-si-calcuri.md](../prieteni-falsi-si-calcuri.md).
- Sursă: DEX ’98 s.v. complet

### RO-INT-025 Starea goală
- Nivel: recomandare
- Greșit: „0 rezultate găsite” → Corect: „Niciun rezultat”
- De ce: Pentru lista goală se preferă „Niciun rezultat” sau „Nu există rezultate”. „Nu au fost găsite rezultate” e corect și nu se marchează ca greșeală (RO-NUG-089 din [ce-nu-e-greseala.md](../ce-nu-e-greseala.md)); varianta scurtă e o preferință. Exemplul „Greșit” nu e o greșeală de limbă, ci un mesaj mai puțin natural.
- Sursă: Mic ghid CE p. 4–6 (Rec. 2, 4)

### RO-INT-022 Zero, cu mesaj propriu
- Nivel: recomandare
- Greșit: „0 fișiere” → Corect: „Niciun fișier”
- De ce: Cazul zero primește mesaj propriu (`=0` în ICU). „0 fișiere” rămâne corect gramatical: CLDR pune zero în categoria `few`. Propui, nu corectezi.
- Sursă: CLDR; ICU; gettext

### RO-INT-013 Mesajul de eroare general
- Nivel: recomandare
- Greșit: „S-a întâmplat ceva prost.” → Corect: „A apărut o problemă. Încearcă din nou peste câteva minute.”
- De ce: Un mesaj de eroare spune ce s-a întâmplat și ce urmează. „S-a întâmplat ceva prost” e corect (RO-NUG-113 din [ce-nu-e-greseala.md](../ce-nu-e-greseala.md)); aici propui doar o formulare mai utilă, nu o corectezi ca greșeală.
- Sursă: DEX ’98 s.v. problemă; Mic ghid CE p. 4 (Rec. 2)

### RO-INT-038 Eroarea are predicat și pasul următor
- Nivel: recomandare
- Greșit: „Eroare procesare fișier” → Corect: „Fișierul nu a putut fi procesat. Verifică formatul.”
- De ce: Mesajul cu predicat spune ce s-a întâmplat, iar al doilea enunț spune ce poate face utilizatorul. Varianta nominală completă („Eroare la procesarea fișierului”) nu e greșeală de limbă, doar mai puțin utilă.
- Sursă: Mic ghid CE p. 4 (Rec. 2)

### RO-INT-043 Mesajul de validare arată pasul următor
- Nivel: recomandare
- Greșit: „Formular incomplet” → Corect: „Mai ai 18 întrebări fără răspuns. Completează-le pe toate.”
- De ce: „Formular incomplet” e corect ca limbă, dar nu ajută. Mesajul de validare e specific și arată cum se rezolvă problema. Numărul și detaliile le iei din formular, nu le inventezi.
- Sursă: Mic ghid CE p. 4 (Rec. 2)

### RO-INT-039 Confirmarea, ca întrebare directă
- Nivel: recomandare
- Greșit: „Sunteți sigur că doriți ștergerea?” → Corect: „Sigur doriți să ștergeți fișierul?”
- De ce: Confirmarea e o întrebare directă către utilizator, în registrul ales. Forma scurtă „Sigur ștergi…?” merge doar într-un produs cu ton intenționat colocvial.
- Sursă: Mozilla RO (uz, nu normă)

### RO-INT-042 Butoanele de apel la acțiune încep cu verbul și sunt specifice
- Nivel: recomandare
- Greșit: „Click aici” → Corect: „Calculează scorul”
- De ce: Un buton de apel la acțiune începe cu un verb și spune ce se întâmplă: „Calculează scorul”, „Vezi rezultatul”. „Click aici” nu spune nimic, iar „Trimite” și „OK” se evită când există o acțiune mai specifică. „OK” e în DOOM 3, deci nu e greșeală de limbă (RO-IMP-012 din [imprumuturi.md](../imprumuturi.md)). Interdicția generală pe „Trimite” și „OK” e o regulă a unui produs, nu o normă.
- Sursă: DOOM 3 s.v. OK; Mic ghid CE p. 9 (Rec. 7)

### RO-INT-037 „Drag-and-drop”
- Nivel: recomandare
- Greșit: „Drag-and-drop fișierul aici” → Corect: „Trage și plasează fișierul aici”
- De ce: Cu „tu”: „Trage și plasează”; cu „dumneavoastră”: „Trageți și plasați”. Ambele verbe sunt în dicționar (DOOM 3 are „trage” și „plasa”), iar formele de imperativ sunt corecte.
- Sursă: DOOM 3 s.v. trage, plasa; DEX ’98 s.v. trage

### RO-INT-027 „Cancel”: „Renunță” sau „Anulează”
- Nivel: recomandare
- Greșit: „Cancel” → Corect: „Renunță”
- De ce: Alegerea „Renunță” este o convenție de terminologie (Mozilla RO o alege pentru consecvență), nu o normă. „Anulează” sau „Anulare” e la fel de corect și nu se marchează ca greșeală. Alegi o formă și o ții în tot produsul.
- Sursă: Mozilla RO (uz, nu normă); DEX ’98 s.v. renunța, anula

### RO-INT-028 „Toggle”
- Nivel: recomandare
- Greșit: „Toggle modul întunecat” → Corect: „Comută modul întunecat”
- De ce: „A comuta” e verbul din DOOM 3 și din DEX ’98 cu sensul tehnic de a schimba conexiunea sau starea; extinderea la starea unei opțiuni din interfață e firească. „Toggle” nu e în dicționare.
- Sursă: DOOM 3 s.v. comuta; DEX ’98 s.v. comuta

## Termeni de interfață

### RO-INT-001 „Share”: „partajează”
- Nivel: recomandare
- Greșit: „Împărtășește fișierul cu echipa.” → Corect: „Partajează fișierul cu echipa.”
- De ce: Pentru resurse (fișiere, linkuri, ecran), „partajează” e echivalentul uzual din localizare. DEX ’98 dă pentru „partaja” doar sensul juridic („a efectua un partaj”), deci sensul informatic e o extensie de terminologie. Pentru opinii, experiențe sau sentimente, „împărtășește” rămâne corect.
- Sursă: DOOM 3 s.v. partaja; DEX ’98 s.v. partaja, împărtăși

### RO-INT-002 „Menu bar”: „bară de meniu”
- Nivel: recomandare
- Greșit: „bară de meniuri” → Corect: „bară de meniu”
- De ce: Convenția Mozilla RO e „bară de meniu”. Motivul invocat („e un singur meniu”) nu se susține, pentru că o bară conține mai multe meniuri (Fișier, Editare, Vizualizare). Niciuna dintre forme nu se marchează drept greșeală; alegi una și o ții.
- Sursă: Mozilla RO (uz, nu normă); DOOM 3 s.v. meniu

### RO-INT-003 „Help”: „Asistență”
- Nivel: recomandare
- Greșit: „Help” → Corect: „Asistență”
- De ce: „Asistență” are în DEX ’98 sensul „sprijin, ajutor”, deci se potrivește. „Ajutor” e la fel de corect.
- Sursă: DEX ’98 s.v. asistență; Mozilla RO (uz, nu normă)

### RO-INT-004 „Breadcrumbs”: „rută de navigare”
- Nivel: recomandare
- Greșit: „breadcrumbs” → Corect: „rută de navigare”
- De ce: „Rută” e în DOOM 3 și în DEX ’98 cu sensul „drum”, deci „rută de navigare” se formează corect și se înțelege. Alte echivalente descriptive nu se marchează ca greșeli.
- Sursă: DOOM 3 s.v. rută; DEX ’98 s.v. rută; Mozilla RO (uz, nu normă)

### RO-INT-005 „Bold”: „aldin”
- Nivel: recomandare
- Greșit: „Text bold” → Corect: „Text aldin”
- De ce: „Aldin” e termenul tipografic românesc (DEX ’98: caracter cu conturul mai gros, „gras”). „Bold” nu e greșit: DOOM 3 are intrarea nouă „bold” pentru font (vezi RO-IMP-048 din [imprumuturi.md](../imprumuturi.md)). În interfață propui „aldin”.
- Sursă: DOOM 3 s.v. bold, aldin; DEX ’98 s.v. aldin

### RO-INT-006 „Timer”: „cronometru”
- Nivel: recomandare
- Greșit: „Setează temporizatorul” → Corect: „Setează cronometrul”
- De ce: În dicționare, „temporizator” este doar adjectiv, „care tărăgănează” (DEX ’98), deci substantivul nu are bază lexicografică. „Cronometru” e instrumentul care măsoară timpul (DEX ’98). Pentru numărătoarea inversă precizezi din context („cronometru cu numărătoare inversă”).
- Sursă: DOOM 3 și DEX ’98 s.v. temporizator, cronometru; Mozilla RO (uz, nu normă)

### RO-INT-007 „Sign in”: „Autentificare”, nu „Înregistrare”
- Nivel: recomandare
- Greșit: „Înregistrare” → Corect: „Autentificare”
- De ce: „A înregistra” înseamnă „a înscrie într-un registru” (DEX ’98), adică a crea un cont, nu a intra într-unul existent. Pentru „sign in” spui „Autentificare” sau „Autentifică-te”. Greșeala apare când butonul „Înregistrare” deschide un cont existent. Sensul „a intra în cont” nu e în DEX ’98 pentru „autentificare”, deci regula e o convenție de terminologie. Judeci după ce face butonul, nu după etichetă; dacă nu știi ce face, întrebi.
- Sursă: DEX ’98 s.v. înregistra; Mozilla RO (uz, nu normă)

### RO-INT-008 „Sign up”: „Înregistrare” sau „Creează un cont”, nu „Autentificare”
- Nivel: recomandare
- Greșit: „Autentificare” → Corect: „Creează un cont”
- De ce: „A autentifica” înseamnă „a legaliza, a face autentic” (DEX ’98), adică a confirma o identitate. Greșeala apare pe butonul care creează un cont nou: „autentificarea” nu denumește crearea unui cont, pentru care se spune „Înregistrare” sau „Creează un cont”.
- Sursă: DEX ’98 s.v. autentifica, înregistra; Mozilla RO (uz, nu normă)

### RO-INT-009 Sensul se deduce din context, nu din cuvânt
- Nivel: recomandare
- Greșit: „Tapetează copii” → Corect: „Ferestre alăturate”
- De ce: Șirul englezesc „Tile children”, tradus după cuvinte, dă un nonsens. Sensul vine din rolul șirului în interfață (aici, aranjarea ferestrelor), nu din traducerea cuvintelor luate separat. Dacă nu ai contextul cheii, întrebi.
- Sursă: Mozilla RO (uz, nu normă)

### RO-INT-010 Mărcile nu devin verbe
- Nivel: recomandare
- Greșit: „Googlează numele produsului.” → Corect: „Caută numele produsului pe Google.”
- De ce: Numele de mărci se declină cât mai puțin, în special în titluri și în elementele proeminente, și nu devin verbe. „Googlează” (formă dedusă, nu e în dicționare) e un derivat de la o marcă. Regula e o convenție de stil, nu o normă.
- Sursă: Mozilla RO (uz, nu normă)

### RO-INT-011 „A se loga”: „Autentifică-te”, „Conectează-te”
- Nivel: recomandare
- Greșit: „Loghează-te pentru a continua.” → Corect: „Autentifică-te pentru a continua.”
- De ce: Verbul „a se loga” e admis în DOOM 3 (intrare nouă), deci nu e greșeală și nu se semnalează ca atare. Ca preferință de stil, în interfață propui „Autentifică-te” sau „Conectează-te”.
- Sursă: DOOM 3 s.v. loga

### RO-INT-017 „Oferă posibilitatea de a”
- Nivel: recomandare
- Greșit: „Aplicația vă oferă posibilitatea de a exporta datele.” → Corect: „Puteți exporta datele.”
- De ce: „Vă oferă posibilitatea de a” e o perifrază nominală. Recomandarea cere să înlocuiești substantivele inutile cu forme verbale. Textul scurt de interfață câștigă cel mai mult din asta.
- Sursă: Mic ghid CE p. 8 (Rec. 6)

### RO-INT-018 „Funcționalitate” sau „funcție”
- Nivel: recomandare
- Greșit: „o nouă funcționalitate” → Corect: „o funcție nouă”
- De ce: „Funcționalitate” e un cuvânt corect (DOOM 3, DEX ’98 – „calitatea de a fi funcțional”) și nu se condamnă. În interfață, „funcție” e mai scurt și mai obișnuit.
- Sursă: DOOM 3 s.v. funcționalitate; DEX ’98 s.v. funcționalitate

### RO-INT-053 „+3 mai multe”
- Nivel: recomandare
- Greșit: „+3 mai multe” → Corect: „încă 3”
- De ce: Calcul după engl. „+3 more”. „Încă” înseamnă în DEX ’98 „în plus, pe deasupra” („Au venit încă două persoane”), exact structura necesară: „încă 3 fișiere”. Substantivul se acordă după regula pluralului (vezi mai jos).
- Sursă: DEX ’98 s.v. încă

### RO-INT-055 Tastele de acces
- Nivel: recomandare
- Greșit: „&Salvează” → Corect: „S&alvează”
- De ce: Exemplul arată acceleratorul pus pe litera s și mutat pe a. Acceleratorii evită literele care coboară sub rând (p, g, q, y), literele înguste (i, l) și s sau t, ușor de confundat cu ș sau ț. E o recomandare de ergonomie, nu o normă a limbii.
- Sursă: Mozilla RO (uz, nu normă)

## Valorile inserate

Un șir cu valoare inserată (`%{nume}`, `{count}`, `%(autor)s`) se scrie astfel încât să fie corect pentru orice valoare posibilă. Valoarea nu primește articol, desinență sau adjectiv acordat, pentru că nu știi ce cuvânt va apărea.

### RO-INT-045 Fără articol sau desinență pe valoare
- Nivel: normă
- Greșit: „Ai terminat %{track_title}-ul.” → Corect: „Ai terminat piesa %{track_title}.”
- De ce: Valoarea nu se vede, deci o desinență sau un articol lipit de ea e greșit pentru o parte din valori. Un substantiv de context care are genul și numărul cunoscut rezolvă acordul. Unde nu e nevoie de flexiune, valoarea stă singură: „%{name} s-a actualizat” e corect când valoarea e la singular (numele unui fișier), pentru că participiul e invariabil, dar auxiliarul nu („s-a”, „s-au”). Pentru valori care pot fi și la plural formulezi fără acord, ca la RO-INT-048.
- Sursă: ICU; DOOM 3 (articolul și desinențele)

### RO-INT-046 Adjectivul nu se acordă cu valoarea
- Nivel: normă
- Greșit: „%{name} este disponibil” → Corect: „Aplicația %{name} este disponibilă”
- De ce: Genul valorii nu se cunoaște, iar un adjectiv sau un participiu variabil acordat cu valoarea va fi greșit pentru o parte din valori. Pui un substantiv cu gen cunoscut și acordezi adjectivul cu el.
- Sursă: DOOM 3 (acordul adjectivului)

### RO-INT-047 Numele marcatorilor rămân neschimbate
- Nivel: normă
- Greșit: „{titlu}” → Corect: „{title}”
- De ce: Numele din acolade e un identificator pe care aplicația îl înlocuiește la afișare, nu un text de tradus. Nu se traduce și nu se înlocuiește cu text: `{title}` și `{domain}` rămân exact ca în sursă.
- Sursă: ICU

### RO-INT-048 Marcatorul de persoană, la singular
- Nivel: recomandare
- Greșit: „%(author)s au răspuns” → Corect: „%(author)s a răspuns”
- De ce: Un marcator de forma `%(nume)s` e un identificator, nu un plural. Sursa engleză are verbul la singular („has replied”), deci verbul românesc rămâne la singular. Dacă sursa admite mai multe persoane, folosești o formulare care nu depinde de număr (de exemplu „Răspuns de la %(author)s”).
- Sursă: Mozilla RO (uz, nu normă); ICU

### RO-INT-049 Marcatorul se reproduce întreg
- Nivel: normă
- Greșit: „%(document_title)” → Corect: „%(document_title)s”
- De ce: Fără „s” la sfârșit, formatul `%(nume)s` nu mai e complet, iar formatarea dă eroare sau afișează greșit. Marcatorul se copiază întreg, nu se rescrie.
- Sursă: ICU; Mozilla RO (uz, nu normă)

### RO-INT-050 Etichetele și variabilele, exact și în același număr
- Nivel: normă
- Greșit: „<0>Citește</0> {{nume}}” → Corect: „<0>Citește</0> {{name}}”
- De ce: Etichetele numerotate (`<0>…</0>`) și variabilele (`{{name}}`) se reproduc exact, în aceeași formă și în același număr ca în sursă. Traduci textul dintre ele, nu marcajul.
- Sursă: regula de fidelitate din SKILL.md (marcajul rămâne neatins); sintaxa din exemplu e a i18next, convenție de localizare, fără normă

### RO-INT-051 Textele comparate de teste nu se traduc
- Nivel: de semnalat
- Greșit: „Salut, lume!” → Corect: „Hello, world!”
- De ce: Exemplul privește un șir pe care un test îl compară literal. Dacă un test automat compară textul literal, traducerea îl strică. Titlul exercițiului se traduce, dar literalul din corpul textului rămâne neschimbat. Nu poți ști din șir dacă e comparat de un test, de aceea semnalezi și întrebi.
- Sursă: ICU; convenție de localizare, fără normă

## Plural

Româna are trei categorii de plural în CLDR și în ICU: `one`, `few` și `other`.

- `one`: exact 1 („1 fișier”).
- `few`: 0, numerele de la 2 la 19 și cele care se termină în 01–19 (101, 119, 1001), plus zecimalele („1,5 fișiere”).
- `other`: restul (20, 100, 120, 1.000), cu „de” după numeral („20 de fișiere”, „1.000 de fișiere”).

Regula lui „de” e în RO-NUG-031 din [ce-nu-e-greseala.md](../ce-nu-e-greseala.md) și în [numere-date-unitati.md](../numere-date-unitati.md).

### RO-INT-020 Cele trei categorii
- Nivel: normă
- Greșit: „{count, plural, one {# fișier} other {# fișiere}}” → Corect: „{count, plural, one {# fișier} few {# fișiere} other {# de fișiere}}”
- De ce: Categoriile englezești (`one`, `other`) dau pentru 20 „20 fișiere”, fără „de”. Verificate cu datele CLDR: 0 → `few`, 1 → `one`, 1,5 → `few`, 2 și 19 → `few`, 20, 99, 100, 120, 200 și 1020 → `other`, 101, 119, 201 și 1001 → `few`. Zecimalele cad în `few`, deci „1,5 fișiere” fără „de”.
- Sursă: CLDR; ICU

### RO-INT-021 Formula gettext
- Nivel: normă
- Greșit: „nplurals=2; plural=(n != 1)” → Corect: „nplurals=3; plural=(n==1 ? 0 : (n==0 || (n%100 > 0 && n%100 < 20)) ? 1 : 2)”
- De ce: Formula cu două forme e cea a englezei. Pentru română sunt trei forme: 0 pentru exact 1, 1 pentru 0 și pentru numerele cu ultimele două cifre între 1 și 19, 2 pentru restul.
- Sursă: gettext; CLDR

### RO-INT-023 Fiecare categorie care acoperă mai multe numere afișează numărul
- Nivel: normă
- Greșit: „{count, plural, one {un fișier} few {câteva fișiere} other {multe fișiere}}” → Corect: „{count, plural, one {# fișier} few {# fișiere} other {# de fișiere}}”
- De ce: În română `few` acoperă 0, 2–19, 101… și `other` acoperă 20, 100, 120…, deci ambele trebuie să conțină numărul (`#` sau `{count}`). Categoriile se aleg după limba țintă, nu după cele din engleză.
- Sursă: ICU; CLDR

### RO-INT-054 Numeralul cu „+”
- Nivel: recomandare
- Greșit: „25+ limbi” → Corect: „cel puțin 25 de limbi”
- De ce: „25+” ascunde locul lui „de”, pe care numeralul îl cere când ultimele două cifre fac 00 sau un număr de la 20 la 99 („25 de limbi”, dar „101 limbi”). „25+” înseamnă „25 sau mai multe”, deci reformularea păstrează limita: „cel puțin 25”. „Peste 25” exclude 25 și schimbă cifra (vezi RO-NUM-048 din [numere-date-unitati.md](../numere-date-unitati.md)).
- Sursă: INJ p. 84–85; DSL, s.v. NUMERAL

## Acordul etichetelor

### RO-INT-052 Etichetele de stare se acordă cu referentul
- Nivel: normă
- Greșit: „Finalizată” → Corect: „Finalizat”
- De ce: Starea „Finalizată” se acordă cu o lecție, „Finalizat” cu un exercițiu; exemplul arată forma feminină pusă pe un exercițiu. O cheie partajată („Completed”) tradusă la feminin devine greșeală când e refolosită pentru un masculin: ceri chei separate pe gen și semnalezi cheile partajate.
- Sursă: DOOM 3 (acordul adjectivului); DSL, s.v. ACORD

### RO-INT-041 Acordul se face după contextul cheii
- Nivel: normă
- Greșit: „(încă nu este validă)” → Corect: „(încă nu este valid)”
- De ce: Exemplul privește cheia „pseudonim”. Șirul englezesc „(not valid yet)” nu marchează genul. Acordul adjectivului vine din substantivul la care se referă cheia („pseudonim” e masculin), nu din formulare. Fără contextul cheii, întrebi.
- Sursă: DOOM 3 (acordul adjectivului)

## Formate, majuscule și punctuație

### RO-INT-024 Formatele ro-RO
- Nivel: recomandare
- Greșit: „Total: 1,250.50 lei, 10/17/2026” → Corect: „Total: 1.250,50 lei, 17.10.2026”
- De ce: Formatele românești: suma înaintea monedei, punctul la mii și virgula la zecimale, data „17.10.2026” sau „17 octombrie 2026”, ora în format de 24 de ore. Regulile complete sunt în [numere-date-unitati.md](../numere-date-unitati.md). În interfață formatul vine de preferat din localizarea „ro-RO”, nu din text scris de mână.
- Sursă: convenție de localizare, fără normă; vezi numere-date-unitati.md

### RO-INT-056 Punctuația finală
- Nivel: recomandare
- Greșit: „Se salvează” → Corect: „Se salvează…”
- De ce: Pentru sursa „Saving…”, traducerea cere elipsa. Punctuația finală a traducerii (punct, două puncte, elipsă) e identică cu a sursei. Nepotrivirile se semnalează pentru verificare, nu se corectează automat. Excepție: ordinea punctului față de ghilimelele de închidere urmează norma românească (punctul după ghilimele când citatul face parte din propoziție).
- Sursă: convenție de localizare, fără normă; ghid gov-md §9.8

Majusculele calchiate din engleză apar des în etichete și în titluri („Limba: Română”, „Luni, 5 Martie”, titluri cu fiecare cuvânt cu majusculă): lunile, zilele și limbile se scriu cu literă mică, iar titlurile și etichetele au majusculă doar la primul cuvânt și la numele proprii. Regulile sunt în [tipografie.md](../tipografie.md).

### RO-INT-057 Unitățile se convertesc doar când numărul e scris în text
- Nivel: recomandare
- Greșit: „Curse de până la 320 de mile” → Corect: „Curse de până la 515 km (320 de mile)”
- De ce: Unitățile imperiale din textul fix pot fi date în sistemul metric, dar numai cu un factor de conversie cunoscut (o milă terestră are circa 1,609 km) și cu valoarea originală păstrată între paranteze. Când numărul e o variabilă, unitatea nu se convertește în șir, pentru că valoarea nu se cunoaște. Vezi RO-NUM-033 și RO-NUM-048 din [numere-date-unitati.md](../numere-date-unitati.md) și [fidelitate.md](../fidelitate.md).
- Sursă: Mozilla RO (uz, nu normă)

## Coerența cu produsul

### RO-INT-014 Elementele de interfață se numesc cum apar pe ecran
- Nivel: de semnalat
- Greșit: „Apăsați «Trimite»” → Corect: „Apăsați «Expediază»”
- De ce: Exemplul presupune că pe ecran butonul se numește „Expediază”. Dacă textul numește un buton, o filă sau un meniu altfel decât apare pe ecran, cititorul caută ceva ce nu există. Semnalezi nepotrivirea; denumirea corectă o poți afla doar de la produs.
- Sursă: convenție de localizare, fără normă

### RO-INT-015 Același termen, aceeași traducere
- Nivel: recomandare
- Greșit: „Preferințe” → Corect: „Setări”
- De ce: Exemplul presupune că aceeași pagină se numește „Setări” în alt ecran. Un termen de interfață rămâne același între ecrane, programe și fișiere de traducere (RO-COE-006 din [coerenta.md](../coerenta.md)). Nu e absolut: același cuvânt sursă poate cere traduceri diferite când are roluri diferite (verb și substantiv).
- Sursă: Mozilla RO (uz, nu normă)

### RO-INT-016 Textul nu contrazice produsul
- Nivel: de semnalat
- Greșit: „Modificările au fost salvate.” → Corect: „[întrebare pentru autor: ce face butonul?]”
- De ce: Exemplul presupune că nu s-a salvat nimic. Mesajul nu afirmă ce produsul nu face: nu scrie „salvate” dacă nu se salvează nimic și nu raportează o eroare de rețea într-o aplicație fără rețea. Aceasta e o regulă de conținut, nu de limbă; o semnalezi și întrebi.
- Sursă: convenție de localizare, fără normă

### RO-INT-044 Fără stări de eroare inventate
- Nivel: de semnalat
- Greșit: „Conexiune pierdută. Încearcă din nou.” → Corect: „[întrebare pentru autor: când poate apărea această eroare?]”
- De ce: Exemplul presupune o aplicație fără conexiune la rețea. Nu inventezi stări de eroare care nu pot apărea în produs. Regula ține de coerența dintre text și produs. Cere confirmarea de la cine cunoaște aplicația.
- Sursă: convenție de localizare, fără normă

## Lista de verificare mecanică

La sfârșit, parcurgi șirurile de interfață cu aceste puncte, fiecare în fișierul lui:

- sedilă în loc de virgulă la ș și ț: [tipografie.md](../tipografie.md);
- punctuația finală identică cu a sursei: RO-INT-056;
- spațiu după punctuație, fără majuscule englezești (title case, luni și zile cu majusculă): [tipografie.md](../tipografie.md);
- marcatori și etichete intacte: RO-INT-047 – RO-INT-050;
- taste de acces fără p, g, q, y, i, l, s, t: RO-INT-055;
- ghilimele „…”: [tipografie.md](../tipografie.md);
- mărci nedeclinate în elementele proeminente: RO-INT-010;
- unități convertite doar când numărul e scris: RO-INT-057.

Spațiul. Se presupune că textul sursă ocupă cam două treimi din spațiul disponibil, ca traducerea să aibă loc să crească. Aceasta e o țintă de proiectare, nu un plafon pentru titlul tradus. La fel, „cel mult 14 cuvinte” pentru un mesaj e o țintă orientativă, nu o limită.
