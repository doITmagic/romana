package evals

import "testing"

func TestNotareNeatins(t *testing.T) {
	c := Caz{ID: "curat-x-01", Mod: "corectura", Intrare: "Ședința a durat două ore.", Asteptat: Asteptat{Neatins: true}}
	n := Noteaza(c, Rezultat{Iesire: "[[TEXT]]\nȘedința a durat două ore.\n[[/TEXT]]\nFără schimbări."})
	if n.FalsPozitive != 0 || n.Probleme != 0 {
		t.Errorf("ieșire identică: %+v", n)
	}
	n = Noteaza(c, Rezultat{Iesire: "[[TEXT]]\nȘedința a durat trei ore.\n[[/TEXT]]"})
	if n.FalsPozitive != 1 {
		t.Errorf("un cuvânt schimbat: %+v", n)
	}
	// Doar spațiile diferă: nu e o schimbare.
	n = Noteaza(c, Rezultat{Iesire: "[[TEXT]]\nȘedința   a durat două ore.\n[[/TEXT]]"})
	if n.FalsPozitive != 0 {
		t.Errorf("spații: %+v", n)
	}
}

func TestNotareNeatinsVerificare(t *testing.T) {
	c := Caz{ID: "curat-x-02", Mod: "verificare", Intrare: "Casa sa se află la margine.", Asteptat: Asteptat{Neatins: true}}
	n := Noteaza(c, Rezultat{Iesire: "Textul nu are probleme. Nu semnalez nimic."})
	if n.FalsPozitive != 0 {
		t.Errorf("fără observații: %+v", n)
	}
	n = Noteaza(c, Rezultat{Iesire: "1. OBS: „Casa sa” · gramatică · normă · P1 · altceva"})
	if n.FalsPozitive != 1 {
		t.Errorf("o observație pe text curat: %+v", n)
	}
}

func TestNotareInterzis(t *testing.T) {
	c := Caz{ID: "context-01", Mod: "verificare", Intrare: "Casa sa se află acolo. Guvernul sa întrunit.",
		Asteptat: Asteptat{
			Observatii: []ObsAsteptata{{Fragment: "Guvernul sa întrunit", Nivel: "normă", Pilon: "gramatică"}},
			Interzis:   []string{"Casa sa"},
		}}
	bun := "- OBS: „Guvernul sa întrunit” · gramatică · normă · P1 · „s-a întrunit”"
	if n := Noteaza(c, Rezultat{Iesire: bun}); n.Probleme != 0 || n.Lipsa != 0 || n.FalsPozitive != 0 {
		t.Errorf("răspuns bun: %+v", n)
	}
	rau := bun + "\n- OBS: „Casa sa” · gramatică · normă · P1 · „Casa s-a”"
	if n := Noteaza(c, Rezultat{Iesire: rau}); n.Probleme != 1 {
		t.Errorf("interzis semnalat: %+v", n)
	}
	// Corectură: schimbarea unui fragment interzis.
	cc := c
	cc.Mod = "corectura"
	n := Noteaza(cc, Rezultat{Iesire: "[[TEXT]]\nCasa s-a află acolo. Guvernul s-a întrunit.\n[[/TEXT]]"})
	if n.Probleme != 1 || n.Lipsa != 0 {
		t.Errorf("corectură, interzis schimbat: %+v", n)
	}
	// Redactare: forma interzisă apare în ieșire.
	cr := Caz{ID: "fidelitate-15", Mod: "redactare", Intrare: "Rescrie: ar putea scurta durata.",
		Asteptat: Asteptat{Interzis: []string{"va scurta"}}}
	if n := Noteaza(cr, Rezultat{Iesire: "Metoda VA SCURTA durata."}); n.Probleme != 1 || n.Fidelitate {
		t.Errorf("redactare interzis: %+v", n)
	}
}

func TestNotareFidelitate(t *testing.T) {
	c := Caz{ID: "pastrare-01", Mod: "corectura", Intrare: "Avem 25 ani de activitate.",
		Asteptat: Asteptat{Observatii: []ObsAsteptata{{Fragment: "25 ani", Nivel: "normă", Pilon: "tipografie"}}}}
	if n := Noteaza(c, Rezultat{Iesire: "[[TEXT]]Avem 25 de ani de activitate.[[/TEXT]]"}); !n.Fidelitate {
		t.Errorf("fără număr nou: %+v", n)
	}
	if n := Noteaza(c, Rezultat{Iesire: "[[TEXT]]Avem 30 de ani de activitate.[[/TEXT]]"}); n.Fidelitate {
		t.Errorf("număr nou: %+v", n)
	}
	cr := Caz{ID: "fidelitate-16", Mod: "redactare", Intrare: "Scrie despre 3 produse.", Asteptat: Asteptat{Cerut: []string{"[de completat"}}}
	n := Noteaza(cr, Rezultat{Iesire: "Oferim 3 produse și 12 servicii."})
	if n.Fidelitate || n.Lipsa != 1 {
		t.Errorf("redactare: %+v", n)
	}
}

func TestNotareObservatii(t *testing.T) {
	c := Caz{ID: "corectura-02", Mod: "verificare",
		Intrare: "Firma are 25 ani de activitate. Colaborăm ca și subcontractor.",
		Asteptat: Asteptat{Observatii: []ObsAsteptata{
			{Fragment: "25 ani", Nivel: "normă", Pilon: "tipografie"},
			{Fragment: "ca și subcontractor", Nivel: "normă", Pilon: "gramatică"},
		}}}
	out := "Rezumat: 1 normă.\nOBS: „25 ani” · tipografie · normă · P1 · „25 de ani” · Numeralul cere „de”.\nFormă corectă: „ca și subcontractor” nu e citat aici.\nneverificat:\n- „Firma” · nu știu"
	n := Noteaza(c, Rezultat{Iesire: out})
	if n.Lipsa != 1 || n.FalsPozitive != 0 {
		t.Errorf("o observație lipsă: %+v", n)
	}
	// Citat parțial din fragmentul așteptat: contează.
	n = Noteaza(c, Rezultat{Iesire: out + "\nOBS: „ca și” · gramatică · normă · P1"})
	_ = n
	out2 := "1. OBS: „25 ani” · normă\n2. OBS: „ca și” · normă\n3. OBS: „Firma” · normă\n4. Rezumat: „Casa” corectă"
	n = Noteaza(c, Rezultat{Iesire: out2})
	if n.Lipsa != 0 || n.FalsPozitive != 1 {
		t.Errorf("potrivire parțială + observație în plus: %+v", n)
	}
}

func TestNotareCerut(t *testing.T) {
	c := Caz{ID: "corectura-02", Mod: "corectura", Intrare: "Are 25 ani.",
		Asteptat: Asteptat{Observatii: []ObsAsteptata{{Fragment: "25 ani", Nivel: "normă", Pilon: "tipografie"}}, Cerut: []string{"25 de ani"}}}
	if n := Noteaza(c, Rezultat{Iesire: "[[TEXT]]Are 25 de ani.[[/TEXT]]"}); n.Lipsa != 0 || n.FalsPozitive != 0 {
		t.Errorf("corectat: %+v", n)
	}
	if n := Noteaza(c, Rezultat{Iesire: "[[TEXT]]Are 25 ani.[[/TEXT]]"}); n.Lipsa != 2 {
		t.Errorf("necorectat (observație + cerut): %+v", n)
	}
}

func TestNotareTrigger(t *testing.T) {
	poz := Caz{ID: "trigger-pozitiv-01", Mod: "trigger", Intrare: "x", Asteptat: Asteptat{Cerut: []string{"romana"}}}
	neg := Caz{ID: "trigger-negativ-01", Mod: "trigger", Intrare: "x", Asteptat: Asteptat{Interzis: []string{"romana"}}}
	if n := Noteaza(poz, Rezultat{PornitCunoscut: true, Pornit: false}); n.Lipsa != 1 {
		t.Errorf("%+v", n)
	}
	if n := Noteaza(poz, Rezultat{PornitCunoscut: true, Pornit: true}); n.Lipsa != 0 {
		t.Errorf("%+v", n)
	}
	if n := Noteaza(neg, Rezultat{PornitCunoscut: true, Pornit: true}); n.Probleme != 1 || n.FalsPozitive != 1 {
		t.Errorf("%+v", n)
	}
	if n := Noteaza(neg, Rezultat{PornitCunoscut: false}); n != (Nota{Fidelitate: true, Nemasurat: true}) {
		t.Errorf("necunoscut: %+v", n)
	}
}

func TestExtrageText(t *testing.T) {
	if got, ok := ExtrageText("a\n[[TEXT]]\nb c\n[[/TEXT]]\nz"); !ok || got != "b c" {
		t.Errorf("%q %v", got, ok)
	}
	if got, ok := ExtrageText("fără marcaje"); ok || got != "fără marcaje" {
		t.Errorf("%q %v", got, ok)
	}
}

// --- atacurile din recenzia 1 ---

func cazDouaObs(mod string) Caz {
	return Caz{ID: "x-01", Mod: mod,
		Intrare: "Casa sa se află acolo, la marginea orașului. Marți, Guvernul sa întrunit într-o ședință extraordinară. Firma are 25 ani de activitate pe piața instalațiilor.",
		Asteptat: Asteptat{
			Observatii: []ObsAsteptata{
				{Fragment: "Guvernul sa întrunit", Nivel: "normă", Pilon: "gramatică"},
				{Fragment: "25 ani", Nivel: "normă", Pilon: "tipografie"},
			},
			Interzis: []string{"Casa sa"},
		}}
}

func TestCitatulTotTextul(t *testing.T) {
	c := cazDouaObs("verificare")
	n := Noteaza(c, Rezultat{Iesire: "OBS: „" + c.Intrare + "” totul"})
	if n.Lipsa != 2 || n.Probleme != 1 || n.FalsPozitive != 1 {
		t.Errorf("citat = tot textul nu trebuie să dea nota perfectă: %+v", n)
	}
}

func TestRescriereTotala(t *testing.T) {
	c := cazDouaObs("corectura")
	n := Noteaza(c, Rezultat{Iesire: "[[TEXT]]\nO locuință liniștită. Executivul a avut o ședință. Compania activează de un sfert de secol, o garanție de seriozitate.\n[[/TEXT]]"})
	if n.Probleme != 1 || n.FalsPozitive < 5 {
		t.Errorf("retranscrierea nu e o corectură: %+v", n)
	}
}

func TestInterzisAdiacentObservatiei(t *testing.T) {
	c := cazDouaObs("corectura")
	// O singură zonă schimbată acoperă și observația, și interzisul.
	n := Noteaza(c, Rezultat{Iesire: "[[TEXT]]\nCasa sa se află acolo, la marginea orașului. Marți, Guvernul s-a întrunit într-o ședință extraordinară. Firma are 25 de ani de activitate pe piața instalațiilor.\n[[/TEXT]]"})
	if n.Probleme != 0 || n.Lipsa != 0 || n.FalsPozitive != 0 {
		t.Errorf("corectură bună: %+v", n)
	}
	n = Noteaza(c, Rezultat{Iesire: "[[TEXT]]\nLocuința lui e la marginea orașului. Marți, Guvernul s-a întrunit într-o ședință extraordinară. Firma are 25 de ani de activitate pe piața instalațiilor.\n[[/TEXT]]"})
	if n.Probleme != 1 {
		t.Errorf("interzis schimbat: %+v", n)
	}
	// Zonă contiguă care cuprinde interzisul și observația vecină.
	c2 := Caz{ID: "x-02", Mod: "corectura", Intrare: "Casa sa a fost vândută. Guvernul sa întrunit.",
		Asteptat: Asteptat{
			Observatii: []ObsAsteptata{{Fragment: "Guvernul sa întrunit", Nivel: "normă", Pilon: "gramatică"}},
			Interzis:   []string{"Casa sa"}}}
	n = Noteaza(c2, Rezultat{Iesire: "[[TEXT]]Locuința a fost vândută. Executivul s-a întrunit.[[/TEXT]]"})
	if n.Probleme != 1 {
		t.Errorf("interzisul atins de o zonă care atinge și observația: %+v", n)
	}
}

func TestNeatinsInTokeni(t *testing.T) {
	c := Caz{ID: "curat-x-03", Mod: "corectura", Intrare: "Ședința a durat două ore și s-a încheiat târziu.", Asteptat: Asteptat{Neatins: true}}
	n := Noteaza(c, Rezultat{Iesire: "[[TEXT]]Reuniunea a avut loc zece ore și s-a închis devreme.[[/TEXT]]"})
	if n.FalsPozitive < 5 {
		t.Errorf("FP se măsoară în tokeni schimbați: %+v", n)
	}
}

func TestNumereMiiCuSpatiu(t *testing.T) {
	c := Caz{ID: "x-04", Mod: "corectura", Intrare: "Costă 1.000 lei, 12,5 % din 8 000 lei.", Asteptat: Asteptat{Neatins: true}}
	if n := Noteaza(c, Rezultat{Iesire: "[[TEXT]]Costă 1000 lei, 12,5% din 8000 lei.[[/TEXT]]"}); !n.Fidelitate {
		t.Errorf("reformatare: %+v", n)
	}
	c.Intrare = "Costă 8\u00a0000 lei."
	if n := Noteaza(c, Rezultat{Iesire: "[[TEXT]]Costă 8000 lei.[[/TEXT]]"}); !n.Fidelitate {
		t.Errorf("NBSP: %+v", n)
	}
	// „12,5” nu se confundă cu „125”.
	c.Intrare = "Costă 12,5 lei."
	if n := Noteaza(c, Rezultat{Iesire: "[[TEXT]]Costă 125 lei.[[/TEXT]]"}); n.Fidelitate {
		t.Errorf("12,5 vs 125: %+v", n)
	}
}

func TestNBSPNuEOSchimbare(t *testing.T) {
	c := Caz{ID: "curat-x-05", Mod: "corectura", Intrare: "Avem 10\u00a0% din total.", Asteptat: Asteptat{Neatins: true}}
	if n := Noteaza(c, Rezultat{Iesire: "[[TEXT]]Avem 10 % din total.[[/TEXT]]"}); n.FalsPozitive != 0 {
		t.Errorf("%+v", n)
	}
}

func TestNemasurat(t *testing.T) {
	cr := Caz{ID: "redactare-blog-01", Mod: "redactare", Intrare: "Scrie.", Asteptat: Asteptat{Interzis: []string{"zzz"}}}
	if n := Noteaza(cr, Rezultat{Iesire: "text"}); !n.Nemasurat {
		t.Errorf("redactare fără judecător: %+v", n)
	}
	n := Noteaza(cr, Rezultat{Iesire: "text", Judecata: &Judecata{Calcuri: 1, Coerenta: 2}})
	if n.Nemasurat || n.Judecator != 3 || n.Probleme != 0 {
		t.Errorf("judecătorul se raportează separat de Interzis: %+v", n)
	}
}

func TestRedactareNotataDoarPeText(t *testing.T) {
	cr := Caz{ID: "redactare-blog-01", Mod: "redactare", Intrare: "Scrie.", Asteptat: Asteptat{Interzis: []string{"va scurta"}}}
	n := Noteaza(cr, Rezultat{Iesire: "Am evitat „va scurta”.\n[[TEXT]]Metoda ar putea scurta durata.[[/TEXT]]"})
	if n.Probleme != 0 {
		t.Errorf("mențiune în comentariu: %+v", n)
	}
}

func TestCerutLaVerificareInObservatii(t *testing.T) {
	c := Caz{ID: "corectura-02", Mod: "verificare", Intrare: "Are 25 ani.",
		Asteptat: Asteptat{Observatii: []ObsAsteptata{{Fragment: "25 ani", Nivel: "normă", Pilon: "tipografie"}}, Cerut: []string{"25 de ani"}}}
	if n := Noteaza(c, Rezultat{Iesire: "Rezumat: ar fi „25 de ani”.\nOBS: „25 ani” · normă"}); n.Lipsa != 1 {
		t.Errorf("cerut doar în rezumat: %+v", n)
	}
	if n := Noteaza(c, Rezultat{Iesire: "OBS: „25 ani” → „25 de ani” · normă"}); n.Lipsa != 0 {
		t.Errorf("%+v", n)
	}
}

func TestCitatScurt(t *testing.T) {
	c := Caz{ID: "x-05", Mod: "verificare", Intrare: "Ai o idee bună.",
		Asteptat: Asteptat{Observatii: []ObsAsteptata{{Fragment: "idee", Nivel: "normă", Pilon: "semantică"}}}}
	n := Noteaza(c, Rezultat{Iesire: "OBS: „a” · normă"})
	if n.Lipsa != 1 {
		t.Errorf("un citat de o literă nu prinde nimic: %+v", n)
	}
}

func TestInterzisPeLimitaDeCuvant(t *testing.T) {
	cr := func(interzis string) Caz {
		return Caz{ID: "redactare-interfata-02", Mod: "redactare", Intrare: "Scrie.", Asteptat: Asteptat{Interzis: []string{interzis}}}
	}
	j := &Judecata{}
	if n := Noteaza(cr("Introdu"), Rezultat{Iesire: "[[TEXT]]Introduceți o adresă validă.[[/TEXT]]", Judecata: j}); n.Probleme != 0 || !n.Fidelitate {
		t.Errorf("Introdu nu se potrivește cu Introduceți: %+v", n)
	}
	if n := Noteaza(cr("Introdu"), Rezultat{Iesire: "[[TEXT]]Introdu o adresă validă.[[/TEXT]]", Judecata: j}); n.Probleme != 1 {
		t.Errorf("forma interzisă: %+v", n)
	}
	if n := Noteaza(cr("a aplica"), Rezultat{Iesire: "[[TEXT]]Primăria va aplica măsura.[[/TEXT]]", Judecata: j}); n.Probleme != 0 {
		t.Errorf("a aplica nu se potrivește cu va aplica: %+v", n)
	}
	if n := Noteaza(cr("a aplica"), Rezultat{Iesire: "[[TEXT]]Poți a aplica acum.[[/TEXT]]", Judecata: j}); n.Probleme != 1 {
		t.Errorf("%+v", n)
	}
	// diacritice și cratimă fac parte din cuvânt
	if n := Noteaza(cr("sigur"), Rezultat{Iesire: "[[TEXT]]Ești asigurat; e un sigur-sigur.[[/TEXT]]", Judecata: j}); n.Probleme != 0 {
		t.Errorf("%+v", n)
	}
	if n := Noteaza(cr("ședință"), Rezultat{Iesire: "[[TEXT]]Reședință nouă.[[/TEXT]]", Judecata: j}); n.Probleme != 0 {
		t.Errorf("diacritice: %+v", n)
	}
	if n := Noteaza(cr("Introdu"), Rezultat{Iesire: "[[TEXT]]Introdu-ți parola.[[/TEXT]]", Judecata: j}); n.Probleme != 0 {
		t.Errorf("cratima leagă cuvântul: %+v", n)
	}
	if n := Noteaza(cr("Introdu-ți"), Rezultat{Iesire: "[[TEXT]]Introdu-ți parola.[[/TEXT]]", Judecata: j}); n.Probleme != 1 {
		t.Errorf("%+v", n)
	}
}

func TestInterzisPeLimitaDeCuvantCorectura(t *testing.T) {
	c := Caz{ID: "x-06", Mod: "corectura", Intrare: "Introduceți datele. Ea va aplica regula.",
		Asteptat: Asteptat{Interzis: []string{"Introdu", "a aplica"}}}
	// Schimbarea lui „Introduceți” nu atinge „Introdu”, nici „va aplica” pe „a aplica”.
	n := Noteaza(c, Rezultat{Iesire: "[[TEXT]]Completați datele. Ea va folosi regula.[[/TEXT]]"})
	if n.Probleme != 0 {
		t.Errorf("%+v", n)
	}
}

func TestInterzisPeLimitaDeCuvantVerificare(t *testing.T) {
	c := Caz{ID: "x-07", Mod: "verificare", Intrare: "Introduceți datele. Ea va aplica regula.",
		Asteptat: Asteptat{Interzis: []string{"Introdu", "a aplica"}}}
	n := Noteaza(c, Rezultat{Iesire: "OBS: „Introduceți datele” · normă\nOBS: „va aplica regula” · normă"})
	if n.Probleme != 0 {
		t.Errorf("%+v", n)
	}
}

// La corectură, o observație care nu e normă (recomandare, de semnalat etc.)
// nu se aplică obligatoriu în text: ajunge s-o propună în afara textului
// corectat (rânduri OBS: sau secțiunea de după [[/TEXT]]). O normă trebuie în
// continuare corectată în text.
func TestCorecturaRecomandareInPropuneri(t *testing.T) {
	c := Caz{ID: "corectura-x-09", Mod: "corectura",
		Intrare: "Are 25 ani. Am făcut un update la aplicație.",
		Asteptat: Asteptat{Observatii: []ObsAsteptata{
			{Fragment: "25 ani", Nivel: "normă", Pilon: "tipografie"},
			{Fragment: "un update", Nivel: "recomandare", Pilon: "semantică"},
		}}}
	text := "[[TEXT]]\nAre 25 de ani. Am făcut un update la aplicație.\n[[/TEXT]]\n"
	// Recomandarea propusă după text, fără schimbare în text: nu lipsește.
	if n := Noteaza(c, Rezultat{Iesire: text + "Propun și: „un update” → „o actualizare”."}); n.Lipsa != 0 || n.FalsPozitive != 0 {
		t.Errorf("recomandare propusă după [[/TEXT]]: %+v", n)
	}
	// Recomandarea pe un rând OBS: nu lipsește.
	if n := Noteaza(c, Rezultat{Iesire: "OBS: „un update” · recomandare · „o actualizare”\n" + text}); n.Lipsa != 0 {
		t.Errorf("recomandare pe rând OBS: %+v", n)
	}
	// Recomandarea nici aplicată, nici propusă: lipsește.
	if n := Noteaza(c, Rezultat{Iesire: text + "Am corectat numeralul."}); n.Lipsa != 1 {
		t.Errorf("recomandare nepropusă: %+v", n)
	}
	// Recomandarea aplicată în text: nu lipsește (ca înainte).
	if n := Noteaza(c, Rezultat{Iesire: "[[TEXT]]\nAre 25 de ani. Am făcut o actualizare la aplicație.\n[[/TEXT]]"}); n.Lipsa != 0 || n.FalsPozitive != 0 {
		t.Errorf("recomandare aplicată: %+v", n)
	}
	// Norma doar propusă, necorectată în text: lipsește.
	n := Noteaza(c, Rezultat{Iesire: "[[TEXT]]\nAre 25 ani. Am făcut un update la aplicație.\n[[/TEXT]]\nOBS: „25 ani” → „25 de ani”\nOBS: „un update” → „o actualizare”"})
	if n.Lipsa != 1 {
		t.Errorf("normă doar propusă: %+v", n)
	}
}
