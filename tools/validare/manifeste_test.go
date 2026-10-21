package validare

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const descriereTest = "Scrie texte în română. Use when writing Romanian."

func scrie(t *testing.T, rad, cale, continut string) {
	t.Helper()
	p := filepath.Join(rad, cale)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(continut), 0o644); err != nil {
		t.Fatal(err)
	}
}

// radacinaBuna scrie un set complet și consecvent de manifeste.
func radacinaBuna(t *testing.T, nume, versiune string) string {
	rad := t.TempDir()
	scrie(t, rad, "skills/romana/SKILL.md", "---\nname: romana\ndescription: "+descriereTest+"\n---\nText.\n")
	plugin := `{"name":"` + nume + `","version":"` + versiune + `","description":"` + descriereTest + `"}`
	scrie(t, rad, ".claude-plugin/plugin.json", plugin)
	scrie(t, rad, "plugin.json", plugin)
	scrie(t, rad, ".claude-plugin/marketplace.json",
		`{"name":"romana","owner":{"name":"X"},"plugins":[{"name":"romana","source":"./","version":"0.1.0"}]}`)
	scrie(t, rad, ".agents/plugins/marketplace.json",
		`{"name":"romana","plugins":[{"name":"romana","source":{"source":"local","path":"./"},"policy":{"installation":"AVAILABLE","authentication":"ON_INSTALL"},"category":"Productivity"}]}`)
	return rad
}

func areProblema(pr []Problema, fragment string) bool {
	for _, p := range pr {
		if strings.Contains(p.Motiv, fragment) {
			return true
		}
	}
	return false
}

func TestManifesteConsecvente(t *testing.T) {
	if pr := Manifeste(radacinaBuna(t, "romana", "0.1.0")); len(pr) != 0 {
		t.Fatalf("set consecvent, dar: %v", pr)
	}

	rad := radacinaBuna(t, "alta", "0.1.0")
	if pr := Manifeste(rad); !areProblema(pr, "name") {
		t.Errorf("nume greșit nedetectat: %v", pr)
	}

	rad = radacinaBuna(t, "romana", "0.1.0")
	scrie(t, rad, "plugin.json", `{"name":"romana","version":"0.2.0","description":"`+descriereTest+`"}`)
	if pr := Manifeste(rad); !areProblema(pr, "version") {
		t.Errorf("versiune diferită nedetectată: %v", pr)
	}

	rad = radacinaBuna(t, "romana", "0.1.0")
	scrie(t, rad, "plugin.json", `{"name":"romana","version":"0.1.0","description":"alta"}`)
	if pr := Manifeste(rad); !areProblema(pr, "description") {
		t.Errorf("descriere diferită nedetectată: %v", pr)
	}

	rad = radacinaBuna(t, "romana", "0.1.0")
	scrie(t, rad, ".agents/plugins/marketplace.json",
		`{"name":"romana","plugins":[{"name":"romana","source":{"source":"local","path":"./nu-exista"}}]}`)
	if pr := Manifeste(rad); !areProblema(pr, "skills/romana") {
		t.Errorf("sursă care nu duce la skills/romana nedetectată: %v", pr)
	}

	rad = radacinaBuna(t, "romana", "0.1.0")
	scrie(t, rad, ".claude-plugin/marketplace.json",
		`{"name":"romana","owner":{"name":"X"},"plugins":[{"name":"romana","source":"../x"}]}`)
	if pr := Manifeste(rad); !areProblema(pr, "source") {
		t.Errorf("source invalid nedetectat: %v", pr)
	}

	rad = radacinaBuna(t, "romana", "0.1.0")
	os.Remove(filepath.Join(rad, "plugin.json"))
	if pr := Manifeste(rad); !areProblema(pr, "plugin.json") {
		t.Errorf("manifest lipsă nedetectat: %v", pr)
	}
}

// Manifestele depozitului real trebuie să fie consecvente.
func TestManifesteDepozit(t *testing.T) {
	if pr := Manifeste("../.."); len(pr) != 0 {
		t.Fatalf("manifeste: %v", pr)
	}
}
