package verifica

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/doITmagic/romana/dex"
)

// fals e un dicționar de probă: o hartă de forme (normalizate).
type fals map[string]bool

func (f fals) Exista(_ context.Context, forma string) (bool, error) {
	return f[dex.Normalizeaza(forma)], nil
}

func (f fals) Sugestii(_ context.Context, forma string, max int) ([]string, error) {
	var rez []string
	for k := range f {
		if dex.Cheie(k) == dex.Cheie(forma) && k != dex.Normalizeaza(forma) {
			rez = append(rez, k)
		}
	}
	sort.Strings(rez)
	if len(rez) > max {
		rez = rez[:max]
	}
	return rez, nil
}

func dictionar(forme string) fals {
	f := fals{}
	for _, w := range strings.Fields(forme) {
		f[w] = true
	}
	return f
}

var baza = dictionar("sau e am lista casă a dă mi nu o în caz limită se el ea este un doar și cuvânt sunt pe acesta text unde")

func cuvinte(r Raport) []string {
	var o []string
	for _, n := range r.Necunoscute {
		o = append(o, n.Cuvant)
	}
	return o
}

func TestNecunoscute(t *testing.T) {
	r, err := Text(context.Background(), baza, "Am sumarizată lista")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Necunoscute) != 1 {
		t.Fatalf("necunoscute: %+v", r.Necunoscute)
	}
	n := r.Necunoscute[0]
	if n.Cuvant != "sumarizată" || n.Neverificat || len(n.Sugestii) != 0 {
		t.Errorf("necunoscut: %+v", n)
	}
	if got := "Am sumarizată lista"[n.Inceput:n.Sfarsit]; got != "sumarizată" {
		t.Errorf("poziție greșită: %q", got)
	}
}

func TestSugestiiDiacritice(t *testing.T) {
	r, err := Text(context.Background(), baza, "Aceasta e casa")
	if err != nil {
		t.Fatal(err)
	}
	var gasit *Necunoscut
	for i := range r.Necunoscute {
		if r.Necunoscute[i].Cuvant == "casa" {
			gasit = &r.Necunoscute[i]
		}
	}
	if gasit == nil || !reflect.DeepEqual(gasit.Sugestii, []string{"casă"}) {
		t.Errorf("casa: %+v", r.Necunoscute)
	}
}

func TestCliticeSiCompuse(t *testing.T) {
	// „s-a” și „într-o” se sprijină pe fragmente elidate; „caz-limită” și
	// „dă-mi” pe părți; „nu-l” pe „nu” și fragmentul „l”.
	for _, text := range []string{
		"s-a", "într-o casă", "dă-mi", "nu-l", "caz-limită", "Nu-l dă-mi și s-a într-o casă",
		"s’a dă’mi", "„nu-l” (caz-limită).",
	} {
		r, err := Text(context.Background(), baza, text)
		if err != nil {
			t.Fatal(err)
		}
		if len(r.Necunoscute) != 0 {
			t.Errorf("%q: %+v", text, r.Necunoscute)
		}
	}
	// o parte necunoscută face compusul necunoscut
	r, _ := Text(context.Background(), baza, "dă-zzz")
	if len(r.Necunoscute) != 1 || r.Necunoscute[0].Cuvant != "dă-zzz" {
		t.Errorf("dă-zzz: %+v", r.Necunoscute)
	}
	// prefixele nu se descompun
	r, _ = Text(context.Background(), baza, "supracaz")
	if len(r.Necunoscute) != 1 {
		t.Errorf("supracaz: %+v", r.Necunoscute)
	}
}

func TestNeverificat(t *testing.T) {
	text := strings.Join([]string{
		"Acesta este un text.",
		"Am vorbit cu Zbigniew despre ONU și despre ceva-ce-nu-exista.",
		"Vezi https://exemplu.ro/zzqx și scrie la zzqx@exemplu.ro în 2024 la mp3 sau x_yz.",
		"Codul `zzqxfunc()` rămâne.",
		"```go",
		"func qqwwee() { zzrr := 1 }",
		"```",
		"Pe GitHub unde Привет e doar un cuvânt.",
	}, "\n")
	r, err := Text(context.Background(), baza, text)
	if err != nil {
		t.Fatal(err)
	}
	for _, n := range r.Necunoscute {
		switch n.Cuvant {
		case "Zbigniew", "ONU", "GitHub", "Привет":
			if !n.Neverificat || n.Motiv == "" {
				t.Errorf("%q trebuia Neverificat: %+v", n.Cuvant, n)
			}
		case "ceva-ce-nu-exista", "vorbit", "ceva", "despre", "cu", "Vezi", "scrie", "Codul", "rămâne", "Am", "la":
			// cuvinte obișnuite lipsite din dicționarul de probă: sunt necunoscute reale
			if n.Neverificat {
				t.Errorf("%q nu trebuia Neverificat: %+v", n.Cuvant, n)
			}
		default:
			t.Errorf("cuvânt neașteptat în raport: %+v", n)
		}
	}
}

func TestLimitaSiCRLF(t *testing.T) {
	_, err := Text(context.Background(), baza, strings.Repeat("a", MaxCaractere+1))
	if !errors.Is(err, ErrPreaMare) {
		t.Errorf("trebuia ErrPreaMare, am primit %v", err)
	}
	// limita exactă e acceptată
	if _, err := Text(context.Background(), baza, strings.Repeat("a ", MaxCaractere/2)); err != nil {
		t.Errorf("limita exactă: %v", err)
	}
	lf := "Am sumarizată lista\nși zzqx  ,\nAm foo."
	crlf := strings.ReplaceAll(lf, "\n", "\r\n")
	a, err1 := Text(context.Background(), baza, lf)
	b, err2 := Text(context.Background(), baza, crlf)
	if err1 != nil || err2 != nil {
		t.Fatal(err1, err2)
	}
	if !reflect.DeepEqual(cuvinte(a), cuvinte(b)) || len(a.Lint) != len(b.Lint) {
		t.Errorf("LF %v/%d, CRLF %v/%d", cuvinte(a), len(a.Lint), cuvinte(b), len(b.Lint))
	}
	if len(a.Lint) == 0 {
		t.Error("linterul trebuia să vadă spațiul înainte de virgulă")
	}
}

func TestLinterInclus(t *testing.T) {
	r, err := Text(context.Background(), baza, "Ei vroiam să plece.")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Lint) == 0 {
		t.Error("lipsesc observațiile linterului")
	}
}

type stricat struct{ fals }

func (stricat) Exista(context.Context, string) (bool, error) {
	return false, errors.New("baza a căzut")
}

func TestEroareDeLaBaza(t *testing.T) {
	if _, err := Text(context.Background(), stricat{}, "casă"); err == nil {
		t.Error("eroarea bazei trebuia propagată")
	}
}

func TestContextAnulat(t *testing.T) {
	ctx, anuleaza := context.WithCancel(context.Background())
	anuleaza()
	if _, err := Text(ctx, baza, "casă"); err == nil {
		t.Error("contextul anulat trebuia respectat")
	}
}

// paragraful e scris de noi; trebuie să nu dea necunoscute nemarcate
// „neverificat” pe baza reală.
const paragraf = `Într-o dimineață de toamnă, bunicul s-a trezit devreme și a ieșit în grădină. ` +
	`Merele coapte atârnau greu în pomi, iar pe alee zăceau frunze galbene. ` +
	`„Nu-l lăsa pe copil singur lângă fântână”, i-a spus bunica, uitându-se pe fereastră. ` +
	`Dă-mi, te rog, coșul cel mare, ca să strâng fructele înainte de ploaie. ` +
	`Vecinii, care munceau de ani buni în aceeași comună, au venit să ajute, ` +
	`iar spre seară toată lumea s-a așezat la masă. ` +
	`Au mâncat o supă fierbinte, au povestit întâmplări vechi și au hotărât că, ` +
	`într-un caz-limită, mai bine ceri ajutor decât să pierzi recolta. ` +
	`Copiii au adormit repede, obosiți după o zi plină de joacă și râsete.`

func TestParagrafRealCorect(t *testing.T) {
	cale := os.Getenv("ROMANA_DEX")
	if cale == "" {
		acasa, err := os.UserHomeDir()
		if err != nil {
			t.Skip("fără director acasă")
		}
		cale = filepath.Join(acasa, ".local", "share", "romana", "dex.sqlite")
	}
	if _, err := os.Stat(cale); err != nil {
		t.Skipf("fără baza reală (%s)", cale)
	}
	d, err := dex.Deschide(cale)
	if err != nil {
		t.Skipf("baza nu se deschide: %v", err)
	}
	defer d.Close()
	if n := len(strings.Fields(paragraf)); n < 70 {
		t.Fatalf("paragraful are doar %d cuvinte", n)
	}
	r, err := Text(context.Background(), d, paragraf)
	if err != nil {
		t.Fatal(err)
	}
	for _, n := range r.Necunoscute {
		if !n.Neverificat {
			t.Errorf("necunoscut nemarcat: %+v", n)
		}
	}
	if len(r.Lint) != 0 {
		t.Errorf("observații de linter pe text corect: %+v", r.Lint)
	}
	// și un cuvânt inventat trebuie să iasă
	r, err = Text(context.Background(), d, "Am sumarizată lista și zzqxvb.")
	if err != nil {
		t.Fatal(err)
	}
	var zz bool
	for _, n := range r.Necunoscute {
		zz = zz || (n.Cuvant == "zzqxvb" && !n.Neverificat)
	}
	if !zz {
		t.Errorf("zzqxvb trebuia găsit: %+v", r.Necunoscute)
	}
}

// numara numără apelurile către dicționar.
type numara struct {
	fals
	sugestii, exista int
}

func (n *numara) Exista(ctx context.Context, f string) (bool, error) {
	n.exista++
	return n.fals.Exista(ctx, f)
}

func (n *numara) Sugestii(ctx context.Context, f string, m int) ([]string, error) {
	n.sugestii++
	return n.fals.Sugestii(ctx, f, m)
}

func TestSugestiiMemorate(t *testing.T) {
	c := &numara{fals: baza}
	r, err := Text(context.Background(), c, strings.Repeat("casa sumarizată ", 50))
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Necunoscute) != 100 {
		t.Fatalf("necunoscute: %d", len(r.Necunoscute))
	}
	if c.sugestii != 2 {
		t.Errorf("Sugestii apelat de %d ori, aștept 2 (casa, sumarizată)", c.sugestii)
	}
	for _, u := range r.Necunoscute {
		if u.Cuvant == "casa" && !reflect.DeepEqual(u.Sugestii, []string{"casă"}) {
			t.Fatalf("sugestii pentru %q: %v", u.Cuvant, u.Sugestii)
		}
	}
}

// Prefixele care se scriu legat nu se descompun după cratimă; „ex-” („fost”)
// se scrie cu cratimă, deci „ex-soț” e corect dacă „soț” există (RO-IMP-013).
func TestPrefixeleNuSeDescompun(t *testing.T) {
	b := dictionar("anti soț ex nici o caz limită nato")
	r, err := Text(context.Background(), b, "ex-soț nici-o anti-nato caz-limită")
	if err != nil {
		t.Fatal(err)
	}
	got := cuvinte(r)
	want := []string{"nici-o", "anti-nato"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("necunoscute %v, aștept %v", got, want)
	}
}

func TestMarkdownSiHtml(t *testing.T) {
	r, err := Text(context.Background(), baza, "Un text <div class=\"foo\"> și ![sumarizată](a.png) și [text](pagina.md) unde")
	if err != nil {
		t.Fatal(err)
	}
	if got := cuvinte(r); !reflect.DeepEqual(got, []string{"sumarizată"}) {
		t.Errorf("necunoscute %v, aștept doar textul-alternativ", got)
	}
}

func TestGrafieVecheFaraSugestieZgomotoasa(t *testing.T) {
	b := dictionar("sunt sint text")
	r, err := Text(context.Background(), b, "Noi sînt acolo")
	if err != nil {
		t.Fatal(err)
	}
	for _, u := range r.Necunoscute {
		if u.Cuvant == "sînt" && len(u.Sugestii) != 0 {
			t.Errorf("sînt are sugestii %v deși linterul propune „sunt”", u.Sugestii)
		}
	}
	if len(r.Lint) == 0 {
		t.Error("linterul ar trebui să marcheze sînt")
	}
}

// Raportul are textul cu greșelile sigure corectate; cuvintele necunoscute
// rămân neatinse, pentru că nu sunt greșeli sigure.
func TestTextCorectat(t *testing.T) {
	r, err := Text(context.Background(), baza, "Aşa  e zzz .")
	if err != nil {
		t.Fatal(err)
	}
	if r.TextCorectat != "Așa e zzz." {
		t.Errorf("text_corectat = %q", r.TextCorectat)
	}
}
