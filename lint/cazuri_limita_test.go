package lint

// Cazurile limită găsite la recenzia fazei 2: zonele HTML și front matter-ul,
// prefixoidele cu „î”, numele proprii și cuvintele străine.

import "testing"

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
