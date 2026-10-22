package lint

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
