package lint

// Cazurile limită găsite la recenzia fazei 2: zonele HTML și front matter-ul,
// prefixoidele cu „î”, numele proprii și cuvintele străine.

import (
	"strings"
	"testing"
)

func TestZoneHtmlSiFrontMatter(t *testing.T) {
	for _, ok := range []string{
		"<p>Text.</p>\n<script>var a = x ? 1 : 2;  f(a , b);</script>\n",
		"<STYLE>a  { color : red }</STYLE>",
		"<script type=\"x\">\nvar a = b ? c : d;\n</script>\n<p>ok</p>",
		"<pre>a  b , c</pre> și <code>x ? y</code>",
		"<div  class=\"a  b\" data-x = \"1\">text</div>",
		"---\ntitle : Un titlu  aici\ndescription: cîndva\n---\n\nText curat.\n",
		"<!-- a  b , c -->",
	} {
		if o := Verifica(ok); len(o) != 0 {
			t.Errorf("%q: aștept zero observații, am %+v", ok, o)
		}
	}
	// textul din afara zonelor se verifică în continuare
	if o := Verifica("<script>x</script>\nun cuvânt , apoi"); len(o) != 1 {
		t.Errorf("textul de după script trebuie verificat: %+v", o)
	}
	// front matter doar la începutul fișierului
	if o := Verifica("Text.\n---\ntitle : x\n---\n"); len(o) == 0 {
		t.Errorf("un „---” de la mijloc nu e front matter")
	}
	if !Ignorat("<script>a  b</script>")(10) {
		t.Error("Ignorat trebuie să acopere conținutul scriptului")
	}
}

func TestPrefixoideCuI(t *testing.T) {
	for _, ok := range []string{
		"videoînregistrare", "teleînregistrare", "bioîngrășăminte", "fotoîntrerupător",
		"portîncărcător", "superîmbogățit", "radioînregistrare", "electroîntrerupător",
		"ecoîngrășământ", "geoînregistrare", "audioînregistrare", "hidroîmpermeabil",
	} {
		if o := Verifica(ok + " aici"); len(o) != 0 {
			t.Errorf("%q: aștept zero observații, am %+v", ok, o)
		}
	}
	// grafia veche rămâne semnalată
	for cuv, prop := range map[string]string{
		"cînd": "când", "sînt": "sunt", "pîine": "pâine", "întîi": "întâi",
		"hotărît": "hotărât", "rîu": "râu", "mîine": "mâine",
	} {
		if p := propuneri(Verifica("Spun " + cuv + " aici")); len(p) != 1 || p[0] != prop {
			t.Errorf("%s: propuneri %v, aștept %s", cuv, p, prop)
		}
	}
	if o := Verifica("Spun coborî aici"); len(o) != 0 {
		t.Errorf("coborî nu se semnalează: %+v", o)
	}
}

func TestNumePropriiCuI(t *testing.T) {
	for _, ok := range []string{
		"Pîrvu a venit.", "(Rîpeanu 1999)", "**Brîndușa** a venit.", "Brîndușa a venit.",
		"Mulțumesc, maître.", "Il faut connaître le maître.", "Maître Dupont a venit.",
		"Işık a venit.", "Spun că Işık e aici.", "Müller și Özgür", "Çelik, Ağa",
		"A venit Rîpeanu.", "\"Pîrvu\" a venit.",
	} {
		if o := Verifica(ok); len(o) != 0 {
			t.Errorf("%q: aștept zero observații, am %+v", ok, o)
		}
	}
	// cuvintele uzuale cu majusculă la început de propoziție rămân semnalate
	for cuv, prop := range map[string]string{
		"Cînd": "Când", "Sînt": "Sunt", "Pînă": "Până", "Mîine": "Mâine", "Cît": "Cât",
	} {
		if p := propuneri(Verifica(cuv + " vine.")); len(p) != 1 || p[0] != prop {
			t.Errorf("%s: propuneri %v, aștept %s", cuv, p, prop)
		}
	}
	// minuscula apare și în altă parte cu „î”: majuscula e aceeași grafie veche
	if p := propuneri(Verifica("Rîul e mare. Un rîu mic.")); len(p) != 2 {
		t.Errorf("rîul/rîu: %v", p)
	}
	// sedila rămâne greșeală la numele românești
	for cuv, prop := range map[string]string{"Ştefan": "Ștefan", "Mureş": "Mureș", "Timişoara": "Timișoara"} {
		if p := propuneri(Verifica("Merg la " + cuv + " mâine.")); len(p) != 1 || p[0] != prop {
			t.Errorf("%s: propuneri %v, aștept %s", cuv, p, prop)
		}
	}
}

func TestMajusculaFlexionataLaInceputDePropozitie(t *testing.T) {
	text := "Întîmplarea a fost ciudată. Tînărul a plecat. Hotărîrea e luată. Mîncarea e gata. Cîinele latră. Pîinea e caldă. Gîndurile zboară."
	if o := Verifica(text); len(o) != 7 {
		t.Errorf("aștept 7 observații, am %d: %+v", len(o), o)
	}
	for _, nume := range []string{"Pîrvu a venit.", "(Rîpeanu 1999)", "**Brîndușa**", "Mîndrescu a venit."} {
		if o := Verifica(nume); len(o) != 0 {
			t.Errorf("%q: numele nu se semnalează: %+v", nume, o)
		}
	}
}

func TestFrontMatterDoarYaml(t *testing.T) {
	o := Verifica("---\n\nSînt aici și vin.\n\n---\n\nCînd? Sînt.")
	if len(o) != 3 {
		t.Errorf("liniile orizontale nu sunt front matter; aștept 3 observații, am %d: %+v", len(o), o)
	}
	if o := Verifica("---\nSînt aici.\n---\n"); len(o) != 1 {
		t.Errorf("proza între separatoare se verifică: %+v", o)
	}
}

func TestEticheteHtmlNuProza(t *testing.T) {
	o := Verifica("x < y și a<b. Cînd c>d mîine")
	if len(o) != 2 {
		t.Errorf("aștept 2 observații (Cînd, mîine), am %d: %+v", len(o), o)
	}
	if o := Verifica(`<a href="x" class=y>link  dublu</a> <br/> <input disabled="">`); len(o) != 1 {
		t.Errorf("etichetele reale rămân ignorate: %+v", o)
	}
}

// O etichetă pomenită în proză sau în cod, fără închidere, nu ascunde restul
// documentului; conținutul unei etichete închise rămâne ignorat.
func TestEtichetaNeinchisaNuAscundeRestul(t *testing.T) {
	for _, text := range []string{
		"Folosește `<code>` pentru cod. Apoi cînd vine pîinea, sînt aici.",
		"Scrie <script> în pagină, apoi cînd vine pîinea.",
		"```html\n<style>\n```\nApoi cînd vine pîinea.",
		"Stă scris „<pre>” aici. Apoi cînd plouă, vine pîinea.",
	} {
		gasit := false
		for _, o := range Verifica(text) {
			if o.Regula == "RO-TIP-128" && text[o.Inceput:o.Sfarsit] == "cînd" {
				gasit = true
			}
		}
		if !gasit {
			t.Errorf("„cînd” trebuie semnalat în %q: %+v", text, Verifica(text))
		}
	}
	text := "<script>const x = ok ? 1 : 2; var cînd = 1;</script>\nAici e bine."
	if obs := Verifica(text); len(obs) != 0 {
		t.Errorf("conținutul din <script> închis trebuie ignorat: %+v", obs)
	}
}

// Grafia veche cu majusculă (început de propoziție, titlu, majuscule întregi)
// se prinde și pentru cuvintele uzuale care nu sunt în lista scurtă, prin
// rădăcină; numele de persoană cu „î” rămân neatinse.
func TestRadaciniVechiCuMajuscula(t *testing.T) {
	text := "Gîndesc la asta. Întîmplător e aici. Vînzătorul a plecat. Mîncăm acum.\n\n# Cîntăm\n\nGÎNDESC MULT. Rîvna lui e mare."
	vrut := map[string]bool{"Gîndesc": true, "Întîmplător": true, "Vînzătorul": true, "Mîncăm": true, "Cîntăm": true, "GÎNDESC": true, "Rîvna": true}
	for _, o := range Verifica(text) {
		if o.Regula == "RO-TIP-128" {
			delete(vrut, text[o.Inceput:o.Sfarsit])
		}
	}
	if len(vrut) != 0 {
		t.Errorf("nesemnalate: %v", vrut)
	}
	for _, nume := range []string{"Pîrvu a venit.", "Mîndrescu a venit.", "Rîpeanu a scris.", "Brîndușa a plecat."} {
		if obs := Verifica(nume); len(obs) != 0 {
			t.Errorf("%q: numele nu se atinge: %+v", nume, obs)
		}
	}
}

// regulile și fragmentele observațiilor, pentru comparații scurte
func fragmente(text string) map[string]string {
	m := map[string]string{}
	for _, o := range Verifica(text) {
		m[text[o.Inceput:o.Sfarsit]] = o.Regula
	}
	return m
}

// Numele de familie cu „î” la început de propoziție sau de rând nu se ating;
// nici în numele compuse cu cratimă.
func TestNumeDeFamilieLaInceput(t *testing.T) {
	for _, text := range []string{
		"Mîndrea a plecat. Cîndea Ion a venit.\nCîntea a spus ceva.\n- Cîndea Maria\n\nGîndea e un nume.",
		"Mîndru a venit acasă.",
		"Au venit Ion și Popescu-Mîndru acasă.",
	} {
		if f := fragmente(text); len(f) != 0 {
			t.Errorf("%q: numele nu se ating: %v", text, f)
		}
	}
	// formele articulate ale cuvintelor uzuale rămân prinse
	f := fragmente("Pîinea e caldă. Mîncarea e gata. Hotărîrea e luată. Întîlnirea a fost lungă.")
	for _, c := range []string{"Pîinea", "Mîncarea", "Hotărîrea", "Întîlnirea"} {
		if f[c] != "RO-TIP-128" {
			t.Errorf("%s trebuie semnalat: %v", c, f)
		}
	}
}

// Gardurile de cod se închid doar cu același caracter, cel puțin la fel de lung.
func TestGarduriDeCod(t *testing.T) {
	text := "````markdown\n```text\nsînt aici\n```\nsînt  două\n````\n\nApoi cînd vine.\n\n~~~\n```\nsînt\n~~~\nUn alt rând și sînt."
	f := fragmente(text)
	gasite := 0
	for _, o := range Verifica(text) {
		if o.Regula == "RO-TIP-128" {
			gasite++
			if o.Inceput < strings.Index(text, "Apoi") {
				t.Errorf("semnalat în cod: %q la %d", text[o.Inceput:o.Sfarsit], o.Inceput)
			}
		}
	}
	if f["cînd"] != "RO-TIP-128" || gasite != 2 {
		t.Errorf("aștept „cînd” și ultimul „sînt”: %v (%d)", f, gasite)
	}
}

// Formulele, coloanele aliniate, codul indentat și tabelele din citat nu
// primesc observații de spațiu.
func TestSpatiiInZoneSpeciale(t *testing.T) {
	for _, text := range []string{
		"Fie $f(x) , g(x)$ și $$a , b$$ funcții, iar drumul e lung pentru că sunt aici.",
		"Iată comanda, iar drumul e lung pentru că sunt aici:\n\n    $ ls    -la   # listare\n    x = a ? 1 : 2\n",
		"Tabelul, iar drumul e lung pentru că sunt aici:\n\n> | a   | b |\n> | c  | x  y |\n",
		"Nume          Cantitate   Preț\nMere          5           3,5\nși așa mai departe, pentru că sunt aici",
	} {
		for _, o := range Verifica(text) {
			if o.Regula == "RO-TIP-068" || o.Regula == "" {
				t.Errorf("%q: observație de spațiu falsă: %q", text, text[o.Inceput:o.Sfarsit])
			}
		}
	}
	// spațiul dublu obișnuit și spațiul înaintea virgulei rămân prinse
	f := fragmente("Drumul e lung  și frumos , iar eu sunt aici pentru că vreau.")
	if _, ok := f["  "]; !ok || f[" ,"] != "RO-TIP-068" {
		t.Errorf("spațiile obișnuite trebuie prinse: %v", f)
	}
}

// Citatul închis cu „…“ e tot citat.
func TestCitatInchisCuGhilimeleSus(t *testing.T) {
	text := "Citat franțuzesc: „toute âme et ses abîmes“ cu final, iar el a zis: „vroiam să plec“ și gata."
	if f := fragmente(text); len(f) != 0 {
		t.Errorf("citatele nu se ating: %v", f)
	}
}

// URL-ul se oprește la ghilimele, paranteze și punctuația de la sfârșit: textul
// lipit de el se verifică.
func TestUrlNuInghiteTextul(t *testing.T) {
	text := `<p>Vezi <a href="https://exemplu.ro/x">Şcoala noastră</a> și sînt aici, [sit](https://a.ro/b)cînd, https://a.ro/c.</p>`
	f := fragmente(text)
	if f["Şcoala"] != "RO-TIP-027" || f["cînd"] != "RO-TIP-128" {
		t.Errorf("textul de lângă URL trebuie verificat: %v", f)
	}
	ign := Ignorat(text)
	if !ign(strings.Index(text, "exemplu.ro")) || ign(strings.Index(text, "c.</p>")+1) {
		t.Error("URL-ul trebuie ignorat până la „.” final, fără el")
	}
}

// Runda a treia: formulele nu înghit proza ($$ de preț, „$” din cod sau din
// variabile), imperfectul cu majusculă se prinde, numele de familie și
// toponimele oficiale cu „î” nu se ating, URL-urile Wikipedia rămân întregi,
// iar punctuația lipită de spațiu se verifică și pe rândurile indentate.
func TestRunda3Lint(t *testing.T) {
	trebuie := func(text string, cuvinte ...string) {
		t.Helper()
		f := fragmente(text)
		for _, c := range cuvinte {
			if _, ok := f[c]; !ok {
				t.Errorf("%q: aștept observație pe %q, am %v", text, c, f)
			}
		}
	}
	nimic := func(text string) {
		t.Helper()
		if f := fragmente(text); len(f) != 0 {
			t.Errorf("%q: nicio observație, am %v", text, f)
		}
	}
	trebuie("Restaurant ieftin: $$ (mediu).\n\nCînd vii, ia pîine, mîine.\n\nRestaurant scump: $$$ (foarte scump).", "Cînd", "pîine", "mîine")
	trebuie("Variabilele $HOME și $PATH; cînd vii: \"$HOME\" mîine \"$PATH\" gata.", "cînd")
	trebuie("Prețul e US$a, cînd vii, nu US$b.", "cînd")
	trebuie("Din `$CERC/a.json` alege și sînt greșeli, apoi `$CERC/b.json` gata.", "sînt")
	nimic("Formula $x , y$ și $$a , b$$ sunt aici, iar drumul e lung pentru că sunt aici.")
	trebuie("Gîndea mult la asta.", "Gîndea")
	nimic("Cîndea a venit. Cîntea a spus. Stîncu a venit. Pînzaru a zis. Tîrzu a râs.")
	nimic("Sîngerei e raion. Rîșcani e oraș.")
	nimic("Vezi https://ro.wikipedia.org/wiki/Iași_(oraș) și https://x.ro/a_(Pîine)/d și https://x.ro/a'b/cînd azi.")
	trebuie("- punct unu\n    continuare cu text , și altceva", " ,")
	trebuie("    Stimate domn , vă scriu.", " ,")
	nimic("Iată, pentru că sunt aici:\n\n    x = a ? 1 : 2\n")
	nimic("Capitolul unu . . . . . . . . . 5\nCapitolul doi . . . . . . . . . 9\nși așa mai departe, pentru că sunt aici")
}

// Runda a patra: imperfectul plural și cuvintele uzuale cu majusculă se prind;
// semnul pomenit („Semnul ! poate”) și descrierea bibliografică nu sunt greșeli.
func TestRunda4Lint(t *testing.T) {
	f := fragmente("Mîncau pîine. Gîndeau mult. Rîdeau toți. Cîntau frumos. Pămîntean sau nu.")
	for _, c := range []string{"Mîncau", "Gîndeau", "Rîdeau", "Cîntau", "Pămîntean"} {
		if f[c] != "RO-TIP-128" {
			t.Errorf("%s trebuie semnalat: %v", c, f)
		}
	}
	for _, text := range []string{"Stîncu a venit. Tîrzu a râs. Pînzaru a zis.", "    Semnul ! poate marca intervenții, iar semnul ? arată îndoiala.",
		"    București : Univers Enciclopedic, 2005."} {
		if f := fragmente(text); len(f) != 0 {
			t.Errorf("%q: nicio observație, am %v", text, f)
		}
	}
	if f := fragmente("Ce faci ? Bine, mulțumesc, pentru că sunt aici."); f[" ?"] != "RO-TIP-068" {
		t.Errorf("„ ?” urmat de propoziție nouă trebuie semnalat: %v", f)
	}
}
