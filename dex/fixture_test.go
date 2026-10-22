package dex

import (
	"strings"
)

// Fixture-ul de test: un export SQL mic în formatul dexonline, cu date
// inventate de noi (nicio definiție nu e copiată din export). Textele
// „TEXT-DOOM3-UNIC” și „TEXT-DEX09-UNIC” nu trebuie să iasă niciodată din
// pachet, pentru că DOOM 3 și DEX ’09 au can_distribute=0.

const (
	subsirDoom3 = "TEXT-DOOM3-UNIC"
	subsirDex09 = "TEXT-DEX09-UNIC"
)

// tuplu face un rând cu n coloane, implicit '0', cu valorile date (deja în
// sintaxa SQL).
func tuplu(n int, v map[int]string) string {
	c := make([]string, n)
	for i := range c {
		c[i] = "'0'"
		if x, ok := v[i]; ok {
			c[i] = x
		}
	}
	return "(" + strings.Join(c, ",") + ")"
}

func insert(tabel string, randuri ...string) string {
	return "INSERT INTO `" + tabel + "` VALUES " + strings.Join(randuri, ",") + ";\n"
}

// dumpMic imită exportul: comentarii, CREATE TABLE, tabele de prisos și
// ordinea reală (Source și Definition la urmă, dumpuri lipite).
func dumpMic() string {
	q := func(s string) string { return "'" + s + "'" }
	src := func(id, nume, dist string) string {
		return tuplu(29, map[int]string{0: id, 1: q(nume), 3: q("Dicționar inventat"), 18: "NULL", 19: dist})
	}
	// Definition: 0 id, 2 sourceId, 3 lexicon, 4 internalRep, 5 status
	def := func(id, sursa, lexicon, status, text string) string {
		return tuplu(14, map[int]string{0: id, 2: sursa, 3: q(lexicon), 4: q(text), 5: status})
	}
	// Lexeme: 0 id, 2 formNoAccent, 6 description, 14 modelType
	lex := func(id, forma, tip string) string {
		return tuplu(23, map[int]string{0: id, 1: q(forma), 2: q(forma), 5: "NULL", 6: "''", 10: "NULL", 14: q(tip)})
	}
	// InflectedForm: 2 formNoAccent, 4 lexemeId, 5 inflectionId, 6 variant, 7 recommended
	ifl := func(forma, lexem, infl, varianta, recomandat string) string {
		return tuplu(10, map[int]string{2: q(forma), 4: lexem, 5: infl, 6: varianta, 7: recomandat})
	}
	el := func(intrare, lexem string) string { return tuplu(8, map[int]string{1: intrare, 2: lexem}) }
	ed := func(intrare, def, rang string) string { return tuplu(7, map[int]string{1: intrare, 2: def, 4: rang}) }
	infl := func(id, descriere string) string { return tuplu(9, map[int]string{0: id, 1: q(descriere)}) }
	tip := func(id, cod, descriere string) string {
		return tuplu(4, map[int]string{0: id, 1: q(cod), 2: q(descriere)})
	}
	lsrc := func(lexem, sursa string) string { return tuplu(7, map[int]string{1: lexem, 2: sursa}) }

	var b strings.Builder
	b.WriteString("-- MySQL dump 10.19\nSET NAMES utf8mb4;\nCREATE TABLE `Lexeme` (\n  `id` int\n);\n")
	b.WriteString(insert("FullTextIndex", "(1,2,3)", "(4,'x),(y',6)"))
	b.WriteString(insert("ModelType", tip("2", "M", "substantiv masculin"), tip("3", "F", "substantiv feminin"),
		tip("8", "V", "verb"), tip("9", "VT", "verb")))
	b.WriteString(insert("Inflection",
		infl("1", "Substantiv feminin, Nominativ-Acuzativ, singular, nearticulat"),
		infl("3", "Substantiv feminin, Nominativ-Acuzativ, plural, nearticulat"),
		infl("5", "Substantiv feminin, Nominativ-Acuzativ, singular, articulat"),
		infl("50", "Verb, Infinitiv"),
		infl("51", "Verb, Indicativ, prezent, persoana I, singular"),
		infl("52", "Verb, Indicativ, imperfect, persoana I, singular")))
	// 1 casă (F), 2 vroi (V), 3 dejivra (VT), 4 caz-limită (F, fără forme în tabelă)
	b.WriteString(insert("Lexeme", lex("1", "casă", "F"), lex("2", "vroi", "V"), lex("3", "dejivra", "VT"), lex("4", "caz-limită", "F")))
	b.WriteString(insert("InflectedForm",
		ifl("casă", "1", "1", "0", "1"), ifl("case", "1", "3", "0", "1"), ifl("casa", "1", "5", "0", "NULL"),
		ifl("Casă", "1", "1", "0", "1"), // dublură cu literă mare: se contopește
		ifl("vroi", "2", "50", "0", "1"), ifl("vroiam", "2", "52", "1", "0"),
		ifl("dejivra", "3", "50", "0", "1"), ifl("dejivrez", "3", "51", "0", "1")))
	b.WriteString(insert("LexemeSource", lsrc("3", "88")))
	// intrările: 10 casă, 20 vroi, 30 dejivra, 40 caz-limită, 41 zi-lumină (fără lexem), 50 din sursa ascunsă
	b.WriteString(insert("EntryLexeme", el("10", "1"), el("20", "2"), el("30", "3"), el("40", "4"), el("50", "1")))
	b.WriteString(insert("EntryDefinition",
		ed("10", "1", "1"), ed("10", "2", "2"), ed("10", "3", "3"), ed("20", "4", "1"), ed("30", "5", "1"),
		ed("40", "6", "1"), ed("41", "7", "1"), ed("10", "8", "4"), ed("50", "9", "1")))
	b.WriteString("-- Dump completed\n-- MySQL dump 10.19\n")
	b.WriteString(insert("Source", src("1", `DEX \'98`, "1"), src("7", "DCR2", "1"), src("27", `DEX \'09`, "0"),
		src("88", "DOOM 3", "0"),
		// sursă ascunsă de dexonline (hidden, coloana 12): nu intră, nici definițiile ei
		tuplu(29, map[int]string{0: "45", 1: "'ASCUNS'", 3: "'x'", 12: "1", 18: "NULL", 19: "1"})))
	b.WriteString(insert("Definition",
		def("1", "1", "casă", "0", `@CASĂ@ #s. f.# Clădire de locuit (a),(b) \"x\". - Din #lat.# @casa@.`),
		def("2", "88", "casă", "0", `@casă@ `+subsirDoom3+` forma corectă`),
		def("3", "27", "casă", "0", `@CASĂ@ `+subsirDex09+` alt text`),
		def("4", "1", "vroi", "0", `@VROI@ #vb.# Formă veche pentru a voi.`),
		def("5", "88", "dejivra", "0", `@dejivra@ `+subsirDoom3+` verb de conjugarea I`),
		def("6", "7", "caz-limită", "0", `@CAZ-LIMITĂ@ #s. n.# Situație la limită, inventată.`),
		def("7", "7", "zi-lumină", "0", `@ZI-LUMINĂ@ #s. f.# Fără lexem legat.`),
		def("8", "1", "casă", "2", "ștearsă"),
		def("9", "45", "casă", "0", "din sursa ascunsă")))
	b.WriteString("-- Dump completed\n")
	return b.String()
}
