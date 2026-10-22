package lint

import "testing"

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
